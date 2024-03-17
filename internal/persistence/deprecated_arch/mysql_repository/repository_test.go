package mysql_repository

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew_Success(t *testing.T) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger: logger.New(context.Background()),
		Db:     db,
	}

	_, err := New(cfg)
	require.NoError(t, err)
}

func TestNew_Fail_Validate(t *testing.T) {
	cfg := &Config{
		Logger: logger.New(context.Background()),
	}

	_, err := New(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "validate:", "expected err to contains 'validate:'")
	require.Contains(t, err.Error(), "validate struct:", "expected err to contains 'validate struct:'")
}

func TestNew_Fail_Ping(t *testing.T) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup(true)

	cfg := &Config{
		Logger: logger.New(context.Background()),
		Db:     db,
	}

	err := db.Close()
	require.NoError(t, err, "expected db to close successfully")

	err = db.Ping()
	require.Error(t, err, "expected to be unable to ping a db")

	_, err = New(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "validate:", "expected err to contains 'validate:'")
	require.Contains(t, err.Error(), "ping:", "expected err to contains 'ping:'")
}
