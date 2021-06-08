package env_reader

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/spf13/afero"

	"github.com/IASamoylov/otus_home_work/hw08_envdir_tool/env_reader/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("wrap err when read dir fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockOS := mocks.NewMockOS(ctrl)
		mockOS.EXPECT().ReadDir("./path").Return(nil, errors.New("fake error"))

		ctx := NewContext(mockOS)
		_, err := ctx.ReadDir("./path")

		var envReaderErr *envReaderErr

		require.ErrorAs(t, err, &envReaderErr)
		require.EqualError(t, envReaderErr, "Errors occurred while getting environment from the directory")
		require.EqualError(t, envReaderErr.Unwrap(), "fake error")
	})

	t.Run("ignore invalid files", func(t *testing.T) {
		memFS := afero.NewMemMapFs()
		memFS.Mkdir("test", 0777)
		mockOS := mocks.NewMockFS(memFS)
		ctx := NewContext(mockOS)

		t.Run("file is dir", func(t *testing.T) {
			memFS.Mkdir("test/invalid_file", 0777)
			env, err := ctx.ReadDir("test")
			require.NoError(t, err)
			require.Empty(t, env)
		})

		t.Run("file has invalid name", func(t *testing.T) {
			afero.WriteFile(memFS, "test/ARG=", []byte("file c"), 0644)
			env, err := ctx.ReadDir("test")
			require.NoError(t, err)
			require.Empty(t, env)
		})
	})

	t.Run("parse environment", func(t *testing.T) {
		memFS := afero.NewMemMapFs()
		memFS.Mkdir("test", 0777)
		mockOS := mocks.NewMockFS(memFS)
		ctx := NewContext(mockOS)

		t.Run("read only first line", func(t *testing.T) {
			afero.WriteFile(memFS, "test/RUNTIME_ENVIRONMENT", []byte("CONTAINER\n\rnew_line"), 0644)
			env, err := ctx.ReadDir("test")

			require.NoError(t, err)
			require.Contains(t, env, "RUNTIME_ENVIRONMENT")
			require.Equal(t, env["RUNTIME_ENVIRONMENT"], EnvValue{
				"RUNTIME_ENVIRONMENT",
				"CONTAINER",
				false,
			})
		})

		t.Run("trim ", func(t *testing.T) {
			runes := []struct {
				b []byte
			}{
				{[]byte("\x00")}, {[]byte("\r")}, {[]byte("\n")},
				{[]byte("\t")}, {[]byte("\v")}, {[]byte("\f")},
				{[]byte(" ")}, {[]byte(string(rune(0x85)))}, {[]byte(string(rune(0x85)))},
			}

			for _, tc := range runes {
				tc := tc
				t.Run(fmt.Sprintf("%X", tc.b), func(t *testing.T) {
					data := [][]byte{
						[]byte("CONTAINER"),
						tc.b,
					}
					afero.WriteFile(memFS, "test/RUNTIME_ENVIRONMENT", bytes.Join(data, []byte("")), 0644)
					env, err := ctx.ReadDir("test")

					require.NoError(t, err)
					require.Contains(t, env, "RUNTIME_ENVIRONMENT")
					require.Equal(t, env["RUNTIME_ENVIRONMENT"], EnvValue{
						"RUNTIME_ENVIRONMENT",
						"CONTAINER",
						false,
					})
				})
			}
		})

		t.Run("remove env when file is empty", func(t *testing.T) {
			afero.WriteFile(memFS, "test/RUNTIME_ENVIRONMENT", []byte{}, 0644)
			env, err := ctx.ReadDir("test")

			require.NoError(t, err)
			require.Contains(t, env, "RUNTIME_ENVIRONMENT")
			require.Equal(t, env["RUNTIME_ENVIRONMENT"], EnvValue{
				Name:       "RUNTIME_ENVIRONMENT",
				NeedRemove: true,
			})
		})
	})
}
