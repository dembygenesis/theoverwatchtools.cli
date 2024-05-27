package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetCapturePages attempts to fetch the capture pages
// entries using the given transaction layer.
func (m *Repository) GetCapturePages(ctx context.Context, tx persistence.TransactionHandler, filters *model.CapturePagesFilters) (*model.PaginatedCapturePages, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCapturePages(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read capture pages: %v", err)
	}

	return res, nil
}

// getCapturePages performs the actual sql-queries
// that fetches capture pages entries.
func (m *Repository) getCapturePages(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CapturePagesFilters,
) (*model.PaginatedCapturePages, error) {
	var (
		paginated  model.PaginatedCapturePages
		pagination = model.NewPagination()
		res        = make([]model.CapturePages, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := []qm.QueryMod{
		qm.InnerJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				mysqlmodel.TableNames.CapturePageSets,
				mysqlmodel.TableNames.CapturePageSets,
				mysqlmodel.CapturePageSetColumns.ID,
				mysqlmodel.TableNames.CapturePages,
				mysqlmodel.CapturePageColumns.CapturePageSetID,
			),
		),
		qm.Select(
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePages,
				mysqlmodel.CapturePageColumns.ID,
				mysqlmodel.CapturePageColumns.ID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePages,
				mysqlmodel.CapturePageColumns.Name,
				mysqlmodel.CapturePageColumns.Name,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePages,
				mysqlmodel.CapturePageColumns.IsControl,
				mysqlmodel.CapturePageColumns.IsControl,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePageSets,
				mysqlmodel.CapturePageSetColumns.ID,
				mysqlmodel.CapturePageColumns.CapturePageSetID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePageSets,
				mysqlmodel.CapturePageSetColumns.Name,
				"organization_type",
			),
		),
	}

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.ID.IN(filters.IdsIn))
		}

		if len(filters.CapturePagesTypeIdIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageSetWhere.ID.IN(filters.CapturePagesTypeIdIn))
		}

		if len(filters.CapturePagesTypeNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageSetWhere.Name.IN(filters.CapturePagesTypeNameIn))
		}

		//if filters.CapturePagesIsControl {
		//	queryMods = append(queryMods, mysqlmodel.CapturePageWhere.IsControl.EQ(filters.CapturePagesIsControl))
		//}

		if len(filters.CapturePagesIsControl) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.IsControl.IN(filters.CapturePagesIsControl))
		}

		if len(filters.CapturePagesNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.Name.IN(filters.CapturePagesNameIn))
		}
	}

	q := mysqlmodel.CapturePages(queryMods...)
	totalCount, err := q.Count(ctx, ctxExec)
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

	pagination.SetQueryBoundaries(page, maxRows, int(totalCount))

	queryMods = append(queryMods, qm.Limit(pagination.MaxRows), qm.Offset(pagination.Offset))
	q = mysqlmodel.CapturePages(queryMods...)

	fmt.Println("the q --- ", strutil.GetAsJson(q))

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get capture pages: %v", err)
	}

	pagination.RowCount = len(res)
	paginated.CapturePages = res
	paginated.Pagination = pagination

	fmt.Println("the res return --- ", strutil.GetAsJson(res))

	return &paginated, nil
}
