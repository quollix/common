package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

var (
	versionValidator = NewVersionValidator(false)

	appName        = "gitea"
	maintainerName = "samplemaintainer"
)

func TestInvalidComposeBytes(t *testing.T) {
	composeFileBytes := []byte("hello")
	err := versionValidator.Validate(composeFileBytes, maintainerName, appName)
	errorString := u.ExtractError(err)
	assert.True(t, strings.Contains(errorString, "yaml: unmarshal errors"))
}

func TestOfficialBranAppIsDenied(t *testing.T) {
	composeFileBytes := []byte("hello")
	err := versionValidator.Validate(composeFileBytes, maintainerName, u.OfficialBrandAppName)
	assert.Equal(t, SystemAppNamesAreAlreadyReserved, u.ExtractError(err))

	err = versionValidator.Validate(composeFileBytes, maintainerName, u.OfficialDatabaseAppName)
	assert.Equal(t, SystemAppNamesAreAlreadyReserved, u.ExtractError(err))
}

func validateComposeFile(composeFileBytes []byte) string {
	err := versionValidator.Validate(composeFileBytes, maintainerName, appName)
	return u.ExtractError(err)
}

func TestValidationSuccess(t *testing.T) {
	compose := `
services:
  gitea:
    image: gitea/gitea:1.20.2
    container_name: samplemaintainer_gitea_gitea
    labels:
      quollix.port: 3000
`
	actualError := validateComposeFile([]byte(compose))
	assert.Equal(t, "", actualError)
}

func TestValidationWrongContainerName(t *testing.T) {
	compose := `
services:
  gitea:
    image: gitea/gitea:1.20.2
    container_name: wrongadmin_gitea_gitea
`
	actualError := validateComposeFile([]byte(compose))
	assert.Equal(t, wrongContainerNameValue, actualError)
}
