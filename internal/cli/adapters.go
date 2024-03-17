package cli

import (
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . stringService
type stringService interface {
	CopyRootPathToClipboard(root string, exclude []string) ([]string, error)
}

//counterfeiter:generate . gptService
type gptService interface {
	ClipCodingStandardsPreface() error
}

//counterfeiter:generate . fileService
type fileService interface {
	CopyDirToAnother(opts *fslib.CopyOptions) error
}
