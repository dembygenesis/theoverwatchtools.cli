package fileutil

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
	"github.com/dembygenesis/local.tools/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type FileUtils interface {
	CopyDirToAnother(opts *fslib.CopyOptions) error
}

//counterfeiter:generate . osLayer
type osLayer interface {
	CopyDirToAnother(opts *fslib.CopyOptions) error
}

func New(conf *config.Config, osLayer osLayer) (FileUtils, error) {
	if conf == nil {
		return nil, model.ErrConfigNil
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

func (g *fileUtils) CopyDirToAnother(opts *fslib.CopyOptions) error {
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
