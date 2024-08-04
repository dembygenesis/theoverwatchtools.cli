package model

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type UpdateOrganization struct {
	Id   int         `json:"id" validate:"required,greater_than_zero"`
	Name null.String `json:"name"`
}

type Organization struct {
	Id            int       `json:"id" boil:"id"`
	Name          string    `json:"name" boil:"name"`
	CreatedBy     null.Int  `json:"created_by" boil:"created_by"`
	LastUpdatedBy null.Int  `json:"last_updated_by" boil:"last_updated_by"`
	CreatedAt     time.Time `json:"created_at" boil:"created_at"`
	LastUpdatedAt null.Time `json:"last_updated_at" boil:"last_updated_at"`
	IsActive      bool      `json:"is_active" boil:"is_active"`
}

type PaginatedOrganizations struct {
	Organizations []Organization `json:"organizations"`
	Pagination    *Pagination    `json:"pagination"`
}

type OrganizationFilters struct {
	OrganizationNameIn     []string  `query:"organization_name_in" json:"organization_name_in"`
	OrganizationIsActive   null.Bool `query:"is_active" json:"is_active"` // Change to null.Bool
	IdsIn                  []int     `query:"ids_in" json:"ids_in"`
	CreatedBy              null.Int  `query:"created_by" json:"created_by"`           // New field
	LastUpdatedBy          null.Int  `query:"last_updated_by" json:"last_updated_by"` // New field
	PaginationQueryFilters `swaggerignore:"true"`
}
