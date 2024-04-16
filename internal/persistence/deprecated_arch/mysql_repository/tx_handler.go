package mysql_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type txHandlerType int

const (
	// _tx indicates this transaction handler is for a single-use transaction.
	_tx txHandlerType = iota
	// _db indicates this is a live connection
	_db txHandlerType = iota
	// _controller indicates this instance's responsibility is to create the previous
	// aforementioned types.
	_controller txHandlerType = iota
)

type txHandler struct {
	hType  txHandlerType
	db     *sqlx.DB
	tx     *sqlx.Tx
	logger *logrus.Entry
}

// getExec returns the context executor.
func (m *txHandler) getExec() boil.ContextExecutor {
	switch m.hType {
	case _tx:
		return m.tx
	case _db:
		return m.db
	default:
		panic("invalid context executor")
	}
}

func (m *txHandler) validate() error {
	if m == nil {
		return errors.New(sysconsts.ErrStructNil)
	}

	if ok := sliceutil.Compare([]txHandlerType{_db, _tx}, func(v txHandlerType) bool {
		return v == m.hType
	}); !ok {
		return errors.New(sysconsts.ErrUnidentifiedTxType)
	}

	return nil
}

func (m *txHandler) Commit(ctx context.Context) error {
	var (
		err error
	)

	// We do nothing
	if m.hType == _db {
		return errors.New(sysconsts.ErrCommitInvalidDB)
	}

	if err = m.tx.Commit(); err != nil {
		return fmt.Errorf("commit: %v", err)
	}

	return nil
}

func (m *txHandler) Rollback(ctx context.Context) {
	var (
		err error
	)

	// We do nothing
	if m.hType == _db {
		return
	}

	if err = m.tx.Rollback(); err != nil {
		m.logger.Warn(logrus.Fields{
			"err": fmt.Errorf("rollback: %v", err),
		})
	}
}

// Tx returns a new instance of the tx handler attached
// with a transaction.
func (m *txHandler) Tx() (Transaction, error) {
	if m.db == nil {
		return nil, errors.New(sysconsts.ErrDatabaseNil)
	}

	tx, err := m.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin _tx: %v", err)
	}

	return &txHandler{
		hType:  _tx,
		tx:     tx,
		logger: m.logger,
	}, nil
}

// Db returns a new instance of the tx handler attached
// with its database instance.
func (m *txHandler) Db() (Transaction, error) {
	if m.db == nil {
		return nil, errors.New(sysconsts.ErrDatabaseNil)
	}

	return &txHandler{
		hType:  _db,
		db:     m.db,
		logger: m.logger,
	}, nil
}
