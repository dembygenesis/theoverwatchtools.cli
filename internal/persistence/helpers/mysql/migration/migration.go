package migration

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql/connection"
	"github.com/golang-migrate/migrate/v4"
	"time"
)

type Tables [][]string

func Migrate(ctx context.Context, cp *mysql.ConnectionParameters, mode mysql.CreateMode) (Tables, error) {
	if err := cp.Validate(true); err != nil {
		return nil, fmt.Errorf("validate conn parameters: %v", err)
	}

	_, err := mysql.GetClient(ctx, &mysql.ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("get client: %v", err)
	}

	_, err = mysql.New(ctx, cp, &mysql.NewOpts{
		Mode:  mode,
		Close: true,
	})
	if err != nil {
		return nil, fmt.Errorf("mysql new: %v", err)
	}

	// Execute migration statements
	m, err := migrate.New(
		"file://internal/database/migrations",
		fmt.Sprintf("mysql://%v", cp.GetConnectionString(false)),
	)
	if err != nil {
		return nil, fmt.Errorf("migrate new: %v", err)
	}

	if err = m.Up(); err != nil {
		return nil, fmt.Errorf("up: %v", err)
	}

	// Show all created tables
	showTablesStmt := fmt.Sprintf(`
			SELECT table_name FROM information_schema.tables
			WHERE table_schema = '%v'
		`, cp.Database)

	conn, err := connection.New(&connection.Settings{
		ConnectTimeout:       time.Second * 10,
		QueryTimeout:         time.Second * 10,
		ExecTimeout:          time.Second * 10,
		ConnectionParameters: cp,
	})
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	tables, _, err := conn.QueryIntoStructSettings(ctx, &connection.QueryAsArrFilter{
		Stmt:       showTablesStmt,
		Pagination: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("query into struct: %v", err)
	}

	return tables, nil
}
