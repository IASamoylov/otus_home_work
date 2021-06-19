package field

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMaxTagValidator(t *testing.T) {
	t.Run("fail cases", func(t *testing.T) {
		tests := []struct {
			tag    reflect.StructTag
			value  interface{}
			errMsg string
		}{
			{
				tag: "time.Time,validate:\"max:2999.07.25\"",
				value: struct {
					Birthday time.Time `validate:"max:2999.07.25"`
				}{},
				errMsg: "tag `validate:\"max:2999.07.25\"` not supported for this type time.Time",
			},
			{
				tag: "int,validate:\"max:Z\"",
				value: struct {
					Birthday int `validate:"max:Z"`
				}{},
				errMsg: "tag `validate:\"max:Z\"` contains an invalid rule value Z for this type int",
			},
			{
				tag: "int8,validate:\"max:-256\"",
				value: struct {
					Age int8 `validate:"max:-256"`
				}{},
				errMsg: "tag `validate:\"max:-256\"` contains an invalid rule value -256 for this type int8",
			},
			{
				tag: "uint8,validate:\"max:-5\"",
				value: struct {
					Age uint8 `validate:"max:-5"`
				}{},
				errMsg: "tag `validate:\"max:-5\"` contains an invalid rule value -5 for this type uint8",
			},
			{
				tag: "uint16,validate:\"max:-5\"",
				value: struct {
					Age uint16 `validate:"max:-5"`
				}{},
				errMsg: "tag `validate:\"max:-5\"` contains an invalid rule value -5 for this type uint16",
			},
			{
				tag: "uint32,validate:\"max:-5\"",
				value: struct {
					Age uint32 `validate:"max:-5"`
				}{},
				errMsg: "tag `validate:\"max:-5\"` contains an invalid rule value -5 for this type uint32",
			},
			{
				tag: "uint64,validate:\"max:-5\"",
				value: struct {
					Age uint64 `validate:"max:-5"`
				}{},
				errMsg: "tag `validate:\"max:-5\"` contains an invalid rule value -5 for this type uint64",
			},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := NewTagValidator(f)
				err := v.validateMax(f.Tags[0])
				require.EqualError(t, err, tc.errMsg)
			})
		}
	})

	t.Run("successful cases", func(t *testing.T) {
		tests := []struct {
			tag   reflect.StructTag
			value interface{}
		}{
			{tag: "int,validate:\"max:-36\"", value: struct {
				ID int `validate:"max:-36"`
			}{}},
			{tag: "int8,validate:\"max:-36\"", value: struct {
				ID int8 `validate:"max:-36"`
			}{}},
			{tag: "int16,validate:\"max:-36\"", value: struct {
				ID int16 `validate:"max:-36"`
			}{}},
			{tag: "int32,validate:\"max:-36\"", value: struct {
				ID int32 `validate:"max:-36"`
			}{}},
			{tag: "int64,validate:\"max:-36\"", value: struct {
				ID int64 `validate:"max:-36"`
			}{}},
			{tag: "uint,validate:\"max:36\"", value: struct {
				ID uint `validate:"max:36"`
			}{}},
			{tag: "uint8,validate:\"max:36\"", value: struct {
				ID uint8 `validate:"max:36"`
			}{}},
			{tag: "uint16,validate:\"max:36\"", value: struct {
				ID uint16 `validate:"max:36"`
			}{}},
			{tag: "uint32,validate:\"max:36\"", value: struct {
				ID uint32 `validate:"max:36"`
			}{}},
			{tag: "uint64,validate:\"max:36\"", value: struct {
				ID uint64 `validate:"max:36"`
			}{}},
		}

		for _, tc := range tests {
			t.Run(string(tc.tag), func(t *testing.T) {
				f := New(reflect.ValueOf(tc.value), 0)
				v := NewTagValidator(f)
				err := v.validateMax(f.Tags[0])
				require.NoError(t, err)
			})
		}
	})
}
