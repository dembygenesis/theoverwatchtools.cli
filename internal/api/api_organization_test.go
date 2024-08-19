package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type testCaseListOrganization struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	mutations       func(t *testing.T, db *sqlx.DB, modules *testassets.Container)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesListOrganizations() []testCaseListOrganization {
	testCases := []testCaseListOrganization{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) {
				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: 1,
				}

				_, err := modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				require.NoError(t, err, "error adding the organization")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")

				var respPaginated model.PaginatedOrganizations
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Len(t, respPaginated.Organizations, 1, "unexpected number of organizations")
				assert.Equal(t, 1, respPaginated.Pagination.RowCount, "unexpected row_count")
				assert.Equal(t, 1, respPaginated.Pagination.TotalCount, "unexpected total_count")
				assert.Equal(t, 1, respPaginated.Pagination.Page, "unexpected page")
			},
		},
	}

	return testCases
}

func Test_ListOrganizations(t *testing.T) {
	for _, testCase := range getTestCasesListOrganizations() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			if testCase.queryParameters == nil {
				testCase.queryParameters = make(map[string]interface{})
			}

			handlers, _ := testCase.getContainer(t)

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

			testCase.mutations(t, db, handlers)

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}

type testCaseDeleteOrganization struct {
	name         string
	getContainer func(t *testing.T) (*testassets.Container, func())
	mutations    func(t *testing.T, db *sqlx.DB, modules *testassets.Container) int
	body         map[string]interface{}
	assertions   func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesDeleteOrganization() []testCaseDeleteOrganization {
	return []testCaseDeleteOrganization{
		{
			name: "success",
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) int {
				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: 1,
				}

				organization, err := modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				require.NoError(t, err, "error adding the organization")

				id := organization.Id
				return id
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			body: map[string]interface{}{
				"id":   1,
				"name": "demby",
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusNoContent, respCode, "unexpected response code")
				assert.Empty(t, resp, "expected empty response body for no-content response")
			},
		},
	}
}

func Test_DeleteOrganization(t *testing.T) {
	for _, testCase := range getTestCasesDeleteOrganization() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()

			handlers, _ := testCase.getContainer(t)

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				Logger:              logger.New(context.TODO()),
			}

			orgId := testCase.mutations(t, db, handlers)

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/organization/"+strconv.Itoa(orgId), bytes.NewBuffer(reqB))
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

type testCaseUpdateOrganization struct {
	name         string
	getContainer func(t *testing.T) (*testassets.Container, func())
	mutations    func(t *testing.T, db *sqlx.DB, modules *testassets.Container)
	body         map[string]interface{}
	assertions   func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesUpdateOrganization() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name: "success",
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) {
				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: 1,
				}

				_, err := modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				require.NoError(t, err, "error adding the organization")
			},
			body: map[string]interface{}{
				"id":   1,
				"name": "demby2",
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				require.NotNil(t, ctn, "unexpected nil container")
				return ctn, func() {
					cleanup()
				}
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
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()

			handlers, _ := testCase.getContainer(t)

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

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/organization", bytes.NewBuffer(reqB))
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			testCase.mutations(t, db, handlers)

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}

type testCaseRestoreOrganization struct {
	name         string
	getContainer func(t *testing.T) (*testassets.Container, func())
	mutations    func(t *testing.T, db *sqlx.DB, modules *testassets.Container) int
	body         map[string]interface{}
	assertions   func(t *testing.T, resp []byte, respCode int, orgID int, db *sqlx.DB)
}

func getTestCasesRestoreOrganization() []testCaseRestoreOrganization {
	return []testCaseRestoreOrganization{
		{
			name: "success",
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) int {
				//organizationModel := &model.CreateOrganization{
				//	Name:   "demby",
				//	UserId: 1,
				//}
				//
				//_, err := modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				//require.NoError(t, err, "error adding the organization")
				//
				//deleteModel := &model.DeleteOrganization{
				//	ID: organizationModel.UserId,
				//}
				//
				//err = modules.OrganizationService.DeleteOrganization(context.Background(), deleteModel)
				//require.NoError(t, err, "error deleting the organization")

				organizationModel := mysqlmodel.Organization{
					ID:        1,
					Name:      "Demby",
					CreatedBy: null.IntFrom(1),
					IsActive:  false,
				}
				err := organizationModel.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "error inserting organization into organization table")

				id := organizationModel.ID
				return id
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			body: map[string]interface{}{
				"id":   1,
				"name": "Demby",
			},
			assertions: func(t *testing.T, resp []byte, respCode int, orgID int, db *sqlx.DB) {
				require.Equal(t, http.StatusNoContent, respCode, "unexpected response code")
				assert.Empty(t, resp, "unexpected non-empty response body")

				organization, err := mysqlmodel.FindOrganization(context.TODO(), db, orgID)
				fmt.Println("the organization returned --- ", strutil.GetAsJson(organization))
				require.NoError(t, err, "unexpected error fetching organization from database")

				assert.True(t, organization.IsActive, "expected is_active to be true")
			},
		},
	}
}

func Test_RestoreOrganization(t *testing.T) {
	for _, testCase := range getTestCasesRestoreOrganization() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()

			handlers, _ := testCase.getContainer(t)

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				Logger:              logger.New(context.TODO()),
			}

			orgId := testCase.mutations(t, db, handlers)

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/organization/"+strconv.Itoa(orgId), bytes.NewBuffer(reqB))
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			fmt.Println("the req --- ", req)

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode, orgId, db)
		})
	}
}
