package dependencies

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

const (
	persistenceMySQL = "persistence_mysql"
)

func GetPersistence() []dingo.Def {
	return []dingo.Def{
		{
			Name: persistenceMySQL,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				db *sqlx.DB,
				txHandler *mysqltx.Handler,
			) (*mysqlstore.Repository, error) {
				store, err := mysqlstore.New(&mysqlstore.Config{
					Logger: logger,
					QueryTimeouts: &persistence.QueryTimeouts{
						Query: cfg.Timeouts.DbQuery,
						Exec:  cfg.Timeouts.DbExec,
					},
				})
				if err != nil {
					return nil, fmt.Errorf("store: %v", err)
				}

				return store, nil
			},
		},
	}
}
