package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
	"gopkg.in/yaml.v3"
)

const wgEasyTestImage = "ghcr.io/wg-easy/wg-easy:15.3"

var wgEasyValidator = NewVersionValidator(false)

func TestVersionValidator_WgEasyAcceptsExpectedDefinition(t *testing.T) {
	assert.Nil(t, wgEasyValidator.Validate(wgEasyComposeWithImage(wgEasyTestImage), u.OfficialMaintainer, wgEasyAppName))
}

func TestVersionValidator_WgEasyAcceptsEquivalentYamlFormatting(t *testing.T) {
	composeMap := map[string]any{}
	err := yaml.Unmarshal([]byte(wgEasyComposeYaml), &composeMap)
	assert.Nil(t, err)
	composeMap["services"].(map[string]any)[wgEasyAppName].(map[string]any)["image"] = wgEasyTestImage
	reformattedContent, err := yaml.Marshal(composeMap)
	assert.Nil(t, err)

	assert.Nil(t, wgEasyValidator.Validate(reformattedContent, u.OfficialMaintainer, wgEasyAppName))
}

func TestVersionValidator_WgEasyAcceptsPinnedDigest(t *testing.T) {
	digest := strings.Repeat("a", 64)
	assert.Nil(t, wgEasyValidator.Validate(wgEasyComposeWithImage(wgEasyTestImage+"@sha256:"+digest), u.OfficialMaintainer, wgEasyAppName))
}

func TestVersionValidator_WgEasyRejectsLatestTag(t *testing.T) {
	assert.NotNil(t, wgEasyValidator.Validate(wgEasyComposeWithImage("ghcr.io/wg-easy/wg-easy:latest"), u.OfficialMaintainer, wgEasyAppName))
}

func TestVersionValidator_WgEasyRejectsDifferentImage(t *testing.T) {
	assert.NotNil(t, wgEasyValidator.Validate(wgEasyComposeWithImage("example.com/wg-easy:15.3"), u.OfficialMaintainer, wgEasyAppName))
}

func TestVersionValidator_WgEasyRejectsChangedDefinition(t *testing.T) {
	content := wgEasyComposeWithImage(wgEasyTestImage)
	content = []byte(strings.Replace(string(content), "      - NET_ADMIN\n", "      - NET_ADMIN\n      - SYS_ADMIN\n", 1))

	assert.NotNil(t, wgEasyValidator.Validate(content, u.OfficialMaintainer, wgEasyAppName))
}

func wgEasyComposeWithImage(image string) []byte {
	return []byte(strings.ReplaceAll(wgEasyComposeYaml, wgEasyImagePlaceholder, image))
}
