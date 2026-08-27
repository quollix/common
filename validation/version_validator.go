package validation

import (
	"bytes"

	u "github.com/quollix/common/utils"
	"gopkg.in/yaml.v3"
)

const (
	AppDefinitionLicenseNotice        = "# This file is licensed under the 0BSD License: https://opensource.org/license/0bsd\\n"
	MissingAppDefinitionLicenseNotice = "app definition must start with: " + AppDefinitionLicenseNotice
)

type VersionValidator interface {
	Validate(composeFileBytes []byte, maintainerName, appName string) error
}

type VersionValidatorImpl struct {
	ComposeValidator       ComposeValidator
	ComposeSyntaxValidator ComposeSyntaxValidator
	RequireLicenseNotice   bool
}

// NewVersionValidator creates a validator; requireLicenseNotice=false is deprecated and kept only for compatibility.
func NewVersionValidator(requireLicenseNotice bool) VersionValidator {
	composeValidator := NewComposeValidator(&ServiceValidatorImpl{})
	return &VersionValidatorImpl{
		ComposeValidator:       composeValidator,
		ComposeSyntaxValidator: NewComposeSyntaxValidator(),
		RequireLicenseNotice:   requireLicenseNotice,
	}
}

func (v *VersionValidatorImpl) Validate(composeFileBytes []byte, maintainerName, appName string) error {
	if err := validateMaintainerAndAppName(maintainerName, appName); err != nil {
		return err
	}

	if u.IsSystemApp(appName) {
		return u.Logger.NewError(SystemAppNamesAreAlreadyReserved)
	}

	if err := v.validateLicenseNotice(composeFileBytes); err != nil {
		return err
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

func (v *VersionValidatorImpl) validateLicenseNotice(composeFileBytes []byte) error {
	if !v.RequireLicenseNotice || bytes.HasPrefix(composeFileBytes, []byte(AppDefinitionLicenseNotice)) {
		return nil
	}
	return u.Logger.NewError(MissingAppDefinitionLicenseNotice)
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

	if err := v.ComposeSyntaxValidator.ValidateComposeSyntax(composeFileBytes); err != nil {
		return err
	}
	return nil
}
