package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
	"github.com/quollix/deepstack"
)

func TestUnknownTopLevelKeyError(t *testing.T) {
	f := &ComposeValidatorImpl{}

	err := f.ValidateComposeMap(map[string]any{
		"services": map[string]any{},
		"extra":    1,
	}, "maint", "app")
	deepstack.AssertDeepStackError(t, err, notAllowedTopLevelKeyword, "key", "extra")
}

func TestGlobalVolumesHaveNoSubKeys(t *testing.T) {
	sv := NewServiceValidatorMock(t)
	f := &ComposeValidatorImpl{ServiceValidator: sv}
	err := f.ValidateComposeMap(map[string]any{
		"services": map[string]any{"app": map[string]any{}},
		"volumes": map[string]any{
			"v1": map[string]any{"driver": "local"},
		},
	}, "maint", "app")
	deepstack.AssertDeepStackError(t, err, globalVolumeShouldNotHaveSubKeywords)
}

func TestGlobalVolumesMustUseMaintainerAppPrefix(t *testing.T) {
	sv := NewServiceValidatorMock(t)
	f := &ComposeValidatorImpl{ServiceValidator: sv}

	err := f.ValidateComposeMap(map[string]any{
		"services": map[string]any{"app": map[string]any{}},
		"volumes": map[string]any{
			"wrongprefix_data": map[string]any{},
		},
	}, "maint", "app")

	deepstack.AssertDeepStackError(t, err, volumeNamePrefixIsWrong, VolumeNameField, "wrongprefix_data", ExpectedPrefixField, "maint_app_")
}

func TestGlobalVolumesMustHaveThreeParts(t *testing.T) {
	sv := NewServiceValidatorMock(t)
	f := &ComposeValidatorImpl{ServiceValidator: sv}

	err := f.ValidateComposeMap(map[string]any{
		"services": map[string]any{"app": map[string]any{}},
		"volumes": map[string]any{
			"maint_app_data_extra": map[string]any{},
		},
	}, "maint", "app")

	deepstack.AssertDeepStackError(t, err, unexpectedUnderscoreCountVolume, VolumeNameField, "maint_app_data_extra", ExpectedPartCount, 3, ActualPartCount, 4)
}

func TestGlobalVolumesMustUseValidSuffixRegex(t *testing.T) {
	sv := NewServiceValidatorMock(t)
	f := &ComposeValidatorImpl{ServiceValidator: sv}

	err := f.ValidateComposeMap(map[string]any{
		"services": map[string]any{"app": map[string]any{}},
		"volumes": map[string]any{
			"maint_app_data-extra": map[string]any{},
		},
	}, "maint", "app")

	assert.True(t, strings.HasPrefix(u.ExtractError(err), "Invalid input."))
}

func TestValidateComposeFile_MainServiceMustExist_CustomError(t *testing.T) {
	f := &ComposeValidatorImpl{}
	err := f.ValidateComposeMap(map[string]any{
		"services": map[string]any{
			"other": map[string]any{},
		},
	}, "maint", "app")
	deepstack.AssertDeepStackError(t, err, mainServiceMustBeDefined)
}

