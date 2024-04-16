package mysqlhelper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
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

	// ConnectionParameters holds the database credentials
	ConnectionParameters *mysqlutil.ConnectionSettings `mapstructure:"parameters"`
}

func (c *Settings) setDefaults() {
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 10 * time.Second
	}

	if c.QueryTimeout == 0 {
		c.QueryTimeout = 10 * time.Second
	}

	if c.ExecTimeout == 0 {
		c.ExecTimeout = 10 * time.Second
	}
}

// validate validates the Connection settings
func (c *Settings) validate() (*sqlx.DB, error) {
	if c != nil {
		c.setDefaults()
	}

	if err := validationutils.Validate(c); err != nil {
		return nil, fmt.Errorf("required: %s", err)
	}

	if err := validationutils.Validate(c.ConnectionParameters); err != nil {
		return nil, fmt.Errorf("required connection parameters: %s", err)
	}

	db, err := NewDbClient(context.TODO(), &ClientOptions{
		PingTimeout: c.ConnectTimeout,
		ConnString:  c.ConnectionParameters.GetConnectionString(false),
		Close:       false,
	})
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	return db, nil
}

// Connection containers the helper method
type Connection struct {
	db       *sqlx.DB
	settings *Settings
}

func NewConnection(settings *Settings) (*Connection, error) {
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
		return nil, errors.New(sysconsts.ErrDatabaseNil)
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

type QueryIntoStructFilter struct {
	Dest       any
	Args       []interface{}
	Stmt       string            `mapstructure:"stmt" json:"stmt" validate:"required"`
	Pagination *model.Pagination `mapstructure:"pagination" json:"pagination"`
}

func (q *QueryIntoStructFilter) Validate() error {
	if err := validationutils.Validate(q); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	if q.Pagination == nil {
		// Assign defaults
		q.Pagination = model.NewPagination()
	} else {
		if err := q.Pagination.ValidatePagination(); err != nil {
			return fmt.Errorf("pagination: %v", err)
		}
	}

	return nil
}

type QueryAsArrFilter struct {
	Stmt       string            `mapstructure:"stmt" json:"stmt" validate:"required"`
	Pagination *model.Pagination `mapstructure:"pagination" json:"pagination"`
}

func (q *QueryAsArrFilter) Validate() error {
	if err := validationutils.Validate(q); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	if q.Pagination == nil {
		// Assign defaults
		q.Pagination = model.NewPagination()
	} else {
		if err := q.Pagination.ValidatePagination(); err != nil {
			return fmt.Errorf("pagination: %v", err)
		}
	}

	return nil
}

// QueryAsArr executes a query with respect to the timeout.
func (c *Connection) QueryAsArr(ctx context.Context, f *QueryAsArrFilter) ([][]string, *model.Pagination, error) {
	if err := f.Validate(); err != nil {
		return nil, nil, fmt.Errorf("filter: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, c.settings.QueryTimeout)
	defer cancel()

	res, err := queryAsStringArrays(ctx, c.db, f.Stmt, f.Pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("query as arr: %v", err)
	}

	return res, f.Pagination, nil
}

type PaginateSettings struct {
	Stmt            string            `json:"stmt" validate:"required"`
	Pagination      *model.Pagination `json:"pagination" validate:"required"`
	Args            []interface{}     `json:"args"`
	timeoutDuration time.Duration
	db              *sqlx.DB
}

func (c *PaginateSettings) setDefaults(conn *Connection) {
	c.db = conn.db
	c.timeoutDuration = conn.settings.QueryTimeout

	if c.Pagination == nil {
		c.Pagination = model.NewPagination()
	}
}

// GetPaginateSettings wraps on PaginateSettings, and applies defaults.
func (c *Connection) GetPaginateSettings(p *PaginateSettings) *PaginateSettings {
	if p == nil {
		return nil
	}
	p.setDefaults(c)
	return p
}

func (c *PaginateSettings) Validate() error {
	err := validationutils.Validate(c)
	if err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	// We have to check internal attributes manually
	if c.db == nil {
		return errors.New(sysconsts.ErrDatabaseNil)
	}
	if c.timeoutDuration == 0 {
		return errors.New(sysconsts.ErrTimeoutNil)
	}

	return nil
}

// QueryIntoStructSettings fetches context, and pagination filters.
func (c *Connection) QueryIntoStructSettings(ctx context.Context, f *QueryAsArrFilter) ([][]string, *model.Pagination, error) {
	if err := f.Validate(); err != nil {
		return nil, nil, fmt.Errorf("filter: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, c.settings.QueryTimeout)
	defer cancel()

	res, err := queryAsStringArrays(ctx, c.db, f.Stmt, f.Pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("query as arr: %v", err)
	}

	return res, f.Pagination, nil
}
