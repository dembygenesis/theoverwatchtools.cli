package clicktrackerlogic

import (
	"context"
	"errors"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/clicktrackerlogic/clicktrackerlogicfakes"
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
	Cleanup    func(ignoreErrors ...bool)
}

type argsGetClickTrackers struct {
	ctx    context.Context
	filter *model.ClickTrackerFilters
}

type testCaseGetClickTrackers struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsGetClickTrackers
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, categories *model.PaginatedClickTrackers, err error)
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

// getGetClickTrackersTestCases returns a list of test cases for the ListClickTrackers method.
func getGetClickTrackersTestCases() []testCaseGetClickTrackers {
	return []testCaseGetClickTrackers{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsGetClickTrackers{
				ctx:    context.Background(),
				filter: &model.ClickTrackerFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedClickTrackers, err error) {
				require.NoError(t, err, "unexpected get click trackers error")
				require.NotNil(t, paginated, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				assert.NotEqual(t, 0, paginated.Pagination.RowCount, "unexpected total count")
			},
		},
		{
			name:            "fail-get-categories",
			getDependencies: getConcreteDependencies,
			args: argsGetClickTrackers{
				ctx:    context.Background(),
				filter: &model.ClickTrackerFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.ClickTrackers)
			},
			assertions: func(t *testing.T, paginated *model.PaginatedClickTrackers, err error) {
				require.Error(t, err, "unexpected get click trackers error")
				require.Contains(t, err.Error(), "get click trackers:")
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
					Persistor:  &clicktrackerlogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			args: argsGetClickTrackers{
				ctx:    context.Background(),
				filter: &model.ClickTrackerFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedClickTrackers, err error) {
				require.Error(t, err, "unexpected get click trackers error")
				require.Contains(t, err.Error(), "get db:")
			},
		},
	}
}

func TestService_GetClickTrackers(t *testing.T) {
	for _, tt := range getGetClickTrackersTestCases() {
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

			paginatedClickTrackers, err := svc.ListClickTrackers(tt.args.ctx, tt.args.filter)
			tt.assertions(t, paginatedClickTrackers, err)
		})
	}
}
