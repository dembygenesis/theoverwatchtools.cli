package model

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type UpdateCapturePage struct {
	Id     int         `json:"id" validate:"required,greater_than_zero"`
	Name   null.String `json:"name"`
	UserId null.Int    `json:"userId"`
}

type CreateCapturePage struct {
	Name   string `json:"name" validate:"required"`
	UserId int    `json:"user_id"`
}

type CapturePage struct {
	Id               int       `json:"id" boil:"id"`
	Name             string    `json:"name" boil:"name"`
	Html             string    `json:"html" boil:"html"`
	Clicks           int       `json:"clicks" boil:"clicks"`
	CapturePageSetId int       `json:"capture_page_set_id" boil:"capture_page_set_id"`
	IsControl        bool      `json:"is_control" boil:"is_control"`
	Impressions      int       `json:"impressions" boil:"impressions"`
	LastImpressionAt time.Time `json:"last_impression_at" boil:"last_impression_at"`
	CreatedBy        string    `json:"created_by" boil:"created_by"`
	LastUpdatedBy    string    `json:"last_updated_by" boil:"last_updated_by"`
	CreatedAt        time.Time `json:"created_at" boil:"created_at"`
	LastUpdatedAt    null.Time `json:"last_updated_at" boil:"last_updated_at"`
	IsActive         bool      `json:"is_active" boil:"is_active"`
}

type PaginatedCapturePages struct {
	CapturePages []CapturePage `json:"capture_pages"`
	Pagination   *Pagination   `json:"pagination"`
}

type CapturePageFilters struct {
	CapturePageNameIn      []string  `query:"capture_page_name_in" json:"capture_page_name_in"`
	CapturePageIsActive    null.Bool `query:"is_active" json:"is_active"`
	IdsIn                  []int     `query:"ids_in" json:"ids_in"`
	CreatedBy              null.Int  `query:"created_by" json:"created_by"`
	LastUpdatedBy          null.Int  `query:"last_updated_by" json:"last_updated_by"`
	PaginationQueryFilters `swaggerignore:"true"`
}