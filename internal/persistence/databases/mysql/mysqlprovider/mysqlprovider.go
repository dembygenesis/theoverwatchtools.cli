package mysqlprovider

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqltxhandler"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Db     *sqlx.DB      `json:"db" validate:"required"`
	Logger *logrus.Entry `json:"logger" validate:"required"`
}

type MySQLProvider struct {
	cfg        *Config
	controller *mysqltxhandler.MySQLTxHandler
}

func (m *MySQLProvider) Db(ctx context.Context) (persistence.TransactionHandler, error) {
	txHandler, err := m.controller.Db(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx handler (db): %v", err)
	}
	return txHandler, nil
}

func (m *MySQLProvider) Tx(ctx context.Context) (persistence.TransactionHandler, error) {
	txHandler, err := m.controller.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx handler (tx %v", err)
	}
	return txHandler, nil
}

func New(cfg *Config) (persistence.TransactionProvider, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	txHandler, err := mysqltxhandler.New(&mysqltxhandler.Config{
		Db:     cfg.Db,
		Logger: cfg.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("tx handler: %v", err)
	}

	prov := &MySQLProvider{
		cfg:        cfg,
		controller: txHandler,
	}

	return prov, nil
}

func (m *Config) Validate() error {
	if err := validationutils.Validate(m); err != nil {
		return err
	}
	return nil
}
