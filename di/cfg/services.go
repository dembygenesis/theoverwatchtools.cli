package cfg

import (
	"github.com/dembygenesis/local.tools/di/cfg/wrappers"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/services"
	"github.com/dembygenesis/local.tools/internal/utils_services/file_utils"
	"github.com/dembygenesis/local.tools/internal/utils_services/gpt_utils"
	"github.com/dembygenesis/local.tools/internal/utils_services/string_utils"
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
			) (*services.Services, error) {
				fileUtils, err := file_utils.New(cfg, wrappers.NewFileUtilsWrapper())
				if err != nil {
					return nil, err
				}

				stringUtils, err := string_utils.New(cfg, wrappers.NewStringUtilsWrapper())
				if err != nil {
					return nil, err
				}

				return services.NewServices(
					stringUtils,
					gpt_utils.New(),
					fileUtils,
				), nil
			},
		},
	}
}
