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
	ctx    context.Context
	filter *model.OrganizationFilters
}

type testCaseGetOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsGetOrganizations
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, organization *model.PaginatedOrganizations, err error)
}

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
				user := mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					CategoryTypeRefID: 1,
					IsActive:          true,
				}
				err := user.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting a user into the user table")

				entry := mysqlmodel.Organization{
					Name:          "Demby",
					CreatedBy:     null.IntFrom(user.ID),
					LastUpdatedBy: null.IntFrom(user.ID),
				}
				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
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
			mutations: func(t *testing.T, db *sqlx.DB) {},
			assertions: func(t *testing.T, paginated *model.PaginatedOrganizations, err error) {
				require.Error(t, err, "unexpected get organizations error")
				require.Contains(t, err.Error(), "get db:")
			},
		},
	}
}

func TestService_GetOrganizations(t *testing.T) {
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

type testCaseUpdateOrganization struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	ctx             context.Context
	mutations       func(t *testing.T, db *sqlx.DB) *model.UpdateOrganization
	assertions      func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error)
}

func getTestCasesUpdateOrganizations() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name:            "success",
			ctx:             context.TODO(),
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) *model.UpdateOrganization {
				user, err := mysqlmodel.Users().One(context.TODO(), db)
				require.NoError(t, err, "error finding user from user table.")

				entry := mysqlmodel.Organization{
					ID:            user.ID,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(user.ID),
					LastUpdatedBy: null.IntFrom(user.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}

				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				toUpdateOrganization := &model.UpdateOrganization{
					Id: user.ID,
					Name: null.String{
						String: "Lawrence",
						Valid:  true,
					},
					UserId: null.Int{
						Int:   user.ID,
						Valid: true,
					},
				}

				return toUpdateOrganization
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
			mutations: func(t *testing.T, db *sqlx.DB) *model.UpdateOrganization {
				newUser := mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					CategoryTypeRefID: 1,
					IsActive:          true,
				}

				err := newUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting new user")

				entry := mysqlmodel.Organization{
					ID:            newUser.ID,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(newUser.ID),
					LastUpdatedBy: null.IntFrom(newUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}

				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample organization data")

				toUpdateOrganization := &model.UpdateOrganization{
					Id: newUser.ID,
					Name: null.String{
						String: "Demby2",
						Valid:  true,
					},
					UserId: null.Int{
						Int:   newUser.ID,
						Valid: true,
					},
				}

				return toUpdateOrganization
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
			mutations: func(t *testing.T, db *sqlx.DB) *model.UpdateOrganization {
				toUpdateOrganization := &model.UpdateOrganization{
					Id: rand.Int(),
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					UserId: null.Int{
						Int:   rand.Int(),
						Valid: true,
					},
				}

				return toUpdateOrganization
			},
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
			mutations: func(t *testing.T, db *sqlx.DB) (toUpdatedOrganization *model.UpdateOrganization) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Organization)
				return nil
			},
			assertions: func(t *testing.T, params *model.UpdateOrganization, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "validate: validate: error, struct is nil", "expected error message to contain 'sql: no rows in result set'")
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

			updateOrganizationStruct := tt.mutations(t, _dependencies.Db)

			organization, err := svc.UpdateOrganization(tt.ctx, updateOrganizationStruct)
			tt.assertions(t, updateOrganizationStruct, organization, err)
		})
	}
}

type testCaseCreateOrganization struct {
	name            string
	ctx             context.Context
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	mutations       func(t *testing.T, db *sqlx.DB) *model.CreateOrganization
	assertions      func(t *testing.T, organization *model.Organization, err error)
}

