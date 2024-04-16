package mysqlconn

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testCtx    = context.TODO()
	testLogger = logger.New(testCtx)
)

func TestNew(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	txHandler, err := mysqltx.New(&mysqltx.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected new error")

	provider, err := New(&Config{
		Logger:    testLogger,
		TxHandler: txHandler,
	})

	require.NoError(t, err, "unexpected new error")
	require.NotNil(t, provider, "unexpected nil provider")
}

func TestNew_Fail(t *testing.T) {
	provider, err := New(&Config{})

	require.Error(t, err, "unexpected new error")
	require.Nil(t, provider, "unexpected nil provider")
}
