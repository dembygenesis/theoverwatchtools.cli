package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getTestCasesMiddleware() []testCaseUpdateCategory {
	return []testCaseUpdateCategory{
		{
			name: "success",
			fnGetApiCfg: func(t *testing.T) (*testApiCfg, *testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return &testApiCfg{catService: ctn.CategoryService, resourceGetter: ctn.ResourceGetter}, ctn, func() {
					cleanup()
				}
			},
			mutations: func(t *testing.T, ctn *testassets.Container) {

			},
			body: map[string]interface{}{
				"id":   1,
				"name": "demby",
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.Equal(t, http.StatusOK, respCode)
				fmt.Println("", string(resp))
			},
		},
	}
}

func Test_Middleware(t *testing.T) {
	for _, testCase := range getTestCasesMiddleware() {
		t.Run(testCase.name, func(t *testing.T) {
			ctn, _, cleanup := testCase.fnGetApiCfg(t)
			defer cleanup()

			cfg := &Config{
				BaseUrl:         testassets.MockBaseUrl,
				Port:            3000,
				CategoryService: ctn.catService,
				Logger:          logger.New(context.TODO()),
				Resource:        ctn.resourceGetter,
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.body)
			require.NoError(t, err, "unexpected error marshalling parameters")

			// req := httptest.NewRequest(http.MethodGet, "/api/v1/test/demby", bytes.NewBuffer(reqB))
			req := httptest.NewRequest(http.MethodGet, "/api/v1/test/Admin", bytes.NewBuffer(reqB))
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
