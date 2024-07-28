package mysqlstore

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore/testhelper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type testCaseGetCategoryType struct {
	name       string
	filter     *model.CategoryTypeFilters
	mutations  func(t *testing.T, db *sqlx.DB)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategoryTypes, err error)
}

func getTestCasesGetCategoryTypes() []testCaseGetCategoryType {
	return []testCaseGetCategoryType{
		{
			name: "success-filtered-ids-in",
			filter: &model.CategoryTypeFilters{
				IdsIn: []int{1, 2},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategoryTypes, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				assert.True(t, len(paginated.Categories) == 2, "unexpected categories length")
			},
		},
		{
			name: "fail-internal-error",
			filter: &model.CategoryTypeFilters{
				IdsIn: []int{1, 2},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.CategoryType)
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategoryTypes, err error) {
				require.Error(t, err, "unexpected error")
				require.Nil(t, paginated, "unexpected nil paginated")
			},
		},
	}
}

func TestGetCategoryTypes(t *testing.T) {
	for _, testCase := range getTestCasesGetCategoryTypes() {
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
			paginated, err := m.GetCategoryTypes(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}
