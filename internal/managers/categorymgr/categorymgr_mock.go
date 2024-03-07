package categorymgr

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
)

type CategoryManager struct {
	cfg *Config
}

func (c *CategoryManager) GetCategories(ctx context.Context, f *model.CategoryFilters) ([]model.Category, error) {
	db, err := c.cfg.TxProvider.Db(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db: %v", err)
	}

	categories, err := c.cfg.Persistor.GetCategories(ctx, db, f)
	if err != nil {
		return nil, fmt.Errorf("persistor: %v", err)
	}

	return categories, nil
}

type Config struct {
	Persistor  categoryPersistor `json:"persistor" validate:"required"`
	TxProvider persistence.TransactionProvider
}

func (c *Config) Validate() error {
	return validationutils.Validate(c)
}

func New(cfg *Config) (*CategoryManager, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	return &CategoryManager{cfg: cfg}, nil
}
