package cli

import (
	"errors"
	"github.com/dembygenesis/local.tools/internal/cli/clifakes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServices_CopyToClipboard_Success(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	srv := Service{
		stringUtils: &mockStringUtils,
		gptUtils:    &mockGptUtils,
		fileUtils:   &mockFileUtils,
	}

	_, err := srv.CopyToClipboard(".", nil)
	require.NoError(t, err, "should have no error")
}

func TestServices_CopyToClipboard_Fail(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	mockStringUtils.CopyRootPathToClipboardReturns(nil, errors.New("mock error"))

	srv := Service{
		stringUtils: &mockStringUtils,
		gptUtils:    &mockGptUtils,
		fileUtils:   &mockFileUtils,
	}

	_, err := srv.CopyToClipboard(".", nil)
	require.Error(t, err, "should have no error")
	require.Contains(t, err.Error(), "mock error")
	require.Contains(t, err.Error(), "copy to clipboard:")
}

func TestServices_ClipCodingStandardsPreface_Success(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	srv := Service{
		stringUtils: &mockStringUtils,
		gptUtils:    &mockGptUtils,
		fileUtils:   &mockFileUtils,
	}

	err := srv.ClipCodingStandardsPreface()
	require.NoError(t, err, "should have no error")
}

func TestServices_New(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	_ = NewService(
		&mockStringUtils,
		&mockGptUtils,
		&mockFileUtils,
	)
}

func TestServices_ClipCodingStandardsPreface_Fail(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	mockGptUtils.ClipCodingStandardsPrefaceReturns(errors.New("mock error"))

	srv := Service{
		stringUtils: &mockStringUtils,
		gptUtils:    &mockGptUtils,
		fileUtils:   &mockFileUtils,
	}

	err := srv.ClipCodingStandardsPreface()
	require.Error(t, err, "should have an error")
	require.Contains(t, err.Error(), "mock error")
	require.Contains(t, err.Error(), "clip coding standards:")
}
