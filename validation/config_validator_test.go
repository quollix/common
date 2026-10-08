package validation

import (
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
	"gopkg.in/yaml.v3"
)

func TestValidateComposeMap_InlineConfigs(t *testing.T) {
	const definition = `
services:
  grafana:
    image: grafana/grafana:13.2.3
    container_name: quollix_grafana_grafana
    labels:
      quollix.port: 3000
    configs:
      - source: datasources
        target: /etc/grafana/provisioning/datasources/victoria.yaml
configs:
  datasources:
    content: |
      url: http://victorialogs:9428
      password: ${SECRET_GRAFANA}
`
	var compose map[string]any
	assert.Nil(t, yaml.Unmarshal([]byte(definition), &compose))
	assert.Nil(t, NewComposeValidator(&ServiceValidatorImpl{}).ValidateComposeMap(compose, "quollix", "grafana"))
}

func TestValidateConfigs_RejectsUnsafeDefinitions(t *testing.T) {
	for _, definition := range []any{
		nil,
		map[string]any{"file": "/etc/passwd"},
		map[string]any{"environment": "HOST_SECRET"},
		map[string]any{"external": true},
		map[string]any{"content": "safe", "name": "other_app"},
		map[string]any{"content": "safe", "file": "/etc/passwd"},
		map[string]any{"content": 42},
	} {
		assert.Equal(t, "config must contain only inline content", u.ExtractError(validateConfigs(map[string]any{"configs": map[string]any{"example": definition}})))
	}
	assert.Equal(t, "configs must be a map", u.ExtractError(validateConfigs(map[string]any{"configs": "invalid"})))
}

func TestValidateConfigs_RejectsInvalidReferences(t *testing.T) {
	for _, tc := range []struct {
		references any
		message    string
	}{
		{"example", "service configs must be a list"},
		{[]any{"example"}, "service config must contain only source and target"},
		{[]any{map[string]any{"source": "example", "target": "/config", "uid": "0"}}, "service config must contain only source and target"},
		{[]any{map[string]any{"source": "missing", "target": "/config"}}, "service config source must be declared globally"},
		{[]any{map[string]any{"source": "example", "target": "relative"}}, "config target must be an absolute container file path"},
		{[]any{map[string]any{"source": "example", "target": "/"}}, "config target must be an absolute container file path"},
		{[]any{map[string]any{"source": "example", "target": "/directory/"}}, "config target must be an absolute container file path"},
		{[]any{map[string]any{"source": "example", "target": "/$HOST_SECRET"}}, "service config requires literal source and target"},
		{[]any{map[string]any{"source": "example", "target": "/config"}, map[string]any{"source": "example", "target": "/dir/../config"}}, "duplicate config target in service"},
	} {
		assert.Equal(t, tc.message, u.ExtractError(validateConfigs(map[string]any{
			"configs":  map[string]any{"example": map[string]any{"content": "safe"}},
			"services": map[string]any{"app": map[string]any{"configs": tc.references}},
		})))
	}
}
