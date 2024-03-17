package dependencies

import (
	"github.com/dembygenesis/local.tools/di/cfg/dependencies/wrappers"
	"github.com/dembygenesis/local.tools/internal/cli"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/services/gptsrv"
	"github.com/sarulabs/dingo/v4"
)

const (
	serviceCli = "service_cli"
)

func GetServicesLayer() []dingo.Def {
	return []dingo.Def{
		{
			Name: serviceCli,
			Build: func(
				cfg *config.App,
			) (*cli.Service, error) {
				fileUtil := wrappers.NewFileUtilsWrapper()
				strUtil := wrappers.NewStringUtilsWrapper()
				gptUtil := gptsrv.New()

				return cli.NewService(strUtil, gptUtil, fileUtil), nil
			},
		},
	}
}
