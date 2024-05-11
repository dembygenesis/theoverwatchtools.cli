package model

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
)

type DeleteCapturePages struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

type RestoreCapturePages struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

// CapturePagesFilters contains the capture pages filters.
type CapturePagesFilters struct {
	CapturePagesNameIn     []string `query:"capture_pages_name_in" json:"capture_pages"`
	CapturePagesTypeNameIn []string `query:"capture_pages_type_name_in" json:"capture_pages_type_name_in"`
	CapturePagesTypeIdIn   []int    `query:"capture_pages_type_id_in" json:"capture_pages_type_id_in"`
	CapturePagesIsControl  bool     `query:"is_control" json:"is_control"`
	IdsIn                  []int    `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}

type CapturePagesType struct {
	Id   int    `json:"id" boil:"id"`
	Name string `json:"name" validate:"required" boil:"name"`
}

type CapturePages struct {
	Id                    int    `json:"id" boil:"id"`
	CapturePagesTypeRefId int    `json:"capture_pages_type_ref_id" boil:"capture_pages_type_ref_id" swaggerignore:"true"`
	Name                  string `json:"name" boil:"name"`
	CapturePagesType      string `json:"capture_pages_type" boil:"capture_pages_type"`
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
