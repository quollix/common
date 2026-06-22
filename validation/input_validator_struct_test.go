package validation

import (
	"testing"
	"time"

	"github.com/quollix/common/assert"
	"github.com/quollix/deepstack"

	u "github.com/quollix/common/utils"
)

type validSimpleRegexStruct struct {
	Value string `validate:"default"`
}

type nestedStruct struct {
	Nested validSimpleRegexStruct
}

type validGenericRegexStruct struct {
	Value string `validate:"email"`
}

type structWithFieldToIgnore struct {
	ByteField    []byte
	BoolField    bool
	IntegerField int
}

type noValidationTag struct {
	Value string
}

type nonPublicField struct {
	value string `validate:"default"`
}

type unknownValidationTag struct {
	Value string `validate:"unknown-validation-tag"`
}

type pointerString struct {
	Value *string `validate:"default"`
}

type stringSliceStruct struct {
	Value []string `validate:"default"`
}

type structWithTimeField struct {
	CreatedAt time.Time
}

type ValidationTestCase struct {
	testName             string
	input                any
	expectedErrorMessage string
	expectedContext      []any
}

func TestValidateStruct(t *testing.T) {
	sampleString := u.OfficialMaintainer
	fieldNameOfTheStructure := "Value"

	invalidInputContext := []any{fieldFieldNameKey, fieldNameOfTheStructure, fieldValidationTag, FieldDefault, fieldAllowedSymbols, "a-z0-9", fieldMinLength, 3, fieldMaxLength, 20}
	testCases := []ValidationTestCase{
		{"valid simple regex struct", validSimpleRegexStruct{"sample"}, "", nil},
		{"invalid value in simple struct", validSimpleRegexStruct{"sample!!"}, buildSimpleRegexErrorMessage(fieldNameOfTheStructure, "a-z0-9", 3, 20), invalidInputContext},

		{"valid simple generic struct", validGenericRegexStruct{"admin@admin.de"}, "", nil},
		{"invalid value in generic struct", validGenericRegexStruct{"no-email"}, "Invalid input. The content of the field Value must be a valid email address.", []any{fieldValidationTag, FieldEmail, fieldRegex, emailRegex, fieldFieldNameKey, "Value"}},

		{"no validation tag", noValidationTag{"asdf"}, noValidationTagError, []any{fieldFieldNameKey, "Value"}},
		{"unknown validation tag", unknownValidationTag{"asdf"}, unknownValidationTagError, []any{fieldValidationTag, "unknown-validation-tag", fieldFieldNameKey, "Value"}},
		{"non-public field", nonPublicField{"asdf"}, canNotValidateNonPublicFieldsError, []any{fieldFieldNameKey, "value"}},

		{"string input fails", "some-string", inputMustBeStructError, []any{fieldInputType, "string"}},
		{"float input fails", 1.23, inputMustBeStructError, []any{fieldInputType, "float64"}},
		{"integer input fails", 123, inputMustBeStructError, []any{fieldInputType, "int"}},

		{"valid struct as pointer", &validSimpleRegexStruct{"sample"}, "", nil},

		{"nil input", nil, inputIsNilError, nil},
		{"check that specific types in fields are skipped", structWithFieldToIgnore{[]byte("asdf"), true, 123}, "", nil},
		{"time fields should be skipped", structWithTimeField{CreatedAt: time.Now()}, "", nil},
		{"pointer string fields should fail", pointerString{&sampleString}, "unsupported field type", []any{fieldFieldNameKey, "Value", FieldType, "*string"}},
		{"nested structs should validate (ok)", nestedStruct{Nested: validSimpleRegexStruct{Value: "sample"}}, "", nil},
		{"nested structs should validate (invalid)", nestedStruct{Nested: validSimpleRegexStruct{Value: "sample!!"}}, buildSimpleRegexErrorMessage(fieldNameOfTheStructure, "a-z0-9", 3, 20), invalidInputContext},
		{"slice of structs should fail", []validSimpleRegexStruct{{"sample"}, {"another"}}, inputMustBeStructError, []any{fieldInputType, "slice"}},

		{"valid string array struct", stringSliceStruct{[]string{"sample", "another"}}, "", nil},
		{"invalid string slice struct", stringSliceStruct{[]string{"sample", "another!!"}}, buildSimpleRegexErrorMessage(fieldNameOfTheStructure, "a-z0-9", 3, 20), invalidInputContext},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			assertValidateResult(t, tc)
		})
	}
}

func assertValidateResult(t *testing.T, tc ValidationTestCase) {
	err := ValidateStruct(tc.input)
	if tc.expectedErrorMessage == "" {
		assert.Nil(t, err)
	} else {
		assert.NotNil(t, err)
		_, ok := err.(*deepstack.DeepStackError)
		assert.True(t, ok)
		deepstack.AssertDeepStackError(t, err, tc.expectedErrorMessage, tc.expectedContext...)
	}
}
