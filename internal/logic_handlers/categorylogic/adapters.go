package categorylogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . persistor
type persistor interface {
	GetCategories(ctx context.Context, tx persistence.TransactionHandler, filters *model.CategoryFilters) (*model.PaginatedCategories, error)
	CreateCategory(ctx context.Context, tx persistence.TransactionHandler, category *model.Category) (*model.Category, error)
	GetCategoryTypeById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.CategoryType, error)
	GetCategoryByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.Category, error)
	UpdateCategory(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateCategory) (*model.Category, error)
}
