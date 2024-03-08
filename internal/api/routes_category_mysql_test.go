package api

import (
	"encoding/json"
	"github.com/dembygenesis/local.tools/internal/api/factory"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/urlutil"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCategories_MySQL(t *testing.T) {
	mysqlCategoryManager, cleanup := factory.TestGetMySQLCategoryManager(t)
	defer cleanup()

	cfg := &Config{
		Port:            3000,
		CategoryManager: mysqlCategoryManager,
	}
	api, err := New(cfg)
	require.NoError(t, err, "unexpected error instantiating api")
	require.NotNil(t, api, "unexpected api nil instance")

	api.Routes()

	reqUrl := "/v1/category"
	reqUrl, err = urlutil.MapToQueryString(reqUrl, map[string][]string{
		"ids_in": {"1", "2", "3"},
	})
	require.NoError(t, err, "unexpected error  trying to build the query url")

	req := httptest.NewRequest("GET", reqUrl, nil)
	req.Header = map[string][]string{
		"Content-Type":    {"application/json"},
		"Accept-Encoding": {"gzip", "deflate", "br"},
	}

	resp, err := api.app.Test(req, 2000)
	require.NoError(t, err, "unexpected error executing test")

	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "unexpected error reading response body")

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respCategories *model.PaginatedCategories
	err = json.Unmarshal(respBytes, &respCategories)
	require.NoError(t, err, "unexpected error unmarshalling respBytes")

	require.True(t, len(respCategories.Categories) > 0, "unexpected empty respCategories")
	require.True(t, len(respCategories.Categories) == 3, "unexpected len not equal to 3")
}

func TestGetCategories_MySQL_Pagination(t *testing.T) {
	mysqlCategoryManager, cleanup := factory.TestGetMySQLCategoryManager(t)
	defer cleanup()

	cfg := &Config{
		Port:            3000,
		CategoryManager: mysqlCategoryManager,
	}
	api, err := New(cfg)
	require.NoError(t, err, "unexpected error instantiating api")
	require.NotNil(t, api, "unexpected api nil instance")

	api.Routes()

	reqUrl := "/v1/category"
	reqUrl, err = urlutil.MapToQueryString(reqUrl, map[string][]string{
		"rows": {"1"},
		"page": {"2"},
	})
	require.NoError(t, err, "unexpected error  trying to build the query url")

	req := httptest.NewRequest("GET", reqUrl, nil)
	req.Header = map[string][]string{
		"Content-Type":    {"application/json"},
		"Accept-Encoding": {"gzip", "deflate", "br"},
	}

	resp, err := api.app.Test(req, 2000)
	require.NoError(t, err, "unexpected error executing test")

	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "unexpected error reading response body")

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respCategories *model.PaginatedCategories
	err = json.Unmarshal(respBytes, &respCategories)
	require.NoError(t, err, "unexpected error unmarshalling respBytes")

	require.True(t, len(respCategories.Categories) > 0, "unexpected empty respCategories")
	require.True(t, len(respCategories.Categories) == 1, "unexpected len not equal to 3")

	require.Equal(t, len(respCategories.Pagination.Pages), 3, "unexpected len not equal to 3")
	require.Equal(t, 2, respCategories.Pagination.Page, "unexpected len not equal to 3")
}
