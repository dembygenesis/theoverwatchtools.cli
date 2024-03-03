package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/lib/file_system"
)

func NewFileUtilsWrapper() *FileWrapper {
	return &FileWrapper{}
}

type FileWrapper struct {
}

func (f *FileWrapper) CopyDirToAnother(opts *file_system.CopyOptions) error {
	return file_system.CopyDirToAnother(opts)
}
