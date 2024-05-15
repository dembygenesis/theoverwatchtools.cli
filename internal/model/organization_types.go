package model

type PaginatedOrganizationTypes struct {
	Organizations []OrganizationType `json:"organization_types"`
	Pagination    *Pagination        `json:"pagination"`
}

// OrganizationTypeFilters contains the organization filters.
type OrganizationTypeFilters struct {
	IdsIn                  []int `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}
