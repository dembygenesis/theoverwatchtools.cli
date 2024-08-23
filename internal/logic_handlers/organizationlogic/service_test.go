package organizationlogic

import (
	"context"
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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"math/rand"
	"strconv"
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

type argsGetOrganizations struct {
	User         mysqlmodel.User
	Organization mysqlmodel.Organization
	filters      *model.OrganizationFilters
}

type testCaseGetOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsGetOrganizations
	mutations       func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations)
	assertions      func(t *testing.T, organization *model.PaginatedOrganizations, err error)
}

func getGetOrganizationTestCases() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsGetOrganizations{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				filters: &model.OrganizationFilters{},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization")
			},
			assertions: func(t *testing.T, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected get organizations error")
				require.NotNil(t, paginated, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				assert.Equal(t, 1, paginated.Pagination.RowCount, "unexpected total count")
			},
		},
		{
			name:            "fail-get-organizations",
			getDependencies: getConcreteDependencies,
			args:            &argsGetOrganizations{},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) {
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
			args:      &argsGetOrganizations{},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) {},
			assertions: func(t *testing.T, paginated *model.PaginatedOrganizations, err error) {
				require.Error(t, err, "unexpected get organizations error")
				require.Contains(t, err.Error(), "get db:")
			},
		},
	}
}

func TestService_GetOrganizations(t *testing.T) {
	for _, testCase := range getGetOrganizationTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			_dependencies, cleanup := testCase.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			testCase.mutations(t, _dependencies.Db, testCase.args)

			paginatedOrganizations, err := svc.ListOrganizations(context.Background(), testCase.args.filters)
			testCase.assertions(t, paginatedOrganizations, err)
		})
	}
}

type argsUpdateOrganizations struct {
	User          mysqlmodel.User
	Organization  mysqlmodel.Organization
	UpdateOrgData *model.UpdateOrganization
}

type testCaseUpdateOrganization struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	ctx             context.Context
	args            *argsUpdateOrganizations
	mutations       func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganizations)
	assertions      func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error)
}

func getTestCasesUpdateOrganizations() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name:            "success",
			ctx:             context.TODO(),
			getDependencies: getConcreteDependencies,
			args: &argsUpdateOrganizations{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Lawrence",
					Lastname:          "Margaja",
					CategoryTypeRefID: 1,
					IsActive:          true,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				UpdateOrgData: &model.UpdateOrganization{
					Id: 1,
					Name: null.String{
						String: "Lawrence",
						Valid:  true,
					},
					UserId: null.Int{
						Int:   4,
						Valid: true,
					},
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganizations) {
				err := args.User.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into user table")

				err = args.Organization.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into organization table")
			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				createdByConvToInt, convErr := strconv.Atoi(organization.CreatedBy)
				require.NoError(t, convErr, "error converting createdBy into an int")
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")
				assert.Equal(t, params.Id, organization.Id, "expected id to be equal")
				assert.Equal(t, params.UserId.Int, createdByConvToInt, "expected created_by id to be equal")
				assert.Equal(t, params.Name.String, "Lawrence", "expected name to be equal")
				assert.Contains(t, organization.Name, "Lawrence", "expected organization name to contain 'Lawrence'")
			},
		},
		{
			name:            "success-name-only",
			ctx:             context.TODO(),
			getDependencies: getConcreteDependencies,
			args: &argsUpdateOrganizations{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Lawrence",
					Lastname:          "Margaja",
					CategoryTypeRefID: 1,
					IsActive:          true,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				UpdateOrgData: &model.UpdateOrganization{
					Id: 1,
					Name: null.String{
						String: "Demby2",
						Valid:  true,
					},
					UserId: null.Int{
						Int:   4,
						Valid: true,
					},
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganizations) {
				err := args.User.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into user table")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")
				assert.Equal(t, params.Name.String, "Demby2", "expected name to be equal")
				assert.Contains(t, organization.Name, "Demby2", "expected organization name to contain 'Demby2'")
			},
		},
		{
			name: "fail-mock",
			ctx:  context.TODO(),
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
			args: &argsUpdateOrganizations{
				UpdateOrgData: &model.UpdateOrganization{
					Id: rand.Int(),
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					UserId: null.Int{
						Int:   rand.Int(),
						Valid: true,
					},
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganizations) {},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Contains(t, err.Error(), "get db: error getting db", "expected error message to contain 'get db: error getting db'")
				assert.Nil(t, organization, "expected organization to be nil due to failure")
				assert.Contains(t, err.Error(), "get db: error getting db", "expected error message to contain 'sql: no rows in result set'")
			},
		},
		{
			name:            "fail-internal-server-error",
			ctx:             context.TODO(),
			getDependencies: getConcreteDependencies,
			args: &argsUpdateOrganizations{
				UpdateOrgData: &model.UpdateOrganization{
					Id: rand.Int(),
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					UserId: null.Int{
						Int:   rand.Int(),
						Valid: true,
					},
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganizations) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Organization)
			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "unable to update organization row", "expected error message to contain 'unable to update organization row'")
			},
		},
	}
}

