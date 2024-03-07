package categorymgr

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqlprovider"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/category/categorymysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testCtx           = context.TODO()
	testLogger        = logger.New(testCtx)
	testQueryTimeouts = &persistence.QueryTimeouts{
		Query: 10 * time.Second,
		Exec:  10 * time.Second,
	}
)

func TestCategoryManager_MySQL_GetCategories(t *testing.T) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	persistor, err := categorymysql.New(&categorymysql.Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	})
	require.NoError(t, err, "unexpected err instantiating category mysql")
	require.NotNil(t, persistor, "unexpected nil category mysql")

	txProvider, err := mysqlprovider.New(&mysqlprovider.Config{
		Db:     db,
		Logger: testLogger,
	})
	require.NoError(t, err, "unexpected err instantiating mysql provider")
	require.NotNil(t, txProvider, "unexpected nil mysql tx provider")

	cm, err := New(&Config{
		Persistor:  persistor,
		TxProvider: txProvider,
	})
	require.NoError(t, err, "unexpected err instantiating category manager")
	require.NotNil(t, cm, "unexpected nil category manager")

	filters := &model.CategoryFilters{}
	categories, err := cm.GetCategories(testCtx, filters)
	assert.NoError(t, err, "unexpected err fetching categories")
	assert.NotNil(t, categories, "unexpected nil categories")
	assert.True(t, len(categories) > 0, "unexpected categories to have a length less that 0")
}
