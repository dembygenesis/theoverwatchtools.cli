package cfg

import (
	"fmt"
	"github.com/dembygenesis/local.tools/di/cfg/dependencies"
	"github.com/sarulabs/dingo/v4"
)

type Provider struct {
	dingo.BaseProvider
}

// getServices is the main configuration func that produces the singleton
func getServices() ([]dingo.Def, error) {
	var services []dingo.Def

	layers := [][]dingo.Def{
		dependencies.GetConfigs(),
		dependencies.GetDatabases(),
		dependencies.GetLoggers(),
		dependencies.GetTransactions(),
		dependencies.GetLogicHandlers(),
		dependencies.GetServicesLayer(),
		dependencies.GetPersistence(),
		dependencies.GetResourceGetter(),
	}

	for _, layer := range layers {
		services = append(services, layer...)
	}

	return services, nil
}

// Load bootstrap the dependencies
func (p *Provider) Load() error {
	services, err := getServices()
	if err != nil {
		return fmt.Errorf("error trying to load the provider: %v", err)
	}

	err = p.AddDefSlice(services)
	if err != nil {
		return fmt.Errorf("error adding dependency definitions: %v", err)
	}

	return nil
}
