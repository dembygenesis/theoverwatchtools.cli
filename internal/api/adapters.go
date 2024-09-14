package api

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . categoryManager
type categoryManager interface {
	GetByName(ctx context.Context, name string, userId int) (*model.Category, error)
	ListCategories(ctx context.Context, filters *model.CategoryFilters) (*model.PaginatedCategories, error)
	CreateCategory(ctx context.Context, category *model.CreateCategory) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.UpdateCategory) (*model.Category, error)
	DeleteCategory(ctx context.Context, params *model.DeleteCategory) error
	RestoreCategory(ctx context.Context, params *model.RestoreCategory) error
}

//counterfeiter:generate . resourceManager
type resourceManager interface {
	Get(ctx context.Context, rType model.ResourceType, param string, userId int) (interface{}, error)
}
