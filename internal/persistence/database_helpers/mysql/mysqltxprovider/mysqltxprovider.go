package mysqltxprovider

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltxhandler"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger    *logrus.Entry                  `json:"logger" validate:"required"`
	TxHandler *mysqltxhandler.MySQLTxHandler `json:"tx_handler" validate:"required"`
}

type Impl struct {
	cfg        *Config
	controller *mysqltxhandler.MySQLTxHandler
}

func (m *Impl) Db(ctx context.Context) (persistence.TransactionHandler, error) {
	txHandler, err := m.controller.Db(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx handler (db): %v", err)
	}
	return txHandler, nil
}

func (m *Impl) Tx(ctx context.Context) (persistence.TransactionHandler, error) {
	txHandler, err := m.controller.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx handler (tx %v", err)
	}
	return txHandler, nil
}

func New(cfg *Config) (*Impl, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	prov := &Impl{
		cfg:        cfg,
		controller: cfg.TxHandler,
	}

	return prov, nil
}

func (m *Config) Validate() error {
	if err := validationutils.Validate(m); err != nil {
		return err
	}
	return nil
}
