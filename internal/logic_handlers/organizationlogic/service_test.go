package organizationlogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
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

type argsDeleteOrganization struct {
	ctx    context.Context
	params *model.DeleteOrganization
}

type testCaseDeleteOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            argsDeleteOrganization
	mutations       func(t *testing.T, db *sqlx.DB) []int
	assertions      func(t *testing.T, db *sqlx.DB, id int)
}

func getTestCasesDeleteOrganizations() []testCaseDeleteOrganizations {
	return []testCaseDeleteOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: argsDeleteOrganization{
				ctx: context.TODO(),
				params: &model.DeleteOrganization{
					ID: 4,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entryOrganization := mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{entryOrganization.ID}

				return Ids
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnOrganization, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Nil(t, returnOrganization, "expected to be nil")
				require.Error(t, err, "error fetching organization from db")
			},
		},
		{
			name:            "fail-organization-not-found",
			getDependencies: getConcreteDependencies,
			args: argsDeleteOrganization{
				ctx: context.TODO(),
				params: &model.DeleteOrganization{
					ID: 1232,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entryOrganization := mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{entryUser.ID}

				return Ids
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnOrganization, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Nil(t, returnOrganization, "expected to be nil")
				require.Error(t, err, "expected an error when deleting a non-existent organization")
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

			OrgIds := testCase.mutations(t, _dependencies.Db)
			require.NoError(t, err, "unexpected new service error")

			err = svc.DeleteOrganization(testCase.args.ctx, testCase.args.params)
			require.NoError(t, err, "unexpected error deleting organization.")
			testCase.assertions(t, db, OrgIds[0])
		})
	}
}

type testCaseRestoreOrganizations struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	ctx             context.Context
	mutations       func(t *testing.T, db *sqlx.DB) *model.RestoreOrganization
	assertions      func(t *testing.T, db *sqlx.DB, id int, err error)
}

func getTestCasesRestoreOrganizations() []testCaseRestoreOrganizations {
	return []testCaseRestoreOrganizations{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			ctx:             context.TODO(),
			mutations: func(t *testing.T, db *sqlx.DB) *model.RestoreOrganization {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entryOrganization := mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      false,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				toBeRestoredOrganization := &model.RestoreOrganization{
					ID: entryOrganization.ID,
				}

				return toBeRestoredOrganization
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.NoError(t, err, "error restoring organization from db")

				returnOrganization, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "error fetching organization from db")
				require.NotNil(t, returnOrganization, "expected organization to be not nil")

				require.True(t, returnOrganization.IsActive, "expected organization to be active")
				assert.Equal(t, returnOrganization.IsActive, true)
			},
		},
		{
			name:            "fail-organization-not-found",
			getDependencies: getConcreteDependencies,
			ctx:             context.TODO(),
			mutations: func(t *testing.T, db *sqlx.DB) *model.RestoreOrganization {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entryOrganization := mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      false,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				toBeRestoredOrganization := &model.RestoreOrganization{
					ID: entryOrganization.ID + 12321,
				}

				return toBeRestoredOrganization
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				returnOrganization, fetchErr := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Error(t, fetchErr, "expected an error when fetching organization from db")
				require.Nil(t, returnOrganization, "expected organization to be nil")
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

			toBeRestoredOrganization := testCase.mutations(t, _dependencies.Db)
			require.NoError(t, err, "unexpected error in mutations")

			err = svc.RestoreOrganization(testCase.ctx, toBeRestoredOrganization)
			testCase.assertions(t, _dependencies.Db, toBeRestoredOrganization.ID, err)
		})
	}
}
