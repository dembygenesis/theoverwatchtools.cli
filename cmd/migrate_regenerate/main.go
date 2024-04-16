package main

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/global"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"os"
	"os/exec"
)

func main() {
	var (
		cfg *config.App
		err error
	)

	ctx := context.TODO()
	log := logger.New(ctx)

	cfg, err = config.New()
	if err != nil {
		log.Fatalf("cfg: %v", err.Error())
	}

	c := &mysqlutil.ConnectionSettings{
		Host:     cfg.MysqlDatabaseCredentials.Host,
		User:     cfg.MysqlDatabaseCredentials.User,
		Pass:     cfg.MysqlDatabaseCredentials.Pass,
		Port:     cfg.MysqlDatabaseCredentials.Port,
		Database: cfg.MysqlDatabaseCredentials.Database,
	}

	db, err := mysqlhelper.NewDbClient(ctx, &mysqlhelper.ClientOptions{
		PingTimeout: 0,
		ConnString:  c.GetConnectionString(true),
		Close:       false,
	})
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	recreateStmts := []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", cfg.MysqlDatabaseCredentials.Database),
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", cfg.MysqlDatabaseCredentials.Database),
	}

	for _, stmt := range recreateStmts {
		_, err = db.ExecContext(ctx, stmt)
		if err != nil {
			log.Fatalf("stmt: %s, err: %v", stmt, err)
		}
	}

	log.Info("Migrating...")
	tables, err := mysqlhelper.Migrate(ctx, c, mysqlhelper.CreateIfNotExists)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Infof("Created tables: %v", tables)

	log.Infof("Generating sqlboiler tables...")
	type cmd struct {
		Name string
		Args []string
	}

	sqlboilerCfg := os.Getenv(global.OsEnvAppDir) + "/sqlboiler.toml"

	commands := []cmd{
		{
			Name: "rm",
			Args: []string{"-rf", "./internal/persistence/database_helpers/mysql/assets/mysqlmodel"},
		},
		{
			Name: "sqlboiler",
			Args: []string{"mysql", "--no-hooks", "--no-tests", "--config", sqlboilerCfg},
		},
	}

	for _, command := range commands {
		execCmd := exec.Command(command.Name, command.Args...)
		_, err = execCmd.CombinedOutput()
		if err != nil {
			log.Fatalf("err for cm: %s, err: %v", command.Name, err)
		}
	}

	log.Infof("Completed generating sqlboiler tables")
}
