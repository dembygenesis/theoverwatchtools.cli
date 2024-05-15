package model

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/volatiletech/null/v8"
	"time"
)

type PaginatedUser struct {
	Pagination *Pagination
	Users      []User
}

// User contains all relevant struct fields
type User struct {
	Id                int         `json:"id" boil:"id"`
	Firstname         string      `json:"firstname" boil:"firstname"`
	Lastname          string      `json:"lastname" boil:"lastname"`
	Email             string      `json:"email" boil:"email"`
	Password          string      `json:"password" boil:"password"`
	CategoryType      string      `json:"category_type" boil:"category_type"`
	CategoryTypeRefId int         `json:"category_type_ref_id" boil:"category_type_ref_id"`
	CreatedBy         null.Int    `json:"created_by" boil:"created_by"`
	LastUpdatedById   null.Int    `json:"last_updated_by" boil:"last_updated_by"`
	LastUpdatedBy     string      `json:"last_updated" boil:"last_updated"`
	CreatedAt         time.Time   `json:"created_at" boil:"created_at"`
	LastUpdatedAt     null.Time   `json:"last_updated_at" boil:"last_updated_at"`
	IsActive          bool        `json:"is_active" boil:"is_active"`
	ResetToken        null.String `json:"reset_token" boil:"reset_token"`
	Address           null.String `json:"address" boil:"address"`
	Birthday          null.Time   `json:"birthday" boil:"birthday"`
	Gender            null.String `json:"gender" boil:"gender"`
	IsSelfRegistered  null.Bool   `json:"is_self_registered" boil:"is_self_registered"`
}

// UserFilters contains the user filters.
type UserFilters struct {
	UserNameIn             []string `query:"user_name_in" json:"user_name_in"`
	UserTypeNameIn         []string `query:"user_type_name_in" json:"user_type_name_in"`
	UserTypeIdIn           []int    `query:"user_type_id_in" json:"user_type_id_in"`
	UserIsActive           []int    `query:"is_active" json:"is_active"`
	IdsIn                  []int    `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}

type PaginatedUsers struct {
	Users      []User      `json:"users"`
	Pagination *Pagination `json:"pagination"`
}

func (c *UserFilters) Validate() error {
	if err := c.ValidatePagination(); err != nil {
		return fmt.Errorf("pagination: %v", err)
	}
	if err := validationutils.Validate(c); err != nil {
		return fmt.Errorf("users filters: %v", err)
	}

	return nil
}

func (u *User) ValidateCreate() error {
	if u == nil {
		return errors.New(sysconsts.ErrStructNil)
	}

	return nil
}

type Users []User
