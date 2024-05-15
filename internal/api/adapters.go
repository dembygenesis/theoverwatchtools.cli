package api

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . organizationService
type organizationService interface {
	ListOrganizations(ctx context.Context, filters *model.OrganizationFilters) (*model.PaginatedOrganizations, error)
	CreateOrganization(ctx context.Context, organization *model.CreateOrganization) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, params *model.DeleteOrganization) error
	RestoreOrganization(ctx context.Context, params *model.RestoreOrganization) error
}

//counterfeiter:generate . userService
type userService interface {
	ListUsers(ctx context.Context, filters *model.UserFilters) (*model.PaginatedUsers, error)
}
