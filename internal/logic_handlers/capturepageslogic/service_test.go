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
	"github.com/volatiletech/null/v8"
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
	assertions      func(t *testing.T, capturepages *model.CapturePages, err error)
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
				},
			},
			assertions: func(t *testing.T, capturepages *model.CapturePages, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturepages, "unexpected nil capture page")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
	}
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

// test update method

type argsUpdateCapturePages struct {
	ctx    context.Context
	params *model.UpdateCapturePages
}

type testCaseUpdateCapturePages struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsUpdateCapturePages
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, params *model.UpdateCapturePages, capturepages *model.CapturePages, err error)
}

func getTestCasesUpdateCapturePages() []testCaseUpdateCapturePages {
	return []testCaseUpdateCapturePages{
		{
			name: "success",
			args: argsUpdateCapturePages{
				ctx: context.TODO(),
				params: &model.UpdateCapturePages{
					Id: 1,
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					CapturePageSetId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateCapturePages, capturepages *model.CapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturepages, "unexpected nil capture pages")
				assert.Equal(t, params.Id, capturepages.Id, "expected id to be equal")
				assert.Equal(t, params.CapturePageSetId.Int, capturepages.CapturePageSetId, "expected capture pages type ref id to be equal")
				assert.Equal(t, params.Name.String, capturepages.Name, "expected name to be equal")
			},
		},
		{
			name: "success-name-only",
			args: argsUpdateCapturePages{
				ctx: context.TODO(),
				params: &model.UpdateCapturePages{
					Id: 1,
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateCapturePages, capturepages *model.CapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturepages, "unexpected nil capture pages")
				assert.Equal(t, params.Name.String, capturepages.Name, "expected name to be equal")
			},
		},
		{
			name: "success-capture-pages-type-ref-id-only",
			args: argsUpdateCapturePages{
				ctx: context.TODO(),
				params: &model.UpdateCapturePages{
					Id: 1,
					CapturePageSetId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateCapturePages, capturepages *model.CapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturepages, "unexpected nil capture pages")
				assert.Equal(t, params.Id, capturepages.Id, "expected id to be equal")
				assert.Equal(t, params.CapturePageSetId.Int, capturepages.CapturePageSetId, "expected name to be equal")
			},
		},
		{
			name: "fail-mock",
			args: argsUpdateCapturePages{
				ctx: context.TODO(),
				params: &model.UpdateCapturePages{
					Id: 1,
					CapturePageSetId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
				cleanup := func(ignoreErrors ...bool) {}

				mockTxProvider := persistencefakes.FakeTransactionProvider{}
				mockTxProvider.TxReturns(nil, errors.New(mockDbReturnsErr))

				return &dependencies{
					Persistor:  &capturepageslogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateCapturePages, capturepages *model.CapturePages, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturepages, "unexpected nil capture pages")
			},
		},
		{
			name: "fail-internal-server-error",
			args: argsUpdateCapturePages{
				ctx: context.TODO(),
				params: &model.UpdateCapturePages{
					Id: 1,
					CapturePageSetId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.CapturePages)
			},
			assertions: func(t *testing.T, params *model.UpdateCapturePages, capturepages *model.CapturePages, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturepages, "unexpected nil capture pages")
			},
		},
	}
}

func TestService_UpdateCapturePages(t *testing.T) {
	for _, tt := range getTestCasesUpdateCapturePages() {
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

			capturepages, err := svc.UpdateCapturePages(tt.args.ctx, tt.args.params)
			tt.assertions(t, tt.args.params, capturepages, err)
		})
	}
}

// test delete method

type argsDeleteCapturePages struct {
	ctx    context.Context
	params *model.DeleteCapturePages
}

type testCaseDeleteCapturePages struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsDeleteCapturePages
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, params *model.DeleteCapturePages, err error)
}

func getTestCasesDeleteCapturePages() []testCaseDeleteCapturePages {
	return []testCaseDeleteCapturePages{
		{
			name: "success",
			args: argsDeleteCapturePages{
				ctx: context.TODO(),
				params: &model.DeleteCapturePages{
					ID: 0,
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.DeleteCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				assert.Equal(t, params.ID, 0)
			},
		},
	}
}

func TestService_DeleteCapturePages(t *testing.T) {
	for _, tt := range getTestCasesDeleteCapturePages() {
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

			err = svc.DeleteCapturePages(tt.args.ctx, tt.args.params)
			tt.assertions(t, tt.args.params, err)
		})
	}
}

// test restore method

type argsRestoreCapturePages struct {
	ctx    context.Context
	params *model.RestoreCapturePages
}

type testCaseRestoreCapturePages struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsRestoreCapturePages
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, params *model.RestoreCapturePages, err error)
}

func getTestCasesRestoreCapturePages() []testCaseRestoreCapturePages {
	return []testCaseRestoreCapturePages{
		{
			name: "success",
			args: argsRestoreCapturePages{
				ctx: context.TODO(),
				params: &model.RestoreCapturePages{
					ID: 1,
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.RestoreCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				assert.Equal(t, params.ID, 1)
			},
		},
	}
}

func TestService_RestoreCapturePages(t *testing.T) {
	for _, tt := range getTestCasesRestoreCapturePages() {
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

			err = svc.RestoreCapturePages(tt.args.ctx, tt.args.params)
			tt.assertions(t, tt.args.params, err)
		})
	}
}
