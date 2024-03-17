package categorylogic

import (
	"context"
	"errors"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic/categorylogicfakes"
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
	Cleanup    func(ignoreErrors ...bool)
}

type argsGetCategories struct {
	ctx    context.Context
	filter *model.CategoryFilters
}

type testCaseGetCategories struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsGetCategories
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, categories *model.PaginatedCategories, err error)
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

// getGetCategoryTestCases returns a list of test cases for the ListCategories method.
func getGetCategoryTestCases() []testCaseGetCategories {
	return []testCaseGetCategories{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsGetCategories{
				ctx:    context.Background(),
				filter: &model.CategoryFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected get categories error")
				require.NotNil(t, paginated, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				assert.NotEqual(t, 0, paginated.Pagination.RowCount, "unexpected total count")
			},
		},
		{
			name:            "fail-get-categories",
			getDependencies: getConcreteDependencies,
			args: argsGetCategories{
				ctx:    context.Background(),
				filter: &model.CategoryFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Category)
			},
			assertions: func(t *testing.T, paginated *model.PaginatedCategories, err error) {
				require.Error(t, err, "unexpected get categories error")
				require.Contains(t, err.Error(), "get categories:")
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
					Persistor:  &categorylogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			args: argsGetCategories{
				ctx:    context.Background(),
				filter: &model.CategoryFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedCategories, err error) {
				require.Error(t, err, "unexpected get categories error")
				require.Contains(t, err.Error(), "get db:")
			},
		},
	}
}

func TestService_GetCategories(t *testing.T) {
	for _, tt := range getGetCategoryTestCases() {
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

			paginatedCategories, err := svc.ListCategories(tt.args.ctx, tt.args.filter)
			tt.assertions(t, paginatedCategories, err)
		})
	}
}

type argsCreateCategory struct {
	ctx      context.Context
	category *model.CreateCategory
}

type testCaseCreateCategory struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsCreateCategory
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, category *model.Category, err error)
}

func TestService_CreateCategory(t *testing.T) {
	for _, tt := range getTestCasesCreateCategory() {
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

			category, err := svc.CreateCategory(tt.args.ctx, tt.args.category)
			tt.assertions(t, category, err)
		})
	}
}

func getTestCasesCreateCategory() []testCaseCreateCategory {
	return []testCaseCreateCategory{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsCreateCategory{
				ctx: context.TODO(),
				category: &model.CreateCategory{
					CategoryTypeRefId: 1,
					Name:              "Example",
				},
			},
			assertions: func(t *testing.T, category *model.Category, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, category, "unexpected nil category")

				require.NotEqual(t, 0, category.Id, "unexpected nil category")
				require.NotEmpty(t, category.Name, "unexpected empty category name")
				require.NotEqual(t, 0, category.CategoryTypeRefId, "unexpected empty category type ref id")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsCreateCategory{
				ctx: context.TODO(),
				category: &model.CreateCategory{
					CategoryTypeRefId: 1,
					Name:              "Example",
				},
			},
			assertions: func(t *testing.T, category *model.Category, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, category, "unexpected nil category")

				require.NotEqual(t, 0, category.Id, "unexpected nil category")
				require.NotEmpty(t, category.Name, "unexpected empty category name")
				require.NotEqual(t, 0, category.CategoryTypeRefId, "unexpected empty category type ref id")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "fail-internal-error",
			getDependencies: getConcreteDependencies,
			args: argsCreateCategory{
				ctx: context.TODO(),
				category: &model.CreateCategory{
					CategoryTypeRefId: 1,
					Name:              "Example",
				},
			},
			assertions: func(t *testing.T, category *model.Category, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, category, "unexpected nil category")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Category)
			},
		},
		{
			name:            "fail-invalid-args",
			getDependencies: getConcreteDependencies,
			args: argsCreateCategory{
				ctx: context.TODO(),
				category: &model.CreateCategory{
					CategoryTypeRefId: 0,
					Name:              "Example",
				},
			},
			assertions: func(t *testing.T, category *model.Category, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, category, "unexpected nil category")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "fail-unique",
			getDependencies: getConcreteDependencies,
			args: argsCreateCategory{
				ctx: context.TODO(),
				category: &model.CreateCategory{
					CategoryTypeRefId: 1,
					Name:              "Example",
				},
			},
			assertions: func(t *testing.T, category *model.Category, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, category, "unexpected nil category")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Category{
					Name:              "Example",
					CategoryTypeRefID: 1,
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
					Persistor:  &categorylogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    fakeCleanup,
				}, fakeCleanup
			},
			args: argsCreateCategory{
				ctx: context.TODO(),
				category: &model.CreateCategory{
					CategoryTypeRefId: 1,
					Name:              "Example",
				},
			},
			assertions: func(t *testing.T, category *model.Category, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, category, "unexpected nil category")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
	}
}

type argsUpdateCategories struct {
	ctx    context.Context
	params *model.UpdateCategory
}

type testCaseUpdateCategory struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsUpdateCategories
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, params *model.UpdateCategory, category *model.Category, err error)
}

func getTestCasesUpdateCategories() []testCaseUpdateCategory {
	return []testCaseUpdateCategory{
		{
			name: "success",
			args: argsUpdateCategories{
				ctx: context.TODO(),
				params: &model.UpdateCategory{
					Id: 1,
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					CategoryTypeRefId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateCategory, category *model.Category, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, category, "unexpected nil category")
				assert.Equal(t, params.Id, category.Id, "expected id to be equal")
				assert.Equal(t, params.CategoryTypeRefId.Int, category.CategoryTypeRefId, "expected category type ref id to be equal")
				assert.Equal(t, params.Name.String, category.Name, "expected name to be equal")
			},
		},
		{
			name: "success-name-only",
			args: argsUpdateCategories{
				ctx: context.TODO(),
				params: &model.UpdateCategory{
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
			assertions: func(t *testing.T, params *model.UpdateCategory, category *model.Category, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, category, "unexpected nil category")
				assert.Equal(t, params.Name.String, category.Name, "expected name to be equal")
			},
		},
		{
			name: "success-category-type-ref-id-only",
			args: argsUpdateCategories{
				ctx: context.TODO(),
				params: &model.UpdateCategory{
					Id: 1,
					CategoryTypeRefId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateCategory, category *model.Category, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, category, "unexpected nil category")
				assert.Equal(t, params.Id, category.Id, "expected id to be equal")
				assert.Equal(t, params.CategoryTypeRefId.Int, category.CategoryTypeRefId, "expected name to be equal")
			},
		},
		{
			name: "fail-mock",
			args: argsUpdateCategories{
				ctx: context.TODO(),
				params: &model.UpdateCategory{
					Id: 1,
					CategoryTypeRefId: null.Int{
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
					Persistor:  &categorylogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateCategory, category *model.Category, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, category, "unexpected nil category")
			},
		},
		{
			name: "fail-internal-server-error",
			args: argsUpdateCategories{
				ctx: context.TODO(),
				params: &model.UpdateCategory{
					Id: 1,
					CategoryTypeRefId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Category)
			},
			assertions: func(t *testing.T, params *model.UpdateCategory, category *model.Category, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, category, "unexpected nil category")
			},
		},
	}
}

func TestService_UpdateCategory(t *testing.T) {
	for _, tt := range getTestCasesUpdateCategories() {
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

			category, err := svc.UpdateCategory(tt.args.ctx, tt.args.params)
			tt.assertions(t, tt.args.params, category, err)
		})
	}
}
