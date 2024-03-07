package cfg

import (
	"github.com/dembygenesis/local.tools/di/cfg/wrappers"
	"github.com/dembygenesis/local.tools/internal/cli"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/services/fileutil"
	"github.com/dembygenesis/local.tools/internal/services/gptutil"
	"github.com/dembygenesis/local.tools/internal/services/strsrvutil"
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
				fileUtils, err := fileutil.New(cfg, wrappers.NewFileUtilsWrapper())
				if err != nil {
					return nil, err
				}

				stringUtils, err := strsrvutil.New(cfg, wrappers.NewStringUtilsWrapper())
				if err != nil {
					return nil, err
				}

				return cli.NewService(
					stringUtils,
					gptutil.New(),
					fileUtils,
				), nil
			},
		},
	}
}
