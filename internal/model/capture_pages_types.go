package model

type PaginatedCapturePagesTypes struct {
	CapturePages []CapturePageType `json:"capture_pages"`
	Pagination   *Pagination       `json:"pagination"`
}

// CapturePagesTypeFilters contains the category filters.
type CapturePagesTypeFilters struct {
	IdsIn                  []int `query:"ids_in" json:"ids_in"`
	PaginationQueryFilters `swaggerignore:"true"`
}
