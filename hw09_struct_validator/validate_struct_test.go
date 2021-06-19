package hw09structvalidator

import (
	"reflect"
	"testing"

	"github.com/IASamoylov/otus_home_work/hw09_struct_validator/field"
	"github.com/stretchr/testify/require"
)

func TestValidateStruct(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
			{
				tag:    "struct,validate:\"nested\"",
				errMsg: "User.Name: string length must be greater or equal 3",
				value: struct {
					User struct {
						Name string `validate:"len:3"`
					} `validate:"nested"`
				}{
					User: struct {
						Name string `validate:"len:3"`
					}{
						Name: "St",
					},
				},
			},
			{
				tag:    "ptr,validate:\"nested,\"",
				errMsg: "User.Name: string length must be greater or equal 3",
				value: struct {
					User struct {
						Name string `validate:"len:3"`
					} `validate:"nested"`
				}{
					User: struct {
						Name string `validate:"len:3"`
					}{
						Name: "St",
					},
				},
			},
			{
				tag:    "ptr,validate:\"nested,\"",
				errMsg: "tag `validate:\"len:Z\"` contains an invalid rule value Z for this type string",
				value: struct {
					User struct {
						Name string `validate:"len:Z"`
					} `validate:"nested"`
				}{
					User: struct {
						Name string `validate:"len:Z"`
					}{
						Name: "St",
					},
				},
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := field.New(reflect.ValueOf(tc.value), 0)
				ok, err := validateStruct(f)
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
				tag: "struct,validate:\"nested,\"",
				value: struct {
					User struct {
						Role string `validate:"in:Admin,User"`
					} `validate:"nested"`
				}{
					User: struct {
						Role string `validate:"in:Admin,User"`
					}{
						Role: "User",
					},
				},
			},
			{
				tag: "ptr,validate:\"nested,\"",
				value: struct {
					User *struct {
						Role string `validate:"in:Admin,User"`
					} `validate:"nested"`
				}{
					User: &struct {
						Role string `validate:"in:Admin,User"`
					}{
						Role: "User",
					},
				},
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := field.New(reflect.ValueOf(tc.value), 0)
				ok, err := validateStruct(f)
				require.True(t, ok)
				require.NoError(t, err)
			})
		}
	})
}
