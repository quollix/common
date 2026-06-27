package validation

import (
	"fmt"
	"regexp"

	u "github.com/quollix/common/utils"
)

var (
	allowedTopLevelKeys             = u.MapOf("services", "volumes")
	composePlaceholderRegex         = regexp.MustCompile(`\$\{([^}]+)\}`)
	allowedComposePlaceholders      = []string{"BASE_DOMAIN", "CLIENT_ID", "CLIENT_SECRET", "IANA_TIMEZONE", "SERVER_HOST"}
	allowedComposePlaceholderLookup = map[string]struct{}{
		"BASE_DOMAIN":   {},
		"CLIENT_ID":     {},
		"CLIENT_SECRET": {},
		"IANA_TIMEZONE": {},
		// SERVER_HOST is kept as a legacy alias for app compose files created before BASE_DOMAIN.
		"SERVER_HOST": {},
	}
)

type ComposeValidator interface {
	ValidateComposePlaceholders(composeFileBytes []byte) error
	ValidateComposeMap(compose map[string]any, maintainerName, appName string) error
}

type ComposeValidatorImpl struct {
	ServiceValidator            ServiceValidator
	ComposeConsistencyValidator ComposeConsistencyValidator
	FileSystemOperator          u.FileSystemOperator
}

func NewComposeValidator(serviceValidator ServiceValidator, fileSystemOperator u.FileSystemOperator) ComposeValidator {
	return &ComposeValidatorImpl{
		ServiceValidator:            serviceValidator,
		ComposeConsistencyValidator: &ComposeConsistencyValidatorImpl{},
		FileSystemOperator:          fileSystemOperator,
	}
}

func (f *ComposeValidatorImpl) ValidateComposePlaceholders(composeFileBytes []byte) error {
	matches := composePlaceholderRegex.FindAllSubmatch(composeFileBytes, -1)
	for _, match := range matches {
		placeholderName := string(match[1])
		if _, isAllowed := allowedComposePlaceholderLookup[placeholderName]; isAllowed {
			continue
		}
		return u.Logger.NewError(
			"unsupported docker compose placeholder",
			"unsupported_placeholder", placeholderName,
			"allowed_placeholders", allowedComposePlaceholders,
		)
	}
	return nil
}

func (f *ComposeValidatorImpl) ValidateComposeMap(compose map[string]any, maintainerName, appName string) error {
	if err := validateTopLevelKeys(compose); err != nil {
		return err
	}
	if err := validateGlobalVolumes(compose, maintainerName, appName); err != nil {
		return err
	}
	if err := f.validateServices(compose, maintainerName, appName); err != nil {
		return err
	}
	composeConsistencyValidator := f.getComposeConsistencyValidator()
	if err := composeConsistencyValidator.ValidateVolumeMappings(compose); err != nil {
		return err
	}
	if err := composeConsistencyValidator.ValidateServiceReferences(compose); err != nil {
		return err
	}
	return nil
}

func (f *ComposeValidatorImpl) getComposeConsistencyValidator() ComposeConsistencyValidator {
	if f.ComposeConsistencyValidator != nil {
		return f.ComposeConsistencyValidator
	}
	return &ComposeConsistencyValidatorImpl{}
}

func validateTopLevelKeys(keys map[string]any) error {
	for key := range keys {
		_, ok := allowedTopLevelKeys[key]
		if !ok {
			return u.Logger.NewError(notAllowedTopLevelKeyword, "key", key)
		}
	}
	return nil
}

func (f *ComposeValidatorImpl) validateServices(compose map[string]any, maintainerName, appName string) error {
	services, ok := compose["services"].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid 'services' section in docker-compose.yml")
	}

	_, doesMainServiceExist := services[appName]
	if !doesMainServiceExist {
		return u.Logger.NewError(mainServiceMustBeDefined)
	}

	for serviceName, serviceValue := range services {
		err := f.validateService(maintainerName, appName, serviceName, serviceValue)
		if err != nil {
			return u.Logger.AddContext(err, ServiceField, serviceName)
		}
	}
	return nil
}

func (f *ComposeValidatorImpl) validateService(maintainerName, appName, serviceName string, serviceValue any) error {
	if err := f.ServiceValidator.ValidateServiceName(serviceName); err != nil {
		return err
	}
	serviceMap, ok := serviceValue.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid service definition for service %v", serviceName)
	}
	if err := f.ServiceValidator.ValidateServiceKeys(serviceName, serviceMap); err != nil {
		return err
	}
	if err := f.ServiceValidator.ValidateImage(serviceName, serviceMap); err != nil {
		return err
	}
	if err := f.ServiceValidator.ValidateContainerName(serviceName, serviceMap, maintainerName, appName); err != nil {
		return err
	}
	if err := f.ServiceValidator.ValidatePorts(serviceName, serviceMap); err != nil {
		return err
	}
	if err := f.ServiceValidator.ValidateServiceVolumes(maintainerName, appName, serviceName, serviceMap); err != nil {
		return err
	}
	if err := f.ServiceValidator.ValidateDeploySection(serviceName, serviceMap); err != nil {
		return err
	}
	if err := f.ServiceValidator.ValidateLabels(appName, serviceName, serviceMap); err != nil {
		return err
	}
	if err := f.ServiceValidator.ValidateNoTzEnvironment(serviceName, serviceMap); err != nil {
		return err
	}
	return nil
}

func validateGlobalVolumes(compose map[string]any, maintainerName, appName string) error {
	volumes, ok := compose["volumes"].(map[string]any)
	if !ok {
		return nil
	}
	for volumeName, volume := range volumes {
		volumesMap, ok := volume.(map[string]any)
		if ok && len(volumesMap) > 0 {
			return u.Logger.NewError(globalVolumeShouldNotHaveSubKeywords)
		}
		if err := validateNamedVolumeName(maintainerName, appName, volumeName); err != nil {
			return err
		}
	}
	return nil
}
