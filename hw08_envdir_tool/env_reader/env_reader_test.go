package env_reader

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/IASamoylov/otus_home_work/hw08_envdir_tool/env_reader/mocks"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
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

	t.Run("wrap err when info of file fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockDirEntry := mocks.NewMockDirEntry(ctrl)
		mockOS := mocks.NewMockOS(ctrl)
		mockDirEntry.EXPECT().Info().Return(nil, errors.New("fake error"))
		mockDirEntry.EXPECT().Name().Return("RUNTIME_ENVIRONMENT").AnyTimes()
		mockDirEntry.EXPECT().IsDir().Return(false)
		mockOS.EXPECT().ReadDir("./path").Return([]os.DirEntry{mockDirEntry}, nil)

		ctx := NewContext(mockOS)
		_, err := ctx.ReadDir("./path")

		var envReaderErr *envReaderErr

		require.ErrorAs(t, err, &envReaderErr)
		require.EqualError(t, envReaderErr, "Error processing file [RUNTIME_ENVIRONMENT]")
		require.EqualError(t, envReaderErr.Unwrap(), "fake error")
	})

	t.Run("wrap err when info of file fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFileInfo := mocks.NewMockFileInfo(ctrl)
		mockFileInfo.EXPECT().Size().Return(int64(1 << 10))
		mockDirEntry := mocks.NewMockDirEntry(ctrl)
		mockDirEntry.EXPECT().Info().Return(mockFileInfo, nil)
		mockDirEntry.EXPECT().Name().Return("RUNTIME_ENVIRONMENT").AnyTimes()
		mockDirEntry.EXPECT().IsDir().Return(false)
		mockOS := mocks.NewMockOS(ctrl)
		mockOS.EXPECT().ReadDir("./path").Return([]os.DirEntry{mockDirEntry}, nil)
		mockOS.EXPECT().Open("path/RUNTIME_ENVIRONMENT").Return(nil, errors.New("fake error"))

		ctx := NewContext(mockOS)
		_, err := ctx.ReadDir("./path")

		var envReaderErr *envReaderErr

		require.ErrorAs(t, err, &envReaderErr)
		require.EqualError(t, envReaderErr, "Error processing file [RUNTIME_ENVIRONMENT]")
		require.EqualError(t, envReaderErr.Unwrap(), "fake error")
	})

	t.Run("ignore invalid files", func(t *testing.T) {
		memFS := afero.NewMemMapFs()
		memFS.Mkdir("test", 0o777)
		mockOS := mocks.NewAferoWrapper(memFS)
		ctx := NewContext(mockOS)

		t.Run("file is dir", func(t *testing.T) {
			memFS.Mkdir("test/invalid_file", 0o777)
			env, err := ctx.ReadDir("test")
			require.NoError(t, err)
			require.Empty(t, env)
		})

		t.Run("file has invalid name", func(t *testing.T) {
			afero.WriteFile(memFS, "test/ARG=", []byte("file c"), 0o644)
			env, err := ctx.ReadDir("test")
			require.NoError(t, err)
			require.Empty(t, env)
		})
	})

	t.Run("parse environment", func(t *testing.T) {
		memFS := afero.NewMemMapFs()
		memFS.Mkdir("test", 0o777)
		mockOS := mocks.NewAferoWrapper(memFS)
		ctx := NewContext(mockOS)

		t.Run("read only first line", func(t *testing.T) {
			tests := []struct {
				name string
				data string
			}{
				{"0x00", "CONTAINER\u0000with new line"},
				{"default line separator", "CONTAINER\n\rnew_line"},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					afero.WriteFile(memFS, "test/RUNTIME_ENVIRONMENT", []byte(tc.data), 0o644)
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

		t.Run("trim invalid symbols", func(t *testing.T) {
			runes := []struct {
				b []byte
			}{
				{[]byte("\r")},
				{[]byte("\n")},
				{[]byte(" ")},
				{[]byte("\t")},
				{[]byte("\v")},
				{[]byte("\f")},
				{[]byte(string(rune(0x85)))},
				{[]byte(string(rune(0x85)))},
			}

			for _, tc := range runes {
				tc := tc
				t.Run(fmt.Sprintf("%X", tc.b), func(t *testing.T) {
					data := [][]byte{
						[]byte("CONTAINER"),
						tc.b,
					}
					afero.WriteFile(memFS, "test/RUNTIME_ENVIRONMENT", bytes.Join(data, []byte("")), 0o644)
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
			tests := []struct {
				name string
				data []byte
			}{
				{"empty file", []byte{}},
				{"only new line symbols", []byte("\u0000\n\r")},
			}
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					afero.WriteFile(memFS, "test/RUNTIME_ENVIRONMENT", tc.data, 0o644)
					env, err := ctx.ReadDir("test")

					require.NoError(t, err)
					require.Contains(t, env, "RUNTIME_ENVIRONMENT")
					require.Equal(t, env["RUNTIME_ENVIRONMENT"], EnvValue{
						Name:       "RUNTIME_ENVIRONMENT",
						NeedRemove: true,
					})
				})
			}
		})

		t.Run("remove env when file has length 0 bytes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockFileInfo := mocks.NewMockFileInfo(ctrl)
			mockFileInfo.EXPECT().Size().Return(int64(0))
			mockDirEntry := mocks.NewMockDirEntry(ctrl)
			mockDirEntry.EXPECT().Info().Return(mockFileInfo, nil)
			mockDirEntry.EXPECT().Name().Return("RUNTIME_ENVIRONMENT").AnyTimes()
			mockDirEntry.EXPECT().IsDir().Return(false)
			mockOS := mocks.NewMockOS(ctrl)
			mockOS.EXPECT().ReadDir("./path").Return([]os.DirEntry{mockDirEntry}, nil)

			ctx := NewContext(mockOS)
			env, err := ctx.ReadDir("./path")

			require.NoError(t, err)
			require.Contains(t, env, "RUNTIME_ENVIRONMENT")
			require.Equal(t, env["RUNTIME_ENVIRONMENT"], EnvValue{
				Name:       "RUNTIME_ENVIRONMENT",
				NeedRemove: true,
			})
		})
	})
}
