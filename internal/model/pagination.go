package model

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
)

type Pagination struct {
	Pages      []int `json:"pages"`
	Rows       int   `mapstructure:"rows" validate:"is_positive" json:"rows"`
	Page       int   `mapstructure:"page" validate:"is_positive" json:"page"`
	Offset     int   `json:"-"`
	TotalCount int   `json:"total_count"`
}

func (p *Pagination) SetData(page, rows, totalCount int) {
	if rows < 1 {
		p.Rows = defaultPaginationRows
	} else if rows > maxPaginationRows {
		p.Rows = maxPaginationRows
	} else {
		p.Rows = rows
	}

	// Ensure the page is within allowed limits
	if page < 1 {
		p.Page = defaultPaginationPage
	} else {
		p.Page = page
	}

	p.TotalCount = totalCount

	// Calculate the total number of pages
	totalPages := totalCount / p.Rows
	if totalCount%p.Rows > 0 {
		totalPages++
	}

	// Populate the Pages slice
	p.Pages = make([]int, totalPages)
	for i := range p.Pages {
		p.Pages[i] = i + 1
	}

	// Adjust the Page number to be within the calculated total pages
	if p.Page > totalPages {
		p.Page = 1 // Reset to page 1 if out of bounds
	}

	p.Offset = (p.Page - 1) * p.Rows
}

func (p *Pagination) GetOffset() int {
	return p.Offset
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
