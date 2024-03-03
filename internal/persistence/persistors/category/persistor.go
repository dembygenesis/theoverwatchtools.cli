package category

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

type Persistor interface {
	// TODO:
	/*CreateCategory(ctx context.Context, tx persistence.TransactionHandler) error
	ReadCategory(ctx context.Context, id int, tx persistence.TransactionHandler) error
	DeleteCategory(ctx context.Context, tx persistence.TransactionHandler) error
	UpdateCategory(ctx context.Context, tx persistence.TransactionHandler) error*/

	GetCategories(
		ctx context.Context,
		tx persistence.TransactionHandler,
		filters *model.CategoryFilters,
	) ([]model.Category, error)
	UpdateCategory(ctx context.Context, tx persistence.TransactionHandler, cat *model.Category) (*model.Category, error)
	/*CreateCategories(ctx context.Context, tx persistence.TransactionHandler) error
	DeleteCategories(ctx context.Context, tx persistence.TransactionHandler) error
	UpdateCategories(ctx context.Context, tx persistence.TransactionHandler) error*/
}
