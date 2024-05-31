package organizationlogic

import (
	"context"
	"errors"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/organizationlogic/organizationlogicfakes"
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

type argsGetOrganizations struct {
	ctx    context.Context
	filter *model.OrganizationFilters
}

type testCaseGetOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsGetOrganizations
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, organizations *model.PaginatedOrganizations, err error)
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

// getGetOrganizationTestCases returns a list of test cases for the ListOrganizations method.
func getGetOrganizationTestCases() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsGetOrganizations{
				ctx:    context.Background(),
				filter: &model.OrganizationFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected get organizations error")
				require.NotNil(t, paginated, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				assert.NotEqual(t, 0, paginated.Pagination.RowCount, "unexpected total count")
			},
		},
		{
			name:            "fail-get-organizations",
			getDependencies: getConcreteDependencies,
			args: argsGetOrganizations{
				ctx:    context.Background(),
				filter: &model.OrganizationFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Organization)
			},
			assertions: func(t *testing.T, paginated *model.PaginatedOrganizations, err error) {
				require.Error(t, err, "unexpected get organizations error")
				require.Contains(t, err.Error(), "get organizations:")
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
					Persistor:  &organizationlogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			args: argsGetOrganizations{
				ctx:    context.Background(),
				filter: &model.OrganizationFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, paginated *model.PaginatedOrganizations, err error) {
				require.Error(t, err, "unexpected get organizations error")
				require.Contains(t, err.Error(), "get db:")
			},
		},
	}
}

func TestService_GetCategories(t *testing.T) {
	for _, tt := range getGetOrganizationTestCases() {
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

			paginatedOrganizations, err := svc.ListOrganizations(tt.args.ctx, tt.args.filter)
			tt.assertions(t, paginatedOrganizations, err)
		})
	}
}

type argsCreateOrganization struct {
	ctx          context.Context
	organization *model.CreateOrganization
}

type testCaseCreateOrganization struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsCreateOrganization
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, category *model.Organization, err error)
}

func TestService_CreateOrganization(t *testing.T) {
	for _, tt := range getTestCasesCreateOrganization() {
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

			category, err := svc.CreateOrganization(tt.args.ctx, tt.args.organization)
			tt.assertions(t, category, err)
		})
	}
}

func getTestCasesCreateOrganization() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsCreateOrganization{
				ctx: context.TODO(),
				organization: &model.CreateOrganization{
					OrganizationTypeRefId: 1,
					Name:                  "Example",
				},
			},
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")

				require.NotEqual(t, 0, organization.Id, "unexpected nil organization")
				require.NotEmpty(t, organization.Name, "unexpected empty organization name")
				require.NotEqual(t, 0, organization.OrganizationTypeRefId, "unexpected empty organization type ref id")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsCreateOrganization{
				ctx: context.TODO(),
				organization: &model.CreateOrganization{
					OrganizationTypeRefId: 1,
					Name:                  "Example",
				},
			},
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")

				require.NotEqual(t, 0, organization.Id, "unexpected nil organization")
				require.NotEmpty(t, organization.Name, "unexpected empty organization name")
				require.NotEqual(t, 0, organization.OrganizationTypeRefId, "unexpected empty organization type ref id")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "fail-internal-error",
			getDependencies: getConcreteDependencies,
			args: argsCreateOrganization{
				ctx: context.TODO(),
				organization: &model.CreateOrganization{
					OrganizationTypeRefId: 1,
					Name:                  "Example",
				},
			},
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Organization)
			},
		},
		{
			name:            "fail-invalid-args",
			getDependencies: getConcreteDependencies,
			args: argsCreateOrganization{
				ctx: context.TODO(),
				organization: &model.CreateOrganization{
					OrganizationTypeRefId: 0,
					Name:                  "Example",
				},
			},
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
		{
			name:            "fail-unique",
			getDependencies: getConcreteDependencies,
			args: argsCreateOrganization{
				ctx: context.TODO(),
				organization: &model.CreateOrganization{
					OrganizationTypeRefId: 1,
					Name:                  "Example",
				},
			},
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Organization{
					Name:                  "Example",
					OrganizationTypeRefID: 1,
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
					Persistor:  &organizationlogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    fakeCleanup,
				}, fakeCleanup
			},
			args: argsCreateOrganization{
				ctx: context.TODO(),
				organization: &model.CreateOrganization{
					OrganizationTypeRefId: 1,
					Name:                  "Example",
				},
			},
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
		},
	}
}

type argsUpdateOrganization struct {
	ctx    context.Context
	params *model.UpdateOrganization
}

type testCaseUpdateOrganization struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsUpdateOrganization
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error)
}

func getTestCasesUpdateOrganizations() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name: "success",
			args: argsUpdateOrganization{
				ctx: context.TODO(),
				params: &model.UpdateOrganization{
					Id: 1,
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					OrganizationTypeRefId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")
				assert.Equal(t, params.Id, organization.Id, "expected id to be equal")
				assert.Equal(t, params.OrganizationTypeRefId.Int, organization.OrganizationTypeRefId, "expected organization type ref id to be equal")
				assert.Equal(t, params.Name.String, organization.Name, "expected name to be equal")
			},
		},
		{
			name: "success-name-only",
			args: argsUpdateOrganization{
				ctx: context.TODO(),
				params: &model.UpdateOrganization{
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
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")
				assert.Equal(t, params.Name.String, organization.Name, "expected name to be equal")
			},
		},
		{
			name: "success-organization-type-ref-id-only",
			args: argsUpdateOrganization{
				ctx: context.TODO(),
				params: &model.UpdateOrganization{
					Id: 1,
					OrganizationTypeRefId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")
				assert.Equal(t, params.Id, organization.Id, "expected id to be equal")
				assert.Equal(t, params.OrganizationTypeRefId.Int, organization.OrganizationTypeRefId, "expected name to be equal")
			},
		},
		{
			name: "fail-mock",
			args: argsUpdateOrganization{
				ctx: context.TODO(),
				params: &model.UpdateOrganization{
					Id: 1,
					OrganizationTypeRefId: null.Int{
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
					Persistor:  &organizationlogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
			},
		},
		{
			name: "fail-internal-server-error",
			args: argsUpdateOrganization{
				ctx: context.TODO(),
				params: &model.UpdateOrganization{
					Id: 1,
					OrganizationTypeRefId: null.Int{
						Int:   3,
						Valid: true,
					},
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Organization)
			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
			},
		},
	}
}

func TestService_UpdateOrganization(t *testing.T) {
	for _, tt := range getTestCasesUpdateOrganizations() {
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

			organization, err := svc.UpdateOrganization(tt.args.ctx, tt.args.params)
			tt.assertions(t, tt.args.params, organization, err)
		})
	}
}
