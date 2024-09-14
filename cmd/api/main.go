package main

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/di/ctn/dic"
	"github.com/dembygenesis/local.tools/internal/api"
	"github.com/dembygenesis/local.tools/internal/api/resource"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/database/migration"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"log"
)

func migrate(cfg *config.App) error {
	c := &mysqlutil.ConnectionSettings{
		Host:     cfg.MysqlDatabaseCredentials.Host,
		User:     cfg.MysqlDatabaseCredentials.User,
		Pass:     cfg.MysqlDatabaseCredentials.Pass,
		Port:     cfg.MysqlDatabaseCredentials.Port,
		Database: cfg.MysqlDatabaseCredentials.Database,
	}

	_log := logger.New(context.Background())
	_log.Info("Migrating...")
	ctx := context.Background()
	_, err := mysqlhelper.MigrateEmbedded(ctx, c, migration.Migrations, mysqlhelper.CreateIfNotExists, false)
	if err != nil {
		return fmt.Errorf("migrate: %v", err)
	}
	return nil
}

func main() {
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("builder: %v", err)
	}

	ctn := builder.Build()

	cfg, err := ctn.SafeGetConfigLayer()
	if err != nil {
		log.Fatalf("cfg: %v", err)
	}

	_logger, err := ctn.SafeGetLoggerLogrus()
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	categoryMgr, err := ctn.SafeGetLogicCategory()
	if err != nil {
		log.Fatalf("category mgr: %v", err)
	}

	resourceGetter, err := resource.New(categoryMgr)
	if err != nil {
		log.Fatalf("resource getter: %v", err)
	}

	apiCfg := &api.Config{
		BaseUrl:         cfg.API.BaseUrl,
		Logger:          _logger,
		Port:            cfg.API.Port,
		CategoryService: categoryMgr,
		Resource:        resourceGetter,
	}

	if err := migrate(cfg); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	_api, err := api.New(apiCfg)
	if err != nil {
		log.Fatalf("api: %v", err)
	}

	if err := _api.Listen(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
