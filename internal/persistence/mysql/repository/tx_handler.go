package repository

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/models"
	"github.com/dembygenesis/local.tools/internal/utility"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type txHandlerType int

const (
	_tx txHandlerType = iota
	_db txHandlerType = iota
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
		return models.ErrStructNil
	}

	if ok := utility.WalkSlice([]txHandlerType{_db, _tx}, func(v txHandlerType) bool {
		return v == m.hType
	}); !ok {
		return models.ErrUnidentifiedTxType
	}

	return nil
}

func (m *txHandler) Commit(ctx context.Context) error {
	var (
		err error
	)

	// We do nothing
	if m.hType == _db {
		return models.ErrCommitInvalidDB
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
