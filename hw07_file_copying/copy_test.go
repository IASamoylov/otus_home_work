package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	testCase    string
	copyOptions CopyOptions
}

type CopyOptions struct {
	fromPath string
	toPath   string
	offset   int64
	limit    int64
	err      error
}

func TestCopy(t *testing.T) {
	t.Run("Argument err when some arguments are invalid", func(t *testing.T) {
		tests := []TestCase{
			{testCase: "fromPath is empty", copyOptions: CopyOptions{"", "./test", 0, 0, NewErrArgumentPath("from")}},
			{testCase: "toPath is empty", copyOptions: CopyOptions{"./test", "", 0, 0, NewErrArgumentPath("to")}},
			{testCase: "offset is negative", copyOptions: CopyOptions{"./test", "./test_", -1, 0, NewErrArgumentNegative("offset")}},
			{testCase: "limit is negative", copyOptions: CopyOptions{"./test", "./test_", 0, -1, NewErrArgumentNegative("limit")}},
			{testCase: "to and from path the same", copyOptions: CopyOptions{"./test.txt", "./test.txt", 0, 0, ErrPathTheSame}},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.testCase, func(t *testing.T) {
				err := Copy(tc.copyOptions.fromPath, tc.copyOptions.toPath, tc.copyOptions.offset, tc.copyOptions.limit)
				require.Error(t, tc.copyOptions.err, err)
			})
		}
	})

	t.Run("", func(t *testing.T) {
		err := Copy("./dirToFile", "./dirToFile/out.txt", 0, 0)
		require.Nil(t, err)
	})
}
