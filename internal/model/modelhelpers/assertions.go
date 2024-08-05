package modelhelpers

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertNonEmptyCategories(t *testing.T, categories []model.Category) {
	for _, category := range categories {
		assert.NotEmpty(t, category.Id, "unexpected empty category id")
		assert.NotEmpty(t, category.Name, "unexpected empty category name")
		assert.NotEmpty(t, category.CategoryType, "unexpected empty category type")
	}
}

func AssertNonEmptyOrganizations(t *testing.T, organizations []model.Organization) {
	for _, organization := range organizations {
		assert.NotEmpty(t, organization.Id, "unexpected empty organization id")
		assert.NotEmpty(t, organization.Name, "unexpected empty organization name")
		assert.NotEmpty(t, organization.CreatedBy, "unexpected empty organization created_by")
	}
}
