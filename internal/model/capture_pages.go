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

// CreateCapturePage struct for creating a new capture page
type CreateCapturePage struct {
	CapturePageSetId int    `json:"capture_page_set_id" validate:"required,greater_than_zero"`
	Name             string `json:"name" validate:"required"`
	//IsControl        int    `json:"is_control" validate:"required,oneof= 1 0"`
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

type CapturePages struct {
	Id                    int    `json:"id" boil:"id"`
	CapturePageSetId      int    `json:"capture_page_set_id" boil:"capture_page_set_id" swaggerignore:"true"`
	Name                  string `json:"name" boil:"name"`
	CapturePageType       string `json:"capture_page_type" boil:"capture_page_type"`
	CapturePagesIsControl int    `boil:"is_control" json:"is_control"`
}

// ToCapturePage converts the CreateCapturePage to a Capture Page.
func (c *CreateCapturePage) ToCapturePage() *CapturePages {
	capturepage := &CapturePages{
		Name:             c.Name,
		CapturePageSetId: c.CapturePageSetId,
		//CapturePagesIsControl: c.IsControl,
	}
	return capturepage
}

func (c *CreateCapturePage) Validate() error {
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("validate: %v", err)
	}
	return nil
}

type CapturePageType struct {
	Id   int    `json:"id" boil:"id"`
	Name string `json:"name" validate:"required" boil:"name"`
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
