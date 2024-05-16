package model

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/volatiletech/null/v8"
)

type UpdateOrganization struct {
	Id                    int         `json:"id" validate:"required,greater_than_zero"`
	OrganizationTypeRefId null.Int    `json:"organization_type_ref_id"`
	Name                  null.String `json:"name"`
}

type DeleteOrganization struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

type RestoreOrganization struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

// OrganizationFilters contains the organization filters.
type OrganizationFilters struct {
	OrganizationNameIn     []string `query:"organization_name_in" json:"organization_name_in"`
	OrganizationTypeNameIn []string `query:"organization_type_name_in" json:"organization_type_name_in"`
	OrganizationTypeIdIn   []int    `query:"organization_type_id_in" json:"organization_type_id_in"`
	OrganizationIsActive   []int    `query:"is_active" json:"is_active"`
	IdsIn                  []int    `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}

type Organization struct {
	Id                    int    `json:"id" boil:"id"`
	OrganizationTypeRefId int    `json:"organization_type_ref_id" boil:"organization_type_ref_id" swaggerignore:"true"`
	Name                  string `json:"name" boil:"name"`
	OrganizationType      string `json:"organization_type" boil:"organization_type"`
}

func (c *OrganizationFilters) Validate() error {
	if err := c.ValidatePagination(); err != nil {
		return fmt.Errorf("pagination: %v", err)
	}
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("organizations filters: %v", err)
	}

	return nil
}

// CreateOrganization struct for creating a new organization
type CreateOrganization struct {
	OrganizationTypeRefId int    `json:"organization_type_ref_id" validate:"required,greater_than_zero"`
	Name                  string `json:"name" validate:"required"`
}

// ToOrganization converts the CreateOrganization to an Organization.
func (c *CreateOrganization) ToOrganization() *Organization {
	organization := &Organization{
		Name:                  c.Name,
		OrganizationTypeRefId: c.OrganizationTypeRefId,
	}
	return organization
}

func (c *CreateOrganization) Validate() error {
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("validate: %v", err)
	}
	return nil
}

type OrganizationType struct {
	Id   int    `json:"id" boil:"id"`
	Name string `json:"name" validate:"required" boil:"name"`
}

type Organizations []Organization

type PaginatedOrganizations struct {
	Organizations []Organization `json:"organizations"`
	Pagination    *Pagination    `json:"pagination"`
}
