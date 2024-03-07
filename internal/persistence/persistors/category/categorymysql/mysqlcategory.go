package categorymysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/databases/mysql/mysqltxhandler"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"sync"
)

var (
	ErrCatNil = errors.New("category provided is nil")
)

type Config struct {
	Logger        *logrus.Entry              `json:"logger" validate:"required"`
	QueryTimeouts *persistence.QueryTimeouts `json:"query_timeouts" validate:"required"`
}

func (m *Config) Validate() error {
	return validationutils.Validate(m)
}

// MySQLCategory is the main struct that
// handles crud operations for MYSQL.
type MySQLCategory struct {
	cfg *Config
}

// New creates a new instance of a MySQLCategory.
func New(cfg *Config) (*MySQLCategory, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %v", err)
	}

	return &MySQLCategory{cfg}, nil
}

// extractCtxExecutor attempts to extract a context executor
// from a persistence.TransactionHandler.
func (m *MySQLCategory) extractCtxExecutor(
	ctx context.Context,
	tx persistence.TransactionHandler,
) (boil.ContextExecutor, error) {
	h, err := mysqltxhandler.ParseExecutorInstance(tx)
	if err != nil {
		return nil, fmt.Errorf("get tx layer: %v", err)
	}

	ctxExec, err := h.GetCtxExecutor()
	if err != nil {
		return nil, fmt.Errorf("get context executor: %v", err)
	}

	return ctxExec, nil
}

func (m *MySQLCategory) UpdateCategory(
	ctx context.Context,
	tx persistence.TransactionHandler,
	cat *model.Category,
) (*model.Category, error) {
	if cat == nil {
		return nil, ErrCatNil
	}
	ctxExec, err := m.extractCtxExecutor(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	cat, err = m.updateCategory(ctx, ctxExec, cat)
	if err != nil {
		return nil, fmt.Errorf("update category: %v", err)
	}

	return cat, nil
}

var (
	mut sync.Mutex
)

func (m *MySQLCategory) updateCategory(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	cat *model.Category,
) (*model.Category, error) {
	entry := mysqlmodel.Category{ID: cat.Id, Name: cat.Name, CategoryTypeRefID: cat.CategoryTypeRefId}
	if _, err := entry.Update(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("update failed: %v", err)
	}
	return cat, nil
}

// GetCategories attempts to fetch the category
// entries using the given transaction layer.
func (m *MySQLCategory) GetCategories(
	ctx context.Context,
	tx persistence.TransactionHandler,
	filters *model.CategoryFilters,
) ([]model.Category, error) {
	ctxExec, err := m.extractCtxExecutor(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCategories(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read categories: %v", err)
	}

	return res, nil
}

// getCategories performs the actual sql-queries
// that fetches category entries.
func (m *MySQLCategory) getCategories(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CategoryFilters,
) ([]model.Category, error) {
	var (
		res []model.Category
		err error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := make([]qm.QueryMod, 0)

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryWhere.ID.IN(filters.IdsIn))
		}
	}

	if err = mysqlmodel.Categories(queryMods...).Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get categories: %v", err)
	}

	return res, nil
}
