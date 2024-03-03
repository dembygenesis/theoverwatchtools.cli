package categorymgr

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . categoryPersistor
type categoryPersistor interface {
	GetCategories(
		ctx context.Context,
		tx persistence.TransactionHandler,
		filters *model.CategoryFilters,
	) ([]model.Category, error)
}
