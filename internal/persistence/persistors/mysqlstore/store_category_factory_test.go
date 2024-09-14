package mysqlstore

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/factory"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"testing"
)

type testCaseFactoryGetCategoriesAssets struct {
	CategoryTypes []*mysqlmodel.CategoryType
	Categories    []*mysqlmodel.Category
	Users         []*mysqlmodel.User
}

type testCaseFactoryGetCategories struct {
	name        string
	getFilterFn func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters
	// mutations   func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssetsr)
	mutations  func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error)
}

func getTestCasesFactoryGetCategories() []testCaseFactoryGetCategories {
	return []testCaseFactoryGetCategories{
		{
			name: "success-filter-ids-in",
			getFilterFn: func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters {
				require.Len(t, assets.CategoryTypes, 1, "unexpected category types length")
				require.Len(t, assets.Categories, 1, "unexpected categories length")

				return &model.CategoryFilters{
					IdsIn: []int{assets.Categories[0].ID},
				}
			},
			mutations: func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets {
				categoryType := factory.NewCategoryTypeFactory(t)
				categoryType.Subject.Name = "Test category type"
				categoryType.Insert(db)

				category := factory.NewCategoryFactory(t)
				category.Insert(db, categoryType.Subject)

				return &testCaseFactoryGetCategoriesAssets{
					CategoryTypes: []*mysqlmodel.CategoryType{categoryType.Subject},
					Categories:    []*mysqlmodel.Category{category.Subject},
					Users:         nil,
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-filter-names-in",
			mutations: func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets {
				categoryType := factory.NewCategoryTypeFactory(t)
				categoryType.Subject.Name = "Test category type"
				categoryType.Insert(db)

				categories := make([]*mysqlmodel.Category, 0)
				for i := 0; i < 2; i++ {
					category := factory.NewCategoryFactory(t)
					category.Subject.Name = fmt.Sprintf("Cat-%v", i)
					category.Insert(db, categoryType.Subject)

					categories = append(categories, category.Subject)
				}

				return &testCaseFactoryGetCategoriesAssets{
					CategoryTypes: []*mysqlmodel.CategoryType{categoryType.Subject},
					Categories:    categories,
					Users:         nil,
				}
			},
			getFilterFn: func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters {
				require.Len(t, assets.CategoryTypes, 1, "unexpected category types length")
				require.Len(t, assets.Categories, 2, "unexpected categories length")

				return &model.CategoryFilters{
					CategoryNameIn: []string{assets.Categories[0].Name, assets.Categories[1].Name},
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 2, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 2, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-category-type-id-in",
			getFilterFn: func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters {
				require.Len(t, assets.CategoryTypes, 1, "unexpected category types length")
				require.Len(t, assets.Categories, 2, "unexpected categories length")

				return &model.CategoryFilters{
					CategoryTypeIdIn: []int{assets.CategoryTypes[0].ID},
				}
			},
			mutations: func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets {
				categoryType := factory.NewCategoryTypeFactory(t)
				categoryType.Subject.Name = "Test category type"
				categoryType.Insert(db)

				categories := make([]*mysqlmodel.Category, 0)
				for i := 0; i < 2; i++ {
					category := factory.NewCategoryFactory(t)
					category.Insert(db, categoryType.Subject)

					categories = append(categories, category.Subject)
				}

				return &testCaseFactoryGetCategoriesAssets{
					CategoryTypes: []*mysqlmodel.CategoryType{categoryType.Subject},
					Categories:    categories,
					Users:         nil,
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) > 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-category-type-name-in",
			getFilterFn: func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters {
				require.Len(t, assets.CategoryTypes, 1, "unexpected category types length")
				require.Len(t, assets.Categories, 2, "unexpected categories length")

				return &model.CategoryFilters{
					CategoryTypeNameIn: []string{assets.CategoryTypes[0].Name},
				}
			},
			mutations: func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets {
				categoryType := factory.NewCategoryTypeFactory(t)
				categoryType.Subject.Name = "User Types"
				categoryType.Insert(db)

				categories := make([]*mysqlmodel.Category, 0)
				for i := 0; i < 2; i++ {
					category := factory.NewCategoryFactory(t)
					category.Insert(db, categoryType.Subject)
					categories = append(categories, category.Subject)
				}

				return &testCaseFactoryGetCategoriesAssets{
					CategoryTypes: []*mysqlmodel.CategoryType{categoryType.Subject},
					Categories:    categories,
					Users:         nil,
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) > 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-multiple-filters",
			getFilterFn: func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters {
				require.Len(t, assets.CategoryTypes, 1, "unexpected category types length")
				require.Len(t, assets.Categories, 1, "unexpected categories length")

				return &model.CategoryFilters{
					CategoryTypeNameIn: []string{assets.CategoryTypes[0].Name},
					CategoryTypeIdIn:   []int{assets.CategoryTypes[0].ID},
					CategoryNameIn:     []string{assets.Categories[0].Name},
				}
			},
			mutations: func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets {
				categoryType := factory.NewCategoryTypeFactory(t)
				categoryType.Subject.Name = "User Types"
				categoryType.Insert(db)

				category := factory.NewCategoryFactory(t)
				category.Subject.Name = "Super Admin"
				category.Insert(db, categoryType.Subject)

				return &testCaseFactoryGetCategoriesAssets{
					CategoryTypes: []*mysqlmodel.CategoryType{categoryType.Subject},
					Categories:    []*mysqlmodel.Category{category.Subject},
					Users:         nil,
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-created-by-filter",
			getFilterFn: func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters {
				require.Len(t, assets.CategoryTypes, 1, "unexpected category types length")
				require.Len(t, assets.Categories, 1, "unexpected categories length")
				require.Len(t, assets.Users, 2, "unexpected users length")

				return &model.CategoryFilters{
					CreatedByIdIn: []int{assets.Users[0].ID},
				}
			},
			mutations: func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets {
				categoryType := factory.NewCategoryTypeFactory(t)
				categoryType.Subject.Name = "User Types"
				categoryType.Insert(db)

				category := factory.NewCategoryFactory(t)
				category.Subject.Name = "Super Admin"
				category.Insert(db, categoryType.Subject)

				user := factory.NewUserFactory(t)
				user.Subject.CategoryTypeRefID = category.Subject.ID
				user.Insert(db)

				category2 := factory.NewCategoryFactory(t)
				category2.Subject.CreatedBy = null.Int{
					Int:   user.Subject.ID,
					Valid: true,
				}
				category2.Subject.Name = "Admin"
				category2.Insert(db, categoryType.Subject)

				user2 := factory.NewUserFactory(t)
				user2.Subject.CategoryTypeRefID = category2.Subject.ID
				user2.Insert(db)

				return &testCaseFactoryGetCategoriesAssets{
					CategoryTypes: []*mysqlmodel.CategoryType{categoryType.Subject},
					Categories:    []*mysqlmodel.Category{category.Subject},
					Users:         []*mysqlmodel.User{user.Subject, user2.Subject},
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "empty-results",
			getFilterFn: func(t *testing.T, assets *testCaseFactoryGetCategoriesAssets) *model.CategoryFilters {
				return &model.CategoryFilters{
					CategoryTypeNameIn: []string{"Saul Goodman"},
					CategoryTypeIdIn:   []int{999},
					CategoryNameIn:     []string{"Super Admin"},
				}
			},
			mutations: func(t *testing.T, db *sqlx.DB) *testCaseFactoryGetCategoriesAssets {
				return &testCaseFactoryGetCategoriesAssets{
					CategoryTypes: []*mysqlmodel.CategoryType{},
					Categories:    []*mysqlmodel.Category{},
					Users:         nil,
				}
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 0, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
	}
}

func Test_Factory_GetCategories(t *testing.T) {
	for _, testCase := range getTestCasesFactoryGetCategories() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t, true)
			defer cleanup()

			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			filter := testCase.getFilterFn(t, testCase.mutations(t, db))
			paginated, err := m.GetCategories(testCtx, txHandlerDb, filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}
