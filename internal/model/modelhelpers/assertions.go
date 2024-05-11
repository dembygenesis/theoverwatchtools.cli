package modelhelpers

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertNonEmptyOrganizations(t *testing.T, organizations []model.Organization) {
	for _, organization := range organizations {
		assert.NotEmpty(t, organization.Id, "unexpected empty organization id")
		assert.NotEmpty(t, organization.Name, "unexpected empty organization name")
		assert.NotEmpty(t, organization.OrganizationType, "unexpected empty organization type")
	}
}
