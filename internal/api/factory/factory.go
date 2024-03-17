package factory

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"testing"
)

var (
	testCtx    = context.TODO()
	testLogger = logger.New(testCtx)
)

func TestMain(m *testing.M) {
	fmt.Println("=== build")
	m.Run()
	fmt.Println("=== teardown")
}

/*func testGetMySQLProvider(t *testing.T) (persistence.TransactionProvider, func(ignoreErrors ...bool)) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)

	txProvider, err := mysqltxprovider.New(&mysqltxprovider.Config{
		Db:     db,
		Logger: testLogger,
	})
	require.NoError(t, err, "error instantiating mysql tx provider")

	return txProvider, cleanup
}

func TestGetMySQLCategoryManager(t *testing.T) (categorysrv.LogicHandler, func(ignoreErrors ...bool)) {
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

	mgr, err := categorysrv.New(&categorysrv.Config{
		Persistor:  persistor,
		TxProvider: txProvider,
	})
	require.NoError(t, err, "error instantiating mysql category manager")

	return mgr, cleanup
}*/
