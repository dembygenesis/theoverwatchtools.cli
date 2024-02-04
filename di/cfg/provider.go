package cfg

import (
	"fmt"
	"github.com/sarulabs/dingo/v4"
)

type Provider struct {
	dingo.BaseProvider
}

// getServices is the main configuration func that produces the singleton
func getServices() ([]dingo.Def, error) {
	var services []dingo.Def

	services = append(services, getServicesLayer()...)
	services = append(services, getConfigLayer()...)

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
