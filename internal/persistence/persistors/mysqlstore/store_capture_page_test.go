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

type argsMutationsGetCapturePage struct {
	User           mysqlmodel.User
	Organization   mysqlmodel.Organization
	CapturePageSet mysqlmodel.CapturePageSet
	CapturePages   []mysqlmodel.CapturePage
}

type testCaseGetCapturePages struct {
	name          string
	argsMutations *argsMutationsGetCapturePage
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsGetCapturePage) *model.CapturePageFilters
	assertions    func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, filters *model.CapturePageFilters)
}

func getTestCasesGetCapturePages() []testCaseGetCapturePages {
	return []testCaseGetCapturePages{
		{
			name: "success-filter-ids-in",
			argsMutations: &argsMutationsGetCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePages: []mysqlmodel.CapturePage{
					{
						ID:               1,
						Name:             "CapturePage A",
						CapturePageSetID: 1,
						CreatedBy:        null.IntFrom(1),
						LastUpdatedBy:    null.IntFrom(1),
						IsActive:         true,
					},
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, filters *model.CapturePageFilters) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil CapturePages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(filters.IdsIn), len(paginated.CapturePages), "unexpected number of capture pages returned")

				for _, capturePage := range paginated.CapturePages {
					assert.Contains(t, filters.IdsIn, capturePage.Id, "returned capture page ID not found in expected IDs")
				}

				assert.Equal(t, len(filters.IdsIn), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetCapturePage) *model.CapturePageFilters {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				var ids []int
				for _, capturePage := range args.CapturePages {
					err = capturePage.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting capture page into the db")
					ids = append(ids, capturePage.ID)
				}

				filters := &model.CapturePageFilters{
					IdsIn: ids,
				}

				return filters
			},
		},
		{
			name: "success-filter-names-in",
			argsMutations: &argsMutationsGetCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePages: []mysqlmodel.CapturePage{
					{
						Name:             "CapturePage A",
						CapturePageSetID: 1,
						CreatedBy:        null.IntFrom(1),
						LastUpdatedBy:    null.IntFrom(1),
						IsActive:         true,
					},
					{
						Name:             "CapturePage B",
						CapturePageSetID: 1,
						CreatedBy:        null.IntFrom(1),
						LastUpdatedBy:    null.IntFrom(1),
						IsActive:         true,
					},
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, filters *model.CapturePageFilters) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")

				for index, capturePage := range paginated.CapturePages {
					assert.Equal(t, filters.CapturePageNameIn[index], capturePage.Name, "unexpected number of capture pages returned")
				}

				for index, capturePage := range paginated.CapturePages {
					assert.Contains(t, filters.CapturePageNameIn[index], capturePage.Name, "returned capture page name not found in expected names")
				}

				modelhelpers.AssertNonEmptyCapturePages(t, paginated.CapturePages)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetCapturePage) *model.CapturePageFilters {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				var names []string
				for _, page := range args.CapturePages {
					err = page.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting capture page into the db")
					names = append(names, page.Name)
				}

				filters := &model.CapturePageFilters{
					CapturePageNameIn: names,
				}

				return filters
			},
		},
		{
			name: "success-multiple-filters",
			argsMutations: &argsMutationsGetCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePages: []mysqlmodel.CapturePage{
					{
						Name:             "CapturePage A",
						CapturePageSetID: 1,
						CreatedBy:        null.IntFrom(4),
						LastUpdatedBy:    null.IntFrom(4),
						IsActive:         true,
					},
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetCapturePage) *model.CapturePageFilters {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				var ids []int
				var names []string
				for _, page := range args.CapturePages {
					err = page.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting capture page into the db")
					ids = append(ids, page.ID)
					names = append(names, page.Name)
				}

				filters := &model.CapturePageFilters{
					CapturePageNameIn: names,
					IdsIn:             ids,
					CreatedBy:         null.IntFrom(args.User.ID),
				}

				return filters
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, filters *model.CapturePageFilters) {
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
			argsMutations: &argsMutationsGetCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePages: []mysqlmodel.CapturePage{
					{
						Name:             "CapturePage A",
						CapturePageSetID: 1,
						CreatedBy:        null.IntFrom(1),
						LastUpdatedBy:    null.IntFrom(1),
						IsActive:         true,
					},
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetCapturePage) *model.CapturePageFilters {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				var ids []int
				var names []string
				for _, page := range args.CapturePages {
					err = page.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting capture page into the db")
					ids = append(ids, page.ID+page.ID*page.ID*page.ID)
					names = append(names, page.Name+page.Name)
				}

				filters := &model.CapturePageFilters{
					CapturePageNameIn: names,
					IdsIn:             ids,
					CreatedBy:         null.IntFrom(args.User.ID + args.User.ID*args.User.ID),
				}

				return filters
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCapturePages, err error, filters *model.CapturePageFilters) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.CapturePages) == 0, "unexpected greater than 1 capture page")
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

			filters := testCase.mutations(t, db, testCase.argsMutations)
			paginated, err := m.GetCapturePages(testCtx, txHandlerDb, filters)

			testCase.assertions(t, db, paginated, err, filters)
		})
	}
}

