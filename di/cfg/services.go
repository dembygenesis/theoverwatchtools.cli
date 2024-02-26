package cfg

import (
	"github.com/dembygenesis/local.tools/di/cfg/wrappers"
	"github.com/dembygenesis/local.tools/internal/cli"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/services/file_utils"
	"github.com/dembygenesis/local.tools/internal/services/gpt_utils"
	"github.com/dembygenesis/local.tools/internal/services/string_utils"
	"github.com/sarulabs/dingo/v4"
)

const (
	servicesLayer = "services_layer"
)

func getServicesLayer() []dingo.Def {
	return []dingo.Def{
		{
			Name: servicesLayer,
			Build: func(
				cfg *config.Config,
			) (*cli.Service, error) {
				fileUtils, err := file_utils.New(cfg, wrappers.NewFileUtilsWrapper())
				if err != nil {
					return nil, err
				}

				stringUtils, err := string_utils.New(cfg, wrappers.NewStringUtilsWrapper())
				if err != nil {
					return nil, err
				}

				return cli.NewService(
					stringUtils,
					gpt_utils.New(),
					fileUtils,
				), nil
			},
		},
	}
}
