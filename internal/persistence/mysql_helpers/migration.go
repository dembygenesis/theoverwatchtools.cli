package helpers

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"strings"
	"time"
)

type Tables [][]string

func Migrate(ctx context.Context, cp *ConnectionParameters, dir string, mode CreateMode) (Tables, error) {
	if err := cp.Validate(true); err != nil {
		return nil, fmt.Errorf("validate conn parameters: %v", err)
	}

	if strings.TrimSpace(dir) == "" {
		return nil, models.ErrEmptyDir
	}

	_, err := NewClient(ctx, &ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("get client: %v", err)
	}

	_, err = NewBuildTeardown(ctx, cp, &NewOpts{
		Mode:  mode,
		Close: true,
	})
	if err != nil {
		return nil, fmt.Errorf("mysql new: %v", err)
	}

	// Execute migration statements
	m, err := migrate.New(
		fmt.Sprintf("file://%s", dir),
		fmt.Sprintf("mysql://%v", cp.GetConnectionString(false)),
	)

	if err != nil {
		return nil, fmt.Errorf("migrate new: %v", err)
	}

	if err = m.Up(); err != nil {
		if !strings.Contains(err.Error(), "no change") {
			return nil, fmt.Errorf("up: %v", err)
		}
	}

	showTablesStmt := fmt.Sprintf(`
			SELECT table_name FROM information_schema.tables
			WHERE table_schema = '%v'
		`, cp.Database)

	conn, err := NewConnection(&Settings{
		ConnectTimeout:       time.Second * 10,
		QueryTimeout:         time.Second * 10,
		ExecTimeout:          time.Second * 10,
		ConnectionParameters: cp,
	})
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	tables, _, err := conn.QueryIntoStructSettings(ctx, &QueryAsArrFilter{
		Stmt:       showTablesStmt,
		Pagination: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("query into struct: %v", err)
	}

	return tables, nil
}
