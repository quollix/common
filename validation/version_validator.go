package validation

import (
	u "github.com/quollix/common/utils"
	"gopkg.in/yaml.v3"
)

type VersionValidator interface {
	Validate(composeFileBytes []byte, maintainerName, appName string) error
}

type VersionValidatorImpl struct {
	FileSystemService u.ComposeSyntaxChecker
	ComposeValidator  ComposeValidator
}

func NewVersionValidator() VersionValidator {
	osWrapper := &u.OsWrapperImpl{}
	fileSystemOperator := u.NewFileSystemOperator(osWrapper)
	return NewVersionValidatorWithDependencies(
		fileSystemOperator,
		NewComposeValidator(&ServiceValidatorImpl{}, fileSystemOperator),
	)
}

func NewVersionValidatorWithDependencies(fileSystemService u.ComposeSyntaxChecker, composeValidator ComposeValidator) VersionValidator {
	return &VersionValidatorImpl{
		FileSystemService: fileSystemService,
		ComposeValidator:  composeValidator,
	}
}

func (v *VersionValidatorImpl) Validate(composeFileBytes []byte, maintainerName, appName string) error {
	if err := validateMaintainerAndAppName(maintainerName, appName); err != nil {
		return err
	}

	if u.IsSystemApp(appName) {
		return u.Logger.NewError(SystemAppNamesAreAlreadyReserved)
	}

	err := v.runValidations(composeFileBytes, maintainerName, appName)
	if err != nil {
		return u.Logger.AddContext(err, MaintainerField, maintainerName, AppField, appName)
	}

	return nil
}

func validateMaintainerAndAppName(maintainerName, appName string) error {
	if err := Validate("Maintainer", FieldDefault, maintainerName); err != nil {
		return u.Logger.AddContext(err, MaintainerField, maintainerName)
	}
	if err := Validate("App", FieldDefault, appName); err != nil {
		return u.Logger.AddContext(err, AppField, appName)
	}
	return nil
}

func (v *VersionValidatorImpl) runValidations(composeFileBytes []byte, maintainerName string, appName string) error {
	var mapData map[string]any
	if err := yaml.Unmarshal(composeFileBytes, &mapData); err != nil {
		return u.Logger.NewError(err.Error())
	}

	if err := v.ComposeValidator.ValidateComposePlaceholders(composeFileBytes); err != nil {
		return err
	}

	if err := v.ComposeValidator.ValidateComposeMap(mapData, maintainerName, appName); err != nil {
		return err
	}

	if err := v.FileSystemService.CheckDockerComposeSyntax(composeFileBytes); err != nil {
		return err
	}
	return nil
}
