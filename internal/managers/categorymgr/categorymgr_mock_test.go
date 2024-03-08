package categorymgr

import (
	"github.com/dembygenesis/local.tools/internal/managers/categorymgr/categorymgrfakes"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/persistencefakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCategoryManager_GetCategories(t *testing.T) {
	fakePersistor := &categorymgrfakes.FakeCategoryPersistor{}
	fakeReturns := &model.PaginatedCategories{
		Categories: []model.Category{
			{
				Id:   1,
				Name: "test",
			},
			{
				Id:   2,
				Name: "test 2",
			},
		},
		Pagination: nil,
	}

	fakePersistor.GetCategoriesReturns(fakeReturns, nil)
	fakeTxProvider := &persistencefakes.FakeTransactionProvider{}

	cm, err := New(&Config{
		Persistor:  fakePersistor,
		TxProvider: fakeTxProvider,
	})
	require.NoError(t, err, "unexpected err instantiating category manager")
	require.NotNil(t, cm, "unexpected nil category manager")

	filters := &model.CategoryFilters{}
	paginatedCategories, err := cm.GetCategories(testCtx, filters)
	assert.NoError(t, err, "unexpected err fetching categories")
	assert.NotNil(t, paginatedCategories.Categories, "unexpected nil categories")
	assert.True(t, len(paginatedCategories.Categories) > 0, "unexpected categories to have a length less that 0")
}
