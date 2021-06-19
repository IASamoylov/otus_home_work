package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexpTagValidator(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
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
				f := New(reflect.ValueOf(tc.value), 0)
				v := NewTagValidator(f)
				err := v.validateRegexp(f.Tags[0])
				require.EqualError(t, err, tc.errMsg)

			})
		}
	})

	t.Run("successful cases", func(t *testing.T) {
		tests := []struct {
			tag   reflect.StructTag
			value interface{}
		}{
			{tag: "string,validate:\"regexp:\\d+\"", value: struct {
				ID string `validate:"regexp:\\d+"`
			}{}},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := NewTagValidator(f)
				err := v.validateRegexp(f.Tags[0])
				require.NoError(t, err)
			})
		}
	})
}
