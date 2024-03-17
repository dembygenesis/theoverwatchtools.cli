package dependencies

import (
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/sarulabs/dingo/v4"
)

const (
	configLayer = "config_layer"
)

func GetConfigs() []dingo.Def {
	return []dingo.Def{
		{
			Name: configLayer,
			Build: func() (*config.App, error) {
				return config.New()
			},
		},
	}
}
