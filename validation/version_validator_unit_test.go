package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
)

type testDependencies struct {
	versionValidator           *VersionValidatorImpl
	composeValidatorMock       *ComposeValidatorMock
	composeSyntaxValidatorMock *ComposeSyntaxValidatorMock
}

func setupTestDependencies(t *testing.T) *testDependencies {
	composeValidator := NewComposeValidatorMock(t)
	composeSyntaxValidator := NewComposeSyntaxValidatorMock(t)
	return &testDependencies{
		versionValidator: &VersionValidatorImpl{
			ComposeValidator:       composeValidator,
			ComposeSyntaxValidator: composeSyntaxValidator,
		},
		composeValidatorMock:       composeValidator,
		composeSyntaxValidatorMock: composeSyntaxValidator,
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
	deps.composeSyntaxValidatorMock.EXPECT().ValidateComposeSyntax(composeFileBytes).Return(nil)

	err := deps.versionValidator.Validate(composeFileBytes, "maintainer", "app")
	assert.Nil(t, err)
}

func TestValidate_StrictLicenseNoticeMissingReturnsError(t *testing.T) {
	deps := setupTestDependencies(t)
	defer assertAllExpectations(t, deps)
	deps.versionValidator.RequireLicenseNotice = true

	err := deps.versionValidator.Validate([]byte("services: {}\n"), "maintainer", "app")
	assert.Equal(t, "app definition must start with: # This file is licensed under the 0BSD License: https://opensource.org/license/0bsd\\n", u.ExtractError(err))
}

func TestValidate_StrictLicenseNoticeSuccess(t *testing.T) {
	deps := setupTestDependencies(t)
	defer assertAllExpectations(t, deps)
	deps.versionValidator.RequireLicenseNotice = true
	composeFileBytes := []byte("# This file is licensed under the 0BSD License: https://opensource.org/license/0bsd\nservices: {}\n")
	expectedComposeMap := map[string]any{"services": map[string]any{}}

	deps.composeValidatorMock.EXPECT().ValidateComposePlaceholders(composeFileBytes).Return(nil)
	deps.composeValidatorMock.EXPECT().ValidateComposeMap(expectedComposeMap, "maintainer", "app").Return(nil)
	deps.composeSyntaxValidatorMock.EXPECT().ValidateComposeSyntax(composeFileBytes).Return(nil)

	err := deps.versionValidator.Validate(composeFileBytes, "maintainer", "app")
	assert.Nil(t, err)
}
