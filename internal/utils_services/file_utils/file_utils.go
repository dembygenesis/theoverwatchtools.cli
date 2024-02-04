package file_utils

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/models"
	"github.com/dembygenesis/local.tools/internal/utils_common"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type FileUtils interface {
	CopyDirToAnother(opts *utils_common.CopyOptions) error
}

//counterfeiter:generate . osLayer
type osLayer interface {
	CopyDirToAnother(opts *utils_common.CopyOptions) error
}

func New(conf *config.Config, osLayer osLayer) (FileUtils, error) {
	if conf == nil {
		return nil, models.ErrConfigNil
	}

	return &fileUtils{
		conf,
		osLayer,
	}, nil
}

type fileUtils struct {
	conf    *config.Config
	osLayer osLayer
}

func (g *fileUtils) CopyDirToAnother(opts *utils_common.CopyOptions) error {
	if opts == nil {
		return fmt.Errorf("opts nil")
	}

	opts.WipeDestination = true
	opts.SourceExclusions = g.conf.FolderAToFolderB.GenericExclusions
	opts.WipeDestinationExclusions = g.conf.FolderAToFolderB.GenericExclusions

	err := g.osLayer.CopyDirToAnother(opts)
	if err != nil {
		return fmt.Errorf("os: %v", err)
	}

	return nil
}
