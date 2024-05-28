package api

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . categoryService
type categoryService interface {
	ListCategories(ctx context.Context, filters *model.CategoryFilters) (*model.PaginatedCategories, error)
	CreateCategory(ctx context.Context, category *model.CreateCategory) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.UpdateCategory) (*model.Category, error)
	DeleteCategory(ctx context.Context, params *model.DeleteCategory) error
	RestoreCategory(ctx context.Context, params *model.RestoreCategory) error
}

//counterfeiter:generate . organizationService
type organizationService interface {
	ListOrganizations(ctx context.Context, filters *model.OrganizationFilters) (*model.PaginatedOrganizations, error)
	CreateOrganization(ctx context.Context, organization *model.CreateOrganization) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, params *model.DeleteOrganization) error
	RestoreOrganization(ctx context.Context, params *model.RestoreOrganization) error
}

//counterfeiter:generate . capturePagesService
type capturePagesService interface {
	ListCapturePages(ctx context.Context, filters *model.CapturePagesFilters) (*model.PaginatedCapturePages, error)
	CreateCapturePages(ctx context.Context, capturepage *model.CreateCapturePage) (*model.CapturePages, error)
	//DeleteCapturePages(ctx context.Context, params *model.DeleteCapturePages) error
	//RestoreCapturePages(ctx context.Context, params *model.RestoreCapturePages) error
}

//counterfeiter:generate . userService
type userService interface {
	ListUsers(ctx context.Context, filters *model.UserFilters) (*model.PaginatedUsers, error)
}
