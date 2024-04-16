package config

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type Timeouts struct {
	DbExec  time.Duration `json:"TIMEOUT_DB_EXEC" mapstructure:"TIMEOUT_DB_EXEC" validate:"required,is_positive_time_duration"`
	DbQuery time.Duration `json:"TIMEOUT_DB_QUERY" mapstructure:"TIMEOUT_DB_QUERY" validate:"required,is_positive_time_duration"`
}

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
	IsProduction bool   `json:"THEOVERWATCHTOOLS_PRODUCTION" mapstructure:"THEOVERWATCHTOOLS_PRODUCTION" validate:"required,is_positive_time_duration"`
	AppDir       string `json:"THEOVERWATCHTOOLS_APP_DIR" mapstructure:"THEOVERWATCHTOOLS_APP_DIR" validate:"required,is_positive_time_duration"`
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
	var err error

	config := App{}
	for _, envVar := range os.Environ() {
		split := strings.SplitN(envVar, "=", 2)
		key := split[0]
		val := split[1]
		viper.Set(key, val)
	}

	err = viper.Unmarshal(&config.Settings)
	if err != nil {
		return &config, fmt.Errorf("error trying to unmarshal the database credentials: %w", err)
	}

	if !config.Settings.IsProduction {
		envFile := fmt.Sprintf("%s/.env", config.Settings.AppDir)
		_, err = os.Stat(envFile)
		if err != nil {
			return nil, fmt.Errorf("file stat: %v", err)
		}

		viper.SetConfigFile(envFile)

		err = viper.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("file state: %v", err)
		}
		viper.AutomaticEnv()
	}

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

	err = viper.Unmarshal(&config.Timeouts)
	if err != nil {
		return &config, fmt.Errorf("error trying to unmarshal the timeouts: %w", err)
	}

	err = viper.Unmarshal(&config.API)
	if err != nil {
		return nil, fmt.Errorf("unmarshal API cfg: %v", err)
	}

	cfgProperties := []interface{}{
		config.API,
	}

	var errs errs.List
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
