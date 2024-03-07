package main

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/api"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/managers/categorymgr"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlprovider"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlutil"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/category/categorymysql"
	"log"
	"time"
)

var (
	testLogger = logger.New(context.TODO())
)

func main() {
	cfg, err := config.New(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	settings := &mysqlhelper.Settings{
		ConnectTimeout: 10 * time.Second,
		QueryTimeout:   10 * time.Second,
		ExecTimeout:    10 * time.Second,
		ConnectionParameters: &mysqlutil.ConnectionSettings{
			Host:     cfg.MysqlDatabaseCredentials.Host,
			User:     cfg.MysqlDatabaseCredentials.User,
			Pass:     cfg.MysqlDatabaseCredentials.Pass,
			Database: cfg.MysqlDatabaseCredentials.Database,
			Port:     cfg.MysqlDatabaseCredentials.Port,
		},
	}
	conn, err := mysqlhelper.NewConnection(settings)
	if err != nil {
		log.Fatalf("new connection: %v", err)
	}

	db, err := conn.GetDBx()
	if err != nil {
		log.Fatalf("new connection: %v", err)
	}

	mysqlTxProvider, err := mysqlprovider.New(&mysqlprovider.Config{
		Db:     db,
		Logger: testLogger,
	})
	if err != nil {
		log.Fatalf("create maria db: %v", err)
	}

	catPersistor, err := categorymysql.New(&categorymysql.Config{
		Logger: logger.New(context.TODO()),
		QueryTimeouts: &persistence.QueryTimeouts{
			Query: 10 * time.Second,
			Exec:  10 * time.Second,
		},
	})
	if err != nil {
		log.Fatalf("category mysql: %v", err)
	}

	categoryMgr, err := categorymgr.New(&categorymgr.Config{
		Persistor:  catPersistor,
		TxProvider: mysqlTxProvider,
	})
	if err != nil {
		log.Fatalf("category mgr: %v", err)
	}

	a, err := api.New(&api.Config{
		Port:            3003,
		CategoryManager: categoryMgr,
	})
	if err != nil {
		log.Fatalf("new: %v", err)
	}

	if err = a.Listen(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
