package api

import (
	"context"
	"encoding/json"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCaseListCapturePages struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	mutations       func(t *testing.T, modules *testassets.Container)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesListCapturePages() []testCaseListCapturePages {
	testCases := []testCaseListCapturePages{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1, 2, 3},
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				var respPaginated model.PaginatedCapturePages
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
				assert.True(t, len(respPaginated.CapturePages) == 3, "unexpected empty capture pages")
				assert.Equal(t, respPaginated.Pagination.RowCount, 3, "unexpected row_count")
				assert.Equal(t, respPaginated.Pagination.TotalCount, 3, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		{
			name: "success-limit-offset-page-1",
			queryParameters: map[string]interface{}{
				"ids_in":   []int{1, 2, 3},
				"page":     1,
				"max_rows": 1,
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equalf(t, http.StatusOK, respCode, "unexpected status code: %v, with resp: %s", respCode, string(resp))

				var respPaginated model.PaginatedCapturePages
				err := json.Unmarshal(resp, &respPaginated)

				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.CapturePages, "unexpected empty capture pages")
				assert.Greater(t, respPaginated.Pagination.TotalCount, 1, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		{
			name: "success-limit-offset-page-2",
			queryParameters: map[string]interface{}{
				"ids_in":   []int{1, 2, 3},
				"page":     2,
				"max_rows": 2,
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equalf(t, http.StatusOK, respCode, "unexpected status code: %v, with resp: %s", respCode, string(resp))

				var respPaginated model.PaginatedCapturePages
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.CapturePages, "unexpected empty capture pages")
				assert.True(t, respPaginated.Pagination.TotalCount > 2, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 2, "unexpected page")
			},
		},
		//{
		//	name: "success-all-filters",
		//	queryParameters: map[string]interface{}{
		//		"ids_in":                    []int{1, 2, 3},
		//		"category_type_id_in":       []int{1},
		//		"capture_page_type_name_in": []string{"User Types"},
		//		"capture_page_name_in":      []string{"Admin"},
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
		{
			name:            "empty_store",
			queryParameters: map[string]interface{}{},
			mutations: func(t *testing.T, modules *testassets.Container) {
				store := modules.MySQLStore
				require.NotNil(t, store, "unexpected nil: store")

				connProvider := modules.ConnProvider
				require.NotNil(t, store, "unexpected nil: txProvider")

				tx, err := connProvider.Tx(context.TODO())
				require.NoError(t, err, "unexpected err for getting tx")

				err = store.DropCategoryTable(context.TODO(), tx)
				require.NoError(t, err, "unexpected err for drop command")

				err = tx.Commit(context.TODO())
				require.NoError(t, err, "unexpected err on commit")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.Equal(t, http.StatusInternalServerError, respCode)
			},
		},
	}

	return testCases
}

func Test_ListCapturePages(t *testing.T) {
	for _, testCase := range getTestCasesListCapturePages() {
		require.NotEmpty(t, testCase.name, "unexpected empty test name")

		t.Run(testCase.name, func(t *testing.T) {
			if testCase.queryParameters == nil {
				testCase.queryParameters = make(map[string]interface{})
			}

			handlers, cleanup := testCase.getContainer(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				CapturePagesService: handlers.CapturePagesService,
				OrganizationService: handlers.OrganizationService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := strutil.AppendQueryToURL("/api/v1/category", testCase.queryParameters)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			testCase.mutations(t, handlers)

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}
