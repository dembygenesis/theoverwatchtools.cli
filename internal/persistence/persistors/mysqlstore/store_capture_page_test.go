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
	"time"
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

type testCaseAddCapturePage struct {
	name            string
	capturePageName string
	assertions      func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage, err error)
	mutations       func(t *testing.T, db *sqlx.DB) (capturePage *model.CreateCapturePage)
}

func getAddCapturePageTestCases() []testCaseAddCapturePage {
	return []testCaseAddCapturePage{
		{
			name:            "success",
			capturePageName: "Example Capture Page",
			mutations: func(t *testing.T, db *sqlx.DB) (capturePage *model.CreateCapturePage) {
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
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				entryCapturePageSet := mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(entryOrganization.ID),
					CreatedBy:         null.IntFrom(entryUser.ID),
					LastUpdatedBy:     null.IntFrom(entryUser.ID),
				}
				err = entryCapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set to db")

				createCapturePage := &model.CreateCapturePage{
					Name:             "DembyCapturePage",
					UserId:           1,
					CapturePageSetId: 1,
				}

				return createCapturePage
			},
			assertions: func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage, err error) {
				assert.NotNil(t, capturePage, "unexpected nil capture page")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePage{*capturePage})
			},
		},
	}
}

func Test_AddCapturePage(t *testing.T) {
	for _, testCase := range getAddCapturePageTestCases() {
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

			capturePageData := testCase.mutations(t, db)

			createdCapturePage, err := m.AddCapturePage(testCtx, txHandlerDb, capturePageData)
			testCase.assertions(t, db, createdCapturePage, err)
		})
	}
}

type testCaseRestoreCapturePage struct {
	name               string
	id                 int
	restoreCapturePage model.CapturePage
	assertions         func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations          func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage) (id int)
}

func getRestoreCapturePageTestCases() []testCaseRestoreCapturePage {
	return []testCaseRestoreCapturePage{
		{
			name: "success",
			restoreCapturePage: model.CapturePage{
				Id:               1,
				Name:             "CapturePage A",
				CapturePageSetId: 1,
				IsActive:         false,
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the capture page")
				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage) (id int) {
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
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				entryCapturePageSet := mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(entryOrganization.ID),
					CreatedBy:         null.IntFrom(entryUser.ID),
					LastUpdatedBy:     null.IntFrom(entryUser.ID),
				}
				err = entryCapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set to db")

				entry := mysqlmodel.CapturePage{
					Name:             capturePage.Name,
					CreatedBy:        entryUser.CreatedBy,
					LastUpdatedBy:    entryUser.CreatedBy,
					CapturePageSetID: capturePage.CapturePageSetId,
					IsActive:         false,
				}
				err = entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				return entry.ID
			},
		},
	}
}

func Test_RestoreCapturePage(t *testing.T) {
	for _, testCase := range getRestoreCapturePageTestCases() {
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

			id := testCase.mutations(t, db, &testCase.restoreCapturePage)
			err = m.RestoreCapturePage(testCtx, txHandlerDb, id)
			testCase.assertions(t, db, id, err)
		})
	}
}

type testCaseDeleteCapturePage struct {
	name              string
	deleteCapturePage model.CapturePage
	assertions        func(t *testing.T, db *sqlx.DB, id int)
	mutations         func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage) (capId int)
}

func getDeleteCapturePageTestCases() []testCaseDeleteCapturePage {
	return []testCaseDeleteCapturePage{
		{
			name: "success",
			deleteCapturePage: model.CapturePage{
				Id:               1,
				Name:             "CapturePage A",
				CapturePageSetId: 1,
				IsActive:         true,
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.Nil(t, err, "unexpected non-nil error")
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, false, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage) (capId int) {
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
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				entryCapturePageSet := mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(entryOrganization.ID),
					CreatedBy:         null.IntFrom(entryUser.ID),
					LastUpdatedBy:     null.IntFrom(entryUser.ID),
				}
				err = entryCapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set to db")

				entry := mysqlmodel.CapturePage{
					Name:             capturePage.Name,
					CreatedBy:        entryUser.CreatedBy,
					LastUpdatedBy:    entryUser.CreatedBy,
					CapturePageSetID: capturePage.CapturePageSetId,
					IsActive:         capturePage.IsActive,
				}
				err = entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				return entry.ID
			},
		},
	}
}

func Test_DeleteCapturePage(t *testing.T) {
	for _, testCase := range getDeleteCapturePageTestCases() {
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

			id := testCase.mutations(t, db, &testCase.deleteCapturePage)
			err = m.DeleteCapturePage(testCtx, txHandlerDb, id)
			require.NoError(t, err, "error deleting organization")
			testCase.assertions(t, db, id)
		})
	}
}

type testCaseUpdateCapturePage struct {
	name              string
	updateCapturePage model.CapturePage
	assertions        func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations         func(t *testing.T, db *sqlx.DB) (updateData model.UpdateCapturePage)
}

func getUpdateCapturePageTestCases() []testCaseUpdateCapturePage {
	return []testCaseUpdateCapturePage{
		{
			name: "success",
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) (updateData model.UpdateCapturePage) {
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
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(entryUser.ID),
					LastUpdatedBy: null.IntFrom(entryUser.ID),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				}
				err = entryOrganization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				entryCapturePageSet := mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(entryOrganization.ID),
					CreatedBy:         null.IntFrom(entryUser.ID),
					LastUpdatedBy:     null.IntFrom(entryUser.ID),
				}
				err = entryCapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set to db")

				entryCapturePage := mysqlmodel.CapturePage{
					ID:               4,
					Name:             "Lawrence Michael",
					CapturePageSetID: 1,
					LastUpdatedAt:    null.TimeFrom(time.Now()),
					CreatedBy:        null.IntFrom(entryUser.ID),
					LastUpdatedBy:    null.IntFrom(entryUser.ID),
				}
				err = entryCapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page into db")

				updateCapturePageData := model.UpdateCapturePage{
					Id:               4,
					Name:             null.StringFrom("Virtue"),
					UserId:           null.IntFrom(1),
					CapturePageSetId: 1,
				}

				return updateCapturePageData
			},
		},
		{
			name: "fail",
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent organization")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 32123")
			},
			mutations: func(t *testing.T, db *sqlx.DB) (updateData model.UpdateCapturePage) {
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

				updateData = model.UpdateCapturePage{
					Id:   32123,
					Name: null.StringFrom("Wrong Name"),
				}
				return updateData
			},
		},
	}
}

func Test_UpdateCapturePage(t *testing.T) {
	for _, testCase := range getUpdateCapturePageTestCases() {
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

		_, err = m.UpdateCapturePage(testCtx, txHandlerDb, &updateData)
		testCase.assertions(t, db, updateData.Id, err)
	}
}
