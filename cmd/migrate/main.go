package main

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/globals"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	helpers2 "github.com/dembygenesis/local.tools/internal/persistence/mysql_helpers"
	"os"
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

	c := &helpers2.ConnectionParameters{
		Host:     cfg.MysqlDatabaseCredentials.Host,
		User:     cfg.MysqlDatabaseCredentials.User,
		Pass:     cfg.MysqlDatabaseCredentials.Pass,
		Database: cfg.MysqlDatabaseCredentials.Database,
		Port:     cfg.MysqlDatabaseCredentials.Port,
	}

	log.Info("Migrating...")
	ctx := context.Background()
	migrationDir := os.Getenv(globals.OsEnvMigrationDir)
	tables, err := helpers2.Migrate(ctx, c, migrationDir, helpers2.CreateIfNotExists)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Infof("Created tables: %v", tables)
}
