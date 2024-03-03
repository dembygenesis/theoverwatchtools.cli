package mysqlprovider

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlhelper"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testCtx    = context.TODO()
	testLogger = logger.New(testCtx)
)

func TestNew(t *testing.T) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	provider, err := New(&Config{
		Db:     db,
		Logger: testLogger,
	})

	require.NoError(t, err, "unexpected new error")
	require.NotNil(t, provider, "unexpected nil provider")
}

func TestNew_Fail(t *testing.T) {
	provider, err := New(&Config{})

	require.Error(t, err, "unexpected new error")
	require.Nil(t, provider, "unexpected nil provider")
}
