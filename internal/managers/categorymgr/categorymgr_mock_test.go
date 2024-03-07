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
	fakeReturns := []model.Category{
		{
			Id:   1,
			Name: "test",
		},
		{
			Id:   2,
			Name: "test 2",
		},
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
	categories, err := cm.GetCategories(testCtx, filters)
	assert.NoError(t, err, "unexpected err fetching categories")
	assert.NotNil(t, categories, "unexpected nil categories")
	assert.True(t, len(categories) > 0, "unexpected categories to have a length less that 0")
}
