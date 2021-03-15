package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCapacityCalculation(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: "ф3a4", expected: 10},
		{input: "日4本2語1", expected: 21},
		{input: "♥4", expected: 12},
		{input: "a4bc2d5e", expected: 13},
		{input: "abccd", expected: 5},
		{input: "", expected: 0},
		{input: "aaa0b", expected: 3},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result := capacity(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "ф3a4", expected: "фффaaaa"},
		{input: "日4本2語1", expected: "日日日日本本語"},
		{input: "♥4", expected: "♥♥♥♥"},
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
