package mysql_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type Repository interface {
	Categories
	TransactionHandler
}

type Categories interface {
	GetCategories(ctx context.Context, t Transaction, f *model.CategoryFilters) (model.Categories, error)
	UpsertCategories(ctx context.Context, t Transaction, categories []model.Category) (model.Categories, error)
}

// TransactionHandler returns transaction handlers.
type TransactionHandler interface {
	Tx() (Transaction, error)
	Db() (Transaction, error)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}

type repository struct {
	cfg *Config
}

func (m *repository) GetCategories(ctx context.Context, t Transaction, f *model.CategoryFilters) (model.Categories, error) {
	//TODO implement me
	panic("implement me")
}

func (m *repository) UpsertCategories(ctx context.Context, t Transaction, categories []model.Category) (model.Categories, error) {
	//TODO implement me
	panic("implement me")
}

// Tx starts an abstract transaction for external resources.
func (m *repository) Tx() (Transaction, error) {
	tx, err := m.cfg.Db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin _tx: %v", err)
	}

	return &txHandler{
		hType:  _tx,
		tx:     tx,
		logger: m.cfg.Logger,
	}, nil
}

func (m *repository) Db() (Transaction, error) {
	if m.cfg.Db == nil {
		return nil, errors.New(sysconsts.ErrDatabaseNil)
	}

	return &txHandler{
		hType:  _db,
		db:     m.cfg.Db,
		logger: m.cfg.Logger,
	}, nil
}

type Config struct {
	// Logger is the logging mechanism
	Logger *logrus.Entry `json:"logger" validate:"required"`

	// Db is the connection instance
	Db *sqlx.DB `json:"_db" validate:"required"`

	// ConnectTimeout is the connect timeout duration
	ConnectTimeout time.Duration `mapstructure:"connect_timeout" validate:"required"`

	// QueryTimeout is the query (READ) timeout duration
	QueryTimeout time.Duration `mapstructure:"query_timeout" validate:"required"`

	// ExecTimeout is the query execution timeout duration
	ExecTimeout time.Duration `mapstructure:"exec_timeout" validate:"required"`
}

func (m *Config) Validate() error {
	if m != nil {
		m.setDefaultTimeouts()
	}

	err := validationutils.Validate(m)
	if err != nil {
		return fmt.Errorf("validate struct: %v", err)
	}

	if err = m.Db.Ping(); err != nil {
		return fmt.Errorf("ping: %v", err)
	}

	return nil
}

func (m *Config) setDefaultTimeouts() {
	if m.ConnectTimeout == 0 {
		m.ConnectTimeout = 10 * time.Second
	}

	if m.ExecTimeout == 0 {
		m.ExecTimeout = 10 * time.Second
	}

	if m.QueryTimeout == 0 {
		m.QueryTimeout = 10 * time.Second
	}
}

func New(cfg *Config) (Repository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}
	return &repository{cfg}, nil
}
