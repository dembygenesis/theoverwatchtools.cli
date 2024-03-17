package mysqltx

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	ErrControllerTypeNotAllowed = "controller type not allowed to serve as executor for db ops"
	ErrDatabaseNil              = "db is nil"
)

type txHandlerType int

const (
	_tx         txHandlerType = iota
	_db         txHandlerType = iota
	_controller txHandlerType = iota
)

type Handler struct {
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

// parseExecutorInstance is a function with a primary use-case of
// fetching a Handler instance from an interface, and also
// validating the instance's responsibility is to serve as a provider
// of data access objects via check its `hType`, and not as a serving controller.
func parseExecutorInstance(i interface{}) (*Handler, error) {
	mySQLTxHandler, ok := i.(*Handler)
	if !ok {
		return nil, errors.New(sysconsts.ErrNonMySQLTxInstance)
	}

	ok, err := mySQLTxHandler.isController()
	if err != nil {
		return nil, fmt.Errorf("is controller: %v", err)
	}

	if ok {
		return nil, errors.New(ErrControllerTypeNotAllowed)
	}

	return mySQLTxHandler, nil
}

// GetCtxExecutor returns the context executor.
func GetCtxExecutor(i interface{}) (boil.ContextExecutor, error) {
	txHandler, err := parseExecutorInstance(i)
	if err != nil {
		return nil, fmt.Errorf("parse executor: %v", err)
	}

	switch txHandler.hType {
	case _tx:
		return txHandler.tx, nil
	case _db:
		return txHandler.db, nil
	case _controller:
		return nil, errors.New(ErrControllerTypeNotAllowed)
	default:
		return nil, errors.New(sysconsts.ErrUnidentifiedTxType)
	}
}

func New(cfg *Config) (*Handler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	if err := cfg.Db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}

	return &Handler{
		hType:  _controller,
		db:     cfg.Db,
		logger: cfg.Logger,
	}, nil
}

type Config struct {
	Logger       *logrus.Entry `json:"logger" validate:"required"`
	Db           *sqlx.DB      `json:"db" validate:"required"`
	DatabaseName string        `json:"database_name" validate:"required"`
}

func (m *Config) Validate() error {
	return validationutils.Validate(m)
}

func (m *Handler) validate() error {
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

func (m *Handler) Commit(ctx context.Context) error {
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

func (m *Handler) Rollback(ctx context.Context) {
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
func (m *Handler) isController() (bool, error) {
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
func (m *Handler) Tx(ctx context.Context) (*Handler, error) {
	if m.db == nil {
		return nil, errors.New(sysconsts.ErrDatabaseNil)
	}

	tx, err := m.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin _tx: %v", err)
	}

	return &Handler{
		hType:  _tx,
		tx:     tx,
		logger: m.logger,
	}, nil
}

// Db returns a new instance of the tx handler attached
// with its database instance.
func (m *Handler) Db(ctx context.Context) (*Handler, error) {
	if m.db == nil {
		return nil, errors.New(sysconsts.ErrDatabaseNil)
	}

	return &Handler{
		hType:  _db,
		db:     m.db,
		logger: m.logger,
	}, nil
}

// Set accepts an active connection, and sets it right.
func (m *Handler) Set(db *sqlx.DB) error {
	if db == nil {
		return errors.New(ErrDatabaseNil)
	}
	m.db = db
	return nil
}
