package bootstrap

import (
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

func TestGetInitialAdminCredentials_UsesEnvCredentials(t *testing.T) {
	t.Setenv(InitialAdminNameEnvVar, "admin")
	t.Setenv(InitialAdminPasswordEnvVar, "password123")

	username, password, err := GetInitialAdminCredentials("quollix")

	assert.Nil(t, err)
	assert.Equal(t, "admin", username)
	assert.Equal(t, "password123", password)
}

func TestGetInitialAdminCredentials_UsesDefaultUsername(t *testing.T) {
	t.Setenv(InitialAdminPasswordEnvVar, "password123")

	username, password, err := GetInitialAdminCredentials("quollix")

	assert.Nil(t, err)
	assert.Equal(t, "quollix", username)
	assert.Equal(t, "password123", password)
}

func TestGetInitialAdminCredentials_GeneratesPassword(t *testing.T) {
	username, password, err := GetInitialAdminCredentials("quollix")

	assert.Nil(t, err)
	assert.Equal(t, "quollix", username)
	assert.Equal(t, GeneratedInitialAdminPasswordSize, len(password))
}

func TestGetInitialAdminCredentials_RejectsInvalidEnvCredentials(t *testing.T) {
	t.Setenv(InitialAdminNameEnvVar, "Admin")
	t.Setenv(InitialAdminPasswordEnvVar, "short")

	username, password, err := GetInitialAdminCredentials("quollix")

	assert.Equal(t, "", username)
	assert.Equal(t, "", password)
	assert.Equal(t, "env variable is not valid", u.ExtractError(err))
}

func TestExtractGeneratedInitialAdminCredentials_UsesDocumentedSearchText(t *testing.T) {
	logs := `{"level":"info","msg":"INITIAL_ADMIN_PASSWORD environment variable is not set, generated random initial admin password","username":"quollix","password":"secret-password"}`

	username, password, err := ExtractGeneratedInitialAdminCredentials(logs)

	assert.Nil(t, err)
	assert.Equal(t, "quollix", username)
	assert.Equal(t, "secret-password", password)
}

func TestExtractGeneratedInitialAdminCredentials_IgnoresOtherPasswordFields(t *testing.T) {
	logs := `{"level":"info","msg":"some other log line","password":"wrong-password"}`

	username, password, err := ExtractGeneratedInitialAdminCredentials(logs)

	assert.Equal(t, "", username)
	assert.Equal(t, "", password)
	assert.Equal(t, "password log line was not found", u.ExtractError(err))
}
