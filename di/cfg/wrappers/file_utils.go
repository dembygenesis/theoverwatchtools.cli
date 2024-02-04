package wrappers

import (
	"github.com/dembygenesis/local.tools/internal/models"
	"github.com/dembygenesis/local.tools/internal/utils_generic"
)

func NewFileUtilsWrapper() *FileWrapper {
	return &FileWrapper{}
}

type FileWrapper struct {
}

func (f *FileWrapper) CopyDir(opts *models.CopyOptions) error {
	return utils_generic.CopyDir(opts)
}
