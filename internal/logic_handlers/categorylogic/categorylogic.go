package categorylogic

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
)

type Config struct {
	TxProvider persistence.TransactionProvider `json:"tx_provider" validate:"required"`
	Logger     *logrus.Entry                   `json:"logger" validate:"required"`
	Persistor  persistor                       `json:"persistor" validate:"required"`
}

func (i *Config) Validate() error {
	return validationutils.Validate(i)
}

type Impl struct {
	cfg *Config
}

func New(cfg *Config) (*Impl, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}
	return &Impl{cfg}, nil
}

// GetCategories returns paginated categories
func (i *Impl) GetCategories(
	ctx context.Context,
	filter *model.CategoryFilters,
) (*model.PaginatedCategories, error) {
	db, err := i.cfg.TxProvider.Db(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx provider: %v", err)
	}

	paginated, err := i.cfg.Persistor.GetCategories(ctx, db, filter)
	if err != nil {
		return nil, fmt.Errorf("get categories: %v", err)
	}

	return paginated, nil
}
