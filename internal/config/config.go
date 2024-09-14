package config

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/errutil"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/spf13/viper"
	"time"
)

var (
	errEnvEmpty = "env '%s' is unset"
)

type Timeouts struct {
	DbExec  time.Duration `json:"DB_QUERY_TIMEOUT" mapstructure:"DB_QUERY_TIMEOUT" validate:"required,is_positive_time_duration"`
	DbQuery time.Duration `json:"DB_EXEC_TIMEOUT" mapstructure:"DB_EXEC_TIMEOUT" validate:"required,is_positive_time_duration"`
}

type MysqlDatabaseCredentials struct {
	Host     string `json:"host" mapstructure:"DB_HOST" validate:"required"`
	User     string `json:"user" mapstructure:"DB_USER" validate:"required"`
	Pass     string `json:"pass" mapstructure:"DB_PASS" validate:"required"`
	Port     int    `json:"port" mapstructure:"DB_PORT" validate:"required"`
	Database string `json:"database" mapstructure:"DB_DATABASE" validate:"required"`
}

type CopyToClipboard struct {
	Exclusions []string `json:"exclusions"`
}

func (c *CopyToClipboard) ParseExclusions(s string) error {
	err := sliceutil.Decode(s, &c.Exclusions)
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
	err := sliceutil.Decode(s, &c.GenericExclusions)
	if err != nil {
		return fmt.Errorf("exclusions decode: %w", err)
	}
	if len(c.GenericExclusions) == 0 {
		return errors.New("exclusions are empty")
	}
	return nil
}

type API struct {
	BaseUrl        string        `json:"base_url" mapstructure:"API_BASE_URL" validate:"required"`
	Port           int           `json:"port" mapstructure:"API_PORT" validate:"required,greater_than_zero"`
	ListenTimeout  time.Duration `json:"listen_timeout" mapstructure:"API_LISTEN_TIMEOUT_SECS" validate:"required,is_positive_time_duration"`
	RequestTimeout time.Duration `json:"request_timeout" mapstructure:"API_REQUEST_TIMEOUT_SECS" validate:"required,is_positive_time_duration"`
}

type Settings struct {
	IsProduction bool   `json:"PRODUCTION" mapstructure:"PRODUCTION" validate:"boolean"`
	AppDir       string `json:"APP_DIR" mapstructure:"APP_DIR" validate:"required"`
}

type App struct {
	Settings                 Settings                 `json:"settings"`
	FolderAToFolderB         FolderAToFolderB         `json:"folder_a_to_folder_b"`
	CopyToClipboard          CopyToClipboard          `json:"copy_to_clipboard"`
	MysqlDatabaseCredentials MysqlDatabaseCredentials `json:"mysql_database_credentials"`
	API                      API                      `json:"API"`
	Timeouts                 Timeouts                 `json:"Timeouts"`
}

func New() (*App, error) {
	viper.Reset()

	// Set env prefix
	viper.SetEnvPrefix("THEOVERWATCHTOOLS")

	// Set app details
	viper.SetDefault("APP_DIR", "/app")
	viper.SetDefault("PRODUCTION", false)

	// Set database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "demby")
	viper.SetDefault("DB_PASS", "secret")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_DATABASE", "example")
	viper.SetDefault("DB_EXEC_TIMEOUT", "10s")
	viper.SetDefault("DB_QUERY_TIMEOUT", "10s")

	// Set API defaults
	viper.SetDefault("API_PORT", 3000)
	viper.SetDefault("API_LISTEN_TIMEOUT_SECS", "10s")
	viper.SetDefault("API_REQUEST_TIMEOUT_SECS", "10s")
	viper.SetDefault("API_BASE_URL", "http://localhost")

	viper.AutomaticEnv()

	// Map configs to struct
	config := App{}

	if err := config.CopyToClipboard.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal copy to clipboard: %v", err)
	}

	if err := config.FolderAToFolderB.ParseExclusions(genericExclusions); err != nil {
		return &config, fmt.Errorf("unmarshal transfer files: %v", err)
	}

	err := viper.Unmarshal(&config.MysqlDatabaseCredentials)
	if err != nil {
		return &config, fmt.Errorf("error trying to unmarshal the database credentials: %w", err)
	}

	err = viper.Unmarshal(&config.Timeouts)
	if err != nil {
		return &config, fmt.Errorf("error trying to unmarshal the timeouts: %w", err)
	}

	err = viper.Unmarshal(&config.API)
	if err != nil {
		return nil, fmt.Errorf("unmarshal API cfg: %v", err)
	}

	err = viper.Unmarshal(&config.Settings)
	if err != nil {
		return nil, fmt.Errorf("unmarshal API cfg: %v", err)
	}

	cfgProperties := []interface{}{
		config.API,
		config.MysqlDatabaseCredentials,
		config.Settings,
		config.Timeouts,
	}

	var errs errutil.List
	for _, cfgProperty := range cfgProperties {
		err = validationutils.Validate(cfgProperty)
		if err != nil {
			errs.AddErr(err)
		}
	}

	if errs.HasErrors() {
		return nil, fmt.Errorf("cfg errors: %v", errs.Single())
	}

	return &config, nil
}
