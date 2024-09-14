package factory

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"testing"
)

type CategoryTypeFactory struct {
	t       *testing.T
	Subject *mysqlmodel.CategoryType
}

func NewCategoryTypeFactory(t *testing.T) *CategoryTypeFactory {
	return &CategoryTypeFactory{
		t:       t,
		Subject: &mysqlmodel.CategoryType{},
	}
}

func (f *CategoryTypeFactory) Insert(exec boil.ContextExecutor) {
	if f.Subject.Name == "" {
		f.Subject.Name = uuid.NewString()
	}
	err := f.Subject.Insert(context.Background(), exec, boil.Infer())
	require.NoError(f.t, err, "unexpected insert error")
}