func getTestCasesCreateOrganization() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			ctx:             context.TODO(),
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, organization, "unexpected nil organization")
				require.NotZero(t, organization.Id, "unexpected zero value for organization id")
				require.NotEmpty(t, organization.Name, "unexpected empty organization name")
				require.NotZero(t, organization.CreatedBy, "unexpected zero value for CreatedBy")
				assert.Contains(t, organization.Name, "Demby", "expected organization name to contain 'Demby'")
			},
			mutations: func(t *testing.T, db *sqlx.DB) *model.CreateOrganization {
				user, err := mysqlmodel.Users().One(context.TODO(), db)
				require.NoError(t, err, "error fetching one user record")

				createOrg := model.CreateOrganization{
					Name:   "Demby",
					UserId: user.ID,
				}

				return &createOrg
			},
		},
		{
			name:            "fail-internal-error",
			getDependencies: getConcreteDependencies,
			ctx:             context.TODO(),
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "validate: validate: error, struct is nil", "expected error message to contain 'sql: no rows in result set'")
			},
			mutations: func(t *testing.T, db *sqlx.DB) *model.CreateOrganization {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.Organization)
				return nil
			},
		},
		{
			name:            "fail-invalid-args",
			getDependencies: getConcreteDependencies,
			ctx:             context.TODO(),
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "create: insert organization: mysqlmodel: unable to insert into organization: Error 1264: Out of range value for column 'created_by' at row 1", "expected error message to contain 'sql: no rows in result set'")
			},
			mutations: func(t *testing.T, db *sqlx.DB) *model.CreateOrganization {
				createOrg := &model.CreateOrganization{
					Name:   "Demby",
					UserId: rand.Int(),
				}

				return createOrg
			},
		},
		{
			name:            "fail-unique",
			getDependencies: getConcreteDependencies,
			ctx:             context.TODO(),
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, organization, "unexpected nil organization")
				assert.Contains(t, err.Error(), "organization already exists", "expected error message to contain 'sql: no rows in result set'")
			},
			mutations: func(t *testing.T, db *sqlx.DB) *model.CreateOrganization {
				user, err := mysqlmodel.Users().One(context.TODO(), db)
				require.NoError(t, err, "error fetching one user record")

				entryOrganization := mysqlmodel.Organization{
					ID:            user.ID,
					Name:          "Demby2",
					CreatedBy:     null.IntFrom(user.ID),
					LastUpdatedBy: null.IntFrom(user.ID),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting into organization table")

				createOrg := &model.CreateOrganization{
					Name:   "Demby2",
					UserId: user.ID,
				}

				return createOrg
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
			ctx: context.TODO(),
			assertions: func(t *testing.T, organization *model.Organization, err error) {
				assert.Error(t, err, "expected error but got none")
				assert.Nil(t, organization, "expected nil organization but got non-nil")
				assert.Contains(t, err.Error(), "get db: mock database error", "expected error message to contain 'sql: no rows in result set'")

			},
			mutations: func(t *testing.T, db *sqlx.DB) *model.CreateOrganization {
				return &model.CreateOrganization{
					Name:   "Demby",
					UserId: rand.Int(),
				}
			},
		},
	}
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

			toBeCreatedOrganization := tt.mutations(t, _dependencies.Db)

			organization, err := svc.AddOrganization(tt.ctx, toBeCreatedOrganization)
			tt.assertions(t, organization, err)
		})
	}
}

type testCaseDeleteOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	mutations       func(t *testing.T, db *sqlx.DB) *model.DeleteOrganization
	assertions      func(t *testing.T, db *sqlx.DB, id int)
}

func getTestCasesDeleteOrganizations() []testCaseDeleteOrganizations {
	return []testCaseDeleteOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB) *model.DeleteOrganization {
				user, err := mysqlmodel.Users().One(context.TODO(), db)
				require.NoError(t, err, "error fetching user from user table")

				entryOrganization := model.DeleteOrganization{
					ID: user.ID,
				}

				return &entryOrganization
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
			mutations: func(t *testing.T, db *sqlx.DB) *model.DeleteOrganization {
				user, err := mysqlmodel.Users().One(context.TODO(), db)
				require.NoError(t, err, "error fetching user from user table")

				entryOrganization := model.DeleteOrganization{
					ID: user.ID,
				}

				return &entryOrganization
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

			delOrganization := testCase.mutations(t, _dependencies.Db)
			require.NoError(t, err, "unexpected new service error")

			err = svc.DeleteOrganization(context.TODO(), delOrganization)
			require.NoError(t, err, "unexpected error deleting organization.")
			testCase.assertions(t, db, delOrganization.ID)
		})
	}
}

type argsRestoreOrganization struct {
	ctx    context.Context
	params *model.RestoreOrganization
}

type testCaseRestoreOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsRestoreOrganization
	mutations       func(t *testing.T, db *sqlx.DB)
	assertions      func(t *testing.T, db *sqlx.DB, id int, err error)
}

func getTestCasesRestoreOrganizations() []testCaseRestoreOrganizations {
	return []testCaseRestoreOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsRestoreOrganization{
				ctx: context.TODO(),
				params: &model.RestoreOrganization{
					ID: 4,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				newUser := mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					CategoryTypeRefID: 1,
					IsActive:          true,
				}

				err := newUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting new user")

				entryOrganization := mysqlmodel.Organization{
					ID:            newUser.ID,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(newUser.ID),
					LastUpdatedBy: null.IntFrom(newUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
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
			args: argsRestoreOrganization{
				ctx: context.TODO(),
				params: &model.RestoreOrganization{
					ID: rand.Int(),
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {},
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

			testCase.mutations(t, _dependencies.Db)
			require.NoError(t, err, "unexpected error in mutations")

			err = svc.RestoreOrganization(testCase.args.ctx, testCase.args.params)
			testCase.assertions(t, _dependencies.Db, testCase.args.params.ID, err)
		})
	}
}
