package mysqlstore

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
)

// ConvertMysqlModelToCategory converts a mysql model category to a model category.
func ConvertMysqlModelToCategory(category *mysqlmodel.Category) *model.Category {
	if category == nil {
		return nil
	}
	return &model.Category{
		Id:                category.ID,
		CategoryTypeRefId: category.CategoryTypeRefID,
		Name:              category.Name,
	}
}
