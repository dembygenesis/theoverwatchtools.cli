package mysqlhelper

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/sliceutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/jmoiron/sqlx"
)

type CreateMode int

func (o CreateMode) Validate() error {
	modes := []CreateMode{
		Recreate,
		CreateIfNotExists,
	}

	if ok := sliceutil.Compare(modes, func(v CreateMode) bool {
		return o == v
	}); !ok {
		return fmt.Errorf("invalid: %v", o)
	}

	return nil
}

const (
	Recreate CreateMode = iota
	CreateIfNotExists
)

type NewOpts struct {
	// Mode determines the creation mode
	Mode CreateMode `json:"required" validate:"required"`

	// Close determines if the database should be closed after creation
	Close bool
}

// Validate validates create options
func (o *NewOpts) Validate() error {
	if o == nil {
		return errors.New(sysconsts.ErrOptsNil)
	}

	if err := o.Mode.Validate(); err != nil {
		return fmt.Errorf("%v: %v", errors.New(sysconsts.ErrInvalidCreateMode).Error(), err)
	}

	return nil
}

// NewBuildTeardown creates a new database with modes:
//
//	"Recreate" - performs a complete teardown, and rebuild of the database
//	"CreateIfNotExists" - creates the database if it doesn't exist
func NewBuildTeardown(ctx context.Context, c *mysqlutil.ConnectionSettings, opts *NewOpts) (*sqlx.DB, error) {
	var (
		db        *sqlx.DB
		errPrefix string
		err       error
	)

	err = validationutils.Validate(c)
	if err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	switch opts.Mode {
	case Recreate:
		errPrefix = "recreate:"
		db, err = recreateDatabase(ctx, c)
	case CreateIfNotExists:
		errPrefix = "create if not exists:"
		db, err = createIfNotExists(ctx, c)
	}

	if err != nil {
		return nil, fmt.Errorf("%s %v", errPrefix, err)
	}

	if opts.Close {
		if err = db.Close(); err != nil {
			return nil, fmt.Errorf("close: %v", err)
		}
	}

	return db, err
}

func recreateDatabase(ctx context.Context, c *mysqlutil.ConnectionSettings) (*sqlx.DB, error) {
	db, err := NewDbClient(ctx, &ClientOptions{
		ConnString: c.GetConnectionString(true),
		Close:      false,
	})
	if err != nil {
		return nil, fmt.Errorf("get client: %v", err)
	}

	dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", c.Database)
	_, err = db.Exec(dropStmt)
	if err != nil {
		return nil, fmt.Errorf("drop if exists: %v", err)
	}

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`", c.Database)
	_, err = db.Exec(createStmt)
	if err != nil {
		return nil, fmt.Errorf("create database: %v", err)
	}

	db, err = NewDbClient(ctx, &ClientOptions{
		ConnString: c.GetConnectionString(false),
		Close:      false,
	})
	if err != nil {
		return nil, fmt.Errorf("reconnect client: %v", err)
	}

	return db, nil
}

func createIfNotExists(ctx context.Context, c *mysqlutil.ConnectionSettings) (*sqlx.DB, error) {
	db, err := NewDbClient(ctx, &ClientOptions{
		ConnString: c.GetConnectionString(true),
		Close:      false,
	})
	if err != nil {
		return nil, fmt.Errorf("get client: %v", err)
	}

	sqlStmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", c.Database)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("create if not exists: %v", err)
	}

	db, err = NewDbClient(ctx, &ClientOptions{
		ConnString: c.GetConnectionString(false),
		Close:      false,
	})
	if err != nil {
		return nil, fmt.Errorf("reconnect client: %v", err)
	}

	return db, nil
}
