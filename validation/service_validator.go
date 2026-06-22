package validation

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	u "github.com/quollix/common/utils"
)

var (
	allowedServiceKeys        = u.MapOf("image", "container_name", "ports", "volumes", "depends_on", "environment", "deploy", "tmpfs", "tty", "user", "command", "entrypoint", "labels")
	portsForbiddenToBeExposed = u.MapOf("22", "53", "80", "443")
)

type ServiceValidator interface {
	ValidateServiceKeys(serviceName string, serviceMap map[string]any) error
	ValidateImage(serviceName string, serviceMap map[string]any) error
	ValidateContainerName(serviceName string, serviceMap map[string]any, maintainerName, appName string) error
	ValidatePorts(serviceName string, serviceMap map[string]any) error
	ValidateServiceVolumes(maintainerName, appName, serviceName string, serviceMap map[string]any) error
	ValidateDeploySection(serviceName string, serviceMap map[string]any) error
	ValidateServiceName(name string) error
	ValidateUrl(url string) bool
	ValidateLabels(appName string, serviceName string, serviceMap map[string]any) error
	ValidateNoTzEnvironment(serviceName string, serviceMap map[string]any) error
}

type ServiceValidatorImpl struct{}

var (
	portLabelKey = "quollix.port"

	mainServiceLabelsMissing            = fmt.Sprintf("main service must define 'labels' including '%s'", portLabelKey)
	mainServiceLabelsMustBeMap          = "'labels' for main service must be a map"
	mainServicePortLabelMissing         = fmt.Sprintf("main service labels must include '%s'", portLabelKey)
	mainServicePortLabelMustBeString    = fmt.Sprintf("'%s' label value must be an integer", portLabelKey)
	mainServicePortLabelMustBeValidPort = fmt.Sprintf("'%s' label value must be an integer between 1 and 65535", portLabelKey)
	tzEnvironmentVariableIsForbidden    = "environment variable 'TZ' is forbidden"
	volumeTargetMustNotBeEmpty          = "volume target must not be empty"
	volumeTargetMustBeAbsolute          = "volume target must be an absolute container path"
	volumeTargetMustNotBeRoot           = "volume target must not be root"
)

func (s *ServiceValidatorImpl) ValidateNoTzEnvironment(serviceName string, serviceMap map[string]any) error {
	environmentValue, hasEnvironment := serviceMap["environment"]
	if !hasEnvironment || environmentValue == nil {
		return nil
	}

	switch typedEnvironment := environmentValue.(type) {
	case []any:
		for _, entry := range typedEnvironment {
			entryString, ok := entry.(string)
			if !ok {
				continue
			}
			if entryString == "TZ" || strings.HasPrefix(entryString, "TZ=") {
				return u.Logger.NewError(tzEnvironmentVariableIsForbidden, ServiceField, serviceName)
			}
		}
	case map[string]any:
		if _, hasTz := typedEnvironment["TZ"]; hasTz {
			return u.Logger.NewError(tzEnvironmentVariableIsForbidden, ServiceField, serviceName)
		}
	case map[any]any:
		if _, hasTz := typedEnvironment["TZ"]; hasTz {
			return u.Logger.NewError(tzEnvironmentVariableIsForbidden, ServiceField, serviceName)
		}
	}

	return nil
}

func (s *ServiceValidatorImpl) ValidateLabels(appName string, serviceName string, serviceMap map[string]any) error {
	if appName != serviceName {
		return nil
	}
	labelsValue, hasLabels := serviceMap["labels"]
	if !hasLabels {
		return u.Logger.NewError(mainServiceLabelsMissing)
	}
	labelsMap, ok := labelsValue.(map[string]any)
	if !ok {
		return u.Logger.NewError(mainServiceLabelsMustBeMap)
	}

	err := extractPort(labelsMap)
	if err != nil {
		return u.Logger.AddContext(err, LabelField, portLabelKey)
	}
	return nil
}

func extractPort(labelsMap map[string]any) error {
	labelValue, hasPort := labelsMap[portLabelKey]
	if !hasPort {
		return u.Logger.NewError(mainServicePortLabelMissing)
	}
	port, ok := labelValue.(int)
	if !ok {
		return u.Logger.NewError(mainServicePortLabelMustBeString)
	}
	if port < 1 || port > 65535 {
		return u.Logger.NewError(mainServicePortLabelMustBeValidPort)
	}
	return nil
}

