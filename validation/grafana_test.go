package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
)

const grafanaNodeExporterCompose = `services:
  grafana:
    image: grafana/grafana:13.2.3
    container_name: quollix_grafana_grafana
    labels:
      quollix.port: 3000
  nodeexporter:
    image: quay.io/prometheus/node-exporter:v1.12.1
    container_name: quollix_grafana_nodeexporter
    pid: host
    volumes:
      - /:/host:ro,rslave
`

func TestComposeValidator_GrafanaWithoutNodeExporter(t *testing.T) {
	compose := parseComposeYml(t, grafanaNodeExporterCompose)
	delete(compose["services"].(map[string]any), "nodeexporter")
	assert.Nil(t, NewComposeValidator(&ServiceValidatorImpl{}).ValidateComposeMap(compose, "quollix", "grafana"))
}

func TestComposeValidator_GrafanaNodeExporter(t *testing.T) {
	for _, test := range []struct {
		name, old, replacement, maintainer, app string
		valid                                   bool
	}{
		{"expected definition", "", "", "quollix", "grafana", true},
		{"different version", "v1.12.1", "v1.12.2", "quollix", "grafana", true},
		{"optional digest", "v1.12.1", "v1.12.1@sha256:" + strings.Repeat("a", 64), "quollix", "grafana", true},
		{"other maintainer", "quollix_", "someone_", "someone", "grafana", false},
		{"other app", "grafana", "another", "quollix", "another", false},
		{"other service", "nodeexporter", "another", "quollix", "grafana", false},
		{"other repository", "quay.io/prometheus/node-exporter", "example.com/prometheus/node-exporter", "quollix", "grafana", false},
		{"repository suffix", "node-exporter:v", "node-exporter-evil:v", "quollix", "grafana", false},
		{"writable mount", ":ro,rslave", ":rw,rslave", "quollix", "grafana", false},
		{"other mount", "/:/host:ro,rslave", "/etc:/host:ro,rslave", "quollix", "grafana", false},
		{"additional mount", "      - /:/host:ro,rslave", "      - /:/host:ro,rslave\n      - /etc:/other:ro", "quollix", "grafana", false},
		{"missing mount", "    volumes:\n      - /:/host:ro,rslave\n", "", "quollix", "grafana", false},
		{"other pid", "pid: host", "pid: container:other", "quollix", "grafana", false},
		{"missing pid", "    pid: host\n", "", "quollix", "grafana", false},
		{"host networking", "pid: host", "pid: host\n    network_mode: host", "quollix", "grafana", false},
		{"other service host access", "  grafana:\n", "  grafana:\n    pid: host\n", "quollix", "grafana", false},
		{"other service host volume", "  grafana:\n", "  grafana:\n    volumes:\n      - /:/host:ro,rslave\n", "quollix", "grafana", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			content := strings.ReplaceAll(grafanaNodeExporterCompose, test.old, test.replacement)
			compose := parseComposeYml(t, content)
			err := NewComposeValidator(&ServiceValidatorImpl{}).ValidateComposeMap(compose, test.maintainer, test.app)
			assert.Equal(t, test.valid, err == nil)
			assert.Equal(t, parseComposeYml(t, content), compose)
		})
	}
}
