package config

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utils_common"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type MysqlDatabaseCredentials struct {
	Host     string `json:"host" mapstructure:"DB_HOST" validate:"required"`
	User     string `json:"user" mapstructure:"DB_USER" validate:"required"`
	Pass     string `json:"pass" mapstructure:"DB_PASS" validate:"required"`
	Port     int    `json:"port" mapstructure:"DB_PORT" validate:"required"`
	Database string `json:"database" mapstructure:"DB_DATABASE" validate:"required"`
}

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
	FolderAToFolderB         FolderAToFolderB         `json:"folder_a_to_folder_b"`
	CopyToClipboard          CopyToClipboard          `json:"copy_to_clipboard"`
	MysqlDatabaseCredentials MysqlDatabaseCredentials `json:"mysq_database_credentials"`
}

// isProduction checks if the `IS_PRODUCTION` envVar isset
func isProduction() bool {
	return os.Getenv("IS_PRODUCTION") != ""
}

func New(envFile string) (*Config, error) {
	var err error
	if isProduction() {
		for _, envVar := range os.Environ() {
			split := strings.SplitN(envVar, "=", 2)
			key := split[0]
			val := split[1]
			viper.Set(key, val)
		}
	} else {
		viper.SetConfigFile(envFile)

		err = viper.ReadInConfig()
		if err != nil {
			log.Fatalf("error reading from config: %v", err)
		}
		viper.AutomaticEnv()
	}

	config := Config{}

	if err = config.CopyToClipboard.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal copy to clipboard: %v", err)
	}

	if err = config.FolderAToFolderB.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal transfer files: %v", err)
	}

	err = viper.Unmarshal(&config.MysqlDatabaseCredentials)
	if err != nil {
		return &config, fmt.Errorf("error trying to unmarshal the database credentials: %w", err)
	}

	return &config, nil
}
