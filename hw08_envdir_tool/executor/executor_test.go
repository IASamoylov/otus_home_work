package executor

import (
	"bytes"
	"errors"
	"io"
	"testing"

	common_mocks "github.com/IASamoylov/otus_home_work/hw08_envdir_tool/common/mocks"
	envreader "github.com/IASamoylov/otus_home_work/hw08_envdir_tool/env_reader"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("failed to remove env", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockOS := common_mocks.NewMockOSFunctions(ctrl)
		var buffer bytes.Buffer
		ctx := NewContext(mockOS, nil, nil, io.Writer(&buffer))
		mockOS.EXPECT().Unsetenv("ENV").Return(errors.New("fake error"))

		evn := envreader.Environment{
			"ENV": {Value: "5", NeedRemove: true, Name: "ENV"},
		}
		result := ctx.RunCmd([]string{
			"TEST",
		}, evn)

		require.EqualValues(t, buffer.Bytes(), []byte("errors when deleting an environment variable ENV\n"))
		require.Equal(t, result, 1)
	})

	t.Run("failed to set env", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockOS := common_mocks.NewMockOSFunctions(ctrl)
		var buffer bytes.Buffer
		ctx := NewContext(mockOS, nil, nil, io.Writer(&buffer))
		mockOS.EXPECT().Unsetenv("ENV").Return(nil)
		mockOS.EXPECT().Setenv("ENV", "5").Return(errors.New("fake error"))

		evn := envreader.Environment{
			"ENV": {Value: "5", NeedRemove: false, Name: "ENV"},
		}
		result := ctx.RunCmd([]string{
			"TEST",
		}, evn)

		require.EqualValues(t, buffer.Bytes(), []byte("errors when setting an environment variable ENV with value 5\n"))
		require.Equal(t, result, 1)
	})
}
