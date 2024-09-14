package factory

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"testing"
)

type CategoryFactory struct {
	t       *testing.T
	Subject *mysqlmodel.Category
}

func NewCategoryFactory(t *testing.T) *CategoryFactory {
	return &CategoryFactory{
		t:       t,
		Subject: &mysqlmodel.Category{},
	}
}

func (f *CategoryFactory) Insert(exec boil.ContextExecutor, categoryType *mysqlmodel.CategoryType) {
	if f.Subject.CategoryTypeRefID == 0 {
		f.Subject.CategoryTypeRefID = categoryType.ID
	}

	if f.Subject.Name == "" {
		f.Subject.Name = uuid.NewString()
	}

	err := f.Subject.Insert(context.Background(), exec, boil.Infer())
	require.NoError(f.t, err, "unexpected insert error")
}
