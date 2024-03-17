package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltxhandler"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (m *Store) UpdateCategory(
	ctx context.Context,
	tx persistence.TransactionHandler,
	cat *model.Category,
) (*model.Category, error) {
	if cat == nil {
		return nil, ErrCatNil
	}
	ctxExec, err := mysqltxhandler.GetCtxExecutor(tx)
	// ctxExec, err := m.extractCtxExecutor(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	cat, err = m.updateCategory(ctx, ctxExec, cat)
	if err != nil {
		return nil, fmt.Errorf("update category: %v", err)
	}

	return cat, nil
}

func (m *Store) updateCategory(
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
func (m *Store) GetCategories(ctx context.Context, tx persistence.TransactionHandler, filters *model.CategoryFilters) (*model.PaginatedCategories, error) {
	ctxExec, err := mysqltxhandler.GetCtxExecutor(tx)
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
func (m *Store) getCategories(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CategoryFilters,
) (*model.PaginatedCategories, error) {
	var (
		paginated  model.PaginatedCategories
		pagination = model.NewPagination()
		res        []model.Category
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := make([]qm.QueryMod, 0)

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryWhere.ID.IN(filters.IdsIn))
		}
	}

	q := mysqlmodel.Categories(queryMods...)
	qCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get categories count: %v", err)
	}

	rows := pagination.Rows
	page := pagination.Page
	if filters != nil {
		rows = filters.Rows
		page = filters.Page
	}

	pagination.SetData(page, rows, int(qCount))

	queryMods = append(queryMods, qm.Limit(pagination.Rows), qm.Offset(pagination.Offset))
	q = mysqlmodel.Categories(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get categories: %v", err)
	}

	paginated.Categories = res
	paginated.Pagination = pagination

	return &paginated, nil
}
