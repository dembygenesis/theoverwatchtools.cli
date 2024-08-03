package model

import "github.com/volatiletech/null/v8"

type UpdateOrganization struct {
	Id   int         `json:"id" validate:"required,greater_than_zero"`
	Name null.String `json:"name"`
}

type Organization struct {
	Id   int    `json:"id" boil:"id"`
	Name string `json:"name" boil:"name"`
}

type PaginatedOrganization struct {
	Organizations []Organization `json:"organizations"`
	Pagination    *Pagination    `json:"pagination"`
}

type OrganizationFilters struct {
	OrganizationNameIn     []string `query:"organization_name_in" json:"organization_name_in"`
	OrganizationIsActive   []int    `query:"is_active" json:"is_active"`
	IdsIn                  []int    `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}
