package validation

import (
	"strings"

	u "github.com/quollix/common/utils"
)

type ComposeConsistencyValidator interface {
	ValidateVolumeMappings(compose map[string]any) error
	ValidateServiceReferences(compose map[string]any) error
}

type ComposeConsistencyValidatorImpl struct{}

var (
	serviceVolumeMustBeDeclaredGlobally = "service volume must be declared as global volume"
	globalVolumeMustBeMounted           = "global volume must be mounted by a service"
	duplicateVolumeTargetInService      = "duplicate volume target in service"
	dependsOnServiceMustExist           = "depends_on service must exist"
	invalidDependsOn                    = "invalid 'depends_on' in service"
)

func (v *ComposeConsistencyValidatorImpl) ValidateVolumeMappings(compose map[string]any) error {
	globalVolumes := extractGlobalVolumes(compose)
	mountedVolumes := map[string]bool{}

	services, ok := compose["services"].(map[string]any)
	if !ok {
		return nil
	}

	for serviceName, serviceValue := range services {
		serviceMap, ok := serviceValue.(map[string]any)
		if !ok {
			continue
		}

		volumeTargets := map[string]bool{}
		for _, volume := range extractServiceVolumeStrings(serviceMap) {
			source, target := parseVolumeSourceAndTarget(volume)
			mountedVolumes[source] = true

			if !globalVolumes[source] {
				return u.Logger.NewError(serviceVolumeMustBeDeclaredGlobally, ServiceField, serviceName, VolumeNameField, source)
			}
			if volumeTargets[target] {
				return u.Logger.NewError(duplicateVolumeTargetInService, ServiceField, serviceName, VolumeField, volume, "volume_target", target)
			}
			volumeTargets[target] = true
		}
	}

	for volumeName := range globalVolumes {
		if !mountedVolumes[volumeName] {
			return u.Logger.NewError(globalVolumeMustBeMounted, VolumeNameField, volumeName)
		}
	}

	return nil
}

func (v *ComposeConsistencyValidatorImpl) ValidateServiceReferences(compose map[string]any) error {
	services, ok := compose["services"].(map[string]any)
	if !ok {
		return nil
	}

	for serviceName, serviceValue := range services {
		serviceMap, ok := serviceValue.(map[string]any)
		if !ok {
			continue
		}

		dependsOn, ok := serviceMap["depends_on"]
		if !ok {
			continue
		}
		dependencies, err := extractDependsOnServices(serviceName, dependsOn)
		if err != nil {
			return err
		}
		for _, dependency := range dependencies {
			if _, exists := services[dependency]; !exists {
				return u.Logger.NewError(dependsOnServiceMustExist, ServiceField, serviceName, "depends_on_service", dependency)
			}
		}
	}

	return nil
}

func extractGlobalVolumes(compose map[string]any) map[string]bool {
	globalVolumes := map[string]bool{}
	volumes, ok := compose["volumes"].(map[string]any)
	if !ok {
		return globalVolumes
	}
	for volumeName := range volumes {
		globalVolumes[volumeName] = true
	}
	return globalVolumes
}

func extractServiceVolumeStrings(serviceMap map[string]any) []string {
	volumesAny, ok := serviceMap["volumes"].([]any)
	if ok {
		volumes := make([]string, 0, len(volumesAny))
		for _, volume := range volumesAny {
			volumeString, ok := volume.(string)
			if ok {
				volumes = append(volumes, volumeString)
			}
		}
		return volumes
	}

	volumesString, ok := serviceMap["volumes"].([]string)
	if ok {
		return volumesString
	}
	return nil
}

func parseVolumeSourceAndTarget(volume string) (string, string) {
	volumeParts := strings.SplitN(volume, ":", 3)
	if len(volumeParts) < 2 {
		return volume, ""
	}
	return volumeParts[0], volumeParts[1]
}

func extractDependsOnServices(serviceName string, dependsOn any) ([]string, error) {
	switch typedDependsOn := dependsOn.(type) {
	case []any:
		dependencies := make([]string, 0, len(typedDependsOn))
		for _, dependency := range typedDependsOn {
			dependencyString, ok := dependency.(string)
			if !ok {
				return nil, u.Logger.NewError(invalidDependsOn, ServiceField, serviceName)
			}
			dependencies = append(dependencies, dependencyString)
		}
		return dependencies, nil
	case []string:
		return typedDependsOn, nil
	default:
		return nil, u.Logger.NewError(invalidDependsOn, ServiceField, serviceName)
	}
}
