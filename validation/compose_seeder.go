package validation

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func CompleteDockerComposeYaml(maintainer, appName string, inputCompose []byte, dataMap map[string]string) ([]byte, error) {
	var composeMap map[string]any
	if err := yaml.Unmarshal(inputCompose, &composeMap); err != nil {
		return nil, err
	}

	networkName := maintainer + "_" + appName

	addExternalNetwork(composeMap, networkName)
	updateServices(composeMap, networkName)
	updateVolumes(composeMap)

	return injectPlaceholders(composeMap, dataMap)
}

func addExternalNetwork(composeMap map[string]any, networkName string) {
	composeMap["networks"] = map[string]any{
		networkName: map[string]any{"external": true},
	}
}

func updateServices(composeMap map[string]any, networkName string) {
	services, _ := composeMap["services"].(map[string]any)
	for serviceName, serviceValue := range services {
		serviceMap, _ := serviceValue.(map[string]any)
		if serviceMap == nil {
			serviceMap = make(map[string]any)
		}

		serviceMap["networks"] = []any{networkName}
		serviceMap["restart"] = "unless-stopped"
		serviceMap["cap_drop"] = []any{"ALL"}
		serviceMap["cap_add"] = []any{
			"CAP_NET_BIND_SERVICE", "CAP_CHOWN", "CAP_FOWNER",
			"CAP_SETGID", "CAP_SETUID", "CAP_DAC_OVERRIDE",
		}

		ensureTzEnvironment(serviceMap)

		services[serviceName] = serviceMap
	}
}

func ensureTzEnvironment(serviceMap map[string]any) {
	const tzEntry = "TZ=${IANA_TIMEZONE}"

	environmentValue, exists := serviceMap["environment"]
	if !exists || environmentValue == nil {
		serviceMap["environment"] = []any{tzEntry}
		return
	}

	switch typedEnvironment := environmentValue.(type) {
	case []any:
		serviceMap["environment"] = ensureStringSliceEntry(typedEnvironment, tzEntry)
	case map[string]any:
		if _, hasTz := typedEnvironment["TZ"]; !hasTz {
			typedEnvironment["TZ"] = "${IANA_TIMEZONE}"
		}
		serviceMap["environment"] = typedEnvironment
	case map[any]any:
		if _, hasTz := typedEnvironment["TZ"]; !hasTz {
			typedEnvironment["TZ"] = "${IANA_TIMEZONE}"
		}
		serviceMap["environment"] = typedEnvironment
	default:
		serviceMap["environment"] = []any{tzEntry}
	}
}

func ensureStringSliceEntry(values []any, newEntry string) []any {
	for _, value := range values {
		if valueString, ok := value.(string); ok && valueString == newEntry {
			return values
		}
	}
	return append(values, newEntry)
}

func updateVolumes(composeMap map[string]any) {
	volumes, _ := composeMap["volumes"].(map[string]any)
	for volumeName, volumeValue := range volumes {
		volumeMap, _ := volumeValue.(map[string]any)
		if volumeMap == nil {
			volumeMap = make(map[string]any)
		}
		volumeMap["name"] = volumeName
		volumes[volumeName] = volumeMap
	}
}

func injectPlaceholders(composeMap map[string]any, dataMap map[string]string) ([]byte, error) {
	composeBytes, err := yaml.Marshal(composeMap)
	if err != nil {
		return nil, err
	}

	content := string(composeBytes)
	for placeholderKey, replacementValue := range dataMap {
		content = strings.ReplaceAll(content, "${"+placeholderKey+"}", replacementValue)
	}
	return []byte(content), nil
}
