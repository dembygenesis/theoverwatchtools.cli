package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	cli "github.com/dembygenesis/local.tools/internal/cli"
	config "github.com/dembygenesis/local.tools/internal/config"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "config_layer",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("config_layer")
				if err != nil {
					var eo *config.Config
					return eo, err
				}
				b, ok := d.Build.(func() (*config.Config, error))
				if !ok {
					var eo *config.Config
					return eo, errors.New("could not cast build function to func() (*config.Config, error)")
				}
				return b()
			},
			Unshared: false,
		},
		{
			Name:  "services_layer",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("services_layer")
				if err != nil {
					var eo *cli.Service
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *cli.Service
					return eo, err
				}
				p0, ok := pi0.(*config.Config)
				if !ok {
					var eo *cli.Service
					return eo, errors.New("could not cast parameter 0 to *config.Config")
				}
				b, ok := d.Build.(func(*config.Config) (*cli.Service, error))
				if !ok {
					var eo *cli.Service
					return eo, errors.New("could not cast build function to func(*config.Config) (*cli.Service, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
	}
}