func (s *ServiceValidatorImpl) ValidateServiceName(serviceName string) error {
	if err := Validate("Service", FieldDefault, serviceName); err != nil {
		return u.Logger.AddContext(err, ServiceField, serviceName)
	}
	return nil
}

func (s *ServiceValidatorImpl) ValidateServiceKeys(serviceName string, serviceMap map[string]any) error {
	for key := range serviceMap {
		_, ok := allowedServiceKeys[key]
		if !ok {
			return u.Logger.NewError(notAllowedKeyInService, ServiceField, serviceName, KeyField, key)
		}
	}
	return nil
}

func (s *ServiceValidatorImpl) ValidateImage(serviceName string, serviceMap map[string]any) error {
	img, ok := serviceMap["image"]
	if !ok {
		return u.Logger.NewError(mustSetImageKey, ServiceField, serviceName)
	}
	imgString, ok := img.(string)
	if !ok {
		return u.Logger.NewError("invalid 'image' in service", ServiceField, serviceName)
	}
	parts := strings.Split(imgString, ":")
	if len(parts) < 2 {
		return u.Logger.NewError(mustSetTheDockerImageTag, ServiceField, serviceName, ImageField, imgString)
	}
	if parts[1] == "latest" {
		return u.Logger.NewError(notAllowedLatestDockerImageTag, ServiceField, serviceName, ImageField, imgString)
	}
	return nil
}

func (s *ServiceValidatorImpl) ValidateContainerName(serviceName string, serviceMap map[string]any, maintainerName, appName string) error {
	actualContainerName, ok := serviceMap["container_name"]
	if !ok {
		return u.Logger.NewError(serviceNeedsContainerNameKeyword, ServiceField, serviceName)
	}
	actualContainerNameString, ok := actualContainerName.(string)
	if !ok {
		return u.Logger.NewError(wrongContainerNameValue, ServiceField, serviceName)
	}
	containerNameParts := strings.Split(actualContainerNameString, "_")
	if len(containerNameParts) != 3 {
		return u.Logger.NewError(
			unexpectedUnderscoreCountContainer,
			ServiceField, serviceName,
			ActualContainerName, actualContainerNameString,
			ExpectedPartCount, 3,
			ActualPartCount, len(containerNameParts),
		)
	}
	expectedContainerName := fmt.Sprintf("%s_%s_%s", maintainerName, appName, serviceName)
	if actualContainerNameString != expectedContainerName {
		return u.Logger.NewError(
			wrongContainerNameValue,
			ServiceField, serviceName,
			ActualContainerName, actualContainerNameString,
			ExpectedContainerName, expectedContainerName,
		)
	}
	return nil
}

func (s *ServiceValidatorImpl) ValidatePorts(serviceName string, serviceMap map[string]any) error {
	ports, ok := serviceMap["ports"]
	if !ok {
		return nil
	}
	portsList, ok := ports.([]any)
	if !ok {
		return u.Logger.NewError("invalid 'ports' in service", ServiceField, serviceName)
	}
	for _, port := range portsList {
		portString, ok := port.(string)
		if !ok {
			return u.Logger.NewError("invalid 'ports' in service", ServiceField, serviceName)
		}
		fields := strings.Split(portString, ":")
		portExposedOnHost := fields[0]
		_, isPortForbidden := portsForbiddenToBeExposed[portExposedOnHost]
		if isPortForbidden {
			return u.Logger.NewError(exposingDefaultPortIsForbidden, ServiceField, serviceName, PortField, portExposedOnHost)
		}
	}
	return nil
}

func (s *ServiceValidatorImpl) ValidateServiceVolumes(maintainerName, appName, serviceName string, serviceMap map[string]any) error {
	volumesList, err := extractVolumesList(serviceName, serviceMap)
	if err != nil {
		return err
	}
	for _, vol := range volumesList {
		if err := validateVolumeEntry(maintainerName, appName, serviceName, vol); err != nil {
			return err
		}
	}
	return nil
}

