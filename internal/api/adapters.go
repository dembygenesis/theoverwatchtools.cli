package api

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . categoryManager
type categoryManager interface {
	GetCategories(ctx context.Context, f *model.CategoryFilters) ([]model.Category, error)
}
