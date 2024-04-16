package strsrv

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"strings"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type StringUtils interface {
	CopyRootPathToClipboard(root string, exclude []string) ([]string, error)
}

//counterfeiter:generate . osLayer
type osLayer interface {
	CopyRootPathToClipboard(root string, exclude []string) ([]string, error)
}

func New(conf *config.App, osLayer osLayer) (StringUtils, error) {
	if conf == nil {
		return nil, errors.New(sysconsts.ErrConfigNil)
	}
	return &stringUtils{conf, osLayer}, nil
}

type stringUtils struct {
	conf    *config.App
	osLayer osLayer
}

func (s *stringUtils) CopyRootPathToClipboard(root string, exclude []string) ([]string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return nil, errors.New(sysconsts.ErrRootMissing)
	}

	if exclude == nil {
		exclude = make([]string, 0)
	}

	exclude = append(exclude, s.conf.CopyToClipboard.Exclusions...)

	files, err := s.osLayer.CopyRootPathToClipboard(root, exclude)
	if err != nil {
		return nil, fmt.Errorf("os: %v", err)
	}

	return files, nil
}
