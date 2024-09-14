package factory

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"testing"
)

type UserFactory struct {
	t       *testing.T
	Subject *mysqlmodel.User
}

func NewUserFactory(t *testing.T) *UserFactory {
	return &UserFactory{
		t:       t,
		Subject: &mysqlmodel.User{},
	}
}

func (f *UserFactory) Insert(exec boil.ContextExecutor) {
	if f.Subject.Firstname == "" {
		f.Subject.Firstname = uuid.NewString()
	}
	if f.Subject.Lastname == "" {
		f.Subject.Lastname = uuid.NewString()
	}

	if f.Subject.Email == "" {
		f.Subject.Email = uuid.NewString()[:5] + "@test.com"
	}

	err := f.Subject.Insert(context.Background(), exec, boil.Infer())
	require.NoError(f.t, err, "unexpected insert error")
}
