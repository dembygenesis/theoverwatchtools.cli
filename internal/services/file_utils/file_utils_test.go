package file_utils

import (
	"errors"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/services/file_utils/file_utilsfakes"
	"github.com/dembygenesis/local.tools/internal/utility"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_fileUtils_CopyDir_Fail_Nil_Opts(t *testing.T) {
	conf := config.Config{}
	fakeFileUtils, _ := New(&conf, &file_utilsfakes.FakeOsLayer{})

	err := fakeFileUtils.CopyDirToAnother(nil)
	require.Error(t, err, "expected opts nil error")
}

func Test_fileUtils_CopyDir_Fail_Os_Layer(t *testing.T) {
	conf := config.Config{}
	fakeOsLayer := file_utilsfakes.FakeOsLayer{}

	fakeFileUtils, _ := New(&conf, &fakeOsLayer)

	fakeOsLayer.CopyDirToAnotherReturns(errors.New("mock error"))

	opts := utility.CopyOptions{}

	err := fakeFileUtils.CopyDirToAnother(&opts)
	require.Error(t, err, "expected opts nil error")
	require.Contains(t, err.Error(), "os:")
	require.Contains(t, err.Error(), "mock error")
}

func Test_fileUtils_CopyDir_Success(t *testing.T) {
	conf := config.Config{}
	fakeOsLayer := file_utilsfakes.FakeOsLayer{}

	fakeFileUtils, _ := New(&conf, &fakeOsLayer)

	opts := utility.CopyOptions{}

	err := fakeFileUtils.CopyDirToAnother(&opts)
	require.NoError(t, err, "expected opts has error")
}
