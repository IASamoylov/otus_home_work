package field

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagParser(t *testing.T) {
	tests := []struct {
		Name string
		tag  reflect.StructTag
		tags []Tag
	}{
		{Name: "ignore not package tags", tag: "json:\"id\"", tags: []Tag{}},
		{Name: "ignore package tag that configured invalid", tag: "validate:", tags: []Tag{}},
		{Name: "zero Value when tag configured without rule", tag: "validate:\"\"", tags: []Tag{
			{ValueIsUndefined: true},
		}},
		{Name: "value will be undefined if incorrectly configured a rule", tag: "validate:\"len:\"", tags: []Tag{
			{Name: LenTagValidation, ValueIsUndefined: true},
		}},
		{Name: "value will be undefined if incorrectly configured a rule", tag: "validate:\"len\"", tags: []Tag{
			{Name: LenTagValidation, ValueIsUndefined: true},
		}},
		{Name: "the tag will be parsed correctly if the rule is configured correctly", tag: "validate:\"len:46\"", tags: []Tag{
			{Name: LenTagValidation, Value: "46"},
		}},
		{Name: "the tag will be parsed correctly if the rule is configured correctly", tag: "validate:\"min:46|max:132\"", tags: []Tag{
			{Name: MinTagValidation, Value: "46"},
			{Name: MaxTagValidation, Value: "132"},
		}},
		{Name: "one of the values will be undefined if incorrectly configured a rule", tag: "validate:\"min:46|max:\"", tags: []Tag{
			{Name: MinTagValidation, Value: "46"},
			{Name: MaxTagValidation, ValueIsUndefined: true},
		}},
		{Name: "one of the values will be undefined if incorrectly configured a rule", tag: "validate:\"min|max:132\"", tags: []Tag{
			{Name: MinTagValidation, ValueIsUndefined: true},
			{Name: MaxTagValidation, Value: "132"},
		}},
		{Name: "one of the values will be undefined if incorrectly configured a rule", tag: "validate:\"min:64|\"", tags: []Tag{
			{Name: MinTagValidation, Value: "64"},
			{ValueIsUndefined: true},
		}},
		{Name: "all values will be undefined if incorrectly configured a rule", tag: "validate:\"|\"", tags: []Tag{
			{ValueIsUndefined: true},
			{ValueIsUndefined: true},
		}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v:%v", tc.Name, tc.tag), func(t *testing.T) {
			result := NewTagParser(tc.tag).Parse()

			require.Len(t, result, len(tc.tags))

			for i, tag := range tc.tags {
				require.Equal(t, tc.tag, result[i].Tag)
				require.Equal(t, tag.Value, result[i].Value)
				require.Equal(t, tag.Name, result[i].Name)
				require.Equal(t, tag.ValueIsUndefined, result[i].ValueIsUndefined)
			}
		})
	}
}
