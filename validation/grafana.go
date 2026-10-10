package validation

import (
	"maps"
	"regexp"
	"slices"

	u "github.com/quollix/common/utils"
)

var nodeExporterImageRegex = regexp.MustCompile(`^quay\.io/prometheus/node-exporter:[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}(@sha256:[a-fA-F0-9]{64})?$`)

func validateGrafanaNodeExporter(compose map[string]any, maintainer, app string) (map[string]any, error) {
	if maintainer != u.OfficialMaintainer || app != "grafana" {
		return compose, nil
	}
	services, _ := compose["services"].(map[string]any)
	exporter, ok := services["nodeexporter"].(map[string]any)
	if !ok {
		return compose, nil
	}
	image, _ := exporter["image"].(string)
	if !nodeExporterImageRegex.MatchString(image) {
		return nil, u.Logger.NewError("invalid Node Exporter image")
	}
	volumes, err := extractVolumesList("nodeexporter", exporter)
	if err != nil {
		return nil, err
	}
	if exporter["pid"] != "host" || !slices.Equal(volumes, []string{"/:/host:ro,rslave"}) {
		return nil, u.Logger.NewError("invalid Node Exporter host access")
	}
	exporter = maps.Clone(exporter)
	delete(exporter, "pid")
	delete(exporter, "volumes")
	services = maps.Clone(services)
	services["nodeexporter"] = exporter
	compose = maps.Clone(compose)
	compose["services"] = services
	return compose, nil
}
