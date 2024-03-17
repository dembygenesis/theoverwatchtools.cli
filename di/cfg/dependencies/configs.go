package dependencies

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/sarulabs/dingo/v4"
	"os"
)

const (
	configLayer = "config_layer"
)

func GetConfigs() []dingo.Def {
	return []dingo.Def{
		{
			Name: configLayer,
			Build: func() (*config.App, error) {
				dir := fmt.Sprintf("%s/%s", os.Getenv("APP_DIR"), ".env")
				return config.New(dir)
			},
		},
	}
}
