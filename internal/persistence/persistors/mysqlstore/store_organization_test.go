package mysqlstore

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore/testhelper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"testing"
	"time"
)

type testCaseGetOrganizations struct {
	name       string
	filter     *model.OrganizationFilters
	mutations  func(t *testing.T, db *sqlx.DB)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error)
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(2),
					LastUpdatedBy: null.IntFrom(3),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err := entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 1, "unexpected number of organizations")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected row count")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.OrganizationFilters{
				OrganizationNameIn: []string{"Organization A", "Organization B"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := []mysqlmodel.Organization{
					{
						ID:            1,
						Name:          "Organization A",
						CreatedBy:     null.IntFrom(2),
						LastUpdatedBy: null.IntFrom(3),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true}, {ID: 2,
						Name:          "Organization B",
						CreatedBy:     null.IntFrom(2),
						LastUpdatedBy: null.IntFrom(3),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true},
				}

				for _, org := range entry {
					err := org.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting sample data")
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organization")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 2, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount == 2, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.OrganizationFilters{
				IdsIn:              []int{1},
				OrganizationNameIn: []string{"Organization A"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Organization{
					ID:            1,
					Name:          "Organization A",
					CreatedBy:     null.IntFrom(2),
					LastUpdatedBy: null.IntFrom(3),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err := entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organization")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 1, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "empty-results",
			filter: &model.OrganizationFilters{
				IdsIn:              []int{2},
				OrganizationNameIn: []string{"Super Admin"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organization")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 0, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
	}
}

func Test_GetOrganizations(t *testing.T) {
	for _, testCase := range getTestCasesGetOrganizations() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			testCase.mutations(t, db)
			paginated, err := m.GetOrganizations(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}

type deleteOrganizationTestCase struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB)
}

func getDeleteOrganizationTestCases() []deleteOrganizationTestCase {
	return []deleteOrganizationTestCase{
		{
			name: "success",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, false, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Organization{
					Name:     "test organization",
					IsActive: true,
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")
			},
		},
	}
}

func Test_DeleteOrganization(t *testing.T) {
	for _, testCase := range getDeleteOrganizationTestCases() {
		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		t.Run(testCase.name, func(t *testing.T) {
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			defer cleanup()
			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			testCase.mutations(t, db)
			err = m.DeleteOrganization(testCtx, txHandlerDb, testCase.id)
			testCase.assertions(t, db, testCase.id, err)
		})
	}
}

func Test_UpdateOrganization_Success(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, m, "unexpected nil")

	txHandler, err := mysqltx.New(&mysqltx.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected error creating the tx handler")

	txHandlerDb, err := txHandler.Db(testCtx)
	require.NoError(t, err, "unexpected error fetching the db from the tx handler")
	require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

	mutations := func(t *testing.T, db *sqlx.DB) int {
		entry := mysqlmodel.Organization{
			ID:            1,
			Name:          "TEST",
			CreatedBy:     null.IntFrom(2),
			LastUpdatedBy: null.IntFrom(3),
			CreatedAt:     time.Now(),
			LastUpdatedAt: null.TimeFrom(time.Now()),
			IsActive:      true,
		}
		err := entry.Insert(context.Background(), db, boil.Infer())
		require.NoError(t, err, "error inserting sample data")
		return entry.ID
	}

	orgID := mutations(t, db)

	paginatedOrganizations, err := m.GetOrganizations(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the organizations from the database")
	require.NotNil(t, paginatedOrganizations, "unexpected nil organizations")
	require.True(t, len(paginatedOrganizations.Organizations) > 0, "unexpected empty organizations")

	updateOrganization := model.UpdateOrganization{
		Id: orgID,
		Name: null.String{
			String: paginatedOrganizations.Organizations[0].Name + " new",
			Valid:  true,
		},
	}

	cat, err := m.UpdateOrganization(testCtx, txHandlerDb, &updateOrganization)
	require.NoError(t, err, "unexpected error updating the organization")
	assert.Equal(t, paginatedOrganizations.Organizations[0].Name+" new", cat.Name)
}

func Test_UpdateOrganization_Fail(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, m, "unexpected nil")

	txHandler, err := mysqltx.New(&mysqltx.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected error creating the tx handler")

	txHandlerDb, err := txHandler.Db(testCtx)
	require.NoError(t, err, "unexpected error fetching the db from the tx handler")
	require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

	updateOrganization := model.UpdateOrganization{
		Id: 3123123123,
		Name: null.String{
			String: "Non-existent Organization",
			Valid:  true,
		},
	}

	_, err = m.UpdateOrganization(testCtx, txHandlerDb, &updateOrganization)
	require.Error(t, err, "expected error when updating a non-existent organization")
	assert.Contains(t, err.Error(), "expected exactly one entry for entity: 3123123123")
}

type createOrganizationTestCase struct {
	name             string
	organizationName string
	createdBy        null.Int
	assertions       func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error)
}

func getAddOrganizationTestCases() []createOrganizationTestCase {
	return []createOrganizationTestCase{
		{
			name:             "success",
			organizationName: "Example Organization",
			createdBy:        null.IntFrom(3),
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.NotNil(t, organization, "unexpected nil organization")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
		{
			name:             "fail-name-exceeds-limit",
			organizationName: strings.Repeat("a", 256),
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.Nil(t, organization, "unexpected non-nil organization")
				assert.Error(t, err, "unexpected nil-error")
			},
		},
	}
}

func Test_AddOrganization(t *testing.T) {
	for _, testCase := range getAddOrganizationTestCases() {
		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		t.Run(testCase.name, func(t *testing.T) {
			defer cleanup()
			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			org := &model.Organization{
				Name:      testCase.organizationName,
				CreatedBy: testCase.createdBy,
			}

			org, err = m.AddOrganization(testCtx, txHandlerDb, org)
			testCase.assertions(t, db, org, err)
		})
	}
}

type restoreOrganizationTestCase struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB)
}

func getRestoreOrganizationTestCases() []restoreOrganizationTestCase {
	return []restoreOrganizationTestCase{
		{
			name: "success",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Organization{
					Name: "test",
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				entry.IsActive = false
				_, err = entry.Update(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected update error")

				err = entry.Reload(context.TODO(), db)
				assert.NoError(t, err, "unexpected reload error")

				assert.Equal(t, false, entry.IsActive)
			},
		},
		{
			name: "fail-missing-entry-to-update",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "unexpected non-nil error")
				assert.Contains(t, err.Error(), "restore:")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, "organization")
			},
		},
	}
}

func Test_RestoreOrganization(t *testing.T) {
	for _, testCase := range getRestoreOrganizationTestCases() {
		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		t.Run(testCase.name, func(t *testing.T) {
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			defer cleanup()
			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			testCase.mutations(t, db)
			err = m.RestoreOrganization(testCtx, txHandlerDb, testCase.id)
			testCase.assertions(t, db, testCase.id, err)
		})
	}
}
