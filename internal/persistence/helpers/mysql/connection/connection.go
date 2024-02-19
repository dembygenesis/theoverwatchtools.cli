package connection

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/dembygenesis/local.tools/internal/utils_common"
	"github.com/jmoiron/sqlx"
	"time"
)

type Settings struct {
	// ConnectTimeout is the connect timeout duration
	ConnectTimeout time.Duration `mapstructure:"connect_timeout" validate:"required"`

	// QueryTimeout is the query (READ) timeout duration
	QueryTimeout time.Duration `mapstructure:"query_timeout" validate:"required"`

	// ExecTimeout is the query execution timeout duration
	ExecTimeout time.Duration `mapstructure:"exec_timeout" validate:"required"`

	// Parameters holds the database credentials
	Parameters *mysql.ConnectionParameters `mapstructure:"parameters"`
}

// validate validates the Connection settings
func (c *Settings) validate() (*sqlx.DB, error) {
	if err := utils_common.ValidateStruct(c); err != nil {
		return nil, fmt.Errorf("required: %s", err)
	}

	db, err := mysql.GetClient(c.Parameters)
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.ConnectTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping ctx: %w", err)
	}

	return db, nil
}

// Connection containers the helper method
type Connection struct {
	db       *sqlx.DB
	settings *Settings
}

func New(settings *Settings) (*Connection, error) {
	if settings == nil {
		return nil, errors.New("settings nil")
	}

	db, err := settings.validate()
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &Connection{
		settings: settings,
		db:       db,
	}, nil
}

func (c *Connection) GetDBx() (*sqlx.DB, error) {
	if c.db == nil {
		return nil, errors.New("hahaha")
	}
	return c.db, nil
}

// Exec executes the query with respect to the timeout.
func (c *Connection) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, c.settings.ExecTimeout)
	defer cancel()
	result, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}
	return result, nil
}

// Query executes a query with respect to the
// timeout.
func (c *Connection) Query(
	ctx context.Context,
	dest interface{},
	query string,
	args ...interface{},
) error {
	ctx, cancel := context.WithTimeout(ctx, c.settings.QueryTimeout)
	defer cancel()
	err := c.db.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return fmt.Errorf("query: %v", err)
	}
	return nil
}
