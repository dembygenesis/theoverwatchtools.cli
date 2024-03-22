package api

import (
	"bytes"
	"encoding/json"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCategories(t *testing.T) {
	handlers, cleanup := testassets.TestGetConcreteLogicHandlers(t)
	defer cleanup()

	cfg := &Config{
		Port:            3000,
		CategoryManager: handlers.Category,
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

	resp, err := api.app.Test(req, 100)
	require.NoError(t, err, "unexpected error executing test")

	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "unexpected error reading response body")

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respPaginated *model.PaginatedCategories
	err = json.Unmarshal(respBytes, &respPaginated)
	require.NoError(t, err, "unexpected error unmarshalling respBytes")

	assert.True(t, len(respPaginated.Categories) > 0, "unexpected empty categories")
	assert.True(t, respPaginated.Pagination.Rows > 0, "unexpected empty rows")
	assert.True(t, respPaginated.Pagination.TotalCount > 0, "unexpected empty count")
	assert.True(t, len(respPaginated.Pagination.Pages) > 0, "unexpected empty pages")
}
