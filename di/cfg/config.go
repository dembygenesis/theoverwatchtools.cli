package cfg

import (
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/sarulabs/dingo/v4"
)

const (
	configLayer = "config_layer"
)

func getConfigLayer() []dingo.Def {
	return []dingo.Def{
		{
			Name: configLayer,
			Build: func() (*config.Config, error) {
				return config.New(".env")
			},
		},
	}
}
