package mysqlstore

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltxhandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testCtx           = context.TODO()
	testLogger        = logger.New(context.TODO())
	testQueryTimeouts = &persistence.QueryTimeouts{
		Query: 10 * time.Second,
		Exec:  10 * time.Second,
	}
)

func TestNew(t *testing.T) {
	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, m, "unexpected nil")
}

func TestNew_Fail(t *testing.T) {
	cfg := &Config{
		Logger:        nil,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.Error(t, err, "unexpected nil error")
	require.Nil(t, m, "unexpected not nil")
}

func Test_ReadCategories(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, m, "unexpected nil")

	txHandler, err := mysqltxhandler.New(&mysqltxhandler.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected error creating the tx handler")

	txHandlerDb, err := txHandler.Db(testCtx)
	require.NoError(t, err, "unexpected error fetching the db from the tx handler")
	require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

	paginatedCategories, err := m.GetCategories(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the categories from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil categories")
	require.True(t, len(paginatedCategories.Categories) > 0, "unexpected empty categories")
}

func Test_UpdateCategories_Success(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, m, "unexpected nil")

	txHandler, err := mysqltxhandler.New(&mysqltxhandler.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected error creating the tx handler")

	txHandlerDb, err := txHandler.Db(testCtx)
	require.NoError(t, err, "unexpected error fetching the db from the tx handler")
	require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

	paginatedCategories, err := m.GetCategories(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the categories from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil categories")
	require.True(t, len(paginatedCategories.Categories) > 0, "unexpected empty categories")

	cat := &paginatedCategories.Categories[0]
	cat.Name = "New Category name"
	cat, err = m.UpdateCategory(testCtx, txHandlerDb, cat)
	require.NoError(t, err, "unexpected error fetching a conflicting category from the database")
	assert.Equal(t, "New Category name", cat.Name)
}

func Test_UpdateCategories_Fail(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, m, "unexpected nil")

	txHandler, err := mysqltxhandler.New(&mysqltxhandler.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected error creating the tx handler")

	txHandlerDb, err := txHandler.Db(testCtx)
	require.NoError(t, err, "unexpected error fetching the db from the tx handler")
	require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

	paginatedCategories, err := m.GetCategories(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the categories from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil categories")
	require.True(t, len(paginatedCategories.Categories) > 0, "unexpected empty categories")

	cat := &paginatedCategories.Categories[0]
	cat.Name = paginatedCategories.Categories[1].Name
	cat, err = m.UpdateCategory(testCtx, txHandlerDb, cat)
	require.Error(t, err, "unexpected error fetching a conflicting category from the database")
	assert.Contains(t, err.Error(), "Duplicate entry")
	assert.Nil(t, cat, "unexpected non nil entry")
}
