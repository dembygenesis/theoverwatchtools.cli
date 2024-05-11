package config

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
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

	fmt.Println("=== flip flip")

	// Set env prefix
	viper.SetEnvPrefix("THEOVERWATCHTOOLS")

	// Set app details
	viper.SetDefault("APP_DIR", "/app")
	viper.SetDefault("PRODUCTION", false)

	// THEOVERWATCHTOOLS

	// Set database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "demby")
	viper.SetDefault("DB_PASS", "secret")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_DATABASE", "example")
	viper.SetDefault("DB_EXEC_TIMEOUT", "10s")
	viper.SetDefault("DB_QUERY_TIMEOUT", "10s")

	// Override top if production == false
	if viper.Get("IS_PRODUCTION") != "1" {
		// Override DB_PORT, and DB_DATABASE.
	}

	// Set API defaults
	viper.SetDefault("API_PORT", 3000)
	viper.SetDefault("API_LISTEN_TIMEOUT_SECS", "10s")
	viper.SetDefault("API_REQUEST_TIMEOUT_SECS", "10s")
	viper.SetDefault("API_BASE_URL", "localhost")

	// Do some additional logic

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

func New3() (*App, error) {
	if os.Getenv(EnvAppDir) == "" {
		return nil, fmt.Errorf(errEnvEmpty, EnvAppDir)
	}

	envDir := fmt.Sprintf("%s/%s", os.Getenv(EnvAppDir), envFile)
	if strings.TrimSpace(envDir) != "" {
		if _, err := os.Stat(envDir); err != nil {
			return nil, fmt.Errorf("env file stat: %v", err)
		}

		if err := godotenv.Load(envDir); err != nil {
			return nil, fmt.Errorf("load env file: %v", err)
		}

		fmt.Println("==== passed this file stat:", envDir)
	}

	app := App{}
	if err := viper.Unmarshal(&app.MysqlDatabaseCredentials); err != nil {
		return nil, fmt.Errorf("unmarshal mysql db credentials: %v", err)
	}

	// THEOVERWATCHTOOLS_TIMEOUT_DB_USER
	a := os.Getenv("THEOVERWATCHTOOLS_DB_HOST")
	b := os.Getenv("DB_USER")

	fmt.Println("======== a:", a)
	fmt.Println("======== b:", b)
	fmt.Println("======== app.MysqlDatabaseCredentials:", strutil.GetAsJson(app.MysqlDatabaseCredentials))
	fmt.Println("======== app.MysqlDatabaseCredentials:", strutil.GetAsJson(app.MysqlDatabaseCredentials))

	return nil, nil
}

func NewOld() (*App, error) {
	fmt.Println("=== hehehe")
	var err error

	config := App{}

	// We should ditch this step, cause we only have our defaults, OR .env file...
	// How about the system env file? Well, too many layers, let's make it simple for now.
	/*for _, envVar := range os.Environ() {
		split := strings.SplitN(envVar, "=", 2)
		key := split[0]
		val := split[1]
		viper.Set(key, val)
	}*/
	viper.AutomaticEnv()

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

	fmt.Println("==== config.Timeouts:", strutil.GetAsJson(config.Timeouts))

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
