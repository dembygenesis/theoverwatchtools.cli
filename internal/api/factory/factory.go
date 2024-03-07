package factory

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/managers/categorymgr"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlprovider"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/category/categorymysql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testCtx    = context.TODO()
	testLogger = logger.New(testCtx)
)

func testGetMySQLProvider(t *testing.T) (persistence.TransactionProvider, func(ignoreErrors ...bool)) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)

	txProvider, err := mysqlprovider.New(&mysqlprovider.Config{
		Db:     db,
		Logger: testLogger,
	})
	require.NoError(t, err, "error instantiating mysql tx provider")

	return txProvider, cleanup
}

func TestGetMySQLCategoryManager(t *testing.T) (*categorymgr.CategoryManager, func(ignoreErrors ...bool)) {
	t.Helper()

	txProvider, cleanup := testGetMySQLProvider(t)

	persistor, err := categorymysql.New(&categorymysql.Config{
		Logger: logger.New(context.TODO()),
		QueryTimeouts: &persistence.QueryTimeouts{
			Query: 10 * time.Second,
			Exec:  10 * time.Second,
		},
	})
	require.NoError(t, err, "error instantiating mysql category provider")

	mgr, err := categorymgr.New(&categorymgr.Config{
		Persistor:  persistor,
		TxProvider: txProvider,
	})
	require.NoError(t, err, "error instantiating mysql category manager")

	return mgr, cleanup
}
