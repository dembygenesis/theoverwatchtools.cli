package resource_test

import (
	"database/sql"
	"github.com/dembygenesis/local.tools/internal/api/resource"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewResourceGetter(t *testing.T) {
	ctn, cleanup := testassets.GetConcreteContainer(t)
	defer cleanup()

	rg, err := resource.New(ctn.CategoryService)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, rg, "unexpected nil")
}

type resourceCfg struct {
	// We need some shit
}

type testCaseResourceGetter struct {
	name            string
	mutations       func(t *testing.T, db *sql.DB)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesResourceGetter() []testCaseResourceGetter {
	return []testCaseResourceGetter{
		{
			name: "get-category",
			mutations: func(t *testing.T, db *sql.DB) {
				// Do something here
				// Create category using a factory function...
				// Now the pain starts HUHUH
			},
		},
	}
}
