package cli

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utility"
)

type Service struct {
	stringUtils stringUtils
	gptUtils    gptUtils
	fileUtils   fileUtils
}

func NewService(
	stringUtils stringUtils,
	gptUtils gptUtils,
	fileUtils fileUtils,
) *Service {
	return &Service{
		stringUtils,
		gptUtils,
		fileUtils,
	}
}

func (s *Service) CopyToClipboard(root string, exclude []string) ([]string, error) {
	files, err := s.stringUtils.CopyRootPathToClipboard(root, exclude)
	if err != nil {
		return nil, fmt.Errorf("copy to clipboard: %v", err)
	}
	return files, nil
}

func (s *Service) ClipCodingStandardsPreface() error {
	err := s.gptUtils.ClipCodingStandardsPreface()
	if err != nil {
		return fmt.Errorf("clip coding standards: %v", err)
	}
	return nil
}

func (s *Service) CopyDirToAnother(opts *utility.CopyOptions) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	err := s.fileUtils.CopyDirToAnother(opts)
	if err != nil {
		return fmt.Errorf("copy folder A to B: %v", err)
	}
	return nil
}