func TestService_UpdateOrganization(t *testing.T) {
	for _, testCase := range getTestCasesUpdateOrganizations() {
		t.Run(testCase.name, func(t *testing.T) {
			_dependencies, cleanup := testCase.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			testCase.mutations(t, _dependencies.Db, testCase.args)

			organization, err := svc.UpdateOrganization(testCase.ctx, testCase.args.UpdateOrgData)
			testCase.assertions(t, testCase.args.UpdateOrgData, organization, err)
		})
	}
}

type argsCreateOrganization struct {
	User                mysqlmodel.User
	CreatedOrganization mysqlmodel.Organization
	Organization        *model.CreateOrganization
}

type testCaseCreateOrganization struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsCreateOrganization
	mutations       func(t *testing.T, db *sqlx.DB, args *argsCreateOrganization)
	assertions      func(t *testing.T, organization *model.Organization, err error)
}

func getTestCasesCreateOrganization() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")
				require.NotZero(t, organization.Id, "unexpected zero value for organization id")
				require.NotEmpty(t, organization.Name, "unexpected empty organization name")
				require.NotZero(t, organization.CreatedBy, "unexpected zero value for CreatedBy")
				assert.Contains(t, organization.Name, "Demby", "expected organization name to contain 'Demby'")
			},
			args: &argsCreateOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: &model.CreateOrganization{
					Name:   "Demby",
					UserId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")
			},
		},
		{
			name:            "fail-internal-error",
			getDependencies: getConcreteDependencies,
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "failed to count organization rows", "expected error message to contain 'failed to count organization rows'")
			},
			args: &argsCreateOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: &model.CreateOrganization{
					Name:   "Demby",
					UserId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateOrganization) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Organization)
			},
		},
		{
			name:            "fail-invalid-args",
			getDependencies: getConcreteDependencies,
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "Cannot add or update a child row: a foreign key constraint fails", "expected error message to contain foreign key constraint failure")
			},
			args: &argsCreateOrganization{
				User: mysqlmodel.User{
					Firstname:         "HAHA",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 2,
				},
				Organization: &model.CreateOrganization{
					Name:   "Younes",
					UserId: 9999,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")
			},
		},
		{
			name:            "fail-unique",
			getDependencies: getConcreteDependencies,
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "organization already exists", "expected error message to contain 'sql: no rows in result set'")
			},
			args: &argsCreateOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				CreatedOrganization: mysqlmodel.Organization{
					Name:          "Demby",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				Organization: &model.CreateOrganization{
					Name:   "Demby",
					UserId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateOrganization) {
				err := args.User.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into user table")

				err = args.CreatedOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into user table")
			},
		},
		{
			name: "fail-mock",
			getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
				fakeCleanup := func(ignoreErrors ...bool) {}
				mockTxProvider := &persistencefakes.FakeTransactionProvider{}
				mockTxProvider.TxReturns(nil, errors.New("mock database error"))

				return &dependencies{
					Persistor:  &organizationlogicfakes.FakePersistor{},
					TxProvider: mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    fakeCleanup,
				}, fakeCleanup
			},
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "expected error but got none")
				assert.Nil(t, organization, "expected nil organization but got non-nil")
				assert.Contains(t, err.Error(), "get db: mock database error", "expected error message to contain 'sql: no rows in result set'")

			},
			args: &argsCreateOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				CreatedOrganization: mysqlmodel.Organization{
					Name:          "Demby",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				Organization: &model.CreateOrganization{
					Name:   "Dembo",
					UserId: 5,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateOrganization) {},
		},
	}
}

