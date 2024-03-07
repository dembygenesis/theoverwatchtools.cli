package mysqltxhandler

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	ErrControllerTypeNotAllowed = errors.New("controller type not allowed to serve as executor for db ops")
)

type txHandlerType int

const (
	_tx         txHandlerType = iota
	_db         txHandlerType = iota
	_controller txHandlerType = iota
)

type MySQLTxHandler struct {
	// hType is the type of responsibility the instance will have:
	//
	// _db        : the instance serves a container for a "db" val that will be utilised in CRUD operations
	// _tx        : the instance serves a container for a "tx" val that will be utilised in CRUD operations
	// _controller: the instance serves a controller living in another domain that spawns instances of _db, and _tx of similar type
	// for returning self-similar instances that will serve as 1 & 2
	hType  txHandlerType
	db     *sqlx.DB
	tx     *sqlx.Tx
	logger *logrus.Entry
}

// ParseExecutorInstance is a function with a primary use-case of
// fetching a MySQLTxHandler instance from an interface, and also
// validating the instance's responsibility is to serve as a provider
// of data access objects via check its `hType`, and not as a serving controller.
func ParseExecutorInstance(i interface{}) (*MySQLTxHandler, error) {
	mySQLTxHandler, ok := i.(*MySQLTxHandler)
	if !ok {
		return nil, model.ErrNonMySQLTxInstance
	}

	ok, err := mySQLTxHandler.isController()
	if err != nil {
		return nil, fmt.Errorf("is controller: %v", err)
	}

	if ok {
		return nil, ErrControllerTypeNotAllowed
	}

	return mySQLTxHandler, nil
}

func New(cfg *Config) (*MySQLTxHandler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	if err := cfg.Db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}

	return &MySQLTxHandler{
		hType:  _controller,
		db:     cfg.Db,
		logger: cfg.Logger,
	}, nil
}

type Config struct {
	Logger *logrus.Entry `json:"logger" validate:"required"`
	Db     *sqlx.DB      `json:"db" validate:"required"`
}

func (m *Config) Validate() error {
	return validationutils.Validate(m)
}

// GetCtxExecutor returns the context executor.
func (m *MySQLTxHandler) GetCtxExecutor() (boil.ContextExecutor, error) {
	switch m.hType {
	case _tx:
		return m.tx, nil
	case _db:
		return m.db, nil
	case _controller:
		return nil, ErrControllerTypeNotAllowed
	default:
		return nil, model.ErrUnidentifiedTxType
	}
}

func (m *MySQLTxHandler) validate() error {
	if m == nil {
		return model.ErrStructNil
	}

	if ok := sliceutil.Compare([]txHandlerType{_db, _tx}, func(v txHandlerType) bool {
		return v == m.hType
	}); !ok {
		return model.ErrUnidentifiedTxType
	}

	return nil
}

func (m *MySQLTxHandler) Commit(ctx context.Context) error {
	var (
		err error
	)

	// We do nothing
	if m.hType == _db {
		return model.ErrCommitInvalidDB
	}

	if err = m.tx.Commit(); err != nil {
		return fmt.Errorf("commit: %v", err)
	}

	return nil
}

func (m *MySQLTxHandler) Rollback(ctx context.Context) {
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

// isController checks if the type is "controller".
func (m *MySQLTxHandler) isController() (bool, error) {
	switch m.hType {
	case _controller:
		return true, nil
	case _tx:
		return false, nil
	case _db:
		return false, nil
	default:
		return false, fmt.Errorf("unidentified hType: %v", m.hType)
	}
}

// Tx returns a new instance of the tx handler attached
// with a transaction.
func (m *MySQLTxHandler) Tx(ctx context.Context) (*MySQLTxHandler, error) {
	if m.db == nil {
		return nil, model.ErrDatabaseNil
	}

	tx, err := m.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin _tx: %v", err)
	}

	return &MySQLTxHandler{
		hType:  _tx,
		tx:     tx,
		logger: m.logger,
	}, nil
}

// Db returns a new instance of the tx handler attached
// with its database instance.
func (m *MySQLTxHandler) Db(ctx context.Context) (*MySQLTxHandler, error) {
	if m.db == nil {
		return nil, model.ErrDatabaseNil
	}

	return &MySQLTxHandler{
		hType:  _db,
		db:     m.db,
		logger: m.logger,
	}, nil
}
