package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/utils_common"
)

func NewFileUtilsWrapper() *FileWrapper {
	return &FileWrapper{}
}

type FileWrapper struct {
}

func (f *FileWrapper) CopyDirToAnother(opts *utils_common.CopyOptions) error {
	return utils_common.CopyDirToAnother(opts)
}
