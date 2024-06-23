package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetCapturePageTypeById attempts to fetch the capture pages.
func (m *Repository) GetCapturePageTypeById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.CapturePageType, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCapturePagesTypes(ctx, ctxExec, &model.CapturePagesTypeFilters{IdsIn: []int{id}})
	if err != nil {
		return nil, fmt.Errorf("read capture pages: %v", err)
	}

	fmt.Println(strutil.GetAsJson(res))
	fmt.Println("this is res.Pagination.RowCount =>>>>>>>>>", strutil.GetAsJson(res.Pagination.RowCount))

	if res.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, mysqlmodel.TableNames.CapturePageSets)
	}

	return &res.CapturePages[0], nil
}

// GetCapturePageTypes attempts to fetch the capture pages
// entries using the given transaction layer.
func (m *Repository) GetCapturePageTypes(ctx context.Context, tx persistence.TransactionHandler, filters *model.CapturePagesTypeFilters) (*model.PaginatedCapturePagesTypes, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCapturePagesTypes(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read capture pages: %v", err)
	}

	return res, nil
}

// getCapturePagesTypes performs the actual sql-queries.
func (m *Repository) getCapturePagesTypes(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CapturePagesTypeFilters,
) (*model.PaginatedCapturePagesTypes, error) {
	var (
		paginated  model.PaginatedCapturePagesTypes
		pagination = model.NewPagination()
		res        = make([]model.CapturePageType, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := make([]qm.QueryMod, 0)

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.ID.IN(filters.IdsIn))
		}
	}

	q := mysqlmodel.CapturePages(queryMods...)
	qCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get capture pages count: %v", err)
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
	q = mysqlmodel.CapturePages(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get capture pages: %v", err)
	}

	pagination.RowCount = len(res)
	paginated.CapturePages = res
	paginated.Pagination = pagination

	return &paginated, nil
}