type argsMutationsAddCapturePage struct {
	User           mysqlmodel.User
	Organization   mysqlmodel.Organization
	CapturePageSet mysqlmodel.CapturePageSet
	CapturePage    model.CreateCapturePage
}

type testCaseAddCapturePage struct {
	name          string
	argsMutations *argsMutationsAddCapturePage
	assertions    func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage, err error)
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsAddCapturePage) *model.CreateCapturePage
}

func getAddCapturePageTestCases() []testCaseAddCapturePage {
	return []testCaseAddCapturePage{
		{
			name: "success",
			argsMutations: &argsMutationsAddCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePage: model.CreateCapturePage{
					Name:             "DembyCapturePage",
					UserId:           1,
					CapturePageSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsAddCapturePage) *model.CreateCapturePage {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				return &args.CapturePage
			},
			assertions: func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage, err error) {
				assert.NotNil(t, capturePage, "unexpected nil capture page")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePage{*capturePage})
			},
		},
		{
			name: "fail-name-too-long",
			argsMutations: &argsMutationsAddCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePage: model.CreateCapturePage{
					Name:   "DembyYounesLawrence",
					UserId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsAddCapturePage) *model.CreateCapturePage {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				return &args.CapturePage
			},
			assertions: func(t *testing.T, db *sqlx.DB, capturePage *model.CapturePage, err error) {
				assert.Nil(t, capturePage, "unexpected non-nil organization")
				assert.Error(t, err, "expected an error due to name exceeding length limit")
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

			capturePageData := testCase.mutations(t, db, testCase.argsMutations)

			createdCapturePage, err := m.AddCapturePage(testCtx, txHandlerDb, capturePageData)
			testCase.assertions(t, db, createdCapturePage, err)
		})
	}
}

type argsMutationsRestoreCapturePage struct {
	User           mysqlmodel.User
	Organization   mysqlmodel.Organization
	CapturePageSet mysqlmodel.CapturePageSet
	CapturePage    mysqlmodel.CapturePage
}

type testCaseRestoreCapturePage struct {
	name          string
	id            int
	argsMutations *argsMutationsRestoreCapturePage
	assertions    func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsRestoreCapturePage) int
}

func getRestoreCapturePageTestCases() []testCaseRestoreCapturePage {
	return []testCaseRestoreCapturePage{
		{
			name: "success",
			argsMutations: &argsMutationsRestoreCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePage: mysqlmodel.CapturePage{
					Name:             "CapturePage A",
					CapturePageSetID: 1,
					CreatedBy:        null.IntFrom(1),
					LastUpdatedBy:    null.IntFrom(1),
					IsActive:         true,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the capture page")
				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsRestoreCapturePage) int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page into the db")

				return args.CapturePage.ID
			},
		},
		{
			name: "fail-missing-entry-to-update",
			argsMutations: &argsMutationsRestoreCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				_, fetchErr := mysqlmodel.FindOrganization(context.TODO(), db, id)
				assert.Error(t, fetchErr, "organization should not exist")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsRestoreCapturePage) int {
				return args.Organization.ID * args.Organization.ID * args.Organization.ID
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

			id := testCase.mutations(t, db, testCase.argsMutations)
			err = m.RestoreCapturePage(testCtx, txHandlerDb, id)
			testCase.assertions(t, db, id, err)
		})
	}
}

