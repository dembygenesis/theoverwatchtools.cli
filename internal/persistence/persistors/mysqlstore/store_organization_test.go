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
	mutations  func(t *testing.T, db *sqlx.DB) []int
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int)
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{4},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.Organizations), "unexpected number of organizations returned")

				for _, org := range paginated.Organizations {
					assert.Contains(t, ids, org.Id, "returned organization ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
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

				entry := mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{entryUser.ID}

				return Ids
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.OrganizationFilters{
				OrganizationNameIn: []string{"Lawrence", "Younes"},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.Organizations), "unexpected number of organizations returned")

				for _, org := range paginated.Organizations {
					assert.Contains(t, ids, org.Id, "returned organization ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
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

				entryOrg1 := mysqlmodel.Organization{
					ID:            1,
					Name:          "Lawrence",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}

				entryOrg2 := mysqlmodel.Organization{
					ID:            2,
					Name:          "Younes",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}

				err = entryOrg1.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting first organization")

				err = entryOrg2.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting second organization")

				Ids := []int{entryOrg1.ID, entryOrg2.ID}

				return Ids
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.OrganizationFilters{
				OrganizationNameIn: []string{"Demby"},
				IdsIn:              []int{1},
				CreatedBy:          null.IntFrom(4),
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

				entry := mysqlmodel.Organization{
					ID:            1,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}

				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting first organization")

				Ids := []int{entry.ID}

				return Ids
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 1, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "empty-results",
			filter: &model.OrganizationFilters{
				OrganizationNameIn:   []string{"Saul Goodman"},
				IdsIn:                []int{1},
				OrganizationIsActive: null.BoolFrom(false),
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {
				return nil
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 0, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 category")

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

			id := testCase.mutations(t, db)
			paginated, err := m.GetOrganizations(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err, id)
		})
	}
}

type testCaseDeleteOrganization struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int)
	mutations  func(t *testing.T, db *sqlx.DB) (orgId int)
}

func getDeleteOrganizationTestCases() []testCaseDeleteOrganization {
	return []testCaseDeleteOrganization{
		{
			name: "success",
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Nil(t, err, "unexpected non-nil error")
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, false, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) (orgId int) {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry := mysqlmodel.Organization{
					ID:            1,
					Name:          "Lawrence",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				return entry.ID
			},
		},
		{
			name: "failure-non-existent-organization",
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				_, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Error(t, err, "expected an error for non-existent organization")
				require.Contains(t, err.Error(), "no rows", "error should indicate that the organization was not found")
			},
			mutations: func(t *testing.T, db *sqlx.DB) (orgId int) {
				entryUser := mysqlmodel.User{
					ID:                931,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				return entryUser.ID
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

			id := testCase.mutations(t, db)
			err = m.DeleteOrganization(testCtx, txHandlerDb, id)
			require.NoError(t, err, "error deleting organization")
			testCase.assertions(t, db, id)
		})
	}
}

type testCaseUpdateOrganization struct {
	name       string
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB) (updateData model.UpdateOrganization)
}

func getUpdateOrganizationTestCases() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name: "success",
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) (updateData model.UpdateOrganization) {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry := mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				updateOrganizationData := model.UpdateOrganization{
					Id:   4,
					Name: null.StringFrom("Organization A new"),
				}

				return updateOrganizationData
			},
		},
		{
			name: "fail",
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent organization")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 32123")
			},
			mutations: func(t *testing.T, db *sqlx.DB) (updateData model.UpdateOrganization) {
				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry := mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				updateData = model.UpdateOrganization{
					Id:   32123,
					Name: null.StringFrom("Wrong Name"),
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

		updateData := testCase.mutations(t, db)

		_, err = m.UpdateOrganization(testCtx, txHandlerDb, &updateData)
		testCase.assertions(t, db, updateData.Id, err)
	}
}

type testCaseCreateOrganization struct {
	name             string
	organizationName string
	createdBy        null.Int
	assertions       func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error)
	mutations        func(t *testing.T, db *sqlx.DB) (organization *model.CreateOrganization)
}

func getAddOrganizationTestCases() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name:             "success",
			organizationName: "Example Organization",
			createdBy:        null.IntFrom(3),
			mutations: func(t *testing.T, db *sqlx.DB) (organization *model.CreateOrganization) {
				createOrganizationData := model.CreateOrganization{
					Name:   "Demby",
					UserId: 1,
				}

				return &createOrganizationData
			},
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.NotNil(t, organization, "unexpected nil organization")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
		{
			name:             "fail-name-too-long",
			organizationName: strings.Repeat("Demby", 256),
			mutations: func(t *testing.T, db *sqlx.DB) (organization *model.CreateOrganization) {
				entry := mysqlmodel.User{
					Firstname:         "DembyYounesLawrence",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry1 := mysqlmodel.Organization{
					ID:            4,
					Name:          "DembyYounesLawrence",
					CreatedBy:     null.IntFrom(entry.ID),
					LastUpdatedBy: null.IntFrom(entry.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entry1.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				createOrganizationData := model.CreateOrganization{
					Name:   "DembyYounesLawrence",
					UserId: 1,
				}

				return &createOrganizationData
			},
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.Nil(t, organization, "unexpected non-nil organization")
				assert.Error(t, err, "expected an error due to name exceeding length limit")
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

			organizationData := testCase.mutations(t, db)

			createOrganization, err := m.AddOrganization(testCtx, txHandlerDb, organizationData)
			testCase.assertions(t, db, createOrganization, err)
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
