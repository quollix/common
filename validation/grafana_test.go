package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
)

const grafanaVectorCompose = `services:
  grafana:
    image: grafana/grafana:13.2.3
    container_name: quollix_grafana_grafana
    labels:
      quollix.port: 3000
  vector:
    image: timberio/vector:0.59.0-debian
    container_name: quollix_grafana_vector
    volumes:
      - /:/host:ro,rslave
      - /var/log/journal:/var/log/journal:ro
      - /run/log/journal:/run/log/journal:ro
      - quollix_grafana_vector:/var/lib/vector
volumes:
  quollix_grafana_vector: null
`

func TestComposeValidator_GrafanaWithoutVector(t *testing.T) {
	compose := parseComposeYml(t, grafanaVectorCompose)
	delete(compose["services"].(map[string]any), "vector")
	delete(compose, "volumes")
	assert.Nil(t, NewComposeValidator(&ServiceValidatorImpl{}).ValidateComposeMap(compose, "quollix", "grafana"))
}

func TestComposeValidator_GrafanaVector(t *testing.T) {
	for _, test := range []struct {
		name, old, replacement, maintainer, app string
		valid                                   bool
	}{
		{"expected definition", "", "", "quollix", "grafana", true},
		{"different version", "0.59.0-debian", "0.59.1-debian", "quollix", "grafana", true},
		{"optional digest", "0.59.0-debian", "0.59.0-debian@sha256:" + strings.Repeat("a", 64), "quollix", "grafana", true},
		{"other maintainer", "quollix_", "someone_", "someone", "grafana", false},
		{"other app", "grafana", "another", "quollix", "another", false},
		{"other service", "vector", "another", "quollix", "grafana", false},
		{"other repository", "timberio/vector", "example.com/timberio/vector", "quollix", "grafana", false},
		{"repository suffix", "vector:0", "vector-evil:0", "quollix", "grafana", false},
		{"writable mount", ":ro,rslave", ":rw,rslave", "quollix", "grafana", false},
		{"other mount", "/:/host:ro,rslave", "/etc:/host:ro,rslave", "quollix", "grafana", false},
		{"additional mount", "      - /:/host:ro,rslave", "      - /:/host:ro,rslave\n      - /etc:/other:ro", "quollix", "grafana", false},
		{"missing mount", "      - /:/host:ro,rslave\n", "", "quollix", "grafana", false},
		{"host pid", "  vector:\n", "  vector:\n    pid: host\n", "quollix", "grafana", false},
		{"writable journal", "/var/log/journal:ro", "/var/log/journal:rw", "quollix", "grafana", false},
		{"missing journal", "      - /run/log/journal:/run/log/journal:ro\n", "", "quollix", "grafana", false},
		{"undeclared checkpoint volume", "volumes:\n  quollix_grafana_vector: null\n", "", "quollix", "grafana", false},
		{"host networking", "  vector:\n", "  vector:\n    network_mode: host\n", "quollix", "grafana", false},
		{"other service host access", "  grafana:\n", "  grafana:\n    pid: host\n", "quollix", "grafana", false},
		{"other service host volume", "  grafana:\n", "  grafana:\n    volumes:\n      - /:/host:ro,rslave\n", "quollix", "grafana", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			content := strings.ReplaceAll(grafanaVectorCompose, test.old, test.replacement)
			compose := parseComposeYml(t, content)
			err := NewComposeValidator(&ServiceValidatorImpl{}).ValidateComposeMap(compose, test.maintainer, test.app)
			assert.Equal(t, test.valid, err == nil)
			assert.Equal(t, parseComposeYml(t, content), compose)
		})
	}
}
