package services

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utils_common"
)

type Services struct {
	stringUtils stringUtils
	gptUtils    gptUtils
	fileUtils   fileUtils
}

func NewServices(
	stringUtils stringUtils,
	gptUtils gptUtils,
	fileUtils fileUtils,
) *Services {
	return &Services{
		stringUtils,
		gptUtils,
		fileUtils,
	}
}

func (s *Services) CopyToClipboard(root string, exclude []string) ([]string, error) {
	files, err := s.stringUtils.CopyRootPathToClipboard(root, exclude)
	if err != nil {
		return nil, fmt.Errorf("copy to clipboard: %v", err)
	}
	return files, nil
}

func (s *Services) ClipCodingStandardsPreface() error {
	err := s.gptUtils.ClipCodingStandardsPreface()
	if err != nil {
		return fmt.Errorf("clip coding standards: %v", err)
	}
	return nil
}

func (s *Services) CopyDirToAnother(opts *utils_common.CopyOptions) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	err := s.fileUtils.CopyDirToAnother(opts)
	if err != nil {
		return fmt.Errorf("copy folder A to B: %v", err)
	}
	return nil
}
