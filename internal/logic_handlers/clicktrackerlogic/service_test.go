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
	"github.com/volatiletech/null/v8"
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

// test update method

type argsUpdateClickTrackers struct {
	ctx    context.Context
	params *model.UpdateClickTracker
}

type testCaseUpdateClickTrackers struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsUpdateClickTrackers
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, params *model.UpdateClickTracker, clicktracker *model.ClickTracker, err error)
}

func getTestCasesUpdateClickTracker() []testCaseUpdateClickTrackers {
	return []testCaseUpdateClickTrackers{
		//{
		//	name: "success",
		//	args: argsUpdateClickTrackers{
		//		ctx: context.TODO(),
		//		params: &model.UpdateClickTracker{
		//			Id: 1,
		//			Name: null.String{
		//				String: "Demby",
		//				Valid:  true,
		//			},
		//			ClickTrackerSetId: null.Int{
		//				Int:   3,
		//				Valid: true,
		//			},
		//		},
		//	},
		//	getDependencies: getConcreteDependencies,
		//	mutations: func(t *testing.T, db *sqlx.DB) {
		//
		//	},
		//	assertions: func(t *testing.T, params *model.UpdateClickTracker, clicktracker *model.ClickTracker, err error) {
		//		require.NoError(t, err, "unexpected error")
		//		require.NotNil(t, clicktracker, "unexpected nil click trackers")
		//		assert.Equal(t, params.Id, clicktracker.Id, "expected id to be equal")
		//		assert.Equal(t, params.ClickTrackerSetId.Int, clicktracker.ClickTrackerSetId, "expected click tracker type ref id to be equal")
		//		assert.Equal(t, params.Name.String, clicktracker.Name, "expected name to be equal")
		//	},
		//},
		//{
		//	name: "success-name-only",
		//	args: argsUpdateClickTrackers{
		//		ctx: context.TODO(),
		//		params: &model.UpdateClickTracker{
		//			Id: 1,
		//			Name: null.String{
		//				String: "Demby",
		//				Valid:  true,
		//			},
		//		},
		//	},
		//	getDependencies: getConcreteDependencies,
		//	mutations: func(t *testing.T, db *sqlx.DB) {
		//
		//	},
		//	assertions: func(t *testing.T, params *model.UpdateClickTracker, clicktracker *model.ClickTracker, err error) {
		//		require.NoError(t, err, "unexpected error")
		//		require.NotNil(t, clicktracker, "unexpected nil click trackers")
		//		assert.Equal(t, params.Name.String, clicktracker.Name, "expected name to be equal")
		//	},
		//},
		{
			name: "success-click-trackers-set-id-only",
			args: argsUpdateClickTrackers{
				ctx: context.TODO(),
				params: &model.UpdateClickTracker{
					Id: 1,
					ClickTrackerSetId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateClickTracker, clicktracker *model.ClickTracker, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, clicktracker, "unexpected nil click trackers")
				assert.Equal(t, params.Id, clicktracker.Id, "expected id to be equal")
				assert.Equal(t, params.ClickTrackerSetId.Int, clicktracker.ClickTrackerSetId, "expected name to be equal")
			},
		},
		//{
		//	name: "fail-mock",
		//	args: argsUpdateClickTrackers{
		//		ctx: context.TODO(),
		//		params: &model.UpdateClickTracker{
		//			Id: 1,
		//			ClickTrackerSetId: null.Int{
		//				Int:   3,
		//				Valid: true,
		//			},
		//		},
		//	},
		//	getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
		//		cleanup := func(ignoreErrors ...bool) {}
		//
		//		mockTxProvider := persistencefakes.FakeTransactionProvider{}
		//		mockTxProvider.TxReturns(nil, errors.New(mockDbReturnsErr))
		//
		//		return &dependencies{
		//			Persistor:  &clicktrackerlogicfakes.FakePersistor{},
		//			TxProvider: &mockTxProvider,
		//			Logger:     mockLogger,
		//			Cleanup:    cleanup,
		//		}, cleanup
		//	},
		//	mutations: func(t *testing.T, db *sqlx.DB) {
		//
		//	},
		//	assertions: func(t *testing.T, params *model.UpdateClickTracker, clicktracker *model.ClickTracker, err error) {
		//		assert.Error(t, err, "unexpected error")
		//		assert.Nil(t, clicktracker, "unexpected nil click trackers")
		//	},
		//},
		//{
		//	name: "fail-internal-server-error",
		//	args: argsUpdateClickTrackers{
		//		ctx: context.TODO(),
		//		params: &model.UpdateClickTracker{
		//			Id: 1,
		//			ClickTrackerSetId: null.Int{
		//				Int:   3,
		//				Valid: true,
		//			},
		//		},
		//	},
		//	getDependencies: getConcreteDependencies,
		//	mutations: func(t *testing.T, db *sqlx.DB) {
		//		testhelper.DropTable(t, db, mysqlmodel.TableNames.CapturePages)
		//	},
		//	assertions: func(t *testing.T, params *model.UpdateClickTracker, clicktracker *model.ClickTracker, err error) {
		//		assert.Error(t, err, "unexpected error")
		//		assert.Nil(t, clicktracker, "unexpected nil click trackers")
		//	},
		//},
	}
}

func TestService_UpdateClickTracker(t *testing.T) {
	for _, tt := range getTestCasesUpdateClickTracker() {
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

			clicktracker, err := svc.UpdateClickTracker(tt.args.ctx, tt.args.params)
			tt.assertions(t, tt.args.params, clicktracker, err)
		})
	}
}

// test restore method

type argsRestoreClickTracker struct {
	ctx    context.Context
	params *model.RestoreClickTracker
}

type testCaseRestoreClickTracker struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsRestoreClickTracker
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, params *model.RestoreClickTracker, err error)
}

func getTestCasesRestoreClickTracker() []testCaseRestoreClickTracker {
	return []testCaseRestoreClickTracker{
		{
			name: "success",
			args: argsRestoreClickTracker{
				ctx: context.TODO(),
				params: &model.RestoreClickTracker{
					ID: 1,
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.RestoreClickTracker, err error) {
				require.NoError(t, err, "unexpected error")
				assert.Equal(t, params.ID, 1)
			},
		},
	}
}

func TestService_RestoreClickTracker(t *testing.T) {
	for _, tt := range getTestCasesRestoreClickTracker() {
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

			err = svc.RestoreClickTracker(tt.args.ctx, tt.args.params)
			tt.assertions(t, tt.args.params, err)
		})
	}
}