func extractVolumesList(serviceName string, serviceMap map[string]any) ([]string, error) {
	volumesValue := serviceMap["volumes"]
	if volumesValue == nil {
		return nil, nil
	}

	if volumesStringSlice, ok := volumesValue.([]string); ok {
		return volumesStringSlice, nil
	}

	volumesAnySlice, ok := volumesValue.([]any)
	if !ok {
		return nil, u.Logger.NewError("invalid 'volumes' in service", ServiceField, serviceName)
	}

	var volumesList []string
	for _, volumeAny := range volumesAnySlice {
		volumeString, ok := volumeAny.(string)
		if !ok {
			return nil, u.Logger.NewError("invalid 'volumes' in service", ServiceField, serviceName)
		}
		volumesList = append(volumesList, volumeString)
	}

	return volumesList, nil
}

func validateVolumeEntry(maintainerName, appName, serviceName, vol string) error {
	parts := strings.Split(vol, ":")
	if len(parts) < 2 {
		return u.Logger.NewError(volumeEntryMissingColonSeparator, ServiceField, serviceName, VolumeField, vol)
	}
	volumeName := parts[0]
	volumeTarget := parts[1]
	if filepath.IsAbs(volumeName) || strings.HasPrefix(volumeName, ".") {
		return u.Logger.NewError(hostDirectoriesMountedForbidden, ServiceField, serviceName)
	}
	if err := validateVolumeTarget(serviceName, vol, volumeTarget); err != nil {
		return err
	}
	if err := validateNamedVolumeName(maintainerName, appName, volumeName); err != nil {
		return u.Logger.AddContext(err, ServiceField, serviceName, VolumeField, vol)
	}
	return nil
}

func validateVolumeTarget(serviceName, volume, target string) error {
	if target == "" {
		return u.Logger.NewError(volumeTargetMustNotBeEmpty, ServiceField, serviceName, VolumeField, volume)
	}
	if !strings.HasPrefix(target, "/") {
		return u.Logger.NewError(volumeTargetMustBeAbsolute, ServiceField, serviceName, VolumeField, volume)
	}
	if target == "/" {
		return u.Logger.NewError(volumeTargetMustNotBeRoot, ServiceField, serviceName, VolumeField, volume)
	}
	return nil
}

func validateNamedVolumeName(maintainerName, appName, volumeName string) error {
	prefix := fmt.Sprintf("%s_%s_", maintainerName, appName)
	if !strings.HasPrefix(volumeName, prefix) {
		return u.Logger.NewError(volumeNamePrefixIsWrong, VolumeNameField, volumeName, ExpectedPrefixField, prefix)
	}
	volumeParts := strings.Split(volumeName, "_")
	if len(volumeParts) != 3 {
		return u.Logger.NewError(
			unexpectedUnderscoreCountVolume,
			VolumeNameField, volumeName,
			ExpectedPartCount, 3,
			ActualPartCount, len(volumeParts),
		)
	}
	if err := Validate("Volume", FieldDefault, volumeParts[2]); err != nil {
		return u.Logger.AddContext(err, VolumeNameField, volumeName)
	}
	return nil
}

func (f *ServiceValidatorImpl) ValidateDeploySection(serviceName string, serviceMap map[string]any) error {
	deploy, has := serviceMap["deploy"]
	if !has {
		return nil
	}
	deployMap, ok := deploy.(map[string]any)
	if !ok {
		return fmt.Errorf("'deploy' keyword in service '%s' must be a map", serviceName)
	}
	for subKey := range deployMap {
		if subKey != "resources" {
			return u.Logger.NewError(deployKeywordMustOnlyContainResources, ServiceField, serviceName, KeyField, subKey)
		}
	}
	resourcesMap, ok := deployMap["resources"].(map[string]any)
	if !ok {
		return nil
	}
	reservations, ok := resourcesMap["reservations"]
	if !ok {
		return nil
	}
	reservationsMap, ok := reservations.(map[string]any)
	if !ok {
		return fmt.Errorf("'reservations' in 'resources' of service '%s' must be a map", serviceName)
	}
	if _, ok = reservationsMap["devices"]; ok {
		return u.Logger.NewError(devicesKeywordIsForbidden, ServiceField, serviceName)
	}
	return nil
}

func (f *ServiceValidatorImpl) ValidateUrl(rawUrl string) bool {
	parsedUrl, parseError := url.Parse(rawUrl)
	if parseError != nil {
		return false
	}
	if parsedUrl.Scheme != "https" {
		return false
	}
	if parsedUrl.Host == "" {
		return false
	}
	if parsedUrl.RawQuery != "" || parsedUrl.Fragment != "" {
		return false
	}
	if parsedUrl.User != nil {
		return false
	}
	return true
}
