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
	"testing"
)

type testCaseGetCapturePages struct {
	name       string
	filter     *model.CapturePageFilters
	mutations  func(t *testing.T, db *sqlx.DB) []int
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, ids []int)
}

func getTestCasesGetCapturePages() []testCaseGetCapturePages {
	return []testCaseGetCapturePages{
		{
			name: "success-filter-ids-in",
			filter: &model.CapturePageFilters{
				IdsIn: []int{4},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil CapturePages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.CapturePages), "unexpected number of capture pages returned")

				for _, capP := range paginated.CapturePages {
					assert.Contains(t, ids, capP.Id, "returned capture page ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {
				entryOrganization := mysqlmodel.Organization{
					Name: "Test Organization",
				}
				err := entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization in the db")

				capturePageSet := mysqlmodel.CapturePageSet{
					ID: 1,
				}
				err = capturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set")

				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
					OrganizationRefID: null.IntFrom(entryOrganization.ID),
				}
				err = entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry := mysqlmodel.CapturePage{
					ID:               4,
					Name:             "TEST",
					Clicks:           4,
					CreatedBy:        null.IntFrom(entryUser.ID),
					LastUpdatedBy:    null.IntFrom(entryUser.ID),
					CapturePageSetID: capturePageSet.ID,
					IsActive:         true,
				}
				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{entryUser.ID}

				return Ids
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.CapturePageFilters{
				CapturePageNameIn: []string{"Lawrence", "Younes"},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.CapturePages), "unexpected number of capture pages returned")

				for _, capP := range paginated.CapturePages {
					assert.Contains(t, ids, capP.Id, "returned capture page ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {

				entryOrganization := mysqlmodel.Organization{
					Name: "Test Organization",
				}
				err := entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization in the db")

				capturePageSet := mysqlmodel.CapturePageSet{
					Name: "lawrence",
					ID:   1,
				}
				err = capturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set")

				capturePageSet2 := mysqlmodel.CapturePageSet{
					Name: "younes",
					ID:   2,
				}
				err = capturePageSet2.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set")

				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
					OrganizationRefID: null.IntFrom(entryOrganization.ID),
				}
				err = entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entryCap1 := mysqlmodel.CapturePage{
					ID:               1,
					Name:             "Lawrence",
					Clicks:           4,
					CreatedBy:        null.IntFrom(entryUser.ID),
					LastUpdatedBy:    null.IntFrom(entryUser.ID),
					CapturePageSetID: capturePageSet.ID,
					IsActive:         true,
				}
				err = entryCap1.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				entryCap2 := mysqlmodel.CapturePage{
					ID:               2,
					Name:             "Younes",
					Clicks:           4,
					CreatedBy:        null.IntFrom(entryUser.ID),
					LastUpdatedBy:    null.IntFrom(entryUser.ID),
					CapturePageSetID: capturePageSet2.ID,
					IsActive:         true,
				}
				err = entryCap2.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{entryCap1.ID, entryCap2.ID}

				return Ids
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.CapturePageFilters{
				CapturePageNameIn: []string{"Demby"},
				IdsIn:             []int{1},
				CreatedBy:         null.IntFrom(4),
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {
				entryOrganization := mysqlmodel.Organization{
					Name: "Test Organization",
				}
				err := entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization in the db")

				capturePageSet := mysqlmodel.CapturePageSet{
					ID: 1,
				}
				err = capturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set")

				entryUser := mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
					OrganizationRefID: null.IntFrom(entryOrganization.ID),
				}
				err = entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry := mysqlmodel.CapturePage{
					ID:               1,
					Name:             "Demby",
					Clicks:           4,
					CreatedBy:        null.IntFrom(entryUser.ID),
					LastUpdatedBy:    null.IntFrom(entryUser.ID),
					CapturePageSetID: capturePageSet.ID,
					IsActive:         true,
				}
				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{entry.ID}

				return Ids
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) == 1, "unexpected greater than 1 capture pages")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
		},
		{
			name: "empty-results",
			filter: &model.CapturePageFilters{
				CapturePageNameIn:   []string{"Saul Goodman"},
				IdsIn:               []int{1},
				CapturePageIsActive: null.BoolFrom(false),
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {
				return nil
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) == 0, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
		},
	}
}

func Test_GetCapturePages(t *testing.T) {
	for _, testCase := range getTestCasesGetCapturePages() {
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
			paginated, err := m.GetCapturePages(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err, id)
		})
	}
}
