package dependencies

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

const (
	txHandler  = "tx_handler"
	txProvider = "tx_provider"
)

func GetTransactions() []dingo.Def {
	return []dingo.Def{
		{
			Name: txProvider,
			Build: func(
				logger *logrus.Entry,
				db *sqlx.DB,
				txHandler *mysqltx.Handler,
			) (*mysqlconn.Provider, error) {
				provider, err := mysqlconn.New(&mysqlconn.Config{
					Logger:    logger,
					TxHandler: txHandler,
				})
				if err != nil {
					return nil, fmt.Errorf("tx provider: %v", err)
				}

				return provider, nil
			},
		},
		{
			Name: txHandler,
			Build: func(
				logger *logrus.Entry,
				db *sqlx.DB,
				config *config.App,
			) (*mysqltx.Handler, error) {
				handler, err := mysqltx.New(&mysqltx.Config{
					Logger:       logger,
					Db:           db,
					DatabaseName: config.MysqlDatabaseCredentials.Database,
				})
				if err != nil {
					return nil, fmt.Errorf("tx handler: %v", err)
				}

				return handler, nil
			},
		},
	}
}
