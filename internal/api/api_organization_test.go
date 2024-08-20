package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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

type testCaseCreateOrganization struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testassets.Container, func())
	mutations         func(t *testing.T, db *sqlx.DB, modules *testassets.Container) *model.CreateOrganization
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesCreateOrganization() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) *model.CreateOrganization {
				tx, err := db.Begin()
				require.NoError(t, err, "unexpected error starting a transaction")
				defer func(tx *sql.Tx) {
					err := tx.Rollback()
					require.NoError(t, err, "error rolling back")
				}(tx)

				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(2),
					LastUpdatedBy:     null.IntFrom(2),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err = entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				createdOrganization := &model.CreateOrganization{
					Name:   "mohamed",
					UserId: entryUser.CreatedBy.Int,
				}

				return createdOrganization
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
	}
}

func Test_CreateOrganization(t *testing.T) {
	for _, testCase := range getTestCasesCreateOrganization() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()

			handlers, _ := testCase.fnGetTestServices(t)

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				Logger:              logger.New(context.TODO()),
			}

			organization := testCase.mutations(t, db, handlers)

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(organization)
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

type testCaseListOrganizations struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	mutations       func(t *testing.T, db *sqlx.DB, modules *testassets.Container)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesListOrganizations() []testCaseListOrganizations {
	testCases := []testCaseListOrganizations{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) {

				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(2),
					LastUpdatedBy:     null.IntFrom(2),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: entryUser.CreatedBy.Int,
				}

				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
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
		{
			name: "success-limit-offset-page-1",
			queryParameters: map[string]interface{}{
				"ids_in":   []int{1, 2, 3},
				"page":     1,
				"max_rows": 1,
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) {
				tx, err := db.Begin()
				require.NoError(t, err, "unexpected error starting a transaction")
				defer func() {
					if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
						require.NoError(t, err, "error rolling back")
					}
				}()

				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(2),
					LastUpdatedBy:     null.IntFrom(2),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err = entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: entryUser.CreatedBy.Int,
				}
				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				require.NoError(t, err, "error inserting organization into organization table")

				organizationModel1 := &model.CreateOrganization{
					Name:   "younes",
					UserId: entryUser.CreatedBy.Int,
				}
				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel1)
				require.NoError(t, err, "error adding the organization")

				err = tx.Commit()
				require.NoError(t, err, "unexpected error committing the transaction")
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
				"page":     2,
				"max_rows": 2,
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) {
				tx, err := db.Begin()
				require.NoError(t, err, "unexpected error starting a transaction")
				defer func() {
					if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
						require.NoError(t, err, "error rolling back")
					}
				}()

				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err = entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: entryUser.CreatedBy.Int,
				}
				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				require.NoError(t, err, "error inserting organization into organization table")

				organizationModel1 := &model.CreateOrganization{
					Name:   "younes",
					UserId: entryUser.CreatedBy.Int,
				}
				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel1)
				require.NoError(t, err, "error inserting organization into organization table")

				organizationModel2 := &model.CreateOrganization{
					Name:   "lawrence",
					UserId: entryUser.CreatedBy.Int,
				}
				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel2)
				require.NoError(t, err, "error adding the organization")

				err = tx.Commit()
				require.NoError(t, err, "unexpected error committing the transaction")
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
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) {
				tx, err := db.Begin()
				require.NoError(t, err, "unexpected error starting a transaction")
				defer func() {
					if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
						require.NoError(t, err, "error rolling back")
					}
				}()

				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err = entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: entryUser.CreatedBy.Int,
				}
				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				require.NoError(t, err, "error adding the organization")

				err = tx.Commit()
				require.NoError(t, err, "unexpected error committing the transaction")
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
				assert.True(t, len(respPaginated.Organizations) > 0, "unexpected empty categories")
				assert.True(t, respPaginated.Pagination.MaxRows > 0, "unexpected empty rows")
				assert.True(t, respPaginated.Pagination.RowCount > 0, "unexpected empty count")
				assert.True(t, len(respPaginated.Pagination.Pages) > 0, "unexpected empty pages")
			},
		},
		{
			name:            "empty_store",
			queryParameters: map[string]interface{}{},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container) {
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

				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(2),
					LastUpdatedBy:     null.IntFrom(2),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: entryUser.CreatedBy.Int,
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
				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(2),
					LastUpdatedBy:     null.IntFrom(2),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: entryUser.CreatedBy.Int,
				}

				_, err = modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
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

				entryUser := mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				}
				err := entryUser.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				organizationModel := mysqlmodel.Organization{
					ID:        1,
					Name:      "Demby",
					CreatedBy: entryUser.CreatedBy,
					IsActive:  false,
				}
				err = organizationModel.Insert(context.TODO(), db, boil.Infer())
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

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode, orgId, db)
		})
	}
}
