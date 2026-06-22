package validation

import (
	"fmt"
	"regexp"

	"github.com/quollix/common/utils"
)

type ValidationFunc func(fieldName, valueToValidate string) error

var (
	emailRegex        = `^[a-zA-Z0-9._%+-]{1,64}@[a-zA-Z0-9.-]{1,253}\.[a-zA-Z]{2,}$`
	simpleSecretRegex = NewSimpleRegex("a-f0-9", 64, 64)

	defaultRegex = NewSimpleRegex("a-z0-9", 3, 20)
)

const (
	FieldDefault        = "default"
	FieldUsername       = "username"
	FieldIgnore         = "ignore"
	FieldVersionName    = "version_name"
	FieldFileName       = "file_name"
	FieldSearchTerm     = "search_term"
	FieldPassword       = "password"
	FieldEmail          = "email"
	FieldNumber         = "number"
	FieldHost           = "host"
	FieldKnownHosts     = "known_hosts"
	FieldResticBackupID = "restic_backup_id"
	FieldRemoteHost     = "remote_host"
	FieldSecret         = "secret"
	FieldType           = "type"
	FieldDefaultOrEmpty = "default_or_empty"
	FieldLoose          = "loose"
)

var ValidationMap = map[string]ValidationFunc{
	FieldDefault:        defaultRegex,
	FieldUsername:       NewSimpleRegex("a-z0-9_-", 3, 20),
	FieldVersionName:    NewSimpleRegex("a-z0-9.", 3, 20),
	FieldFileName:       NewSimpleRegex("a-zA-Z0-9._-", 1, 100),
	FieldSearchTerm:     NewSimpleRegex("a-z0-9", 0, 20),
	FieldPassword:       NewSimpleRegex("a-zA-Z0-9._-", 8, 30),
	FieldEmail:          NewGenericRegex(emailRegex, "Invalid input. The content of the field %s must be a valid email address."),
	FieldNumber:         NewSimpleRegex("0-9", 1, 20),
	FieldHost:           NewSimpleRegex("a-zA-Z0-9:._-", 0, 64),
	FieldKnownHosts:     NewGenericRegex(`^[A-Za-z0-9.:,/_+=#@\[\]| \r\n-]{0,}$`, "Invalid input. The content of the field %s must contain valid SSH known_hosts entries."),
	FieldResticBackupID: NewSimpleRegex("a-f0-9", 64, 64),
	FieldRemoteHost:     NewSimpleRegex("a-zA-Z0-9._-", 0, 64),
	FieldSecret:         simpleSecretRegex,
	FieldDefaultOrEmpty: NewGenericRegex(`^$|^[a-z0-9]{3,20}$`, "Invalid input. The content of the field %s must be empty or be between 3 and 20 characters long. Allowed symbols are: a-z0-9."),
	FieldLoose:          NewGenericRegex(`^[A-Za-z0-9!@#$%^&*()_\-+=\.,:;\/\?\[\]\{\}\|~<>]{1,128}$`, "Invalid input. The content of the field %s contains unsupported symbols."),
}

func NewSimpleRegex(allowedSymbols string, minLength, maxLength int) ValidationFunc {
	compiledRegex := regexp.MustCompile(fmt.Sprintf("^[%s]{%d,%d}$", allowedSymbols, minLength, maxLength))
	return func(fieldName, valueToValidate string) error {
		if len(valueToValidate) < minLength || len(valueToValidate) > maxLength || !compiledRegex.MatchString(valueToValidate) {
			return utils.Logger.NewError(
				buildSimpleRegexErrorMessage(fieldName, allowedSymbols, minLength, maxLength),
				fieldFieldNameKey, fieldName,
				fieldAllowedSymbols, allowedSymbols,
				fieldMinLength, minLength,
				fieldMaxLength, maxLength,
			)
		}
		return nil
	}
}

func NewGenericRegex(regexString, humanReadableMessageTemplate string) ValidationFunc {
	compiledRegex := regexp.MustCompile(regexString)
	return func(fieldName, valueToValidate string) error {
		if !compiledRegex.MatchString(valueToValidate) {
			return utils.Logger.NewError(
				fmt.Sprintf(humanReadableMessageTemplate, fieldName),
				fieldFieldNameKey, fieldName,
				fieldRegex, regexString,
			)
		}
		return nil
	}
}

func buildSimpleRegexErrorMessage(fieldName, allowedSymbols string, minLength, maxLength int) string {
	var lengthRequirement string
	switch minLength {
	case maxLength:
		lengthRequirement = fmt.Sprintf("must be exactly %d characters long", minLength)
	case 0:
		lengthRequirement = fmt.Sprintf("must be at most %d characters long", maxLength)
	default:
		lengthRequirement = fmt.Sprintf("must be between %d and %d characters long", minLength, maxLength)
	}

	return fmt.Sprintf(
		"Invalid input. The content of the field %s %s. Allowed symbols are: %s.",
		fieldName,
		lengthRequirement,
		allowedSymbols,
	)
}
