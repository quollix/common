package validation

import (
	"strings"
	"testing"

	"github.com/quollix/common/assert"
	u "github.com/quollix/common/utils"
	"github.com/quollix/deepstack"
)

func TestValidateServiceKeys(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}
	err := serviceValidator.ValidateServiceKeys("svc", map[string]any{
		"image": "repo:1.0.0",
	})
	assert.Nil(t, err)

	err = serviceValidator.ValidateServiceKeys("svc", map[string]any{"privileged": true})
	assert.Equal(t, notAllowedKeyInService, u.ExtractError(err))

	err = serviceValidator.ValidateServiceKeys("svc", map[string]any{
		"image": "repo:1.0.0",
		"gpus":  "all",
	})
	assert.Equal(t, notAllowedKeyInService, u.ExtractError(err))
}

func TestValidateImage(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}
	err := serviceValidator.ValidateImage("svc", map[string]any{"image": "repo:1.2.3"})
	assert.Nil(t, err)

	err = serviceValidator.ValidateImage("svc", map[string]any{"image": "repo:latest"})
	deepstack.AssertDeepStackError(t, err, notAllowedLatestDockerImageTag, ServiceField, "svc", ImageField, "repo:latest")

	err = serviceValidator.ValidateImage("svc", map[string]any{"image": "repo"})
	deepstack.AssertDeepStackError(t, err, mustSetTheDockerImageTag, ServiceField, "svc", ImageField, "repo")

	err = serviceValidator.ValidateImage("svc", map[string]any{"image": map[string]any{"name": "repo:1.2.3"}})
	deepstack.AssertDeepStackError(t, err, "invalid 'image' in service", ServiceField, "svc")

	err = serviceValidator.ValidateImage("svc", map[string]any{})
	deepstack.AssertDeepStackError(t, err, mustSetImageKey, ServiceField, "svc")
}

func TestValidateContainerName(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}
	err := serviceValidator.ValidateContainerName("svc", map[string]any{
		"container_name": "maint_app_svc",
	}, "maint", "app")
	assert.Nil(t, err)

	err = serviceValidator.ValidateContainerName("svc", map[string]any{}, "maint", "app")
	deepstack.AssertDeepStackError(t, err, serviceNeedsContainerNameKeyword, ServiceField, "svc")

	err = serviceValidator.ValidateContainerName("svc", map[string]any{
		"container_name": "maint_app_other",
	}, "maint", "app")
	deepstack.AssertDeepStackError(t, err, wrongContainerNameValue, ServiceField, "svc", ActualContainerName, "maint_app_other", ExpectedContainerName, "maint_app_svc")

	err = serviceValidator.ValidateContainerName("svc", map[string]any{
		"container_name": "maint_app",
	}, "maint", "app")
	deepstack.AssertDeepStackError(t, err, unexpectedUnderscoreCountContainer, ServiceField, "svc", ActualContainerName, "maint_app", ExpectedPartCount, 3, ActualPartCount, 2)

	err = serviceValidator.ValidateContainerName("svc", map[string]any{
		"container_name": "maint_app_svc_extra",
	}, "maint", "app")
	deepstack.AssertDeepStackError(t, err, unexpectedUnderscoreCountContainer, ServiceField, "svc", ActualContainerName, "maint_app_svc_extra", ExpectedPartCount, 3, ActualPartCount, 4)
}

func TestValidatePorts(t *testing.T) {
	svc := &ServiceValidatorImpl{}

	assert.Nil(t, svc.ValidatePorts("svc", map[string]any{}))
	assert.Nil(t, svc.ValidatePorts("svc", map[string]any{"ports": []any{"8080:80"}}))
	assert.Nil(t, svc.ValidatePorts("svc", map[string]any{"ports": []any{"123"}}))

	forbidden := []string{"22", "53", "80", "443"}
	for _, p := range forbidden {
		t.Run("forbidden_"+p, func(t *testing.T) {
			err := svc.ValidatePorts("svc", map[string]any{"ports": []any{p + ":9999"}})
			deepstack.AssertDeepStackError(t, err, exposingDefaultPortIsForbidden, ServiceField, "svc", PortField, p)
		})
	}
	deepstack.AssertDeepStackError(t, svc.ValidatePorts("svc", map[string]any{"ports": "22"}), "invalid 'ports' in service", ServiceField, "svc")
	deepstack.AssertDeepStackError(t, svc.ValidatePorts("svc", map[string]any{"ports": []any{22}}), "invalid 'ports' in service", ServiceField, "svc")
}

