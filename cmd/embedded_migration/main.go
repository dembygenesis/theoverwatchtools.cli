package main

import (
	"context"
	"github.com/dembygenesis/local.tools/di/ctn/dic"
	"github.com/dembygenesis/local.tools/internal/database/migration"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"log"
)

func migrate() {
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("builder: %v", err)
	}

	ctn := builder.Build()
	cfg, err := ctn.SafeGetConfigLayer()
	if err != nil {
		log.Fatalf("get cfg: %v", err)
	}

	_log := logger.New(context.Background())

	c := &mysqlutil.ConnectionSettings{
		Host:     cfg.MysqlDatabaseCredentials.Host,
		User:     cfg.MysqlDatabaseCredentials.User,
		Pass:     cfg.MysqlDatabaseCredentials.Pass,
		Port:     cfg.MysqlDatabaseCredentials.Port,
		Database: cfg.MysqlDatabaseCredentials.Database,
	}

	_log.Info("Migrating")
	ctx := context.Background()
	tables, err := mysqlhelper.MigrateEmbedded(ctx, c, migration.Migrations, mysqlhelper.CreateIfNotExists)
	if err != nil {
		_log.Fatalf("migrate: %v", err)
	}

	_log.Infof("Created tables: %v", tables)
}

func main() {
	migrate()
}
