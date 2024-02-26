package main

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql/migration"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	var (
		cfg *config.Config
		err error
	)

	ctx := context.Background()
	logger := common.GetLogger(ctx)

	cfg, err = config.New(".env")
	if err != nil {
		log.Fatalf("error trying to initialize the builder: %v", err.Error())
	}

	c := mysql.ConnectionParameters{
		Host:     cfg.MysqlDatabaseCredentials.Host,
		User:     cfg.MysqlDatabaseCredentials.User,
		Pass:     cfg.MysqlDatabaseCredentials.Pass,
		Database: cfg.MysqlDatabaseCredentials.Database,
		Port:     cfg.MysqlDatabaseCredentials.Port,
	}

	tables, err := migration.Migrate(ctx, &c, mysql.Recreate)
	if err != nil {
		logger.Fatalf("migrate: %v", err)
	}

	logger.Info("Created tables", tables)
}
