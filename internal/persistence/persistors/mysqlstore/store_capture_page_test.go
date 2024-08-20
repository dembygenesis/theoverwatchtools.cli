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
				assert.Equal(t, len(ids), len(paginated.CapturePages), "unexpected number of organizations returned")

				for _, cap := range paginated.CapturePages {
					assert.Contains(t, ids, cap.Id, "returned CapturePages ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
			mutations: func(t *testing.T, db *sqlx.DB) []int {
				entryUser := mysqlmodel.CapturePageSet{
					Name:              "Demby",
					OrganizationRefID: null.IntFrom(1),
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				entry := mysqlmodel.CapturePage{
					ID:               4,
					Name:             "TEST",
					Clicks:           4,
					CapturePageSetID: entryUser.OrganizationRefID.Int,
				}
				err = entry.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{entryUser.ID}

				return Ids
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
