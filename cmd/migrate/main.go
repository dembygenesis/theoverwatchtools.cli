package main

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/globals"
	"github.com/dembygenesis/local.tools/internal/persistence/mysql/helpers"
	"os"
)

func main() {
	var (
		cfg *config.Config
		err error
	)

	logger := common.GetLogger(context.TODO())

	cfg, err = config.New(".env")
	if err != nil {
		logger.Fatalf("cfg: %v", err.Error())
	}

	c := &helpers.ConnectionParameters{
		Host:     cfg.MysqlDatabaseCredentials.Host,
		User:     cfg.MysqlDatabaseCredentials.User,
		Pass:     cfg.MysqlDatabaseCredentials.Pass,
		Database: cfg.MysqlDatabaseCredentials.Database,
		Port:     cfg.MysqlDatabaseCredentials.Port,
	}

	logger.Info("Migrating...")
	ctx := context.Background()
	migrationDir := os.Getenv(globals.OsEnvMigrationDir)
	tables, err := helpers.Migrate(ctx, c, migrationDir, helpers.CreateIfNotExists)
	if err != nil {
		logger.Fatalf("migrate: %v", err)
	}

	logger.Infof("Created tables: %v", tables)
}
