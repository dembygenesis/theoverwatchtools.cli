package mysql_repository

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlhelper"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testContext = context.TODO()
	testLogger  = logger.New(testContext)
)

func Test_mySQLRepository_TxHandler_DB(t *testing.T) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger: testLogger,
		Db:     db,
	}

	r, err := New(cfg)
	require.NoError(t, err, "unexpected error initialising repository")

	tDB, err := r.Db()
	require.NoError(t, err, "unexpected error fetching database")

	f := &model.CategoryFilters{}

	categories, err := r.GetCategories(testContext, tDB, f)
	require.NoError(t, err, "unexpected error fetching categories")
	require.NotEmpty(t, categories, "unexpected empty categories")
}

func Test_mySQLRepository_TxHandler_Tx(t *testing.T) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger: testLogger,
		Db:     db,
	}

	r, err := New(cfg)
	require.NoError(t, err, "unexpected error initialising repository")

	tTx, err := r.Tx()
	defer tTx.Rollback(testContext)
	require.NoError(t, err, "unexpected error initialising transaction")

	tDb, err := r.Db()
	require.NoError(t, err, "unexpected error initialising db")

	f := &model.CategoryFilters{}

	categories, err := r.GetCategories(testContext, tTx, f)
	require.NoError(t, err, "unexpected error fetching existing categories")
	require.NotEmpty(t, categories, "unexpected empty existing categories")

	toUpsert := []model.Category{
		{
			Name:              "New entry 1",
			CategoryTypeRefId: 1,
		},
		{
			Name:              "New entry 2",
			CategoryTypeRefId: 1,
		},
	}

	_, err = r.UpsertCategories(testContext, tTx, toUpsert)
	require.NoError(t, err, "unexpected error upserting categories")

	categoriesPreCommit, err := r.GetCategories(testContext, tDb, f)
	require.NoError(t, err, "unexpected error fetching categories pre-commit categories")
	require.NotEmpty(t, categoriesPreCommit, "unexpected empty categories pre-commit categories")
	require.Equal(t, len(categories), len(categoriesPreCommit), "unexpected unequal categories, and pre-commit new categories")

	err = tTx.Commit(testContext)
	require.NoError(t, err, "unexpected error during commit")

	categoriesPostCommit, err := r.GetCategories(testContext, tDb, f)

	require.NoError(t, err, "unexpected error fetching categories post-commit categories")
	require.NotEmpty(t, categoriesPostCommit, "unexpected empty categories post-commit categories")

	require.True(t, len(categoriesPostCommit) > len(categoriesPreCommit))

	expectedEntryCountDiff := len(categories) + len(toUpsert)
	actualEntryCountDiff := len(categoriesPostCommit)
	require.Equal(t, expectedEntryCountDiff, actualEntryCountDiff, "unexpected diffs are unequal")
}
