package models

type Category struct {
	Id                int    `json:"id" boil:"id"`
	CategoryTypeRefId int    `json:"category_type_ref_id" boil:"category_type_ref_id"`
	Name              string `json:"name" validate:"required" boil:"name"`
}

type Categories []Category

// CategoryFilter contains the category filters.
type CategoryFilter struct {
	IdsIn []int
}
