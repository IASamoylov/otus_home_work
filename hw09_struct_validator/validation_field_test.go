package hw09structvalidator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidationField(t *testing.T) {
	t.Run("has validation tags works fine", func(t *testing.T) {
		tests := []struct {
			name     string
			excepted bool
			tags     []validationFieldTag
		}{
			{name: "tags is empty", excepted: false, tags: []validationFieldTag{}},
			{name: "tags is not empty", excepted: true, tags: make([]validationFieldTag, 1)},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				field := &validationField{tags: tc.tags}
				require.Equal(t, tc.excepted, field.hasValidationTags())
			})
		}
	})

	t.Run("tag parsing", func(t *testing.T) {
		tests := []struct {
			name string
			tag  reflect.StructTag
			tags []validationFieldTag
		}{
			{name: "ignore not package tags", tag: "json:\"id\"", tags: []validationFieldTag{}},
			{name: "ignore package tag that configured invalid", tag: "validate:", tags: []validationFieldTag{}},
			{name: "zero value when tag configuration with error", tag: "validate:\"\"", tags: []validationFieldTag{
				{valueIsUndefined: true},
			}},
			{name: "value is undefined when tag configuration with error", tag: "validate:\"len:\"", tags: []validationFieldTag{
				{name: lenValidation, valueIsUndefined: true},
			}},
			{name: "value is undefined when tag configuration with error", tag: "validate:\"len\"", tags: []validationFieldTag{
				{name: lenValidation, valueIsUndefined: true},
			}},
			{name: "tag contain one valid rule", tag: "validate:\"len:46\"", tags: []validationFieldTag{
				{name: lenValidation, value: "46"},
			}},
			{name: "tag contains more than one rules", tag: "validate:\"min:46|max:132\"", tags: []validationFieldTag{
				{name: minValidation, value: "46"},
				{name: maxValidation, value: "132"},
			}},
			{name: "value is undefined when tag configuration with error", tag: "validate:\"min:46|max:\"", tags: []validationFieldTag{
				{name: minValidation, value: "46"},
				{name: maxValidation, valueIsUndefined: true},
			}},
			{name: "value is undefined when tag configuration with error", tag: "validate:\"min|max:132\"", tags: []validationFieldTag{
				{name: minValidation, valueIsUndefined: true},
				{name: maxValidation, value: "132"},
			}},
			{name: "configuration error when tag contains more than one rules", tag: "validate:\"min:64|\"", tags: []validationFieldTag{
				{name: minValidation, value: "64"},
				{valueIsUndefined: true},
			}},
			{name: "configuration error when tag contains more than one rules", tag: "validate:\"|\"", tags: []validationFieldTag{
				{valueIsUndefined: true},
				{valueIsUndefined: true},
			}},
		}

		for _, tc := range tests {
			t.Run(fmt.Sprintf("%v:%v", tc.name, tc.tag), func(t *testing.T) {
				result := parseTags(tc.tag)

				require.Len(t, result, len(tc.tags))

				for i, tag := range tc.tags {
					require.Equal(t, tc.tag, result[i].tag)
					require.Equal(t, tag.value, result[i].value)
					require.Equal(t, tag.name, result[i].name)
					require.Equal(t, tag.valueIsUndefined, result[i].valueIsUndefined)
				}
			})
		}
	})
}

