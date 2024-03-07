package model

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
)

type Pagination struct {
	Pages      []int `json:"pages"`
	Rows       int   `mapstructure:"rows" validate:"is_positive" json:"rows"`
	Page       int   `mapstructure:"page" validate:"is_positive" json:"page"`
	TotalCount int   `json:"total_count"`
}

func (p *Pagination) SetPaginationDefaults() {
	if p.Page < 1 || p.Page == 0 {
		p.Page = defaultPaginationPage
	}
	if p.Rows < 1 || p.Rows == 0 {
		p.Rows = defaultPaginationRows
	}
}

func NewPagination() *Pagination {
	return &Pagination{
		Rows: defaultPaginationRows,
		Page: defaultPaginationPage,
	}
}

func (p *Pagination) ValidatePagination() error {
	// Check if negative values, and crash em, else all is fine.

	err := validationutils.Validate(p)
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
