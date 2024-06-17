package model

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/volatiletech/null/v8"
	"strings"
	"time"
)

// ClickTracker represents the ClickTracker model
type ClickTracker struct {
	Id                  int         `json:"id" boil:"id"`
	Name                string      `json:"name" boil:"name" validate:"required"`
	ClickTrackerSetName string      `json:"click_tracker_set_name" boil:"click_tracker_set_name" validate:"required"`
	UrlName             null.String `json:"url_name" boil:"url_name" validate:"required"`
	RedirectUrl         null.String `json:"redirect_url" boil:"redirect_url" validate:"required"`
	Clicks              null.Int    `json:"clicks" boil:"clicks" validate:"min=0"`
	UniqueClicks        null.Int    `json:"unique_clicks" boil:"unique_clicks" validate:"min=0"`
	CreatedBy           int         `json:"created_by" boil:"created_by" validate:"required,greater_than_zero"`
	UpdatedBy           int         `json:"updated_by" boil:"updated_by" validate:"required,greater_than_zero"`
	ClickTrackerSetId   int         `json:"click_tracker_set_id" boil:"click_tracker_set_id" validate:"required,greater_than_zero"`
	CountryId           null.Int    `json:"country_id" boil:"country_id"`
	CreatedAt           time.Time   `json:"created_at" boil:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at" boil:"updated_at"`
	DeletedAt           null.Time   `json:"deleted_at" boil:"deleted_at"`
}

// CreateClickTracker struct for creating a new ClickTracker
type CreateClickTracker struct {
	Name              string      `json:"name" validate:"required"`
	UrlName           null.String `json:"url_name" validate:"required"`
	RedirectUrl       null.String `json:"redirect_url" validate:"required"`
	CreatedBy         int         `json:"created_by" validate:"required,greater_than_zero"`
	UpdatedBy         int         `json:"updated_by" validate:"required,greater_than_zero"`
	ClickTrackerSetId int         `json:"click_tracker_set_id" validate:"required,greater_than_zero"`
	CountryId         null.Int    `json:"country_id"`
}

// UpdateClickTracker struct for updating an existing ClickTracker
type UpdateClickTracker struct {
	Id           int         `json:"id" validate:"required,greater_than_zero"`
	Name         null.String `json:"name"`
	UrlName      null.String `json:"url_name"`
	RedirectUrl  null.String `json:"redirect_url"`
	Clicks       null.Int    `json:"clicks"`
	UniqueClicks null.Int    `json:"unique_clicks"`
	UpdatedBy    int         `json:"updated_by" validate:"required,greater_than_zero"`
	CountryId    null.Int    `json:"country_id"`
}

// DeleteClickTracker struct for deleting a ClickTracker
type DeleteClickTracker struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

// RestoreClickTracker struct for restoring a deleted ClickTracker
type RestoreClickTracker struct {
	ID int `json:"id" validate:"required,greater_than_zero"`
}

// ClickTrackerFilters contains the ClickTracker filters
type ClickTrackerFilters struct {
	IdsIn                 []int    `query:"ids_in" json:"ids_in"`
	NameIn                []string `query:"name_in" json:"name_in"`
	ClickTrackerSetNameIn []string `query:"click_tracker_set_name_in" json:"click_tracker_set_name_in"`
	UrlNameIn             []string `query:"url_name_in" json:"url_name_in"`
	RedirectUrlIn         []string `query:"redirect_url_in" json:"redirect_url_in"`
	ClicksIn              []int    `query:"clicks_in" json:"clicks_in"`
	UniqueClicksIn        []int    `query:"unique_clicks_in" json:"unique_clicks_in"`
	CreatedByIn           []int    `query:"created_by_in" json:"created_by_in"`
	UpdatedByIn           []int    `query:"updated_by_in" json:"updated_by_in"`
	ClickTrackerSetIdIn   []int    `query:"click_tracker_set_id_in" json:"click_tracker_set_id_in"`
	CountryIdIn           []int    `query:"country_id_in" json:"country_id_in"`
	PaginationQueryFilters
}

type PaginatedClickTrackers struct {
	ClickTrackers []ClickTracker `json:"click_trackers"`
	Pagination    *Pagination    `json:"pagination"`
}

// Validate validates the ClickTracker filters
func (c *ClickTrackerFilters) Validate() error {
	if err := c.ValidatePagination(); err != nil {
		return fmt.Errorf("pagination: %v", err)
	}
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("click tracker filters: %v", err)
	}
	return nil
}

// Validate method for CreateClickTracker
func (c *CreateClickTracker) Validate() error {
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("validate: %v", err)
	}
	return nil
}

// Validate method for UpdateClickTracker
func (c *UpdateClickTracker) Validate() error {
	var errList errs.List
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	hasAtLeastOneUpdateParameter := false
	if c.Name.Valid {
		if strings.TrimSpace(c.Name.String) != "" {
			hasAtLeastOneUpdateParameter = true
		} else {
			errList.Add(sysconsts.ErrEmptyName)
		}
	}

	if c.UrlName.Valid {
		if strings.TrimSpace(c.UrlName.String) != "" {
			hasAtLeastOneUpdateParameter = true
		} else {
			errList.Add(sysconsts.ErrInvalidUrlName)
		}
	}

	if c.RedirectUrl.Valid {
		if strings.TrimSpace(c.RedirectUrl.String) != "" {
			hasAtLeastOneUpdateParameter = true
		} else {
			errList.Add(sysconsts.ErrInvalidRedirectUrl)
		}
	}

	if c.Clicks.Valid {
		if c.Clicks.Int >= 0 {
			hasAtLeastOneUpdateParameter = true
		} else {
			errList.Add(sysconsts.ErrInvalidClicks)
		}
	}

	if c.UniqueClicks.Valid {
		if c.UniqueClicks.Int >= 0 {
			hasAtLeastOneUpdateParameter = true
		} else {
			errList.Add(sysconsts.ErrInvalidUniqueClicks)
		}
	}

	if !hasAtLeastOneUpdateParameter {
		return errors.New(sysconsts.ErrHasNotASingleValidateUpdateParameter)
	}

	return nil
}

// ToClickTracker converts the CreateClickTracker to a ClickTracker.
func (c *CreateClickTracker) ToClickTracker() *ClickTracker {
	return &ClickTracker{
		Name:              c.Name,
		UrlName:           c.UrlName,
		RedirectUrl:       c.RedirectUrl,
		CreatedBy:         c.CreatedBy,
		UpdatedBy:         c.UpdatedBy,
		ClickTrackerSetId: c.ClickTrackerSetId,
		CountryId:         c.CountryId,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// ToClickTracker converts the UpdateClickTracker to a ClickTracker.
func (c *UpdateClickTracker) ToClickTracker(existing *ClickTracker) *ClickTracker {
	if c.Name.Valid {
		existing.Name = c.Name.String
	}
	if c.UrlName.Valid {
		existing.UrlName = null.StringFrom(c.UrlName.String) // Convert string to null.String
	}
	if c.RedirectUrl.Valid {
		existing.RedirectUrl = null.StringFrom(c.RedirectUrl.String) // Convert string to null.String
	}
	if c.Clicks.Valid {
		existing.Clicks = null.IntFrom(c.Clicks.Int) // Convert int to null.Int
	}
	if c.UniqueClicks.Valid {
		existing.UniqueClicks = null.IntFrom(c.UniqueClicks.Int) // Convert int to null.Int
	}
	if c.CountryId.Valid {
		existing.CountryId = c.CountryId
	}
	existing.UpdatedBy = c.UpdatedBy
	existing.UpdatedAt = time.Now()
	return existing
}