func TestValidationFieldTagValidateSuccessful(t *testing.T) {
	tests := []struct {
		tag   reflect.StructTag
		value interface{}
	}{
		{tag: "string,validate:\"len:36\"", value: struct {
			ID string `validate:"len:36"`
		}{}},
		{tag: "string,validate:\"regexp:\\d+\"", value: struct {
			ID string `validate:"regexp:\\d+"`
		}{}},
		{tag: "string,validate:\"in:admin,user\"", value: struct {
			Role string `validate:"in:admin,user"`
		}{}},
		{tag: "string,validate:\"regexp:\\d+|len:36\"", value: struct {
			Role string `validate:"regexp:\\d+|len:36"`
		}{}},
		{tag: "int8,validate:\"in:-3,44\"", value: struct {
			ID int8 `validate:"in:-3,44"`
		}{}},
		{tag: "int8,validate:\"min:-36\"", value: struct {
			ID int8 `validate:"min:-36"`
		}{}},
		{tag: "int8,validate:\"max:36\"", value: struct {
			ID int8 `validate:"max:36"`
		}{}},
		{tag: "int16,validate:\"in:-3,44\"", value: struct {
			ID int16 `validate:"in:-3,44"`
		}{}},
		{tag: "int16,validate:\"min:-36\"", value: struct {
			ID int16 `validate:"min:-36"`
		}{}},
		{tag: "int16,validate:\"max:36\"", value: struct {
			ID int16 `validate:"max:36"`
		}{}},
		{tag: "int32,validate:\"in:-3,44\"", value: struct {
			ID int32 `validate:"in:-3,44"`
		}{}},
		{tag: "int32,validate:\"min:-36\"", value: struct {
			ID int32 `validate:"min:-36"`
		}{}},
		{tag: "int32,validate:\"max:36\"", value: struct {
			ID int32 `validate:"max:36"`
		}{}},
		{tag: "int64,validate:\"in:-3,44\"", value: struct {
			ID int64 `validate:"in:-3,44"`
		}{}},
		{tag: "int64,validate:\"min:-36\"", value: struct {
			ID int64 `validate:"min:-36"`
		}{}},
		{tag: "int64,validate:\"max:36\"", value: struct {
			ID int64 `validate:"max:36"`
		}{}},
		{tag: "uint8,validate:\"in:3,44\"", value: struct {
			ID uint8 `validate:"in:3,44"`
		}{}},
		{tag: "uint8,validate:\"min:36\"", value: struct {
			ID uint8 `validate:"min:36"`
		}{}},
		{tag: "uint8,validate:\"max:36\"", value: struct {
			ID uint8 `validate:"max:36"`
		}{}},
		{tag: "uint16,validate:\"in:3,44\"", value: struct {
			ID uint16 `validate:"in:3,44"`
		}{}},
		{tag: "uint16,validate:\"min:36\"", value: struct {
			ID uint16 `validate:"min:36"`
		}{}},
		{tag: "uint16,validate:\"max:36\"", value: struct {
			ID uint16 `validate:"max:36"`
		}{}},
		{tag: "uint32,validate:\"in:3,44\"", value: struct {
			ID uint32 `validate:"in:3,44"`
		}{}},
		{tag: "uint32,validate:\"min:36\"", value: struct {
			ID uint32 `validate:"min:36"`
		}{}},
		{tag: "uint32,validate:\"max:36\"", value: struct {
			ID uint32 `validate:"max:36"`
		}{}},
		{tag: "uint64,validate:\"in:3,44\"", value: struct {
			ID uint64 `validate:"in:3,44"`
		}{}},
		{tag: "uint64,validate:\"min:36\"", value: struct {
			ID uint64 `validate:"min:36"`
		}{}},
		{tag: "uint64,validate:\"max:36\"", value: struct {
			ID uint64 `validate:"max:36"`
		}{}},
	}

	for _, tc := range tests {
		t.Run(string(tc.tag), func(t *testing.T) {
			ok, err := newValidationField(reflect.ValueOf(tc.value), 0).validateTags()
			require.True(t, ok)
			require.NoError(t, err)
		})
	}
}

func TestValidationFieldValidateFail(t *testing.T) {
	tests := []struct {
		tag   reflect.StructTag
		value interface{}
	}{}

	for _, tc := range tests {
		t.Run(string(tc.tag), func(t *testing.T) {
			ok, err := newValidationField(reflect.ValueOf(tc.value), 0).validateTags()
			require.True(t, ok)
			require.NoError(t, err)
		})
	}
}
