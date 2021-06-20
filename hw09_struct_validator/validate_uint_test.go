package hw09structvalidator

import (
	"reflect"
	"testing"

	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
	"github.com/stretchr/testify/require"
)

func TestValidateUint(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
			{
				tag: "int,validate:\"in:5,9\"",
				value: struct {
					ID uint `validate:"in:5,9"`
				}{
					ID: 13,
				},
				errMsg: "ID: value must be in range [5,9]",
			},
			{
				tag: "int,validate:\"in:5,9,\"",
				value: struct {
					ID uint `validate:"in:5,9"`
				}{
					ID: 3,
				},
				errMsg: "ID: value must be in range [5,9]",
			},
			{
				tag: "int,validate:\"min:9,\"",
				value: struct {
					ID uint `validate:"min:9"`
				}{
					ID: 3,
				},
				errMsg: "ID: value must be greater or equal 9",
			},
			{
				tag: "int,validate:\"max:9,\"",
				value: struct {
					ID uint `validate:"max:9"`
				}{
					ID: 13,
				},
				errMsg: "ID: value must be les or equal 9",
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := field.New(reflect.ValueOf(tc.value), 0)
				ok, err := validateInt(f)
				require.False(t, ok)
				require.EqualError(t, err, tc.errMsg)
			})
		}
	})

	t.Run("successful cases", func(t *testing.T) {
		tests := []struct {
			tag   reflect.StructTag
			value interface{}
		}{
			{
				tag: "string,validate:\"in:5,9,\"",
				value: struct {
					ID uint `validate:"in:5,9"`
				}{
					ID: 7,
				},
			},
			{
				tag: "string,validate:\"max:7\"",
				value: struct {
					ID uint `validate:"max:7"`
				}{
					ID: 5,
				},
			},
			{
				tag: "string,validate:\"min:7\"",
				value: struct {
					ID uint `validate:"min:7"`
				}{
					ID: 13,
				},
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := field.New(reflect.ValueOf(tc.value), 0)
				ok, err := validateInt(f)
				require.True(t, ok)
				require.NoError(t, err)
			})
		}
	})
}