func TestValidateServiceVolumes(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}

	testCases := []struct {
		name        string
		serviceMap  map[string]any
		expectedMsg string
		errorArgs   []any
	}{
		{"no volumes", map[string]any{}, "", nil},
		{"named volume string", map[string]any{"volumes": []string{"maint_app_data:/data"}}, "", nil},
		{"named volume any", map[string]any{"volumes": []any{"maint_app_data:/data"}}, "", nil},

		{"host path string", map[string]any{"volumes": []string{"/host/path:/data"}}, hostDirectoriesMountedForbidden, []any{ServiceField, "svc"}},
		{"host path any", map[string]any{"volumes": []any{"/host/path:/data"}}, hostDirectoriesMountedForbidden, []any{ServiceField, "svc"}},

		{"relative path string", map[string]any{"volumes": []string{"./rel:/data"}}, hostDirectoriesMountedForbidden, []any{ServiceField, "svc"}},
		{"relative path any", map[string]any{"volumes": []any{"./rel:/data"}}, hostDirectoriesMountedForbidden, []any{ServiceField, "svc"}},

		{"wrong prefix string", map[string]any{"volumes": []string{"wrongprefix_data:/data"}}, volumeNamePrefixIsWrong, []any{ServiceField, "svc", VolumeField, "wrongprefix_data:/data", VolumeNameField, "wrongprefix_data", ExpectedPrefixField, "maint_app_"}},
		{"wrong prefix any", map[string]any{"volumes": []any{"wrongprefix_data:/data"}}, volumeNamePrefixIsWrong, []any{ServiceField, "svc", VolumeField, "wrongprefix_data:/data", VolumeNameField, "wrongprefix_data", ExpectedPrefixField, "maint_app_"}},
		{"missing trailing underscore in prefix", map[string]any{"volumes": []string{"maint_app:/data"}}, volumeNamePrefixIsWrong, []any{ServiceField, "svc", VolumeField, "maint_app:/data", VolumeNameField, "maint_app", ExpectedPrefixField, "maint_app_"}},
		{"missing colon in volume entry", map[string]any{"volumes": []string{"maint_app_data"}}, volumeEntryMissingColonSeparator, []any{ServiceField, "svc", VolumeField, "maint_app_data"}},
		{"empty container target", map[string]any{"volumes": []string{"maint_app_data:"}}, volumeTargetMustNotBeEmpty, []any{ServiceField, "svc", VolumeField, "maint_app_data:"}},
		{"relative container target", map[string]any{"volumes": []string{"maint_app_data:data"}}, volumeTargetMustBeAbsolute, []any{ServiceField, "svc", VolumeField, "maint_app_data:data"}},
		{"root container target", map[string]any{"volumes": []string{"maint_app_data:/"}}, volumeTargetMustNotBeRoot, []any{ServiceField, "svc", VolumeField, "maint_app_data:/"}},
		{"too many volume name parts", map[string]any{"volumes": []string{"maint_app_data_extra:/data"}}, unexpectedUnderscoreCountVolume, []any{ServiceField, "svc", VolumeField, "maint_app_data_extra:/data", VolumeNameField, "maint_app_data_extra", ExpectedPartCount, 3, ActualPartCount, 4}},

		{"non string volume entry", map[string]any{"volumes": []any{123}}, "invalid 'volumes' in service", []any{ServiceField, "svc"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := serviceValidator.ValidateServiceVolumes("maint", "app", "svc", testCase.serviceMap)
			if testCase.expectedMsg == "" {
				assert.Nil(t, err)
				return
			}
			deepstack.AssertDeepStackError(t, err, testCase.expectedMsg, testCase.errorArgs...)
		})
	}
}

func TestValidateServiceVolumes_InvalidVolumeSuffixRegex(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}

	testCases := []string{
		"maint_app_:/data",
		"maint_app_data-extra:/data",
	}

	for _, volume := range testCases {
		err := serviceValidator.ValidateServiceVolumes("maint", "app", "svc", map[string]any{
			"volumes": []string{volume},
		})
		assert.True(t, strings.HasPrefix(u.ExtractError(err), "Invalid input."))
	}
}

