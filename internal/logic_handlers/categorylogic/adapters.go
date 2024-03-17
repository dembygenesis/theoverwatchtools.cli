package categorylogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

type persistor interface {
	GetCategories(ctx context.Context, tx persistence.TransactionHandler, filters *model.CategoryFilters) (*model.PaginatedCategories, error)
}
