package validation

import (
	"maps"
	"regexp"
	"slices"

	u "github.com/quollix/common/utils"
)

var vectorImageRegex = regexp.MustCompile(`^timberio/vector:[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}(@sha256:[a-fA-F0-9]{64})?$`)

func validateGrafanaVector(compose map[string]any, maintainer, app string) (map[string]any, error) {
	if maintainer != u.OfficialMaintainer || app != "grafana" {
		return compose, nil
	}
	services, _ := compose["services"].(map[string]any)
	collector, ok := services["vector"].(map[string]any)
	if !ok {
		return compose, nil
	}
	image, _ := collector["image"].(string)
	if !vectorImageRegex.MatchString(image) {
		return nil, u.Logger.NewError("invalid Vector image")
	}
	volumes, err := extractVolumesList("vector", collector)
	if err != nil {
		return nil, err
	}
	if !slices.Equal(volumes, []string{
		"/:/host:ro,rslave",
		"/var/log/journal:/var/log/journal:ro",
		"/run/log/journal:/run/log/journal:ro",
		"quollix_grafana_vector:/var/lib/vector",
	}) {
		return nil, u.Logger.NewError("invalid Vector host access")
	}
	collector = maps.Clone(collector)
	collector["volumes"] = []string{"quollix_grafana_vector:/var/lib/vector"}
	services = maps.Clone(services)
	services["vector"] = collector
	compose = maps.Clone(compose)
	compose["services"] = services
	return compose, nil
}
