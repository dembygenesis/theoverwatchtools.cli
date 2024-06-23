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
				assert.True(t, len(respPaginated.CapturePages) == 1, "unexpected empty capture pages")
				assert.Equal(t, respPaginated.Pagination.RowCount, 1, "unexpected row_count")
				assert.Equal(t, respPaginated.Pagination.TotalCount, 1, "unexpected total_count")
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
				assert.Greater(t, respPaginated.Pagination.TotalCount, 0, "unexpected total_count")
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
				assert.True(t, respPaginated.Pagination.TotalCount > 0, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		{
			name: "success-all-filters",
			queryParameters: map[string]interface{}{
				"ids_in":              []int{1, 2},
				"html":                []string{"<html><body><h1>Example 1</h1></body></html>"},
				"capture_page_set_id": []int{1},
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
				assert.True(t, len(respPaginated.CapturePages) > 0, "unexpected empty capture pages")
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

				err = store.DropCapturePagesTable(context.TODO(), tx)
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
			fmt.Println("===============>", strutil.GetAsJson(handlers))
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				CapturePagesService: handlers.CapturePagesService,
				ClickTrackerService: handlers.ClickTrackerService,
				Logger:              logger.New(context.TODO()),
			}
			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := strutil.AppendQueryToURL("/api/v1/capturepages", testCase.queryParameters)
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

// test for post method

type testCaseCreateCapturePages struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	body              map[string]interface{}
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesCreateCapturePages() []testCaseCreateCapturePages {
	return []testCaseCreateCapturePages{
		{
			name: "success",
			body: map[string]interface{}{
				"name":                "Example",
				"capture_page_set_id": 1,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{catService: container.CategoryService, orgService: container.OrganizationService, capPagesService: container.CapturePagesService, ctService: container.ClickTrackerService}, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				respStr := string(resp)
				require.NotNilf(t, resp, "unexpected nil response: %s", respStr)
				require.Equal(t, http.StatusCreated, respCode, "unexpected non-equal response code: %s", respStr)

				var capturepages *model.CapturePages
				err := json.Unmarshal(resp, &capturepages)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				require.NotNil(t, capturepages, "unexpected nil capture pages")

				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePages{*capturepages})
			},
		},
		{
			name: "fail-empty-body",
			body: map[string]interface{}{},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{catService: container.CategoryService, orgService: container.OrganizationService, capPagesService: container.CapturePagesService, ctService: container.ClickTrackerService}, func() {
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
			name: "fail-invalid-capture-page-set-id",
			body: map[string]interface{}{
				"name":                "Example",
				"capture_page_set_id": 19999,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{catService: container.CategoryService, orgService: container.OrganizationService, capPagesService: container.CapturePagesService, ctService: container.ClickTrackerService}, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {

				assert.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, respCode, http.StatusBadRequest)
				require.Contains(t, string(resp), "capture_page_sets")
			},
		},
		{
			name: "fail-mock-server-error",
			body: map[string]interface{}{
				"name":                "Example",
				"capture_page_set_id": 1,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				fakeCategoryService := apifakes.FakeCategoryService{}
				fakeCategoryService.CreateCategoryReturns(nil, errors.New("mock error"))
				fakeOrganizationService := apifakes.FakeOrganizationService{}
				fakeOrganizationService.CreateOrganizationReturns(nil, errors.New("mock error"))
				fakeCapturePagesService := apifakes.FakeCapturePagesService{}
				fakeCapturePagesService.CreateCapturePagesReturns(nil, errors.New("mock error"))
				fakeClickTrackerService := apifakes.FakeClickTrackerService{}
				fakeClickTrackerService.CreateClickTrackerReturns(nil, errors.New("mock error"))
				return &testServices{catService: &fakeCategoryService, orgService: &fakeOrganizationService, capPagesService: &fakeCapturePagesService, ctService: &fakeClickTrackerService}, func() {}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, http.StatusInternalServerError, respCode)
			},
		},
	}
}

func Test_CreateCapturePages(t *testing.T) {
	for _, testCase := range getTestCasesCreateCapturePages() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				CapturePagesService: handlers.capPagesService,
				ClickTrackerService: handlers.ctService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			fmt.Println("this is REQB =========================> ", reqB)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/capturepages", bytes.NewBuffer(reqB))
			fmt.Println("this is REQ =========================> ", req)
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

// test for update method

type testCaseUpdateCapturePages struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	body              map[string]interface{}
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesUpdateCapturePages() []testCaseUpdateCapturePages {
	return []testCaseUpdateCapturePages{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						orgService:      container.OrganizationService,
						capPagesService: container.CapturePagesService,
						catService:      container.CategoryService,
						ctService:       container.ClickTrackerService,
					}, func() {
						cleanup()
					}
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			body: map[string]interface{}{
				"id":   1,
				"name": "demby",
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var capturepages *model.CapturePages
				err := json.Unmarshal(resp, &capturepages)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePages{*capturepages})
			},
		},
	}
}

func Test_UpdateCapturePages(t *testing.T) {
	for _, testCase := range getTestCasesUpdateCapturePages() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CapturePagesService: handlers.capPagesService,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				ClickTrackerService: handlers.ctService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/capturepages", bytes.NewBuffer(reqB))
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

// test for delete method

type testCaseDeleteCapturePages struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	capID             int
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesDeleteCapturePages() []testCaseDeleteCapturePages {
	return []testCaseDeleteCapturePages{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						catService:      container.CategoryService,
						orgService:      container.OrganizationService,
						capPagesService: container.CapturePagesService,
						ctService:       container.ClickTrackerService,
					}, func() {
						cleanup()
					}
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			capID: 1,
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var capturepages *model.CapturePages
				err := json.Unmarshal(resp, &capturepages)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePages{*capturepages})
			},
		},
	}
}

func Test_DeleteCapturePages(t *testing.T) {
	for _, testCase := range getTestCasesDeleteCapturePages() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CapturePagesService: handlers.capPagesService,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				ClickTrackerService: handlers.ctService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := fmt.Sprintf("/api/v1/capturepages/%d", testCase.capID)
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

// test for restore method

type testCaseRestoreCapturePages struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	capID             int
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesRestoreCapturePages() []testCaseRestoreCapturePages {
	return []testCaseRestoreCapturePages{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						catService:      container.CategoryService,
						orgService:      container.OrganizationService,
						capPagesService: container.CapturePagesService,
						ctService:       container.ClickTrackerService,
					}, func() {
						cleanup()
					}
			},
			mutations: func(t *testing.T, modules *testassets.Container) {

			},
			capID: 1,
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var capturepages *model.CapturePages
				err := json.Unmarshal(resp, &capturepages)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePages{*capturepages})
			},
		},
	}
}

func Test_RestoreCapturePages(t *testing.T) {
	for _, testCase := range getTestCasesRestoreCapturePages() {
		t.Run(testCase.name, func(t *testing.T) {
			handlers, cleanup := testCase.fnGetTestServices(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CapturePagesService: handlers.capPagesService,
				CategoryService:     handlers.catService,
				OrganizationService: handlers.orgService,
				ClickTrackerService: handlers.ctService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := fmt.Sprintf("/api/v1/capturepages/%d", testCase.capID)
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
