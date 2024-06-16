package mysqlstore

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type testCaseGetClickTrackers struct {
	name       string
	filter     *model.ClickTrackerFilters
	mutations  func(t *testing.T, db *sqlx.DB)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error)
}

func getTestCasesGetClickTrackers() []testCaseGetClickTrackers {
	return []testCaseGetClickTrackers{
		{
			name: "success-filter-ids-in",
			filter: &model.ClickTrackerFilters{
				IdsIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) == 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.ClickTrackerFilters{
				NameIn: []string{"Tracker 1", "Tracker 2"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click tracker")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) == 2, "unexpected greater than 1 click tracker")
				assert.True(t, paginated.Pagination.RowCount == 2, "unexpected count to be greater than 1 click tracker")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
		{
			name: "success-category-type-id-in",
			filter: &model.ClickTrackerFilters{
				ClickTrackerSetIdIn: []int{2},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error) {

				fmt.Println("Number of categories retrieved ---- ", strutil.GetAsJson(paginated.ClickTrackers))

				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) > 1, "unexpected greater than 1 click trackers")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 click trackers")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
		{
			name: "success-category-type-name-in",
			filter: &model.ClickTrackerFilters{
				ClickTrackerSetNameIn: []string{"Tracker Set 1"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click tracker")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) > 1, "unexpected greater than 1 click tracker")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.ClickTrackerFilters{
				ClickTrackerSetNameIn: []string{"Tracker Set 1"},
				ClickTrackerSetIdIn:   []int{1},
				NameIn:                []string{"Tracker 1"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) == 1, "unexpected greater than 1 click tracker")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 click tracker")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
		{
			name: "empty-results",
			filter: &model.ClickTrackerFilters{
				ClickTrackerSetNameIn: []string{"Saul Goodman"},
				ClickTrackerSetIdIn:   []int{1},
				NameIn:                []string{"Super Admin"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) == 0, "unexpected greater than 1 click tracker")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 click tracker")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
	}
}

func Test_GetClickTrackers(t *testing.T) {
	for _, testCase := range getTestCasesGetClickTrackers() {
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

			fmt.Println("the click trackers filters --- ", strutil.GetAsJson(testCase.filter))

			testCase.mutations(t, db)
			paginated, err := m.GetClickTrackers(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}
