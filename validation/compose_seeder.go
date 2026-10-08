package validation

import (
	"gopkg.in/yaml.v3"
	"strings"
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

func injectPlaceholders(composeMap map[string]any, dataMap map[string]string) ([]byte, error) {
	composeBytes, err := yaml.Marshal(composeMap)
	if err != nil {
		return nil, err
	}
	content := string(composeBytes)
	for name, value := range dataMap {
		content = strings.ReplaceAll(content, "${"+name+"}", value)
	}
	return []byte(content), nil
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

		services[serviceName] = serviceMap
	}
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
