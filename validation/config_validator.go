package validation

import (
	"path"
	"strings"

	u "github.com/quollix/common/utils"
)

func validateConfigs(compose map[string]any) error {
	configs := map[string]any{}
	if value, exists := compose["configs"]; exists {
		var ok bool
		configs, ok = value.(map[string]any)
		if !ok {
			return u.Logger.NewError("configs must be a map")
		}
	}
	for name, value := range configs {
		definition, ok := value.(map[string]any)
		if !ok || len(definition) != 1 {
			return u.Logger.NewError("config must contain only inline content", "config", name)
		}
		_, ok = definition["content"].(string)
		if !ok {
			return u.Logger.NewError("config must contain only inline content", "config", name)
		}
	}
	services, _ := compose["services"].(map[string]any)
	for serviceName, value := range services {
		service, _ := value.(map[string]any)
		if err := validateServiceConfigs(service, configs); err != nil {
			return u.Logger.AddContext(err, ServiceField, serviceName)
		}
	}
	return nil
}

func validateServiceConfigs(service, configs map[string]any) error {
	value, exists := service["configs"]
	if !exists {
		return nil
	}
	references, ok := value.([]any)
	if !ok {
		return u.Logger.NewError("service configs must be a list")
	}
	targets := map[string]bool{}
	for _, value := range references {
		reference, ok := value.(map[string]any)
		if !ok || len(reference) != 2 {
			return u.Logger.NewError("service config must contain only source and target")
		}
		source, sourceOK := reference["source"].(string)
		target, targetOK := reference["target"].(string)
		if !sourceOK || !targetOK || strings.ContainsAny(source+target, "$\x00") {
			return u.Logger.NewError("service config requires literal source and target")
		}
		if _, exists := configs[source]; !exists {
			return u.Logger.NewError("service config source must be declared globally", "config", source)
		}
		if !path.IsAbs(target) || path.Clean(target) == "/" || strings.HasSuffix(target, "/") {
			return u.Logger.NewError("config target must be an absolute container file path", "target", target)
		}
		target = path.Clean(target)
		if targets[target] {
			return u.Logger.NewError("duplicate config target in service", "target", target)
		}
		targets[target] = true
	}
	return nil
}