type argsMutationsDeleteCapturePage struct {
	User           mysqlmodel.User
	Organization   mysqlmodel.Organization
	CapturePageSet mysqlmodel.CapturePageSet
	CapturePage    mysqlmodel.CapturePage
}

type testCaseDeleteCapturePage struct {
	name          string
	argsMutations *argsMutationsDeleteCapturePage
	assertions    func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations     func(t *testing.T, db *sqlx.DB, capturePage *argsMutationsDeleteCapturePage) int
}

func getDeleteCapturePageTestCases() []testCaseDeleteCapturePage {
	return []testCaseDeleteCapturePage{
		{
			name: "success",
			argsMutations: &argsMutationsDeleteCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePage: mysqlmodel.CapturePage{
					Name:             "CapturePage A",
					CapturePageSetID: 1,
					CreatedBy:        null.IntFrom(1),
					LastUpdatedBy:    null.IntFrom(1),
					IsActive:         true,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the capture page")

				assert.Equal(t, false, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsDeleteCapturePage) int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set into the db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page into the db")

				return args.CapturePage.ID
			},
		},
		{
			name: "failure-non-existent-capture-page",
			argsMutations: &argsMutationsDeleteCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				_, errFind := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.Error(t, errFind, "expected an error for non-existent capture page")
				require.Contains(t, errFind.Error(), "no rows", "error should indicate that the capture page was not found")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsDeleteCapturePage) int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting user into the db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into the db")

				return args.Organization.ID + args.Organization.ID*args.Organization.ID
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

			id := testCase.mutations(t, db, testCase.argsMutations)
			err = m.DeleteCapturePage(testCtx, txHandlerDb, id)
			require.NoError(t, err, "error deleting capture page")
			testCase.assertions(t, db, id, err)
		})
	}
}

type argsMutationsUpdateCapturePage struct {
	User              mysqlmodel.User
	Organization      mysqlmodel.Organization
	CapturePageSet    mysqlmodel.CapturePageSet
	CapturePage       mysqlmodel.CapturePage
	UpdateCapturePage model.UpdateCapturePage
}

type testCaseUpdateCapturePage struct {
	name              string
	updateCapturePage model.CapturePage
	argsMutations     *argsMutationsUpdateCapturePage
	assertions        func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations         func(t *testing.T, db *sqlx.DB, args *argsMutationsUpdateCapturePage) model.UpdateCapturePage
}

func getUpdateCapturePageTestCases() []testCaseUpdateCapturePage {
	return []testCaseUpdateCapturePage{
		{
			name: "success",
			argsMutations: &argsMutationsUpdateCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:              "Example Capture Page Set",
					OrganizationRefID: null.IntFrom(1),
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:               4,
					Name:             "Lawrence Michael",
					CapturePageSetID: 1,
					LastUpdatedAt:    null.TimeFrom(time.Now()),
					CreatedBy:        null.IntFrom(1),
					LastUpdatedBy:    null.IntFrom(1),
				},
				UpdateCapturePage: model.UpdateCapturePage{
					Id:               4,
					Name:             null.StringFrom("Virtue"),
					UserId:           null.IntFrom(1),
					CapturePageSetId: 1,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the capture page")

				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsUpdateCapturePage) model.UpdateCapturePage {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page set")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting capture page")

				return args.UpdateCapturePage
			},
		},
		{
			name: "fail",
			argsMutations: &argsMutationsUpdateCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				UpdateCapturePage: model.UpdateCapturePage{
					Id:   32123,
					Name: null.StringFrom("Wrong Name"),
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent capture page")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 32123")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsUpdateCapturePage) model.UpdateCapturePage {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization")

				return args.UpdateCapturePage
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

		updateData := testCase.mutations(t, db, testCase.argsMutations)

		_, err = m.UpdateCapturePage(testCtx, txHandlerDb, &updateData)
		testCase.assertions(t, db, updateData.Id, err)
	}
}
