package hw09structvalidator

import (
	"reflect"
	"testing"

	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
	"github.com/stretchr/testify/require"
)

func TestValidateString(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
			{
				tag: "string,validate:\"in:Admin,User,\"",
				value: struct {
					Role string `validate:"in:Admin,User"`
				}{
					Role: "Developer",
				},
				errMsg: "Role: value is not contains in range Admin,User",
			},
			{
				tag: "string,validate:\"regexp:\\d+\"",
				value: struct {
					Role string `validate:"regexp:\\d+"`
				}{
					Role: "User",
				},
				errMsg: "Role: value does not match the mask \\d+",
			},
			{
				tag: "string,validate:\"len:3\"",
				value: struct {
					Name string `validate:"len:3"`
				}{
					Name: "St",
				},
				errMsg: "Name: string length must be equal 3",
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := field.New(reflect.ValueOf(tc.value), 0)
				ok, err := validateString(f)
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
				tag: "string,validate:\"in:Admin,User,\"",
				value: struct {
					Role string `validate:"in:Admin,User"`
				}{
					Role: "User",
				},
			},
			{
				tag: "string,validate:\"regexp:\\d+\"",
				value: struct {
					ID string `validate:"regexp:\\d+"`
				}{
					ID: "2134",
				},
			},
			{
				tag: "string,validate:\"len:3\"",
				value: struct {
					Name string `validate:"len:3"`
				}{
					Name: "Iva",
				},
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := field.New(reflect.ValueOf(tc.value), 0)
				ok, err := validateString(f)
				require.True(t, ok)
				require.NoError(t, err)
			})
		}
	})
}
