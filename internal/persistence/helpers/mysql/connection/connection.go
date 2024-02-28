package connection

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/models"
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

	// ConnectionParameters holds the database credentials
	ConnectionParameters *mysql.ConnectionParameters `mapstructure:"parameters"`
}

// validate validates the Connection settings
func (c *Settings) validate() (*sqlx.DB, error) {

	if err := utils_common.ValidateStruct(c); err != nil {
		return nil, fmt.Errorf("required: %s", err)
	}

	db, err := mysql.GetClient(context.TODO(), &mysql.ClientOptions{
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
		return nil, models.ErrDatabaseNil
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
	Stmt       string             `mapstructure:"stmt" json:"stmt" validate:"required"`
	Pagination *models.Pagination `mapstructure:"pagination" json:"pagination"`
}

func (q *QueryIntoStructFilter) Validate() error {
	if err := utils_common.ValidateStruct(q); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	if q.Pagination == nil {
		// Assign defaults
		q.Pagination = models.NewPagination()
	} else {
		if err := q.Pagination.Validate(); err != nil {
			return fmt.Errorf("pagination: %v", err)
		}
	}

	return nil
}

type QueryAsArrFilter struct {
	Stmt       string             `mapstructure:"stmt" json:"stmt" validate:"required"`
	Pagination *models.Pagination `mapstructure:"pagination" json:"pagination"`
}

func (q *QueryAsArrFilter) Validate() error {
	if err := utils_common.ValidateStruct(q); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	if q.Pagination == nil {
		// Assign defaults
		q.Pagination = models.NewPagination()
	} else {
		if err := q.Pagination.Validate(); err != nil {
			return fmt.Errorf("pagination: %v", err)
		}
	}

	return nil
}

// QueryAsArr executes a query with respect to the timeout.
func (c *Connection) QueryAsArr(ctx context.Context, f *QueryAsArrFilter) ([][]string, *models.Pagination, error) {
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
	Stmt            string             `json:"stmt" validate:"required"`
	Pagination      *models.Pagination `json:"pagination" validate:"required"`
	Args            []interface{}      `json:"args"`
	timeoutDuration time.Duration
	db              *sqlx.DB
}

func (c *PaginateSettings) setDefaults(conn *Connection) {
	c.db = conn.db
	c.timeoutDuration = conn.settings.QueryTimeout

	if c.Pagination == nil {
		c.Pagination = models.NewPagination()
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
	err := utils_common.ValidateStruct(c)
	if err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	// We have to check internal attributes manually
	if c.db == nil {
		return models.ErrDatabaseNil
	}
	if c.timeoutDuration == 0 {
		return models.ErrTimeoutNil
	}

	return nil
}

func Paginate[T any](
	ctx context.Context,
	dest *[]T,
	settings *PaginateSettings,
) error {
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, settings.timeoutDuration)
	defer cancel()

	return paginate(ctx, settings.db, dest, settings.Stmt, settings.Pagination, settings.Args...)
}

// QueryIntoStructSettings fetches context, and pagination filters.
func (c *Connection) QueryIntoStructSettings(ctx context.Context, f *QueryAsArrFilter) ([][]string, *models.Pagination, error) {
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
