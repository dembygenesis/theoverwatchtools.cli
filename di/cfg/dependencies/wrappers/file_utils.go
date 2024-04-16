package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
)

func NewFileUtilsWrapper() *FileWrapper {
	return &FileWrapper{}
}

type FileWrapper struct {
}

func (f *FileWrapper) CopyDirToAnother(opts *fslib.CopyOptions) error {
	return fslib.CopyDirToAnother(opts)
}
