package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/utility"
)

func NewFileUtilsWrapper() *FileWrapper {
	return &FileWrapper{}
}

type FileWrapper struct {
}

func (f *FileWrapper) CopyDirToAnother(opts *utility.CopyOptions) error {
	return utility.CopyDirToAnother(opts)
}
