package model

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/volatiletech/null/v8"
)

type DeleteCapturePages struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

type RestoreCapturePages struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

// CapturePagesFilters contains the capture pages filters.
type CapturePagesFilters struct {
	CapturePagesNameIn     []string  `query:"capture_pages_name_in" json:"capture_pages"`
	CapturePagesTypeNameIn []string  `query:"capture_pages_type_name_in" json:"capture_pages_type_name_in"`
	CapturePagesTypeIdIn   []int     `query:"capture_page_type_id_in" json:"capture_page_type_id_in"`
	CapturePagesIsControl  null.Bool `query:"is_control" json:"is_control"`
	IdsIn                  []int     `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}

type CapturePagesType struct {
	Id   int    `json:"id" boil:"id"`
	Name string `json:"name" validate:"required" boil:"name"`
}

type CapturePages struct {
	Id               int    `json:"id" boil:"id"`
	CapturePageSetId int    `json:"capture_page_set_id" boil:"capture_page_set_id" swaggerignore:"true"`
	Name             string `json:"name" boil:"name"`
}

func (c *CapturePagesFilters) Validate() error {
	if err := c.ValidatePagination(); err != nil {
		return fmt.Errorf("pagination: %v", err)
	}
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("capture pages filters: %v", err)
	}

	return nil
}

type PaginatedCapturePages struct {
	CapturePages []CapturePages `json:"capture_pages"`
	Pagination   *Pagination    `json:"pagination"`
}
