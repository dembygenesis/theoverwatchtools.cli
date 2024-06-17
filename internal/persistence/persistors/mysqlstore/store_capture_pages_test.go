package mysqlstore

import (
	"context"
	"errors"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistencefakes"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore/testhelper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"testing"
)

type testCaseGetCapturePages struct {
	name       string
	filter     *model.CapturePagesFilters
	mutations  func(t *testing.T, db *sqlx.DB)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error)
}

func getTestCasesGetCapturePages() []testCaseGetCapturePages {
	return []testCaseGetCapturePages{
		{
			name: "success-filter-is-control",
			filter: &model.CapturePagesFilters{
				CapturePagesIsControl: null.Bool{
					Bool:  true,
					Valid: true,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture page")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) == 1, "unexpected greater than 1 capture page")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
		},
		{
			name: "fail-ctx",
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error) {
				assert.Error(t, err, "test error")
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.CapturePagesFilters{
				CapturePagesNameIn: []string{"Example Page 1", "Example Page 2"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture page")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) == 2, "unexpected greater than 1 capture page")
				assert.True(t, paginated.Pagination.RowCount == 2, "unexpected count to be greater than 1 capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
		},
		{
			name: "success-capture-page-set-id",
			filter: &model.CapturePagesFilters{
				CapturePagesTypeIdIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture page")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) > 1, "unexpected greater than 1 capture page")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
		},
		{
			name: "success-capture-page-type-name-in",
			filter: &model.CapturePagesFilters{
				CapturePagesTypeNameIn: []string{"Test Page Set 1"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture page")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) > 1, "unexpected greater than 1 capture page")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.CapturePagesFilters{
				CapturePagesTypeNameIn: []string{"Test Page Set 1"},
				CapturePagesTypeIdIn:   []int{1},
				CapturePagesNameIn:     []string{"Example Page 1"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture page")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) == 1, "unexpected greater than 1 capture page")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
		},
		{
			name: "empty-results",
			filter: &model.CapturePagesFilters{
				CapturePagesTypeNameIn: []string{"Saul Goodman"},
				CapturePagesTypeIdIn:   []int{13},
				CapturePagesNameIn:     []string{"Super Admin"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) == 0, "unexpected greater than 1 capture pages")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 capture pages")

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

			defer cleanup()
			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			var txHandler persistence.TransactionHandler

			if testCase.name == "fail-ctx" {
				fakeTxHandler := &persistencefakes.FakeTransactionHandler{}
				fakeTxHandler.GetCtxExecutorReturns(nil, errors.New("mock error"))
				txHandler = fakeTxHandler
				testCase.filter = nil
			} else {
				txHandlerInstance, err := mysqltx.New(&mysqltx.Config{
					Logger:       testLogger,
					Db:           db,
					DatabaseName: cp.Database,
				})
				require.NoError(t, err, "unexpected error creating the tx handler")

				txHandlerDb, err := txHandlerInstance.Db(testCtx)
				txHandler = txHandlerDb
				require.NoError(t, err, "unexpected error fetching the db from the tx handler")
				require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")
			}

			testCase.mutations(t, db)
			paginated, err := m.GetCapturePages(testCtx, txHandler, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}

//test for the post method

type createCapturePagesTestCase struct {
	name             string
	capturePageName  string
	capturePageSetId int
	assertions       func(t *testing.T, db *sqlx.DB, capturepages *model.CapturePages, err error)
}

func getAddCapturePageTestCases() []createCapturePagesTestCase {
	return []createCapturePagesTestCase{
		{
			name:             "success",
			capturePageName:  "Example Capture Page",
			capturePageSetId: 1,
			assertions: func(t *testing.T, db *sqlx.DB, capturepages *model.CapturePages, err error) {
				assert.NotNil(t, capturepages, "unexpected nil capture pages")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePages{*capturepages})
			},
		},
		{
			name:             "fail-name-exceeds-limit",
			capturePageName:  strings.Repeat("a", 256),
			capturePageSetId: 1,
			assertions: func(t *testing.T, db *sqlx.DB, capturepages *model.CapturePages, err error) {
				assert.Nil(t, capturepages, "unexpected non-nil capture pages")
				assert.Error(t, err, "unexpected nil-error")
			},
		},
		{
			name:             "fail-invalid-capture-pages-set-id",
			capturePageName:  strings.Repeat("a", 255),
			capturePageSetId: 199,
			assertions: func(t *testing.T, db *sqlx.DB, capturepages *model.CapturePages, err error) {
				assert.Nil(t, capturepages, "unexpected non-nil capture pages")
				assert.Error(t, err, "unexpected nil-error")
			},
		},
	}
}

func Test_AddCapturePage(t *testing.T) {
	for _, testCase := range getAddCapturePageTestCases() {
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

			capt := &model.CapturePages{
				Name:             testCase.capturePageName,
				CapturePageSetId: testCase.capturePageSetId,
			}

			capt, err = m.AddCapturePage(testCtx, txHandlerDb, capt)
			testCase.assertions(t, db, capt, err)
		})
	}
}

//test for the update method

func Test_UpdateCapturePages_Success(t *testing.T) {
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

	paginatedCapturePages, err := m.GetCapturePages(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the capture pages from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil capture pages")
	require.True(t, len(paginatedCapturePages.CapturePages) > 0, "unexpected empty capture pages")

	updateCapturePages := model.UpdateCapturePages{
		Id: 1,
		CapturePageSetId: null.Int{
			Int:   paginatedCapturePages.CapturePages[0].CapturePageSetId,
			Valid: true,
		},
		Name: null.String{
			String: paginatedCapturePages.CapturePages[0].Name + " new",
			Valid:  true,
		},
	}

	capt, err := m.UpdateCapturePages(testCtx, txHandlerDb, &updateCapturePages)
	require.NoError(t, err, "unexpected error updating a conflicting capture pages from the database")
	assert.Equal(t, paginatedCapturePages.CapturePages[0].Name+" new", capt.Name)
}

func Test_UpdateCapturePages_Fail(t *testing.T) {
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

	paginatedCapturePages, err := m.GetCapturePages(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the capture pages from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil capture pages")
	require.True(t, len(paginatedCapturePages.CapturePages) > 0, "unexpected empty capture pages")

	updateCapturePages := model.UpdateCapturePages{
		Id: paginatedCapturePages.CapturePages[1].Id,
		CapturePageSetId: null.Int{
			Int:   paginatedCapturePages.CapturePages[0].CapturePageSetId,
			Valid: true,
		},
		Name: null.String{
			String: paginatedCapturePages.CapturePages[0].Name,
			Valid:  true,
		},
	}

	capt, err := m.UpdateCapturePages(testCtx, txHandlerDb, &updateCapturePages)
	require.Error(t, err, "unexpected nil error fetching a conflicting capture pages from the database")
	assert.Contains(t, err.Error(), "Duplicate entry")
	assert.Nil(t, capt, "unexpected non nil entry")
}

//test for the delete method

type deleteCapturePagesTestCase struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB)
}

func getDeleteCapturePagesTestCases() []deleteCapturePagesTestCase {
	return []deleteCapturePagesTestCase{
		{
			name: "success",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the capture page")

				assert.Equal(t, 0, entry.IsControl)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.CapturePage{
					CapturePageSetID: 1,
					Name:             "test",
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")
			},
		},
	}
}

func Test_DeleteCapturePage(t *testing.T) {
	for _, testCase := range getDeleteCapturePagesTestCases() {
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
			err = m.DeleteCapturePages(testCtx, txHandlerDb, testCase.id)
			testCase.assertions(t, db, testCase.id, err)
		})
	}
}

//test for the restore method

type restoreCapturePagesTestCase struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB)
}

func getRestoreCapturePagesTestCases() []restoreCapturePagesTestCase {
	return []restoreCapturePagesTestCase{
		{
			name: "success",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the capture pages")

				assert.Equal(t, 1, entry.IsControl)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.CapturePage{
					CapturePageSetID: 1,
					Name:             "test",
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				entry.IsControl = 0
				_, err = entry.Update(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected update error")

				err = entry.Reload(context.TODO(), db)
				assert.NoError(t, err, "unexpected reload error")

				assert.Equal(t, 0, entry.IsControl)
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
				testhelper.DropTable(t, db, "capture_pages")
			},
		},
	}
}

func Test_RestoreCapturePages(t *testing.T) {
	for _, testCase := range getRestoreCapturePagesTestCases() {
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
			err = m.RestoreCapturePages(testCtx, txHandlerDb, testCase.id)
			testCase.assertions(t, db, testCase.id, err)
		})
	}
}
