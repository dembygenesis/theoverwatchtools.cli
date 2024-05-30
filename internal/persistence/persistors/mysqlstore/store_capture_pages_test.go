package mysqlstore

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
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
			paginated, err := m.GetCapturePages(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}
