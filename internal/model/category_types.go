package model

type PaginatedCategoryTypes struct {
	Categories []CategoryType `json:"category_types"`
	Pagination *Pagination    `json:"pagination"`
}

// CategoryTypeFilters contains the category filters.
type CategoryTypeFilters struct {
	IdsIn                  []int `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}
