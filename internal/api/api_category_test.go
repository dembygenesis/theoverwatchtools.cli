package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/api/apifakes"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServices struct {
	catService      categoryService
	orgService      organizationService
	capPagesService capturePagesService
	ctService       clickTrackerService
}

type testCaseCreateCategory struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	body              map[string]interface{}
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesCreateCategory() []testCaseCreateCategory {
	return []testCaseCreateCategory{
		{
			name: "success",
			body: map[string]interface{}{
				"name":                 "Example",
				"category_type_ref_id": 1,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{catService: container.CategoryService, orgService: container.OrganizationService, capPagesService: container.CapturePagesService}, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				respStr := string(resp)
				require.NotNilf(t, resp, "unexpected nil response: %s", respStr)
				require.Equal(t, http.StatusCreated, respCode, "unexpected non-equal response code: %s", respStr)

				var category *model.Category
				err := json.Unmarshal(resp, &category)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				require.NotNil(t, category, "unexpected nil category")

				modelhelpers.AssertNonEmptyCategories(t, []model.Category{*category})
			},
		},
		{
			name: "fail-empty-body",
			body: map[string]interface{}{},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{catService: container.CategoryService, orgService: container.OrganizationService, capPagesService: container.CapturePagesService}, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, respCode, http.StatusBadRequest)
				require.Contains(t, string(resp), "validate:")
			},
		},
		{
			name: "fail-invalid-category-ref-id",
			body: map[string]interface{}{
				"name":                 "Example",
				"category_type_ref_id": 19999,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{catService: container.CategoryService, orgService: container.OrganizationService, capPagesService: container.CapturePagesService}, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, respCode, http.StatusBadRequest)
				require.Contains(t, string(resp), "invalid category_type_id")
			},
		},
		{
			name: "fail-mock-server-error",
			body: map[string]interface{}{
				"name":                 "Example",
				"category_type_ref_id": 1,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				fakeCategoryService := apifakes.FakeCategoryService{}
				fakeCategoryService.CreateCategoryReturns(nil, errors.New("mock error"))
				fakeOrganizationService := apifakes.FakeOrganizationService{}
				fakeOrganizationService.CreateOrganizationReturns(nil, errors.New("mock error"))
				fakeCapturePagesService := apifakes.FakeCapturePagesService{}
				fakeCapturePagesService.CreateCapturePagesReturns(nil, errors.New("mock error"))
				return &testServices{catService: &fakeCategoryService, orgService: &fakeOrganizationService, capPagesService: &fakeCapturePagesService}, func() {}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, http.StatusInternalServerError, respCode)
			},
		},
	}
}

func Test_CreateCategory(t *testing.T) {
	for _, testCase := range getTestCasesCreateCategory() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				CapturePagesService: handlers.capPagesService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPost, "/api/v1/category", bytes.NewBuffer(reqB))
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}

type testCaseListCategory struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	mutations       func(t *testing.T, modules *testassets.Container)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesListCategories() []testCaseListCategory {
	testCases := []testCaseListCategory{
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
				var respPaginated model.PaginatedCategories
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
				assert.True(t, len(respPaginated.Categories) == 3, "unexpected empty categories")
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

				var respPaginated model.PaginatedCategories
				err := json.Unmarshal(resp, &respPaginated)

				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.Categories, "unexpected empty categories")
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

				var respPaginated model.PaginatedCategories
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.Categories, "unexpected empty categories")
				assert.True(t, respPaginated.Pagination.TotalCount > 2, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 2, "unexpected page")
			},
		},
		{
			name: "success-all-filters",
			queryParameters: map[string]interface{}{
				"ids_in":                []int{1, 2, 3},
				"category_type_id_in":   []int{1},
				"category_type_name_in": []string{"User Types"},
				"category_name_in":      []string{"Admin"},
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
				var respPaginated model.PaginatedCategories
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
				assert.True(t, len(respPaginated.Categories) > 0, "unexpected empty categories")
				assert.True(t, respPaginated.Pagination.MaxRows > 0, "unexpected empty rows")
				assert.True(t, respPaginated.Pagination.RowCount > 0, "unexpected empty count")
				assert.True(t, len(respPaginated.Pagination.Pages) > 0, "unexpected empty pages")
			},
		},
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

func Test_ListCategories(t *testing.T) {
	for _, testCase := range getTestCasesListCategories() {
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
				OrganizationService: handlers.OrganizationService,
				CapturePagesService: handlers.CapturePagesService,
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

			fmt.Println("testing the cat tests")
		})
	}
}

type testCaseUpdateCategory struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	body              map[string]interface{}
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesUpdateCategory() []testCaseUpdateCategory {
	return []testCaseUpdateCategory{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						orgService:      container.OrganizationService,
						capPagesService: container.CapturePagesService,
						catService:      container.CategoryService,
					}, func() {
						cleanup()
					}
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			body: map[string]interface{}{
				"id":                   1,
				"category_type_ref_id": 1,
				"name":                 "Lawrence",
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var category *model.Category
				err := json.Unmarshal(resp, &category)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyCategories(t, []model.Category{*category})
			},
		},
	}
}

func Test_UpdateCategory(t *testing.T) {
	for _, testCase := range getTestCasesUpdateCategory() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CapturePagesService: handlers.capPagesService,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/category", bytes.NewBuffer(reqB))
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}

type testCaseDeleteCategory struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	catID             int
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesDeleteCategory() []testCaseDeleteCategory {
	return []testCaseDeleteCategory{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						catService:      container.CategoryService,
						orgService:      container.OrganizationService,
						capPagesService: container.CapturePagesService,
					}, func() {
						cleanup()
					}
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			catID: 1,
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var organization *model.Category
				err := json.Unmarshal(resp, &organization)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyCategories(t, []model.Category{*organization})
			},
		},
	}
}

func Test_DeleteCategory(t *testing.T) {
	for _, testCase := range getTestCasesDeleteCategory() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CapturePagesService: handlers.capPagesService,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := fmt.Sprintf("/api/v1/category/%d", testCase.catID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			_, err = api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")
		})
	}
}

type testCaseRestoreCategory struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	catID             int
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesRestoreCategory() []testCaseRestoreCategory {
	return []testCaseRestoreCategory{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						catService:      container.CategoryService,
						orgService:      container.OrganizationService,
						capPagesService: container.CapturePagesService,
					}, func() {
						cleanup()
					}
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			catID: 1,
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var categories *model.Category
				err := json.Unmarshal(resp, &categories)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyCategories(t, []model.Category{*categories})
			},
		},
	}
}

func Test_RestoreCategory(t *testing.T) {
	for _, testCase := range getTestCasesRestoreCategory() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CapturePagesService: handlers.capPagesService,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := fmt.Sprintf("/api/v1/category/%d", testCase.catID)
			req := httptest.NewRequest(http.MethodPatch, url, nil)
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			_, err = api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")
		})
	}
}
