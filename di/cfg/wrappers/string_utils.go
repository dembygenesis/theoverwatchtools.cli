package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/utils_common"
)

func NewStringUtilsWrapper() *StringWrapper {
	return &StringWrapper{}
}

type StringWrapper struct {
}

func (f *StringWrapper) CopyRootPathToClipboard(root string, exclude []string) ([]string, error) {
	return utils_common.CopyRootPathToClipboard(root, exclude)
}
