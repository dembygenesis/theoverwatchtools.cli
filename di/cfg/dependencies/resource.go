package dependencies

import (
	apirsrc "github.com/dembygenesis/local.tools/internal/api/resource"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

const (
	resourceGetter = "resource_getter"
)

func GetResourceGetter() []dingo.Def {
	return []dingo.Def{
		{
			Name: resourceGetter,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				cat *categorylogic.Service,
			) (*apirsrc.Provider, error) {
				return apirsrc.New(cat)
			},
		},
	}
}
