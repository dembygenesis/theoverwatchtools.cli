package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetCategoryTypeById attempts to fetch the category.
func (m *Repository) GetCategoryTypeById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.CategoryType, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCategoryTypes(ctx, ctxExec, &model.CategoryTypeFilters{IdsIn: []int{id}})
	if err != nil {
		return nil, fmt.Errorf("read categories: %v", err)
	}

	if res.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, mysqlmodel.TableNames.CategoryType)
	}

	return &res.Categories[0], nil
}

// GetCategoryTypes attempts to fetch the category
// entries using the given transaction layer.
func (m *Repository) GetCategoryTypes(ctx context.Context, tx persistence.TransactionHandler, filters *model.CategoryTypeFilters) (*model.PaginatedCategoryTypes, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCategoryTypes(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read categories: %v", err)
	}

	return res, nil
}

// getCategoryTypes performs the actual sql-queries.
func (m *Repository) getCategoryTypes(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CategoryTypeFilters,
) (*model.PaginatedCategoryTypes, error) {
	var (
		paginated  model.PaginatedCategoryTypes
		pagination = model.NewPagination()
		res        = make([]model.CategoryType, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := make([]qm.QueryMod, 0)

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryTypeWhere.ID.IN(filters.IdsIn))
		}
	}

	q := mysqlmodel.CategoryTypes(queryMods...)
	qCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get categories count: %v", err)
	}

	page := pagination.Page
	maxRows := pagination.MaxRows
	if filters != nil {
		if filters.Page.Valid {
			page = filters.Page.Int
		}
		if filters.MaxRows.Valid {
			maxRows = filters.MaxRows.Int
		}
	}

	pagination.SetQueryBoundaries(page, maxRows, int(qCount))

	queryMods = append(queryMods, qm.Limit(pagination.MaxRows), qm.Offset(pagination.Offset))
	q = mysqlmodel.CategoryTypes(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get categories: %v", err)
	}

	pagination.RowCount = len(res)
	paginated.Categories = res
	paginated.Pagination = pagination

	return &paginated, nil
}
