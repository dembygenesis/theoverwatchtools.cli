package mysqlstore

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
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
	mutations  func(t *testing.T, db *sqlx.DB, organization *model.Organization)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error)
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) {
				entry := mysqlmodel.Organization{
					ID:            organization.Id,
					Name:          "TEST",
					CreatedBy:     organization.CreatedBy,
					LastUpdatedBy: organization.LastUpdatedBy,
					CreatedAt:     organization.CreatedAt,
					LastUpdatedAt: organization.LastUpdatedAt,
					IsActive:      organization.IsActive,
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
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) {

				entry := []mysqlmodel.Organization{
					{
						ID:            organization.Id,
						Name:          "Organization A",
						CreatedBy:     organization.CreatedBy,
						LastUpdatedBy: organization.LastUpdatedBy,
						CreatedAt:     organization.CreatedAt,
						LastUpdatedAt: organization.LastUpdatedAt,
						IsActive:      true,
					},
					{
						ID:            organization.Id + 1,
						Name:          "Organization B",
						CreatedBy:     organization.CreatedBy,
						LastUpdatedBy: organization.LastUpdatedBy,
						CreatedAt:     organization.CreatedAt,
						LastUpdatedAt: organization.LastUpdatedAt,
						IsActive:      true,
					},
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
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) {
				entry := mysqlmodel.Organization{
					ID:            organization.Id,
					Name:          "Organization A",
					CreatedBy:     organization.CreatedBy,
					LastUpdatedBy: organization.LastUpdatedBy,
					CreatedAt:     organization.CreatedAt,
					LastUpdatedAt: organization.LastUpdatedAt,
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
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) {
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

			organization := &model.Organization{
				Id:            1,
				CreatedBy:     null.IntFrom(1),
				LastUpdatedBy: null.IntFrom(1),
				IsActive:      true,
			}

			testCase.mutations(t, db, organization)
			paginated, err := m.GetOrganizations(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}

type testCaseDeleteOrganization struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int)
}

func getDeleteOrganizationTestCases() []testCaseDeleteOrganization {
	return []testCaseDeleteOrganization{
		{
			name: "success",
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, false, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int) {
				entry := mysqlmodel.Organization{
					ID:        organization.Id,
					Name:      organization.Name,
					CreatedBy: organization.CreatedBy,
					IsActive:  organization.IsActive,
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")
				return organization.Id
			},
		},
	}
}

func Test_DeleteOrganization(t *testing.T) {
	for _, testCase := range getDeleteOrganizationTestCases() {
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

			organization := &model.Organization{
				Id:        1,
				Name:      "Organization A",
				CreatedBy: null.IntFrom(1),
				IsActive:  true,
			}

			id := testCase.mutations(t, db, organization)

			err = m.DeleteOrganization(testCtx, txHandlerDb, id)
			testCase.assertions(t, db, id, err)
		})
	}
}

type testCaseUpdateOrganization struct {
	name                string
	initialOrganization mysqlmodel.Organization
	updateOrganization  model.UpdateOrganization
	assertions          func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations           func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization)
}

func getUpdateOrganizationTestCases() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name: "success",
			initialOrganization: mysqlmodel.Organization{
				ID:            1,
				Name:          "Organization A",
				CreatedBy:     null.IntFrom(1),
				LastUpdatedBy: null.IntFrom(1),
				CreatedAt:     time.Now(),
				LastUpdatedAt: null.TimeFrom(time.Now()),
				IsActive:      true,
			},
			updateOrganization: model.UpdateOrganization{
				Id:   1,
				Name: null.StringFrom("Organization A new"),
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization) {
				updateData = model.UpdateOrganization{
					Id:   organization.Id,
					Name: organization.Name,
				}
				return updateData
			},
		},
		{
			name: "fail",
			initialOrganization: mysqlmodel.Organization{
				ID:            1,
				Name:          "Organization A",
				CreatedBy:     null.IntFrom(1),
				LastUpdatedBy: null.IntFrom(1),
				CreatedAt:     time.Now(),
				LastUpdatedAt: null.TimeFrom(time.Now()),
				IsActive:      true,
			},
			updateOrganization: model.UpdateOrganization{
				Id:   32123,
				Name: null.StringFrom("Wrong Name"),
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent organization")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 32123")
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.UpdateOrganization) (updateData model.UpdateOrganization) {
				updateData = model.UpdateOrganization{
					Id:   organization.Id,
					Name: organization.Name,
				}
				return updateData
			},
		},
	}
}

func Test_UpdateOrganization(t *testing.T) {
	for _, testCase := range getUpdateOrganizationTestCases() {
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

		err = testCase.initialOrganization.Insert(context.Background(), db, boil.Infer())
		require.NoError(t, err, "unexpected error inserting initial organization")

		updateData := testCase.mutations(t, db, &testCase.updateOrganization)

		_, err = m.UpdateOrganization(testCtx, txHandlerDb, &updateData)
		testCase.assertions(t, db, updateData.Id, err)
	}
}

type testCaseCreateOrganization struct {
	name             string
	organizationName string
	createdBy        null.Int
	assertions       func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error)
}

func getAddOrganizationTestCases() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
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
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
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

			org := &model.Organization{
				Name:      testCase.organizationName,
				CreatedBy: testCase.createdBy,
			}

			org, err = m.AddOrganization(testCtx, txHandlerDb, org)
			testCase.assertions(t, db, org, err)
		})
	}
}

type testCaseRestoreOrganization struct {
	name                string
	id                  int
	restoreOrganization model.Organization
	assertions          func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations           func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int)
}

func getRestoreOrganizationTestCases() []testCaseRestoreOrganization {
	return []testCaseRestoreOrganization{
		{
			name: "success",
			restoreOrganization: model.Organization{
				Id:       1,
				Name:     "Organization A",
				IsActive: false,
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")
				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int) {
				entry := mysqlmodel.Organization{
					ID:       organization.Id,
					Name:     organization.Name,
					IsActive: true,
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")
				assert.Equal(t, true, entry.IsActive)
				return organization.Id
			},
		},
		{
			name: "fail-missing-entry-to-update",
			restoreOrganization: model.Organization{
				Id:       99999,
				Name:     "Non-existent Organization",
				IsActive: false,
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				_, fetchErr := mysqlmodel.FindOrganization(context.TODO(), db, id)
				assert.Error(t, fetchErr, "organization should not exist")
			},
			mutations: func(t *testing.T, db *sqlx.DB, organization *model.Organization) (id int) {
				return organization.Id
			},
		},
	}
}

func Test_RestoreOrganization(t *testing.T) {
	for _, testCase := range getRestoreOrganizationTestCases() {
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

			id := testCase.mutations(t, db, &testCase.restoreOrganization)
			err = m.RestoreOrganization(testCtx, txHandlerDb, id)
			testCase.assertions(t, db, id, err)
		})
	}
}
