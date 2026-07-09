package frontend

import (
	"testing"

	"github.com/quollix/common/assert"
	"github.com/quollix/deepstack"
)

func TestConfig_Resolve_RequiresFrontendFolderPath(t *testing.T) {
	_, err := Config{}.resolve()

	deepstack.AssertDeepStackError(t, err, "frontend folder path must not be empty")
}

func TestConfig_Resolve_ReturnsErrorForMissingResourceFolderPath(t *testing.T) {
	_, err := Config{FrontendFolderPath: "testdata/missing"}.resolve()

	assert.NotNil(t, err)
	deepstack.AssertDeepStackError(t, err, "stat testdata/missing/resources: no such file or directory", "resource_folder_path", "testdata/missing/resources")
}

func TestConfig_Resolve_ReturnsErrorWhenResourceFolderPathIsNotDirectory(t *testing.T) {
	_, err := Config{FrontendFolderPath: "testdata/frontend-with-resource-file"}.resolve()

	deepstack.AssertDeepStackError(t, err, "resource folder path must be a directory", "resource_folder_path", "testdata/frontend-with-resource-file/resources")
}

func TestConfig_Resolve_UsesSpecificFrontendFolderPath(t *testing.T) {
	resolvedConfig, err := Config{FrontendFolderPath: "testdata/frontend/"}.resolve()

	assert.Nil(t, err)
	assert.Equal(t, "/frontend/resources", resolvedConfig.publicResourcePath)
	assert.Equal(t, framePath, resolvedConfig.framePath)
}
