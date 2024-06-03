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

type testCaseCreateOrganization struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	body              map[string]interface{}
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesCreateOrganization() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name: "success",
			body: map[string]interface{}{
				"name":                     "Example",
				"organization_type_ref_id": 1,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						orgService:      container.OrganizationService,
						catService:      container.CategoryService,
						capPagesService: container.CapturePagesService,
					}, func() {
						cleanup()
					}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				respStr := string(resp)
				require.NotNilf(t, resp, "unexpected nil response: %s", respStr)
				require.Equal(t, http.StatusCreated, respCode, "unexpected non-equal response code: %s", respStr)

				var organization *model.Organization
				err := json.Unmarshal(resp, &organization)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				require.NotNil(t, organization, "unexpected nil organization")

				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
		{
			name: "fail-empty-body",
			body: map[string]interface{}{},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						catService:      container.CategoryService,
						capPagesService: container.CapturePagesService,
						orgService:      container.OrganizationService,
					}, func() {
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
			name: "fail-invalid-organization-ref-id",
			body: map[string]interface{}{
				"name":                     "Example",
				"organization_type_ref_id": 19999,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				container, cleanup := testassets.GetConcreteContainer(t)
				return &testServices{
						catService:      container.CategoryService,
						capPagesService: container.CapturePagesService,
						orgService:      container.OrganizationService,
					}, func() {
						cleanup()
					}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, respCode, http.StatusBadRequest)
				require.Contains(t, string(resp), "invalid organization_type_id")
			},
		},
		{
			name: "fail-mock-server-error",
			body: map[string]interface{}{
				"name":                     "Example",
				"organization_type_ref_id": 1,
			},
			fnGetTestServices: func(t *testing.T) (*testServices, func()) {
				fakeOrganizationService := apifakes.FakeOrganizationService{}
				fakeOrganizationService.CreateOrganizationReturns(nil, errors.New("mock error"))

				fakeCategoriesService := apifakes.FakeCategoryService{}
				fakeCategoriesService.CreateCategoryReturns(nil, errors.New("mock error"))

				fakeCapturePagesService := apifakes.FakeCapturePagesService{}
				fakeCapturePagesService.CreateCapturePagesReturns(nil, errors.New("mock error"))
				return &testServices{
					orgService:      &fakeOrganizationService,
					catService:      &fakeCategoriesService,
					capPagesService: &fakeCapturePagesService,
				}, func() {}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, http.StatusInternalServerError, respCode)
			},
		},
	}
}

func Test_CreateOrganization(t *testing.T) {
	for _, testCase := range getTestCasesCreateOrganization() {
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

			req := httptest.NewRequest(http.MethodPost, "/api/v1/organization", bytes.NewBuffer(reqB))
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

type testCaseListOrganization struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	mutations       func(t *testing.T, modules *testassets.Container)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesListOrganizations() []testCaseListOrganization {
	testCases := []testCaseListOrganization{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1, 2},
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
				var respPaginated model.PaginatedOrganizations
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
				assert.True(t, len(respPaginated.Organizations) == 2, "unexpected empty organizations")
				assert.Equal(t, respPaginated.Pagination.RowCount, 2, "unexpected row_count")
				assert.Equal(t, respPaginated.Pagination.TotalCount, 2, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		{
			name: "success-limit-offset-page-1",
			queryParameters: map[string]interface{}{
				"ids_in":   []int{1, 2},
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

				var respPaginated model.PaginatedOrganizations
				err := json.Unmarshal(resp, &respPaginated)

				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.Organizations, "unexpected empty organizations")
				assert.Greater(t, respPaginated.Pagination.TotalCount, 1, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		{
			name: "success-limit-offset-page-2",
			queryParameters: map[string]interface{}{
				"ids_in":   []int{1, 2, 3},
				"page":     1,
				"max_rows": 3,
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

				var respPaginated model.PaginatedOrganizations

				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.Organizations, "unexpected empty organizations")
				assert.True(t, respPaginated.Pagination.TotalCount > 2, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		{
			name: "success-all-filters",
			queryParameters: map[string]interface{}{
				"ids_in":                     []int{1, 2, 3},
				"organizations_type_id_in":   []int{1},
				"organizations_type_name_in": []string{"User Types"},
				"organizations_name_in":      []string{"Admin"},
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
				var respPaginated model.PaginatedOrganizations
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
				assert.True(t, len(respPaginated.Organizations) > 0, "unexpected empty organizations")
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

				err = store.DropOrganizationTable(context.TODO(), tx)
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

func Test_ListOrganizations(t *testing.T) {
	for _, testCase := range getTestCasesListOrganizations() {
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

			url := strutil.AppendQueryToURL("/api/v1/organization", testCase.queryParameters)
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

type testCaseUpdateOrganization struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	body              map[string]interface{}
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesUpdateOrganization() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
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
			body: map[string]interface{}{
				"id":   1,
				"name": "demby",
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var organization *model.Organization
				err := json.Unmarshal(resp, &organization)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
	}
}

func Test_UpdateOrganization(t *testing.T) {
	for _, testCase := range getTestCasesUpdateOrganization() {
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

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/organization", bytes.NewBuffer(reqB))
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

type testCaseDeleteOrganization struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	orgID             int
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesDeleteOrganization() []testCaseDeleteOrganization {
	return []testCaseDeleteOrganization{
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
			orgID: 1,
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var organization *model.Organization
				err := json.Unmarshal(resp, &organization)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
	}
}

func Test_DeleteOrganization(t *testing.T) {
	for _, testCase := range getTestCasesDeleteOrganization() {
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

			url := fmt.Sprintf("/api/v1/organization/%d", testCase.orgID)
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

type testCaseRestoreOrganization struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testServices, func())
	mutations         func(t *testing.T, modules *testassets.Container)
	orgID             int
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesRestoreOrganization() []testCaseRestoreOrganization {
	return []testCaseRestoreOrganization{
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
				// Optional: Perform any mutations needed before the test
			},
			orgID: 1, // Ensure the organization ID to delete is provided
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var organization *model.Organization
				err := json.Unmarshal(resp, &organization)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
	}
}

func Test_RestoreOrganization(t *testing.T) {
	for _, testCase := range getTestCasesRestoreOrganization() {
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

			url := fmt.Sprintf("/api/v1/organization/%d", testCase.orgID)
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
