package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/utils/clipboard_util"
)

func NewStringUtilsWrapper() *StringWrapper {
	return &StringWrapper{}
}

type StringWrapper struct {
}

func (f *StringWrapper) CopyRootPathToClipboard(root string, exclude []string) ([]string, error) {
	return clipboard_util.CopyRootPathToClipboard(root, exclude)
}
