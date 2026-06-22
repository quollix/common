package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

type testDependencies struct {
	versionValidator     VersionValidator
	composeCheckerMock   *u.ComposeSyntaxCheckerMock
	composeValidatorMock *ComposeValidatorMock
}

func setupTestDependencies(t *testing.T) *testDependencies {
	composeChecker := u.NewComposeSyntaxCheckerMock(t)
	composeValidator := NewComposeValidatorMock(t)
	return &testDependencies{
		versionValidator:     NewVersionValidatorWithDependencies(composeChecker, composeValidator),
		composeCheckerMock:   composeChecker,
		composeValidatorMock: composeValidator,
	}
}

func assertAllExpectations(t *testing.T, d *testDependencies) {
	d.composeValidatorMock.AssertExpectations(t)
}

func TestValidate_DeniesReservedNames(t *testing.T) {
	deps := setupTestDependencies(t)
	defer assertAllExpectations(t, deps)

	err := deps.versionValidator.Validate([]byte("x"), "maintainer", u.OfficialMaintainer)
	assert.Equal(t, SystemAppNamesAreAlreadyReserved, u.ExtractError(err))

	err = deps.versionValidator.Validate([]byte("x"), "maintainer", u.OfficialDatabaseAppName)
	assert.Equal(t, SystemAppNamesAreAlreadyReserved, u.ExtractError(err))
}

func TestValidate_DeniesInvalidMaintainerName(t *testing.T) {
	deps := setupTestDependencies(t)
	defer assertAllExpectations(t, deps)

	err := deps.versionValidator.Validate([]byte("x"), "invalid_name", "app")
	assert.True(t, strings.Contains(u.ExtractError(err), "Invalid input."))
}

func TestValidate_DeniesInvalidAppName(t *testing.T) {
	deps := setupTestDependencies(t)
	defer assertAllExpectations(t, deps)

	err := deps.versionValidator.Validate([]byte("x"), "maintainer", "invalid_name")
	assert.True(t, strings.Contains(u.ExtractError(err), "Invalid input."))
}

func TestValidate_SuccessfulValidation(t *testing.T) {
	deps := setupTestDependencies(t)
	defer assertAllExpectations(t, deps)
	composeFileBytes := []byte("services: {}\n")
	expectedComposeMap := map[string]any{"services": map[string]any{}}

	deps.composeValidatorMock.EXPECT().ValidateComposePlaceholders(composeFileBytes).Return(nil)
	deps.composeValidatorMock.EXPECT().ValidateComposeMap(expectedComposeMap, "maintainer", "app").Return(nil)
	deps.composeCheckerMock.EXPECT().CheckDockerComposeSyntax(composeFileBytes).Return(nil)

	err := deps.versionValidator.Validate(composeFileBytes, "maintainer", "app")
	assert.Nil(t, err)
}
