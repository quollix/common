package validation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/quollix/common/utils"
	"github.com/quollix/deepstack"
)

var (
	inputIsNilError                    = "input is nil"
	inputMustBeStructError             = "input must be a data structure"
	canNotValidateNonPublicFieldsError = "cannot validate non-public fields"
	noValidationTagError               = "no validation tag found for field"
	unknownValidationTagError          = "unknown validation tag"

	fieldInputType      = "input_type"
	fieldFieldNameKey   = "field_name"
	fieldMinLength      = "min_length"
	fieldMaxLength      = "max_length"
	fieldAllowedSymbols = "allowed_symbols"
	fieldValidationTag  = "validation_tag"
	fieldRegex          = "regex"
)

func ValidateStruct(inputStructure any) error {
	workStructure := reflect.Indirect(reflect.ValueOf(inputStructure))
	if !workStructure.IsValid() {
		return utils.Logger.NewError(inputIsNilError)
	}
	if workStructure.Kind() != reflect.Struct {
		return utils.Logger.NewError(inputMustBeStructError, fieldInputType, workStructure.Kind().String())
	}

	for _, structField := range reflect.VisibleFields(workStructure.Type()) {
		fieldValue := workStructure.FieldByIndex(structField.Index)
		if err := validateField(fieldValue, structField); err != nil {
			return err
		}
	}
	return nil
}

func validateField(field reflect.Value, structField reflect.StructField) error {
	if !field.CanInterface() {
		return utils.Logger.NewError(canNotValidateNonPublicFieldsError, fieldFieldNameKey, structField.Name)
	}

	switch field.Kind() {
	case reflect.String:
		return validateString(field, structField)
	case reflect.Bool:
		return nil
	case reflect.Int:
		return nil
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			return nil
		}
		return ValidateStruct(field.Interface())
	case reflect.Slice:
		switch field.Type().Elem().Kind() {
		case reflect.Uint8:
			return nil
		case reflect.String:
			return validateSlice(field, structField)
		}
	}

	return utils.Logger.NewError("unsupported field type", fieldFieldNameKey, structField.Name, FieldType, field.Type().String())
}

func validateSlice(field reflect.Value, structField reflect.StructField) error {
	stringSlice := field.Interface().([]string)
	for _, value := range stringSlice {
		if err := validateString(reflect.ValueOf(value), structField); err != nil {
			return err
		}
	}
	return nil
}

func validateString(field reflect.Value, structField reflect.StructField) error {
	fieldTag := structField.Tag.Get("validate")
	if fieldTag == "" {
		return utils.Logger.NewError(noValidationTagError, fieldFieldNameKey, structField.Name)
	}
	err := Validate(structField.Name, fieldTag, field.String())
	if err != nil {
		return utils.Logger.AddContext(err, fieldFieldNameKey, structField.Name, fieldValidationTag, fieldTag)
	}
	return nil
}

func Validate(fieldName, validationTag, fieldValue string) error {
	if validationTag == FieldIgnore {
		return nil
	}

	validationFunc, found := ValidationMap[validationTag]
	if !found {
		return utils.Logger.NewError(unknownValidationTagError, fieldValidationTag, validationTag)
	}

	return validationFunc(fieldName, fieldValue)
}

func ReadBody[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	var result T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Warn(err)
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return nil, false
	}
	defer utils.Close(r.Body)

	if err = json.Unmarshal(body, &result); err != nil {
		utils.Logger.Warn(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return nil, false
	}

	if err = ValidateStruct(result); err != nil {
		utils.Logger.Warn(err)
		deepStackError, ok := err.(*deepstack.DeepStackError)
		if ok {
			http.Error(w, deepStackError.Message, http.StatusBadRequest)
		} else {
			http.Error(w, "invalid input", http.StatusBadRequest)
		}
		return nil, false
	}

	return &result, true
}

func formatErrorWithContext(errorMessage string, context map[string]any) string {
	if len(context) == 0 {
		return errorMessage
	}
	keys := make([]string, 0, len(context))
	for k := range context {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s: %v", k, context[k]))
	}
	return fmt.Sprintf("%s (%s)", errorMessage, strings.Join(parts, ", "))
}
