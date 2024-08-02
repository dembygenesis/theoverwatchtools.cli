package model

import "github.com/volatiletech/null"

type UpdateOrganization struct {
	Id                    int         `json:"id" validate:"required,greater_than_zero"`
	OrganizationTypeRefId null.Int    `json:"organization_type_ref_id"`
	Name                  null.String `json:"name"`
}

type Organization struct {
	Id                    int    `json:"id" boil:"id"`
	OrganizationTypeRefId int    `json:"organization_type_ref_id" boil:"organization_type_ref_id" swaggerignore:"true"`
	Name                  string `json:"name" boil:"name"`
	OrganizationType      string `json:"organization_type" boil:"organization_type"`
}

type PaginatedOrganization struct {
	Organizations []Organization `json:"organizations"`
	Pagination    *Pagination    `json:"pagination"`
}

type OrganizationFilters struct {
	OrganizationNameIn     []string `query:"organization_name_in" json:"organization_name_in"`
	OrganizationTypeNameIn []string `query:"organization_type_name_in" json:"organization_type_name_in"`
	OrganizationTypeIdIn   []int    `query:"organization_type_id_in" json:"organization_type_id_in"`
	OrganizationIsActive   []int    `query:"is_active" json:"is_active"`
	IdsIn                  []int    `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}