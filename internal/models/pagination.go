package models

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utils_common"
)

type Pagination struct {
	Pages      []int `json:"pages"`
	Rows       int   `mapstructure:"rows" validate:"required,int_greater_than_zero" json:"rows"`
	Page       int   `mapstructure:"page" validate:"required,int_greater_than_zero" json:"page"`
	TotalCount int   `json:"total_count"`
}

func NewPagination() *Pagination {
	return &Pagination{
		Rows: defaultPaginationRows,
		Page: defaultPaginationPage,
	}
}

func (p *Pagination) Validate() error {
	err := utils_common.ValidateStruct(p)
	if err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	if p.Rows > maxPaginationRows {
		return fmt.Errorf("max pagination rows exceeded: %v", p.Rows)
	}

	return err
}

func (p *Pagination) SetData(
	pages []int,
	rows int,
	page int,
	totalCount int,
) {
	p.Pages = pages
	p.Rows = rows
	p.Page = page
	p.TotalCount = totalCount
}
