package capturepageslogic

import (
	"context"
	"errors"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/capturepageslogic/capturepageslogicfakes"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistencefakes"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore/testhelper"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	mockTimeout      = 5 * time.Second
	mockLogger       = logger.New(context.TODO())
	mockDbReturnsErr = "error getting db"
)

type dependencies struct {
	Persistor  persistor
	Logger     *logrus.Entry
	TxProvider persistence.TransactionProvider
	Db         *sqlx.DB
	Cleanup    func(ignoreError ...bool)
}

type argsGetCapturePages struct {
	ctx    context.Context
	filter *model.CapturePagesFilters
}

type testCaseGetCapturePages struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsGetCapturePages
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, capturePages *model.PaginatedCapturePages, err error)
}

func getConcreteDependencies(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)

	store, err := mysqlstore.New(&mysqlstore.Config{
		Logger: mockLogger,
		QueryTimeouts: &persistence.QueryTimeouts{
			Query: mockTimeout,
			Exec:  mockTimeout,
		},
	})
	require.NoError(t, err, "unexpected new mysqlstore error")

	tx, err := mysqltx.New(&mysqltx.Config{
		Logger:       mockLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected new mysqltx error")

	prov, err := mysqlconn.New(&mysqlconn.Config{
		Logger:    mockLogger,
		TxHandler: tx,
	})
	require.NoError(t, err, "unexpected new mysqlconn error")

	return &dependencies{
		Persistor:  store,
		TxProvider: prov,
		Logger:     mockLogger,
		Cleanup:    cleanup,
		Db:         db,
	}, cleanup
}

// getGetCapturePagesTestCases returns a list of test cases for the ListCapturePages method.
func getGetCapturePagesTestCases() []testCaseGetCapturePages {
	return []testCaseGetCapturePages{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsGetCapturePages{
				ctx:    context.Background(),
				filter: &model.CapturePagesFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected get capture pages error")
				require.NotNil(t, paginated, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
				assert.NotEqual(t, 0, paginated.Pagination.RowCount, "unexpected total count")
			},
		},
		{
			name:            "fail-get-capture-pages",
			getDependencies: getConcreteDependencies,
			args: argsGetCapturePages{
				ctx:    context.Background(),
				filter: &model.CapturePagesFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.CapturePages)
			},
			assertions: func(t *testing.T, paginated *model.PaginatedCapturePages, err error) {
				require.Error(t, err, "unexpected get capture pages error")
				require.Contains(t, err.Error(), "get capture pages:")
			},
		},
		{
			name: "fail-mock-get-db",
			getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
				cleanup := func(ignoreErrors ...bool) {

				}

				mockTxProvider := persistencefakes.FakeTransactionProvider{}
				mockTxProvider.DbReturns(nil, errors.New(mockDbReturnsErr))

				return &dependencies{
					Persistor:  &capturepageslogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			args: argsGetCapturePages{
				ctx:    context.Background(),
				filter: &model.CapturePagesFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedCapturePages, err error) {
				require.Error(t, err, "unexpected get capture pages error")
				require.Contains(t, err.Error(), "get db:")
			},
		},
	}
}

func TestService_GetCapturePages(t *testing.T) {
	for _, tt := range getGetCapturePagesTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			_dependencies, cleanup := tt.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			tt.mutations(t, _dependencies.Db)

			paginatedCapturePages, err := svc.ListCapturePages(tt.args.ctx, tt.args.filter)
			tt.assertions(t, paginatedCapturePages, err)
		})
	}
}
