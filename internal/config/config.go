package config

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utils_generic"
)

type CopyToClipboard struct {
	Exclusions []string `json:"exclusions" mapstructure:"exclusions"`
}

func (c *CopyToClipboard) ParseExclusions(s string) error {
	err := utils_generic.DecodeToStruct(s, &c.Exclusions)
	if err != nil {
		return fmt.Errorf("exclusions decode: %v", err)
	}
	if len(c.Exclusions) == 0 {
		return errors.New("exclusions are empty")
	}
	return nil
}

type TransferFiles struct {
	Exclusions []string `json:"exclusions" mapstructure:"exclusions"`
}

func (c *TransferFiles) ParseExclusions(s string) error {
	err := utils_generic.DecodeToStruct(s, &c.Exclusions)
	if err != nil {
		return fmt.Errorf("exclusions decode: %v", err)
	}
	if len(c.Exclusions) == 0 {
		return errors.New("exclusions are empty")
	}
	return nil
}

type Config struct {
	TransferFiles   TransferFiles   `json:"transfer_files"`
	CopyToClipboard CopyToClipboard `json:"copy_to_clipboard"`
}

func New() (*Config, error) {
	var err error

	config := Config{}

	if err = config.CopyToClipboard.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal copy to clipboard: %v", err)
	}

	if err = config.TransferFiles.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal transfer files: %v", err)
	}

	return &config, nil
}
