package fileutil

import (
	"errors"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
	"github.com/dembygenesis/local.tools/internal/services/fileutil/fileutilfakes"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_fileUtils_CopyDir_Fail_Nil_Opts(t *testing.T) {
	conf := config.Config{}
	fakeFileUtils, _ := New(&conf, &fileutilfakes.FakeOsLayer{})

	err := fakeFileUtils.CopyDirToAnother(nil)
	require.Error(t, err, "expected opts nil error")
}

func Test_fileUtils_CopyDir_Fail_Os_Layer(t *testing.T) {
	conf := config.Config{}
	fakeOsLayer := fileutilfakes.FakeOsLayer{}

	fakeFileUtils, _ := New(&conf, &fakeOsLayer)

	fakeOsLayer.CopyDirToAnotherReturns(errors.New("mock error"))

	opts := fslib.CopyOptions{}

	err := fakeFileUtils.CopyDirToAnother(&opts)
	require.Error(t, err, "expected opts nil error")
	require.Contains(t, err.Error(), "os:")
	require.Contains(t, err.Error(), "mock error")
}

func Test_fileUtils_CopyDir_Success(t *testing.T) {
	conf := config.Config{}
	fakeOsLayer := fileutilfakes.FakeOsLayer{}

	fakeFileUtils, _ := New(&conf, &fakeOsLayer)

	opts := fslib.CopyOptions{}

	err := fakeFileUtils.CopyDirToAnother(&opts)
	require.NoError(t, err, "expected opts has error")
}