func TestValidateDeploySection(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}
	err := serviceValidator.ValidateDeploySection("svc", map[string]any{})
	assert.Nil(t, err)

	err = serviceValidator.ValidateDeploySection("svc", map[string]any{
		"deploy": map[string]any{"resources": map[string]any{}},
	})
	assert.Nil(t, err)

	err = serviceValidator.ValidateDeploySection("svc", map[string]any{
		"deploy": map[string]any{"replicas": 1},
	})
	deepstack.AssertDeepStackError(t, err, deployKeywordMustOnlyContainResources, ServiceField, "svc", KeyField, "replicas")

	err = serviceValidator.ValidateDeploySection("svc", map[string]any{
		"deploy": map[string]any{
			"resources": map[string]any{
				"reservations": map[string]any{
					"devices": []any{},
				},
			},
		},
	})
	deepstack.AssertDeepStackError(t, err, devicesKeywordIsForbidden, ServiceField, "svc")
}

func TestServiceNameValidation(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}

	assert.Nil(t, serviceValidator.ValidateServiceName("validname123"))
	err := serviceValidator.ValidateServiceName("invalid_name")
	deepstack.AssertDeepStackError(t, err, buildSimpleRegexErrorMessage("Service", "a-z0-9", 3, 20), ServiceField, "invalid_name", fieldFieldNameKey, "Service", "allowed_symbols", "a-z0-9", "min_length", 3, "max_length", 20)
}

func TestIsValidHttpsUrlWithoutExtras_AcceptsValidUrls(t *testing.T) {
	s := ServiceValidatorImpl{}
	for _, rawUrl := range []string{
		"https://example.com",
		"https://example.com/",
		"https://example.com/path",
		"https://example.com/path/",
		"https://example.com/path/path2",
		"https://sub.example.com",
		"https://example.com:8443",
	} {
		assert.True(t, s.ValidateUrl(rawUrl))
	}
}

func TestIsValidHttpsUrlWithoutExtras_RejectsUrlsWithExtrasOrNonHttps(t *testing.T) {
	s := ServiceValidatorImpl{}
	for _, rawUrl := range []string{
		"http://example.com",
		"ftp://example.com",
		"example.com",
		"https://",
		"https://example.com?query=1",
		"https://example.com/#fragment",
		"https://user:pass@example.com",
		"https://example.com?query=1#fragment",
		" https://example.com",
		"https://example.com ",
		"https://exa mple.com",
		"://example.com",
	} {
		assert.False(t, s.ValidateUrl(rawUrl))
	}
}

func TestValidateLabels(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}

	testCases := []struct {
		appName     string
		serviceName string
		serviceMap  map[string]any
		expectedMsg string
	}{
		{"app", "worker", map[string]any{}, ""},
		{"app", "app", map[string]any{}, mainServiceLabelsMissing},
		{"app", "app", map[string]any{"labels": []any{portLabelKey + "=8080"}}, mainServiceLabelsMustBeMap},
		{"app", "app", map[string]any{"labels": map[string]any{}}, mainServicePortLabelMissing},
		{"app", "app", map[string]any{"labels": map[string]any{portLabelKey: "abc"}}, mainServicePortLabelMustBeString},
		{"app", "app", map[string]any{"labels": map[string]any{portLabelKey: 0}}, mainServicePortLabelMustBeValidPort},
		{"app", "app", map[string]any{"labels": map[string]any{portLabelKey: 70000}}, mainServicePortLabelMustBeValidPort},
		{"app", "app", map[string]any{"labels": map[string]any{portLabelKey: 8080}}, ""},
	}

	for _, testCase := range testCases {
		actualError := serviceValidator.ValidateLabels(testCase.appName, testCase.serviceName, testCase.serviceMap)
		if testCase.expectedMsg == "" {
			assert.Nil(t, actualError)
		} else {
			assert.Equal(t, testCase.expectedMsg, u.ExtractError(actualError))
		}
	}
}

func TestValidateNoTzEnvironment_AllowsWhenNotPresent(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}
	err := serviceValidator.ValidateNoTzEnvironment("svc", map[string]any{
		"environment": []any{"A=1"},
	})
	assert.Nil(t, err)
}

func TestValidateNoTzEnvironment_DeniesListStyle(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}
	err := serviceValidator.ValidateNoTzEnvironment("svc", map[string]any{
		"environment": []any{"A=1", "TZ=Europe/Berlin"},
	})
	deepstack.AssertDeepStackError(t, err, tzEnvironmentVariableIsForbidden, ServiceField, "svc")
}

func TestValidateNoTzEnvironment_DeniesMapStyle(t *testing.T) {
	serviceValidator := &ServiceValidatorImpl{}
	err := serviceValidator.ValidateNoTzEnvironment("svc", map[string]any{
		"environment": map[string]any{"TZ": "Europe/Berlin"},
	})
	deepstack.AssertDeepStackError(t, err, tzEnvironmentVariableIsForbidden, ServiceField, "svc")
}
