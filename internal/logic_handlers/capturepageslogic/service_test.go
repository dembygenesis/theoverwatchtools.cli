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
	"github.com/volatiletech/sqlboiler/v4/boil"
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

// test post method

type argsCreateCapturePages struct {
	ctx          context.Context
	capturepages *model.CreateCapturePage
}

type testCaseCreateCapturePages struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsCreateCapturePages
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, category *model.CapturePages, err error)
}

func TestService_CreateCapturePages(t *testing.T) {
	for _, tt := range getTestCasesCreateCapturePages() {
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

			capturepages, err := svc.CreateCapturePages(tt.args.ctx, tt.args.capturepages)
			tt.assertions(t, capturepages, err)
		})
	}
}

func getTestCasesCreateCapturePages() []testCaseCreateCapturePages {
	return []testCaseCreateCapturePages{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsCreateCapturePages{
				ctx: context.TODO(),
				capturepages: &model.CreateCapturePage{
					CapturePageSetId: 1,
					Name:             "Example",
					IsControl:        1,
				},
			},
			assertions: func(t *testing.T, capturepages *model.CapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturepages, "unexpected nil capture page")

				require.NotEqual(t, 0, capturepages.Id, "unexpected nil capture page")
				require.NotEmpty(t, capturepages.Name, "unexpected empty capture page name")
				require.NotEqual(t, 0, capturepages.CapturePageSetId, "unexpected empty capture page set id")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsCreateCapturePages{
				ctx: context.TODO(),
				capturepages: &model.CreateCapturePage{
					CapturePageSetId: 1,
					Name:             "Example",
					IsControl:        1,
				},
			},
			assertions: func(t *testing.T, capturepages *model.CapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturepages, "unexpected nil capture page")

				require.NotEqual(t, 0, capturepages.Id, "unexpected nil capture page")
				require.NotEmpty(t, capturepages.Name, "unexpected empty capture page name")
				require.NotEqual(t, 0, capturepages.CapturePageSetId, "unexpected empty capture page set id")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "fail-internal-error",
			getDependencies: getConcreteDependencies,
			args: argsCreateCapturePages{
				ctx: context.TODO(),
				capturepages: &model.CreateCapturePage{
					CapturePageSetId: 1,
					Name:             "Example",
					IsControl:        1,
				},
			},
			assertions: func(t *testing.T, capturepages *model.CapturePages, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturepages, "unexpected nil capture page")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.CapturePages)
			},
		},
		{
			name:            "fail-invalid-args",
			getDependencies: getConcreteDependencies,
			args: argsCreateCapturePages{
				ctx: context.TODO(),
				capturepages: &model.CreateCapturePage{
					CapturePageSetId: 0,
					Name:             "Example",
					IsControl:        0,
				},
			},
			assertions: func(t *testing.T, capturepages *model.CapturePages, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturepages, "unexpected nil capture page")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "fail-unique",
			getDependencies: getConcreteDependencies,
			args: argsCreateCapturePages{
				ctx: context.TODO(),
				capturepages: &model.CreateCapturePage{
					CapturePageSetId: 1,
					Name:             "Example",
					IsControl:        1,
				},
			},
			assertions: func(t *testing.T, capturepages *model.CapturePages, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturepages, "unexpected nil capture page")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.CapturePage{
					Name:             "Example",
					CapturePageSetID: 1,
					IsControl:        1,
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "unexpected insert error")
			},
		},
		{
			name: "fail-mock",
			getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
				fakeCleanup := func(ignoreErrors ...bool) {}
				mockTxProvider := persistencefakes.FakeTransactionProvider{}
				mockTxProvider.TxReturns(&persistencefakes.FakeTransactionHandler{}, errors.New(mockDbReturnsErr))

				return &dependencies{
					Persistor:  &capturepageslogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    fakeCleanup,
				}, fakeCleanup
			},
			args: argsCreateCapturePages{
				ctx: context.TODO(),
				capturepages: &model.CreateCapturePage{
					CapturePageSetId: 1,
					Name:             "Example",
					IsControl:        1,
				},
			},
			assertions: func(t *testing.T, category *model.CapturePages, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, category, "unexpected nil capture page")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
	}
}
