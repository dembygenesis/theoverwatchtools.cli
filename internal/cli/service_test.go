package cli

import (
	"errors"
	"github.com/dembygenesis/local.tools/internal/cli/clifakes"
	"github.com/dembygenesis/local.tools/internal/utils_common"
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

func TestServices_CopyDirToAnother_Success(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	srv := Service{
		stringUtils: &mockStringUtils,
		gptUtils:    &mockGptUtils,
		fileUtils:   &mockFileUtils,
	}

	opts := &utils_common.CopyOptions{
		Source:                    "./testA",
		SourceExclusions:          nil,
		Destination:               "./testB",
		WipeDestination:           false,
		WipeDestinationExclusions: nil,
	}

	err := srv.CopyDirToAnother(opts)

	require.NoError(t, err, "should have no error")
}

func TestServices_CopyDirToAnother_Validate_Fail(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	srv := Service{
		stringUtils: &mockStringUtils,
		gptUtils:    &mockGptUtils,
		fileUtils:   &mockFileUtils,
	}

	opts := &utils_common.CopyOptions{}

	err := srv.CopyDirToAnother(opts)

	require.Error(t, err, "expected an error due to missing source and destination")
	require.Contains(t, err.Error(), "validate:")
	require.Contains(t, err.Error(), "'source' must have a value")
	require.Contains(t, err.Error(), "'destination' must have a value")
}

func TestServices_CopyDirToAnother_Opts_Fail(t *testing.T) {
	mockStringUtils := clifakes.FakeStringUtils{}
	mockGptUtils := clifakes.FakeGptUtils{}
	mockFileUtils := clifakes.FakeFileUtils{}

	mockFileUtils.CopyDirToAnotherReturns(errors.New("forced error in copy operation"))
	srv := Service{
		stringUtils: &mockStringUtils,
		gptUtils:    &mockGptUtils,
		fileUtils:   &mockFileUtils,
	}

	opts := &utils_common.CopyOptions{
		Source:                    "./testA",
		SourceExclusions:          nil,
		Destination:               "./testB",
		WipeDestination:           false,
		WipeDestinationExclusions: nil,
	}

	err := srv.CopyDirToAnother(opts)

	require.Error(t, err, "expected an error from copy operation")
}
