package hw09structvalidator

import (
	"fmt"
	"reflect"
	"testing"
	"time"

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
				{name: lenTagValidation, valueIsUndefined: true},
			}},
			{name: "value is undefined when tag configuration with error", tag: "validate:\"len\"", tags: []validationFieldTag{
				{name: lenTagValidation, valueIsUndefined: true},
			}},
			{name: "tag contain one valid rule", tag: "validate:\"len:46\"", tags: []validationFieldTag{
				{name: lenTagValidation, value: "46"},
			}},
			{name: "tag contains more than one rules", tag: "validate:\"min:46|max:132\"", tags: []validationFieldTag{
				{name: minTagValidation, value: "46"},
				{name: maxTagValidation, value: "132"},
			}},
			{name: "value is undefined when tag configuration with error", tag: "validate:\"min:46|max:\"", tags: []validationFieldTag{
				{name: minTagValidation, value: "46"},
				{name: maxTagValidation, valueIsUndefined: true},
			}},
			{name: "value is undefined when tag configuration with error", tag: "validate:\"min|max:132\"", tags: []validationFieldTag{
				{name: minTagValidation, valueIsUndefined: true},
				{name: maxTagValidation, value: "132"},
			}},
			{name: "configuration error when tag contains more than one rules", tag: "validate:\"min:64|\"", tags: []validationFieldTag{
				{name: minTagValidation, value: "64"},
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
		{tag: "struct,validate:\"nested\"", value: struct {
			User struct{ ID int } `validate:"nested"`
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
		tag        reflect.StructTag
		value      interface{}
		wrapErrMsg string
		errMsg     string
	}{
		{
			tag: "string,validate:\"uuid\"",
			value: struct {
				ID string `validate:"uuid"`
			}{},
			errMsg: "tag `validate:\"uuid\"` configured incorrectly for field ID",
		},
		{
			tag: "string,validate:\"len:\"",
			value: struct {
				ID string `validate:"len:"`
			}{},
			errMsg: "tag `validate:\"len:\"` configured incorrectly validation rule for field ID",
		},
		{
			tag: "int32,validate:\"in:3,\"",
			value: struct {
				Age int32 `validate:"in:3,"`
			}{},
			errMsg: "tag `validate:\"in:3,\"` configured incorrectly validation rule for field Age",
		},
		{
			tag: "int32,validate:\"in:3\"",
			value: struct {
				Age int32 `validate:"in:3"`
			}{},
			errMsg: "tag `validate:\"in:3\"` configured incorrectly validation rule for field Age",
		},
		{
			tag: "int32,validate:\"in:3,7,18\"",
			value: struct {
				Age int32 `validate:"in:3,7,18"`
			}{},
			errMsg: "tag `validate:\"in:3,7,18\"` configured incorrectly validation rule for field Age",
		},
		{
			tag: "time.Time,validate:\"in:2021.07.25,2021.09.25\"",
			value: struct {
				Birthday time.Time `validate:"in:2021.07.25,2021.09.25"`
			}{},
			errMsg: "tag `validate:\"in:2021.07.25,2021.09.25\"` not supported for this type time.Time",
		},
		{
			tag: "int32,validate:\"len:36\"",
			value: struct {
				ID int32 `validate:"len:36"`
			}{},
			errMsg: "tag `validate:\"len:36\"` not supported for this type int32",
		},
		{
			tag: "string,validate:\"min:13\"",
			value: struct {
				ID string `validate:"min:13"`
			}{},
			errMsg: "tag `validate:\"min:13\"` not supported for this type string",
		},
		{
			tag: "time.Time,validate:\"max:2999.07.25\"",
			value: struct {
				Birthday time.Time `validate:"max:2999.07.25"`
			}{},
			errMsg: "tag `validate:\"max:2999.07.25\"` not supported for this type time.Time",
		},
		{
			tag: "int32,validate:\"max:Z\"",
			value: struct {
				Birthday int32 `validate:"max:Z"`
			}{},
			errMsg: "tag `validate:\"max:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "int64,validate:\"max:Z\"",
			value: struct {
				Birthday int64 `validate:"max:Z"`
			}{},
			errMsg: "tag `validate:\"max:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "int8,validate:\"min:Z\"",
			value: struct {
				Birthday int8 `validate:"min:Z"`
			}{},
			errMsg: "tag `validate:\"min:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "int16,validate:\"min:Z\"",
			value: struct {
				Birthday int16 `validate:"min:Z"`
			}{},
			errMsg: "tag `validate:\"min:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "int32,validate:\"min:Z\"",
			value: struct {
				Birthday int32 `validate:"min:Z"`
			}{},
			errMsg: "tag `validate:\"min:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "uint8,validate:\"max:Z\"",
			value: struct {
				Birthday uint8 `validate:"max:Z"`
			}{},
			errMsg: "tag `validate:\"max:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "uint16,validate:\"max:Z\"",
			value: struct {
				Birthday uint16 `validate:"max:Z"`
			}{},
			errMsg: "tag `validate:\"max:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "uint16,validate:\"min:Z\"",
			value: struct {
				Birthday uint16 `validate:"min:Z"`
			}{},
			errMsg: "tag `validate:\"min:Z\"` contains an invalid rule value Z",
		},
		{
			tag: "int,validate:\"regexp:\\d+\"",
			value: struct {
				ID int `validate:"regexp:\\d+"`
			}{},
			errMsg: "tag `validate:\"regexp:\\\\d+\"` not supported for this type int",
		},
	}

	for _, tc := range tests {
		t.Run(string(tc.tag), func(t *testing.T) {
			ok, err := newValidationField(reflect.ValueOf(tc.value), 0).validateTags()
			require.False(t, ok)
			require.EqualError(t, err, tc.errMsg)
		})
	}
}
