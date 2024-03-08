package model

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
)

type Category struct {
	Id                int    `json:"id" boil:"id"`
	CategoryTypeRefId int    `json:"category_type_ref_id" boil:"category_type_ref_id"`
	Name              string `json:"name" validate:"required" boil:"name"`
}

type Categories []Category

type PaginatedCategories struct {
	Categories []Category  `json:"categories"`
	Pagination *Pagination `json:"pagination"`
}

// CategoryFilters contains the category filters.
type CategoryFilters struct {
	IdsIn []int `query:"ids_in" json:"ids_in"`
	Pagination
}

// Validate validates the pagination parameters,
// and the filters provided.
func (c *CategoryFilters) Validate() error {
	if err := c.ValidatePagination(); err != nil {
		return fmt.Errorf("pagination: %v", err)
	}

	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("category filters: %v", err)
	}

	return nil
}
