package config

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utils_common"
)

type CopyToClipboard struct {
	Exclusions []string `json:"exclusions" mapstructure:"exclusions"`
}

func (c *CopyToClipboard) ParseExclusions(s string) error {
	err := utils_common.DecodeToStruct(s, &c.Exclusions)
	if err != nil {
		return fmt.Errorf("exclusions decode: %v", err)
	}
	if len(c.Exclusions) == 0 {
		return errors.New("exclusions are empty")
	}
	return nil
}

type FolderAToFolderB struct {
	GenericExclusions []string `json:"generic_exclusions" mapstructure:"generic_exclusions"`
}

func (c *FolderAToFolderB) ParseExclusions(s string) error {
	err := utils_common.DecodeToStruct(s, &c.GenericExclusions)
	if err != nil {
		return fmt.Errorf("exclusions decode: %v", err)
	}
	if len(c.GenericExclusions) == 0 {
		return errors.New("exclusions are empty")
	}
	return nil
}

type Config struct {
	FolderAToFolderB FolderAToFolderB `json:"folder_a_to_folder_b"`
	CopyToClipboard  CopyToClipboard  `json:"copy_to_clipboard"`
}

func New() (*Config, error) {
	var err error

	config := Config{}

	if err = config.CopyToClipboard.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal copy to clipboard: %v", err)
	}

	if err = config.FolderAToFolderB.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal transfer files: %v", err)
	}

	return &config, nil
}
