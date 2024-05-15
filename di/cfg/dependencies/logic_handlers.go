package dependencies

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/authlogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

const (
	logicCategory = "logic_category"
	logicAuth     = "logic_auth"
)

func GetLogicHandlers() []dingo.Def {
	return []dingo.Def{
		{
			Name: logicCategory,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
				store *mysqlstore.Repository,
			) (*categorylogic.Service, error) {
				logic, err := categorylogic.New(&categorylogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
					Persistor:  store,
				})
				if err != nil {
					return nil, fmt.Errorf("logicategory: %v", err)
				}
				return logic, nil
			},
		},
		{
			Name: logicAuth,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
			) (*authlogic.Impl, error) {
				logic, err := authlogic.New(&authlogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
				})
				if err != nil {
					return nil, fmt.Errorf("logicauth: %v", err)
				}
				return logic, nil
			},
		},
	}
}
