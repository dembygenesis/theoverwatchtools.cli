package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/utilities/cliputil"
)

func NewStringUtilsWrapper() *StringWrapper {
	return &StringWrapper{}
}

type StringWrapper struct {
}

func (f *StringWrapper) CopyRootPathToClipboard(root string, exclude []string) ([]string, error) {
	return cliputil.CopyRootPathToClipboard(root, exclude)
}
