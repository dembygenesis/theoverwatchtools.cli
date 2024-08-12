package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testCaseListOrganization struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	mutations       func(t *testing.T, db *sqlx.DB, modules *testassets.Container) (organizations []interface{})
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int, organization []interface{})
}

func getTestCasesListOrganizations() []testCaseListOrganization {
	testCases := []testCaseListOrganization{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1, 2, 3},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) (returnedOrganizations []interface{}) {
				entryUser := mysqlmodel.User{
					ID:                5,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				orgs := []mysqlmodel.Organization{
					{
						ID:            1,
						Name:          "TEST1",
						CreatedBy:     null.IntFrom(entryUser.ID),
						LastUpdatedBy: null.IntFrom(entryUser.ID),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true,
					},
					{
						ID:            2,
						Name:          "TEST2",
						CreatedBy:     null.IntFrom(entryUser.ID),
						LastUpdatedBy: null.IntFrom(entryUser.ID),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true,
					},
					{
						ID:            3,
						Name:          "TEST3",
						CreatedBy:     null.IntFrom(entryUser.ID),
						LastUpdatedBy: null.IntFrom(entryUser.ID),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true,
					},
				}

				for _, org := range orgs {
					err := org.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting organization")
				}

				organizations, err := mysqlmodel.Organizations(
					qm.Where("id = ? OR id = ? OR id = ?", 1, 2, 3),
				).All(context.Background(), db)

				fmt.Println("the organization query?? --- ", strutil.GetAsJson(organizations))

				require.NoError(t, err, "unexpected error fetching organizations from db")
				require.Equal(t, 3, len(organizations), "unexpected number of organizations in the database")

				returnedOrganizations = make([]interface{}, len(organizations))
				for i, org := range organizations {
					returnedOrganizations[i] = org
				}
				return returnedOrganizations
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int, organizations []interface{}) {
				var respPaginated model.PaginatedOrganizations
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				var expectedOrganizations []mysqlmodel.Organization
				for _, org := range organizations {
					if org, ok := org.(mysqlmodel.Organization); ok {
						expectedOrganizations = append(expectedOrganizations, org)
					}
				}

				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
				assert.True(t, len(respPaginated.Organizations) == len(expectedOrganizations), "unexpected number of organizations in the response")
				assert.Equal(t, respPaginated.Pagination.RowCount, len(expectedOrganizations), "unexpected row_count")
				assert.Equal(t, respPaginated.Pagination.TotalCount, len(expectedOrganizations), "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		//{
		//	name: "success-limit-offset-page-1",
		//	queryParameters: map[string]interface{}{
		//		"ids_in":   []int{1, 2, 3},
		//		"page":     1,
		//		"max_rows": 1,
		//	},
		//	mutations: func(t *testing.T, modules *testassets.Container) {
		//
		//	},
		//	getContainer: func(t *testing.T) (*testassets.Container, func()) {
		//		ctn, cleanup := testassets.GetConcreteContainer(t)
		//		return ctn, func() {
		//			cleanup()
		//		}
		//	},
		//	assertions: func(t *testing.T, resp []byte, respCode int) {
		//		require.Equalf(t, http.StatusOK, respCode, "unexpected status code: %v, with resp: %s", respCode, string(resp))
		//
		//		var respPaginated model.PaginatedCategories
		//		err := json.Unmarshal(resp, &respPaginated)
		//
		//		require.NoError(t, err, "unexpected error unmarshalling the response")
		//		assert.NotNil(t, respPaginated.Categories, "unexpected empty categories")
		//		assert.Greater(t, respPaginated.Pagination.TotalCount, 1, "unexpected total_count")
		//		assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
		//	},
		//},
		//{
		//	name: "success-limit-offset-page-2",
		//	queryParameters: map[string]interface{}{
		//		"ids_in":   []int{1, 2, 3},
		//		"page":     2,
		//		"max_rows": 2,
		//	},
		//	mutations: func(t *testing.T, modules *testassets.Container) {
		//
		//	},
		//	getContainer: func(t *testing.T) (*testassets.Container, func()) {
		//		ctn, cleanup := testassets.GetConcreteContainer(t)
		//		return ctn, func() {
		//			cleanup()
		//		}
		//	},
		//	assertions: func(t *testing.T, resp []byte, respCode int) {
		//		require.Equalf(t, http.StatusOK, respCode, "unexpected status code: %v, with resp: %s", respCode, string(resp))
		//
		//		var respPaginated model.PaginatedCategories
		//		err := json.Unmarshal(resp, &respPaginated)
		//		require.NoError(t, err, "unexpected error unmarshalling the response")
		//		assert.NotNil(t, respPaginated.Categories, "unexpected empty categories")
		//		assert.True(t, respPaginated.Pagination.TotalCount > 2, "unexpected total_count")
		//		assert.Equal(t, respPaginated.Pagination.Page, 2, "unexpected page")
		//	},
		//},
		//{
		//	name: "success-all-filters",
		//	queryParameters: map[string]interface{}{
		//		"ids_in":                []int{1, 2, 3},
		//		"category_type_id_in":   []int{1},
		//		"category_type_name_in": []string{"User Types"},
		//		"category_name_in":      []string{"Admin"},
		//	},
		//	mutations: func(t *testing.T, modules *testassets.Container) {
		//
		//	},
		//	getContainer: func(t *testing.T) (*testassets.Container, func()) {
		//		ctn, cleanup := testassets.GetConcreteContainer(t)
		//		return ctn, func() {
		//			cleanup()
		//		}
		//	},
		//	assertions: func(t *testing.T, resp []byte, respCode int) {
		//		var respPaginated model.PaginatedCategories
		//		err := json.Unmarshal(resp, &respPaginated)
		//		require.NoError(t, err, "unexpected error unmarshalling the response")
		//
		//		assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
		//		assert.True(t, len(respPaginated.Categories) > 0, "unexpected empty categories")
		//		assert.True(t, respPaginated.Pagination.MaxRows > 0, "unexpected empty rows")
		//		assert.True(t, respPaginated.Pagination.RowCount > 0, "unexpected empty count")
		//		assert.True(t, len(respPaginated.Pagination.Pages) > 0, "unexpected empty pages")
		//	},
		//},
		//{
		//	name:            "empty_store",
		//	queryParameters: map[string]interface{}{},
		//	mutations: func(t *testing.T, modules *testassets.Container) {
		//		store := modules.MySQLStore
		//		require.NotNil(t, store, "unexpected nil: store")
		//
		//		connProvider := modules.ConnProvider
		//		require.NotNil(t, store, "unexpected nil: txProvider")
		//
		//		tx, err := connProvider.Tx(context.TODO())
		//		require.NoError(t, err, "unexpected err for getting tx")
		//
		//		err = store.DropCategoryTable(context.TODO(), tx)
		//		require.NoError(t, err, "unexpected err for drop command")
		//
		//		err = tx.Commit(context.TODO())
		//		require.NoError(t, err, "unexpected err on commit")
		//	},
		//	getContainer: func(t *testing.T) (*testassets.Container, func()) {
		//		ctn, cleanup := testassets.GetConcreteContainer(t)
		//		return ctn, func() {
		//			cleanup()
		//		}
		//	},
		//	assertions: func(t *testing.T, resp []byte, respCode int) {
		//		assert.Equal(t, http.StatusInternalServerError, respCode)
		//	},
		//},
	}

	return testCases
}

func Test_ListOrganizations(t *testing.T) {
	for _, testCase := range getTestCasesListOrganizations() {
		db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		defer cleanup()
		require.NotEmpty(t, testCase.name, "unexpected empty test name")

		t.Run(testCase.name, func(t *testing.T) {
			if testCase.queryParameters == nil {
				testCase.queryParameters = make(map[string]interface{})
			}

			handlers, _ := testCase.getContainer(t)
			//defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := strutil.AppendQueryToURL("/api/v1/organization", testCase.queryParameters)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			createdOrganizations := testCase.mutations(t, db, handlers)

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode, createdOrganizations)
		})
	}
}
