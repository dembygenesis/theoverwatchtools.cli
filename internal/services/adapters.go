package services

import "github.com/dembygenesis/local.tools/internal/models"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . stringUtils
type stringUtils interface {
	CopyRootPathToClipboard(root string, exclude []string) ([]string, error)
}

//counterfeiter:generate . gptUtils
type gptUtils interface {
	ClipCodingStandardsPreface() error
}

//counterfeiter:generate . fileUtils
type fileUtils interface {
	CopyDir(opts *models.CopyOptions) error
}
