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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (m *Repository) DropCapturePageTable(
	ctx context.Context,
	tx persistence.TransactionHandler,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("extract context executor: %v", err)
	}

	dropStmts := []string{
		"SET FOREIGN_KEY_CHECKS = 0;",
		"DROP TABLE capture_page;",
		"SET FOREIGN_KEY_CHECKS = 1;",
	}

	for _, stmt := range dropStmts {
		if _, err := queries.Raw(stmt).Exec(ctxExec); err != nil {
			return fmt.Errorf("dropping capture page table sql stmt: %v", err)
		}
	}

	return nil
}

func (m *Repository) DeleteCapturePage(
	ctx context.Context,
	tx persistence.TransactionHandler,
	id int,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("get ctx exec: %w", err)
	}

	entry := &mysqlmodel.CapturePage{
		ID:       id,
		IsActive: false,
	}
	if _, err = entry.Update(ctx, ctxExec, boil.Whitelist(mysqlmodel.CapturePageColumns.IsActive)); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (m *Repository) UpdateCapturePage(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateCapturePage) (*model.CapturePage, error) {
	if params == nil {
		return nil, ErrCapNil
	}
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	entry := &mysqlmodel.CapturePage{ID: params.Id}
	cols := []string{mysqlmodel.CapturePageColumns.ID}

	if params.Name.Valid {
		entry.Name = params.Name.String
		cols = append(cols, mysqlmodel.CapturePageColumns.Name)
	}

	_, err = entry.Update(ctx, ctxExec, boil.Whitelist(cols...))
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	capturePage, err := m.GetCapturePageById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get capture page by id: %w", err)
	}

	return capturePage, nil
}

func (m *Repository) GetCapturePages(ctx context.Context, tx persistence.TransactionHandler, filters *model.CapturePageFilters) (*model.PaginatedCapturePages, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	res, err := m.getCapturePages(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read capture pages: %w", err)
	}

	fmt.Println("the res at store cp concrete --- ", strutil.GetAsJson(res))

	return res, nil
}

func (m *Repository) getCapturePages(ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CapturePageFilters) (*model.PaginatedCapturePages, error) {

	var (
		paginated  model.PaginatedCapturePages
		pagination = model.NewPagination()
		res        = make([]model.CapturePage, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := []qm.QueryMod{
		qm.InnerJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				mysqlmodel.TableNames.User,
				mysqlmodel.TableNames.User,
				mysqlmodel.UserColumns.ID,
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.CreatedBy,
			),
		),
		qm.Select(
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.ID,
				mysqlmodel.CapturePageColumns.ID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.Name,
				mysqlmodel.CapturePageColumns.Name,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.CreatedBy,
				mysqlmodel.CapturePageColumns.CreatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.LastUpdatedBy,
				mysqlmodel.CapturePageColumns.LastUpdatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.CreatedAt,
				mysqlmodel.CapturePageColumns.CreatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.LastUpdatedAt,
				mysqlmodel.CapturePageColumns.LastUpdatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CapturePage,
				mysqlmodel.CapturePageColumns.IsActive,
				mysqlmodel.CapturePageColumns.IsActive,
			),
		),
	}

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			fmt.Println("the filters --- ", filters.IdsIn)
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.ID.IN(filters.IdsIn))
		}

		if filters.CapturePageIsActive.Valid {
			fmt.Println("the filters valid --- ", filters.CapturePageIsActive)
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.IsActive.EQ(filters.CapturePageIsActive.Bool))
		}

		if len(filters.CapturePageNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.Name.IN(filters.CapturePageNameIn))
		}

		if filters.CreatedBy.Valid {
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.CreatedBy.EQ(filters.CreatedBy))
		}

		if filters.LastUpdatedBy.Valid {
			queryMods = append(queryMods, mysqlmodel.CapturePageWhere.LastUpdatedBy.EQ(filters.LastUpdatedBy))
		}
	}

	q := mysqlmodel.CapturePages(queryMods...)
	totalCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get capture pages count: %w", err)
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

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get capture pages: %w", err)
	}

	pagination.RowCount = len(res)
	paginated.CapturePages = res
	paginated.Pagination = pagination

	return &paginated, nil
}

func (m *Repository) GetCapturePageById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.CapturePage, error) {
	paginated, err := m.GetCapturePages(ctx, tx, &model.CapturePageFilters{
		IdsIn: []int{id},
	})

	fmt.Println("the paginated --- ", strutil.GetAsJson(paginated))

	if err != nil {
		return nil, fmt.Errorf("capture page filtered by id: %w", err)
	}

	if paginated.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, id)
	}

	return &paginated.CapturePages[0], nil
}

func (m *Repository) AddCapturePage(ctx context.Context, tx persistence.TransactionHandler, capturePage *model.CreateCapturePage) (*model.CapturePage, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	entry := &mysqlmodel.CapturePage{
		Name:             capturePage.Name,
		CreatedBy:        null.IntFrom(capturePage.UserId),
		LastUpdatedBy:    null.IntFrom(capturePage.UserId),
		CapturePageSetID: capturePage.CapturePageSetId,
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert capture page: %w", err)
	}

	createdCapturePage, err := m.GetCapturePageById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get capture page by id: %w", err)
	}

	return createdCapturePage, nil
}

func (m *Repository) RestoreCapturePage(ctx context.Context, tx persistence.TransactionHandler, id int) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("get ctx exec: %w", err)
	}

	entry := &mysqlmodel.CapturePage{ID: id, IsActive: true}

	if _, err = entry.Update(ctx, ctxExec, boil.Whitelist(mysqlmodel.CapturePageColumns.IsActive)); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	return nil
}
