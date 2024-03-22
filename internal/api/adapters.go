package api

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . categoryService
type categoryService interface {
	GetCategories(ctx context.Context, filters *model.CategoryFilters) (*model.PaginatedCategories, error)
}
