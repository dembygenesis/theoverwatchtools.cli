package factory

import "github.com/dembygenesis/local.tools/internal/model"

var (
	MockPaginatedCategories = &model.PaginatedCategories{
		Categories: []model.Category{
			{
				Id:   1,
				Name: "test",
			},
			{
				Id:   2,
				Name: "test 2",
			},
		},
	}
)
