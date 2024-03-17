package dependencies

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/dingo/v4"
	"time"
)

const (
	dbMySQL = "db_mysql"
)

func GetDatabases() []dingo.Def {
	return []dingo.Def{
		{
			Name: dbMySQL,
			Build: func(cfg *config.App) (*sqlx.DB, error) {
				connStr := mysqlutil.ConnectionSettings{
					Host:     cfg.MysqlDatabaseCredentials.Host,
					User:     cfg.MysqlDatabaseCredentials.User,
					Pass:     cfg.MysqlDatabaseCredentials.Pass,
					Database: cfg.MysqlDatabaseCredentials.Database,
					Port:     cfg.MysqlDatabaseCredentials.Port,
				}

				opts := &mysqlhelper.ClientOptions{
					PingTimeout: time.Second * 10,
					ConnString:  connStr.GetConnectionString(false),
					Close:       false,
				}

				db, err := mysqlhelper.NewDbClient(context.TODO(), opts)
				if err != nil {
					return nil, fmt.Errorf("mysql client: %v", err)
				}

				return db, nil
			},
		},
	}
}
