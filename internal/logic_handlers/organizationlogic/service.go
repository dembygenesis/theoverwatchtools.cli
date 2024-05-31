package organizationlogic

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Config struct {
	TxProvider persistence.TransactionProvider `json:"tx_provider" validate:"required"`
	Logger     *logrus.Entry                   `json:"Logger" validate:"required"`
	Persistor  persistor                       `json:"Persistor" validate:"required"`
}

func (i *Config) Validate() error {
	return validationutils.Validate(i)
}

type Service struct {
	cfg *Config
}

func New(cfg *Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}
	return &Service{cfg}, nil
}

func (i *Service) validateOrganizationTypeId(ctx context.Context, handler persistence.TransactionHandler, id int) error {
	_, err := i.cfg.Persistor.GetOrganizationTypeById(ctx, handler, id)
	if err != nil {
		return fmt.Errorf("invalid organization_type_id: %v", err)
	}
	return nil
}

// ListOrganizations returns paginated organizations
func (i *Service) ListOrganizations(
	ctx context.Context,
	filter *model.OrganizationFilters,
) (*model.PaginatedOrganizations, error) {
	db, err := i.cfg.TxProvider.Db(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}

	paginated, err := i.cfg.Persistor.GetOrganizations(ctx, db, filter)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get organizations: %v", err),
		})
	}

	return paginated, nil
}

// CreateOrganization creates a new organization.
func (i *Service) CreateOrganization(ctx context.Context, params *model.CreateOrganization) (*model.Organization, error) {
	if err := params.Validate(); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("validate: %v", err),
		})
	}

	tx, err := i.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	organization := params.ToOrganization()

	_, err = i.cfg.Persistor.GetOrganizationTypeById(ctx, tx, organization.OrganizationTypeRefId)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid organization_type_id: %v", err),
		})
	}

	exists, err := i.cfg.Persistor.GetOrganizationByName(ctx, tx, organization.Name)
	if err != nil {
		if !strings.Contains(err.Error(), sysconsts.ErrExpectedExactlyOneEntry) {
			return nil, errs.New(&errs.Cfg{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("check organization unique: %v", err),
			})
		}
	}
	if exists != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("organization already exists"),
		})
	}

	organization, err = i.cfg.Persistor.CreateOrganization(ctx, tx, organization)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("create: %v", err),
		})
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit: %v", err),
		})
	}

	return organization, nil
}

// DeleteOrganization deletes a organization by ID.
func (s *Service) DeleteOrganization(ctx context.Context, params *model.DeleteOrganization) error {
	tx, err := s.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	fmt.Println("the params id ---- ", params.ID)
	err = s.cfg.Persistor.DeleteOrganization(ctx, tx, params.ID)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("delete organization: %v", err),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit transaction: %v", err),
		})
	}

	return nil
}

// RestoreOrganization restores a deleted organization by ID.
func (s *Service) RestoreOrganization(ctx context.Context, params *model.RestoreOrganization) error {
	tx, err := s.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	err = s.cfg.Persistor.RestoreOrganization(ctx, tx, params.ID)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("restore organization: %v", err),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit transaction: %v", err),
		})
	}

	return nil
}

// UpdateOrganization updates an existing organization.
func (i *Service) UpdateOrganization(ctx context.Context, params *model.UpdateOrganization) (*model.Organization, error) {
	tx, err := i.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if params.OrganizationTypeRefId.Valid {
		if err := i.validateOrganizationTypeId(ctx, tx, params.OrganizationTypeRefId.Int); err != nil {
			return nil, fmt.Errorf("organization_type_id: %w", err)
		}
	}

	category, err := i.cfg.Persistor.UpdateOrganization(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}

	return category, nil
}
