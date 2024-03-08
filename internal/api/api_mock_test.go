package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dembygenesis/local.tools/internal/api/apifakes"
	"github.com/dembygenesis/local.tools/internal/api/factory"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testCtx    = context.TODO()
	testPort   = 3000
	testLogger = logger.New(testCtx)
)

func TestGetCategories_Mock(t *testing.T) {
	fakeCategoryManager := &apifakes.FakeCategoryManager{}
	fakeCategoryManager.GetCategoriesReturns(factory.MockPaginatedCategories, nil)

	cfg := &Config{
		Port:            3000,
		CategoryManager: fakeCategoryManager,
	}
	api, err := New(cfg)
	require.NoError(t, err, "unexpected error instantiating api")
	require.NotNil(t, api, "unexpected api nil instance")

	api.Routes()

	reqB, err := json.Marshal(map[string]interface{}{
		"ids_in": []int{1, 2, 3},
	})
	require.NoError(t, err, "unexpected error marshalling parameters")

	req := httptest.NewRequest("GET", "/v1/category", bytes.NewBuffer(reqB))
	req.Header = map[string][]string{
		"Content-Type":    {"application/json"},
		"Accept-Encoding": {"gzip", "deflate", "br"},
	}

	resp, err := api.app.Test(req, 1)
	require.NoError(t, err, "unexpected error executing test")

	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "unexpected error reading response body")

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respCategories *model.PaginatedCategories
	err = json.Unmarshal(respBytes, &respCategories)
	require.NoError(t, err, "unexpected error unmarshalling respBytes")

	require.Equal(t, factory.MockPaginatedCategories, respCategories)
}
