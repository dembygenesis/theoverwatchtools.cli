package cli

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
)

type Service struct {
	stringUtils stringService
	gptUtils    gptService
	fileUtils   fileService
}

func NewService(
	stringUtils stringService,
	gptUtils gptService,
	fileUtils fileService,
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

func (s *Service) CopyDirToAnother(opts *fslib.CopyOptions) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	err := s.fileUtils.CopyDirToAnother(opts)
	if err != nil {
		return fmt.Errorf("copy folder A to B: %v", err)
	}
	return nil
}