func TestValidateComposeFile_HappyPath(t *testing.T) {
	sv := NewServiceValidatorMock(t)
	cv := NewComposeConsistencyValidatorMock(t)
	defer sv.AssertExpectations(t)
	defer cv.AssertExpectations(t)
	f := &ComposeValidatorImpl{ServiceValidator: sv, ComposeConsistencyValidator: cv}

	service := map[string]any{
		"image":          "repo/app:1.0.0",
		"container_name": "maint_app_app",
		"ports":          []any{"8080:8080"},
		"volumes":        []any{"maint_app_vol:/data"},
		"networks":       []any{"default"},
		"deploy":         map[string]any{"resources": map[string]any{}},
	}
	services := map[string]any{
		"app": service,
	}
	compose := map[string]any{
		"services": services,
		"volumes": map[string]any{
			"maint_app_vol": map[string]any{},
		},
	}

	sv.EXPECT().ValidateServiceName("app").Return(nil)
	sv.EXPECT().ValidateServiceKeys(service).Return(nil)
	sv.EXPECT().ValidateImage("app", service).Return(nil)
	sv.EXPECT().ValidateContainerName("app", service, "maint", "app").Return(nil)
	sv.EXPECT().ValidatePorts(service).Return(nil)
	sv.EXPECT().ValidateServiceVolumes("maint", "app", "app", service).Return(nil)
	sv.EXPECT().ValidateDeploySection(service).Return(nil)
	sv.EXPECT().ValidateLabels("app", "app", service).Return(nil)
	sv.EXPECT().ValidateNoTzEnvironment(service).Return(nil)
	cv.EXPECT().ValidateVolumeMappings(compose).Return(nil)
	cv.EXPECT().ValidateServiceReferences(compose).Return(nil)

	assert.Nil(t, f.ValidateComposeMap(compose, "maint", "app"))
}

func TestValidateComposeMap_AddsServiceContextToServiceValidatorError(t *testing.T) {
	sv := NewServiceValidatorMock(t)
	defer sv.AssertExpectations(t)
	f := &ComposeValidatorImpl{ServiceValidator: sv}

	service := map[string]any{
		"image":          "repo/app:1.0.0",
		"container_name": "maint_app_app",
		"ports":          []any{"80:8080"},
	}
	compose := map[string]any{
		"services": map[string]any{
			"app": service,
		},
	}

	sv.EXPECT().ValidateServiceName("app").Return(nil)
	sv.EXPECT().ValidateServiceKeys(service).Return(nil)
	sv.EXPECT().ValidateImage("app", service).Return(nil)
	sv.EXPECT().ValidateContainerName("app", service, "maint", "app").Return(nil)
	sv.EXPECT().ValidatePorts(service).Return(u.Logger.NewError(exposingDefaultPortIsForbidden, PortField, "80"))

	err := f.ValidateComposeMap(compose, "maint", "app")

	deepstack.AssertDeepStackError(t, err, exposingDefaultPortIsForbidden, ServiceField, "app", PortField, "80")
}

func TestValidateComposePlaceholders_WithoutPlaceholders(t *testing.T) {
	f := &ComposeValidatorImpl{}
	err := f.ValidateComposePlaceholders([]byte(""))
	assert.Nil(t, err)
}

func TestValidateComposePlaceholders_WithAllowedPlaceholder(t *testing.T) {
	f := &ComposeValidatorImpl{}
	err := f.ValidateComposePlaceholders([]byte("${APP_SECRET}${BASE_DOMAIN}${SERVER_HOST}${CLIENT_ID}"))
	assert.Nil(t, err)
}

func TestValidateComposePlaceholders_WithAllowedSecretPlaceholder(t *testing.T) {
	f := &ComposeValidatorImpl{}
	err := f.ValidateComposePlaceholders([]byte("${SECRET_POSTGRES_PASSWORD}${SECRET_APP_2_TOKEN}"))
	assert.Nil(t, err)
}

func TestValidateComposePlaceholders_WithUnexpectedPlaceholder(t *testing.T) {
	f := &ComposeValidatorImpl{}
	err := f.ValidateComposePlaceholders([]byte("${CLIENTID}"))
	deepstack.AssertDeepStackError(
		t,
		err,
		"unsupported docker compose placeholder",
		"unsupported_placeholder", "CLIENTID",
		"allowed_placeholders", allowedComposePlaceholders,
	)
}

func TestValidateComposePlaceholders_WithInvalidSecretPlaceholder(t *testing.T) {
	f := &ComposeValidatorImpl{}
	err := f.ValidateComposePlaceholders([]byte("${SECRET_postgres_PASSWORD}"))
	deepstack.AssertDeepStackError(
		t,
		err,
		"unsupported docker compose placeholder",
		"unsupported_placeholder", "SECRET_postgres_PASSWORD",
		"allowed_placeholders", allowedComposePlaceholders,
	)
}
