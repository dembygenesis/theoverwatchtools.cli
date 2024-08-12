package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
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
				tx, err := db.Begin()
				require.NoError(t, err, "unexpected error starting a transaction")
				defer tx.Rollback()

				organizationModel := &model.CreateOrganization{
					Name:   "demby",
					UserId: 1,
				}

				value, err := modules.OrganizationService.AddOrganization(context.Background(), organizationModel)
				require.NoError(t, err, "error adding the organization")
				fmt.Println("the value ---- ", strutil.GetAsJson(value))

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
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()

			require.NotEmpty(t, testCase.name, "unexpected empty test name")

			if testCase.queryParameters == nil {
				testCase.queryParameters = make(map[string]interface{})
			}

			handlers, _ := testCase.getContainer(t)

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       logger.New(context.TODO()),
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(context.TODO())
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				Logger:              logger.New(context.TODO()),
				TxHandler:           txHandlerDb,
			}

			testCase.mutations(t, db, handlers)

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := strutil.AppendQueryToURL("/api/v1/organization", testCase.queryParameters)
			fmt.Printf("Constructed URL: %s\n", url)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			fmt.Printf("Request: %+v\n", req)

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")

			var responseMap map[string]interface{}
			if err := json.Unmarshal(respBytes, &responseMap); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			fmt.Println("Response --------------- ", strutil.GetAsJson(responseMap))

			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}
