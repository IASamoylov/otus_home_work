package field

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestField(t *testing.T) {
	t.Run("has validation tags works fine", func(t *testing.T) {
		tests := []struct {
			name     string
			excepted bool
			tags     []Tag
		}{
			{name: "tags is empty", excepted: false, tags: []Tag{}},
			{name: "tags is not empty", excepted: true, tags: make([]Tag, 1)},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				field := &Field{Tags: tc.tags}
				require.Equal(t, tc.excepted, field.HasValidationTags())
			})
		}
	})
}
