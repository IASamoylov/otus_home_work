package env_reader

import (
	"errors"
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
}
