package field

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInTagValidator(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
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
				tag: "int32,validate:\"in:3,a\"",
				value: struct {
					Age int32 `validate:"in:3,a"`
				}{},
				errMsg: "tag `validate:\"in:3,a\"` contains an invalid rule value 3,a for this type int32",
			},
			{
				tag: "int32,validate:\"in:a,z\"",
				value: struct {
					Age int32 `validate:"in:a,z"`
				}{},
				errMsg: "tag `validate:\"in:a,z\"` contains an invalid rule value a,z for this type int32",
			},
			{
				tag: "int8,validate:\"in:-256,78\"",
				value: struct {
					Age int8 `validate:"in:-256,78"`
				}{},
				errMsg: "tag `validate:\"in:-256,78\"` contains an invalid rule value -256,78 for this type int8",
			},
			{
				tag: "uint8,validate:\"in:-5,78\"",
				value: struct {
					Age uint8 `validate:"in:-5,78"`
				}{},
				errMsg: "tag `validate:\"in:-5,78\"` contains an invalid rule value -5,78 for this type uint8",
			},
			{
				tag: "uint16,validate:\"in:5,65935\"",
				value: struct {
					Age uint16 `validate:"in:5,65935"`
				}{},
				errMsg: "tag `validate:\"in:5,65935\"` contains an invalid rule value 5,65935 for this type uint16",
			},
			{
				tag: "uint32,validate:\"in:235,42954457295\"",
				value: struct {
					Age uint32 `validate:"in:235,42954457295"`
				}{},
				errMsg: "tag `validate:\"in:235,42954457295\"` contains an invalid rule value 235,42954457295 for this type uint32",
			},
			{
				tag: "uint64,validate:\"in:-245,8436346\"",
				value: struct {
					Age uint64 `validate:"in:-245,8436346"`
				}{},
				errMsg: "tag `validate:\"in:-245,8436346\"` contains an invalid rule value -245,8436346 for this type uint64",
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := newTagValidator(f)
				err := v.validateIn(f.Tags[0])
				require.EqualError(t, err, tc.errMsg)
			})
		}
	})

	t.Run("successful cases", func(t *testing.T) {
		tests := []struct {
			tag   reflect.StructTag
			value interface{}
		}{
			{tag: "string,validate:\"in:admin,user,dev\"", value: struct {
				Role string `validate:"in:admin,user,dev"`
			}{}},
			{tag: "string,validate:\"in:admin,user\"", value: struct {
				Role string `validate:"in:admin,user"`
			}{}},
			{tag: "string,validate:\"in:admin\"", value: struct {
				Role string `validate:"in:admin"`
			}{}},
			{tag: "int,validate:\"in:-3,44\"", value: struct {
				ID int `validate:"in:-3,44"`
			}{}},
			{tag: "int8,validate:\"in:-3,44\"", value: struct {
				ID int8 `validate:"in:-3,44"`
			}{}},
			{tag: "int16,validate:\"in:-3,44\"", value: struct {
				ID int16 `validate:"in:-3,44"`
			}{}},
			{tag: "int32,validate:\"in:-3,44\"", value: struct {
				ID int32 `validate:"in:-3,44"`
			}{}},
			{tag: "int64,validate:\"in:-3,44\"", value: struct {
				ID int64 `validate:"in:-3,44"`
			}{}},
			{tag: "uint,validate:\"in:3,44\"", value: struct {
				ID uint `validate:"in:3,44"`
			}{}},
			{tag: "uint8,validate:\"in:3,44\"", value: struct {
				ID uint8 `validate:"in:3,44"`
			}{}},
			{tag: "uint16,validate:\"in:3,44\"", value: struct {
				ID uint16 `validate:"in:3,44"`
			}{}},
			{tag: "uint32,validate:\"in:3,44\"", value: struct {
				ID uint32 `validate:"in:3,44"`
			}{}},
			{tag: "uint64,validate:\"in:3,44\"", value: struct {
				ID uint64 `validate:"in:3,44"`
			}{}},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := newTagValidator(f)
				err := v.validateIn(f.Tags[0])
				require.NoError(t, err)
			})
		}
	})
}
