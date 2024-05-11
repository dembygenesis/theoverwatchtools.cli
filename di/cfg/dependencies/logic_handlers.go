package dependencies

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/organizationlogic"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

const (
	logicOrganization = "logic_organization"
)

func GetLogicHandlers() []dingo.Def {
	return []dingo.Def{
		{
			Name: logicOrganization,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
				store *mysqlstore.Repository,
			) (*organizationlogic.Service, error) {
				logic, err := organizationlogic.New(&organizationlogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
					Persistor:  store,
				})
				if err != nil {
					return nil, fmt.Errorf("logicorganization: %v", err)
				}
				return logic, nil
			},
		},
	}
}
