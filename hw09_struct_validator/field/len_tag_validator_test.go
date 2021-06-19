package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLenTagValidator(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
			{
				tag: "string,validate:\"len:\"",
				value: struct {
					ID string `validate:"len:"`
				}{},
				errMsg: "tag `validate:\"len:\"` contains an invalid rule value  for this type string",
			},
			{
				tag: "string,validate:\"len:afasf\"",
				value: struct {
					ID string `validate:"len:afasf"`
				}{},
				errMsg: "tag `validate:\"len:afasf\"` contains an invalid rule value afasf for this type string",
			},
			{
				tag: "int32,validate:\"len:36\"",
				value: struct {
					ID int32 `validate:"len:36"`
				}{},
				errMsg: "tag `validate:\"len:36\"` not supported for this type int32",
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := NewTagValidator(f)
				err := v.validateLen(f.Tags[0])
				require.EqualError(t, err, tc.errMsg)
			})
		}
	})

	t.Run("successful cases", func(t *testing.T) {
		tests := []struct {
			tag   reflect.StructTag
			value interface{}
		}{
			{tag: "string,validate:\"len:36\"", value: struct {
				ID string `validate:"len:36"`
			}{}},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := NewTagValidator(f)
				err := v.validateLen(f.Tags[0])
				require.NoError(t, err)
			})
		}
	})
}
