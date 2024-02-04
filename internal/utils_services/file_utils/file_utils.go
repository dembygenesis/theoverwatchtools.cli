package file_utils

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/models"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type FileUtils interface {
	CopyDir(opts *models.CopyOptions) error
}

//counterfeiter:generate . osLayer
type osLayer interface {
	CopyDir(opts *models.CopyOptions) error
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

func (g *fileUtils) CopyDir(opts *models.CopyOptions) error {
	if opts == nil {
		return fmt.Errorf("opts nil")
	}

	opts.WipeDestination = true
	opts.WipeDestinationExclusions = g.conf.TransferFiles.Exclusions

	err := g.osLayer.CopyDir(opts)
	if err != nil {
		return fmt.Errorf("os: %v", err)
	}

	return nil
}
