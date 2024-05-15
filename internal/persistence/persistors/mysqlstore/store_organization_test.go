package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore/testhelper"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"testing"
	"time"
)

var (
	testCtx           = context.TODO()
	testLogger        = logger.New(context.TODO())
	testQueryTimeouts = &persistence.QueryTimeouts{
		Query: 10 * time.Second,
		Exec:  10 * time.Second,
	}
)

type testCaseGetOrganizations struct {
	name       string
	filter     *model.OrganizationFilters
	mutations  func(t *testing.T, db *sqlx.DB)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error)
}

type createOrganizationTestCase struct {
	name              string
	organizationName  string
	organizationRefId int
	assertions        func(t *testing.T, db *sqlx.DB, category *model.Organization, err error)
}

func getAddOrganizationTestCases() []createOrganizationTestCase {
	return []createOrganizationTestCase{
		{
			name:              "success",
			organizationName:  "Example Organization",
			organizationRefId: 1,
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.NotNil(t, organization, "unexpected nil organization")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
		{
			name:              "fail-name-exceeds-limit",
			organizationName:  strings.Repeat("a", 256),
			organizationRefId: 1,
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.Nil(t, organization, "unexpected non-nil organization")
				assert.Error(t, err, "unexpected nil-error")
			},
		},
		{
			name:              "fail-invalid-organization-type-id",
			organizationName:  strings.Repeat("a", 255),
			organizationRefId: 199,
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

			organization := &model.Organization{
				Name:                  testCase.organizationName,
				OrganizationTypeRefId: testCase.organizationRefId,
			}

			organization, err = m.AddOrganization(testCtx, txHandlerDb, organization)
			testCase.assertions(t, db, organization, err)
		})
	}
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {

				fmt.Println("Number of organizations retrieved ---- ", len(paginated.Organizations))

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
			name: "success-filter-names-in",
			filter: &model.OrganizationFilters{
				OrganizationNameIn: []string{"TEST DATA 1", "TEST DATA 2"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 2, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount == 2, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "success-organization-type-id-in",
			filter: &model.OrganizationFilters{
				OrganizationTypeIdIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {

				fmt.Println("Number of organizations retrieved ---- ", strutil.GetAsJson(paginated.Organizations))

				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) > 1, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "success-organization-type-name-in",
			filter: &model.OrganizationFilters{
				OrganizationTypeNameIn: []string{"Organization Type 1"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) > 1, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.OrganizationFilters{
				OrganizationTypeNameIn: []string{"Organization Type 1"},
				OrganizationTypeIdIn:   []int{1},
				OrganizationNameIn:     []string{"TEST DATA 1"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
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
				OrganizationTypeNameIn: []string{"Organization Type 1"},
				OrganizationTypeIdIn:   []int{2},
				OrganizationNameIn:     []string{"TEST DATA 1"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
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
			fmt.Println("the organs filters ---- ", strutil.GetAsJson(testCase.filter))
			paginated, err := m.GetOrganizations(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}

func Test_ReadOrganizations(t *testing.T) {
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

	paginatedOrganizations, err := m.GetOrganizations(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the organizations from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil organizations")
	require.True(t, len(paginatedOrganizations.Organizations) > 0, "unexpected empty organizations")
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

				assert.Equal(t, 0, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Organization{
					OrganizationTypeRefID: 1,
					Name:                  "Organization Type 1",
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

func Test_UpdateOrganizations_Success(t *testing.T) {
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

	paginatedOrganizations, err := m.GetOrganizations(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the organizations from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil organizations")
	require.True(t, len(paginatedOrganizations.Organizations) > 0, "unexpected empty organizations")

	updateOrganization := model.UpdateOrganization{
		Id: 1,
		OrganizationTypeRefId: null.Int{
			Int:   paginatedOrganizations.Organizations[0].OrganizationTypeRefId,
			Valid: true,
		},
		Name: null.String{
			String: paginatedOrganizations.Organizations[0].Name + " new",
			Valid:  true,
		},
	}

	cat, err := m.UpdateOrganization(testCtx, txHandlerDb, &updateOrganization)
	require.NoError(t, err, "unexpected error updating a conflicting organization from the database")
	assert.Equal(t, paginatedOrganizations.Organizations[0].Name+" new", cat.Name)
}

func Test_UpdateOrganizations_Fail(t *testing.T) {
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

	paginatedOrganizations, err := m.GetOrganizations(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the organizations from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil organizations")
	require.True(t, len(paginatedOrganizations.Organizations) > 0, "unexpected empty organizations")

	updateOrganization := model.UpdateOrganization{
		Id: paginatedOrganizations.Organizations[1].Id,
		OrganizationTypeRefId: null.Int{
			Int:   paginatedOrganizations.Organizations[0].OrganizationTypeRefId,
			Valid: true,
		},
		Name: null.String{
			String: paginatedOrganizations.Organizations[0].Name,
			Valid:  true,
		},
	}

	cat, err := m.UpdateOrganization(testCtx, txHandlerDb, &updateOrganization)
	require.Error(t, err, "unexpected nil error fetching a conflicting category from the database")
	assert.Contains(t, err.Error(), "Duplicate entry")
	assert.Nil(t, cat, "unexpected non nil entry")
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

				assert.Equal(t, 1, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Organization{
					OrganizationTypeRefID: 1,
					Name:                  "test",
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				entry.IsActive = 0
				_, err = entry.Update(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected update error")

				err = entry.Reload(context.TODO(), db)
				assert.NoError(t, err, "unexpected reload error")

				assert.Equal(t, 0, entry.IsActive)
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
