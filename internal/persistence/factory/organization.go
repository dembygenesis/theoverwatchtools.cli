package factory

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"testing"
)

type OrganizationFactory struct {
	t       *testing.T
	Subject *mysqlmodel.Organization
}

func NewOrganizationFactory(t *testing.T) *OrganizationFactory {
	return &OrganizationFactory{
		t:       t,
		Subject: &mysqlmodel.Organization{},
	}
}

func (f *OrganizationFactory) Insert(exec boil.ContextExecutor) {
	if f.Subject.ID == 0 {
		f.Subject.Name = uuid.NewString()
	}

	if f.Subject.Name == "" {
		f.Subject.Name = uuid.NewString()
	}

	err := f.Subject.Insert(context.Background(), exec, boil.Infer())
	require.NoError(f.t, err, "unexpected insert error")
}
