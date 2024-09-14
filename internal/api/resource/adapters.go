package resource

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
)

type categoryGetter interface {
	GetByName(ctx context.Context, name string, userId int) (*model.Category, error)
}