func TestService_CreateOrganization(t *testing.T) {
	for _, testCase := range getTestCasesCreateOrganization() {
		t.Run(testCase.name, func(t *testing.T) {
			_dependencies, cleanup := testCase.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			testCase.mutations(t, _dependencies.Db, testCase.args)

			organization, err := svc.CreateOrganization(context.Background(), testCase.args.Organization)
			testCase.assertions(t, organization, err)
		})
	}
}

type argsDeleteOrganization struct {
	User            mysqlmodel.User
	Organization    mysqlmodel.Organization
	DelOrganization *model.DeleteOrganization
}

type testCaseDeleteOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsDeleteOrganization
	mutations       func(t *testing.T, db *sqlx.DB, args *argsDeleteOrganization)
	assertions      func(t *testing.T, db *sqlx.DB, id int)
}

func getTestCasesDeleteOrganizations() []testCaseDeleteOrganizations {
	return []testCaseDeleteOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsDeleteOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				DelOrganization: &model.DeleteOrganization{
					ID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteOrganization) {
				err := args.User.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into user model")

				err = args.Organization.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into organization table")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnOrganization, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Nil(t, returnOrganization, "expected to be nil")
				require.Error(t, err, "error fetching organization from db")
				assert.Contains(t, err.Error(), "sql: no rows in result set", "expected error message to contain 'sql: no rows in result set'")
			},
		},
		{
			name:            "fail-organization-not-found",
			getDependencies: getConcreteDependencies,
			args: &argsDeleteOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            3,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				DelOrganization: &model.DeleteOrganization{
					ID: 2,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into user model")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnOrganization, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Nil(t, returnOrganization, "expected to be nil")
				require.Error(t, err, "expected an error when deleting a non-existent organization")
				assert.Contains(t, err.Error(), "sql: no rows in result set", "expected error message to contain 'sql: no rows in result set'")
			},
		},
	}
}

func TestDeleteOrganization(t *testing.T) {
	for _, testCase := range getTestCasesDeleteOrganizations() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			_dependencies, cleanup := testCase.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			testCase.mutations(t, _dependencies.Db, testCase.args)
			require.NoError(t, err, "unexpected new service error")

			err = svc.DeleteOrganization(context.TODO(), testCase.args.DelOrganization)
			require.NoError(t, err, "unexpected error deleting organization.")
			testCase.assertions(t, db, testCase.args.DelOrganization.ID)
		})
	}
}

type argsRestoreOrganization struct {
	User         mysqlmodel.User
	Organization mysqlmodel.Organization
	params       *model.RestoreOrganization
}

type testCaseRestoreOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsRestoreOrganization
	mutations       func(t *testing.T, db *sqlx.DB, args *argsRestoreOrganization)
	assertions      func(t *testing.T, db *sqlx.DB, id int, err error)
}

func getTestCasesRestoreOrganizations() []testCaseRestoreOrganizations {
	return []testCaseRestoreOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsRestoreOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				params: &model.RestoreOrganization{
					ID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.NoError(t, err, "error restoring organization from db")

				returnOrganization, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "error fetching organization from db")
				require.NotNil(t, returnOrganization, "expected organization to be not nil")

				require.True(t, returnOrganization.IsActive, "expected organization to be active")
				assert.Equal(t, returnOrganization.IsActive, true)
				assert.Contains(t, returnOrganization.Name, "Demby", "expected organization name to contain 'Demby'")
			},
		},
		{
			name:            "fail-organization-not-found",
			getDependencies: getConcreteDependencies,
			args: &argsRestoreOrganization{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            2,
					Name:          "Younes",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				params: &model.RestoreOrganization{
					ID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				returnOrganization, fetchErr := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Error(t, fetchErr, "expected an error when fetching organization from db")
				require.Nil(t, returnOrganization, "expected organization to be nil")
				assert.Contains(t, fetchErr.Error(), "sql: no rows in result set", "expected error message to contain 'sql: no rows in result set'")
			},
		},
	}
}

func TestRestoreOrganization(t *testing.T) {
	for _, testCase := range getTestCasesRestoreOrganizations() {
		t.Run(testCase.name, func(t *testing.T) {
			_dependencies, cleanup := testCase.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			testCase.mutations(t, _dependencies.Db, testCase.args)
			require.NoError(t, err, "unexpected error in mutations")

			err = svc.RestoreOrganization(context.Background(), testCase.args.params)
			testCase.assertions(t, _dependencies.Db, testCase.args.params.ID, err)
		})
	}
}
