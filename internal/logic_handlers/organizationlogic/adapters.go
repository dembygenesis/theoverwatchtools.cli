package organizationlogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . persistor
type persistor interface {
	GetOrganizations(ctx context.Context, tx persistence.TransactionHandler, filters *model.OrganizationFilters) (*model.PaginatedOrganizations, error)
	CreateOrganization(ctx context.Context, tx persistence.TransactionHandler, organization *model.Organization) (*model.Organization, error)
	GetOrganizationById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.Organization, error)
	UpdateOrganization(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateOrganization) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, tx persistence.TransactionHandler, id int) error
	RestoreOrganization(ctx context.Context, tx persistence.TransactionHandler, id int) error
}
