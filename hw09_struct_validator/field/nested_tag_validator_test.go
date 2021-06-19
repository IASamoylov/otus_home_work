package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNestedTagValidator(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
			{
				tag: "int32,validate:\"nested\"",
				value: struct {
					ID int32 `validate:"nested"`
				}{},
				errMsg: "tag `validate:\"nested\"` not supported for this type int32",
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
			{tag: "struct,validate:\"nested\"", value: struct {
				User struct{ ID int } `validate:"nested"`
			}{}},
			{tag: "struct,validate:\"nested\"", value: struct {
				User *struct{ ID int } `validate:"nested"`
			}{}},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := NewTagValidator(f)
				err := v.validateNested(f.Tags[0])
				require.NoError(t, err)
			})
		}
	})
}
