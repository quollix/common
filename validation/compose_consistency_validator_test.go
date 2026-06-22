package validation

import (
	"testing"

	"github.com/quollix/common/assert"
	"github.com/quollix/deepstack"
	"gopkg.in/yaml.v3"
)

var validator = &ComposeConsistencyValidatorImpl{}

func TestValidateVolumeMappings_AllowsMatchingGlobalAndServiceVolumes(t *testing.T) {
	err := validator.ValidateVolumeMappings(parseComposeYml(t, `
services:
  app:
    volumes:
      - maint_app_data:/data
volumes:
  maint_app_data: {}
`))

	assert.Nil(t, err)
}

func TestValidateVolumeMappings_AllowsSameGlobalVolumeMountedByMultipleServices(t *testing.T) {
	err := validator.ValidateVolumeMappings(parseComposeYml(t, `
services:
  app:
    volumes:
      - maint_app_data:/data
  worker:
    volumes:
      - maint_app_data:/worker-data
volumes:
  maint_app_data: {}
`))

	assert.Nil(t, err)
}

func TestValidateVolumeMappings_DeniesServiceVolumeMissingGlobalDeclaration(t *testing.T) {
	err := validator.ValidateVolumeMappings(parseComposeYml(t, `
services:
  app:
    volumes:
      - maint_app_data:/data
volumes: {}
`))

	deepstack.AssertDeepStackError(t, err, serviceVolumeMustBeDeclaredGlobally, ServiceField, "app", VolumeNameField, "maint_app_data")
}

func TestValidateVolumeMappings_DeniesUnusedGlobalVolume(t *testing.T) {
	err := validator.ValidateVolumeMappings(parseComposeYml(t, `
services:
  app:
    volumes: []
volumes:
  maint_app_data: {}
`))

	deepstack.AssertDeepStackError(t, err, globalVolumeMustBeMounted, VolumeNameField, "maint_app_data")
}

func TestValidateVolumeMappings_DeniesDuplicateVolumeTargetWithinService(t *testing.T) {
	err := validator.ValidateVolumeMappings(parseComposeYml(t, `
services:
  app:
    volumes:
      - maint_app_data:/data:ro
      - maint_app_cache:/data:rw
volumes:
  maint_app_cache: {}
  maint_app_data: {}
`))

	deepstack.AssertDeepStackError(t, err, duplicateVolumeTargetInService, ServiceField, "app", VolumeField, "maint_app_cache:/data:rw", "volume_target", "/data")
}

func TestValidateServiceReferences_AllowsDependsOnListExistingService(t *testing.T) {
	err := validator.ValidateServiceReferences(parseComposeYml(t, `
services:
  app:
    depends_on:
      - database
  database: {}
`))

	assert.Nil(t, err)
}

func TestValidateServiceReferences_DeniesDependsOnListMissingService(t *testing.T) {
	err := validator.ValidateServiceReferences(parseComposeYml(t, `
services:
  app:
    depends_on:
      - database
`))

	deepstack.AssertDeepStackError(t, err, dependsOnServiceMustExist, ServiceField, "app", "depends_on_service", "database")
}

func TestValidateServiceReferences_DeniesDependsOnMap(t *testing.T) {
	err := validator.ValidateServiceReferences(parseComposeYml(t, `
services:
  app:
    depends_on:
      database:
        condition: service_started
  database: {}
`))

	deepstack.AssertDeepStackError(t, err, invalidDependsOn, ServiceField, "app")
}

func TestValidateServiceReferences_DeniesInvalidDependsOnType(t *testing.T) {
	err := validator.ValidateServiceReferences(parseComposeYml(t, `
services:
  app:
    depends_on: database
  database: {}
`))

	deepstack.AssertDeepStackError(t, err, invalidDependsOn, ServiceField, "app")
}

func parseComposeYml(t *testing.T, composeYml string) map[string]any {
	compose := map[string]any{}
	err := yaml.Unmarshal([]byte(composeYml), &compose)
	assert.Nil(t, err)
	return compose
}
