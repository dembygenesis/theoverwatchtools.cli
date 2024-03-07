package main

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlutil"
)

func main() {
	var (
		cfg *config.Config
		err error
	)

	log := logger.New(context.TODO())

	cfg, err = config.New(".env")
	if err != nil {
		log.Fatalf("cfg: %v", err.Error())
	}

	c := &mysqlutil.ConnectionSettings{
		Host:     cfg.MysqlDatabaseCredentials.Host,
		User:     cfg.MysqlDatabaseCredentials.User,
		Pass:     cfg.MysqlDatabaseCredentials.Pass,
		Database: cfg.MysqlDatabaseCredentials.Database,
		Port:     cfg.MysqlDatabaseCredentials.Port,
	}

	log.Info("Migrating...")
	ctx := context.Background()
	tables, err := mysqlhelper.Migrate(ctx, c, mysqlhelper.CreateIfNotExists)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Infof("Created tables: %v", tables)
}
