package model

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/volatiletech/null/v8"
	"strings"
)

// ClickTrackerSet represents the ClickTrackerSet model
type ClickTrackerSet struct {
	ID                     int       `json:"id"`
	Name                   string    `json:"name"`
	UrlName                string    `json:"url_name"`
	CreatedBy              null.Int  `json:"created_by"`
	UpdatedBy              null.Int  `json:"updated_by"`
	OrganizationID         null.Int  `json:"organization_id"`
	AnalyticsNumberOfLinks null.Int  `json:"analytics_number_of_links"`
	AnalyticsLastUpdatedAt null.Time `json:"analytics_last_updated_at"`
	CreatedAt              null.Time `json:"created_at"`
	UpdatedAt              null.Time `json:"updated_at"`
	DeletedAt              null.Time `json:"deleted_at"`
}

// CreateClickTrackerSet struct for creating a new ClickTrackerSet
type CreateClickTrackerSet struct {
	Name                   string `json:"name" validate:"required"`
	UrlName                string `json:"url_name" validate:"required"`
	CreatedBy              int    `json:"created_by" validate:"required,greater_than_zero"`
	UpdatedBy              int    `json:"updated_by" validate:"required,greater_than_zero"`
	OrganizationID         int    `json:"organization_id" validate:"required,greater_than_zero"`
	AnalyticsNumberOfLinks int    `json:"analytics_number_of_links" validate:"min=0"`
}

// UpdateClickTrackerSet struct for updating an existing ClickTrackerSet
type UpdateClickTrackerSet struct {
	ID                     int    `json:"id" validate:"required,greater_than_zero"`
	Name                   string `json:"name"`
	UrlName                string `json:"url_name"`
	UpdatedBy              int    `json:"updated_by" validate:"required,greater_than_zero"`
	AnalyticsNumberOfLinks int    `json:"analytics_number_of_links"`
}

// DeleteClickTrackerSet struct for deleting a ClickTrackerSet
type DeleteClickTrackerSet struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

// RestoreClickTrackerSet struct for restoring a deleted ClickTrackerSet
type RestoreClickTrackerSet struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

// ClickTrackerSetFilters contains the ClickTrackerSet filters
type ClickTrackerSetFilters struct {
	IdsIn       []int    `query:"ids_in" json:"ids_in"`
	NameIn      []string `query:"name_in" json:"name_in"`
	UrlNameIn   []string `query:"url_name_in" json:"url_name_in"`
	CreatedByIn []int    `query:"created_by_in" json:"created_by_in"`
	UpdatedByIn []int    `query:"updated_by_in" json:"updated_by_in"`
	PaginationQueryFilters
}

type PaginatedClickTrackerSets struct {
	ClickTrackerSets []ClickTrackerSet `json:"click_tracker_sets"`
	Pagination       *Pagination       `json:"pagination"`
}

// Validate validates the ClickTrackerSet filters
func (c *ClickTrackerSetFilters) Validate() error {
	if err := c.ValidatePagination(); err != nil {
		return fmt.Errorf("pagination: %v", err)
	}
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("click tracker set filters: %v", err)
	}
	return nil
}

// Validate method for CreateClickTrackerSet
func (c *CreateClickTrackerSet) Validate() error {
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("validate: %v", err)
	}
	return nil
}

// Validate method for UpdateClickTrackerSet
func (c *UpdateClickTrackerSet) Validate() error {
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	hasAtLeastOneUpdateParameter := false
	if strings.TrimSpace(c.Name) != "" {
		hasAtLeastOneUpdateParameter = true
	}

	if strings.TrimSpace(c.UrlName) != "" {
		hasAtLeastOneUpdateParameter = true
	}

	if c.AnalyticsNumberOfLinks != 0 {
		hasAtLeastOneUpdateParameter = true
	}

	if !hasAtLeastOneUpdateParameter {
		return errors.New(sysconsts.ErrHasNotASingleValidateUpdateParameter)
	}

	return nil
}
