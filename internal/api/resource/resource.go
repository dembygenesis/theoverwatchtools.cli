package resource

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errutil"
)

const (
	CategoryResource model.ResourceType = "resource_type_category"
)

type Provider struct {
	categoryGetter categoryGetter
}

func (a *Provider) Get(ctx context.Context, rType model.ResourceType, param string, userId int) (interface{}, error) {
	switch rType {
	case CategoryResource:
		if category, err := a.categoryGetter.GetByName(ctx, param, userId); err != nil {
			return nil, fmt.Errorf("get category by name: %s, user id: %v, %w", param, userId, err)
		} else {
			return category, nil
		}
	default:
		return nil, fmt.Errorf("unsupported Resource type: %v", rType)
	}
}

func New(categoryGetter categoryGetter) (*Provider, error) {
	var errs errutil.List
	if categoryGetter == nil {
		errs.Add("categoryGetter is required")
	}

	if errs.HasErrors() {
		return nil, fmt.Errorf("validate: %w", errs.Single())
	}

	return &Provider{
		categoryGetter: categoryGetter,
	}, nil
}
