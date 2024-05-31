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

// GetOrganizationTypeById attempts to fetch the organization.
func (m *Repository) GetOrganizationTypeById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.OrganizationType, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getOrganizationTypes(ctx, ctxExec, &model.OrganizationTypeFilters{IdsIn: []int{id}})
	if err != nil {
		return nil, fmt.Errorf("read organizations: %v", err)
	}

	if res.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, mysqlmodel.TableNames.OrganizationType)
	}

	return &res.Organizations[0], nil
}

// GetOrganizationTypes attempts to fetch the organization
// entries using the given transaction layer.
func (m *Repository) GetOrganizationTypes(ctx context.Context, tx persistence.TransactionHandler, filters *model.OrganizationTypeFilters) (*model.PaginatedOrganizationTypes, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getOrganizationTypes(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read organizations: %v", err)
	}

	return res, nil
}

// getOrganizationTypes performs the actual sql-queries.
func (m *Repository) getOrganizationTypes(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.OrganizationTypeFilters,
) (*model.PaginatedOrganizationTypes, error) {
	var (
		paginated  model.PaginatedOrganizationTypes
		pagination = model.NewPagination()
		res        = make([]model.OrganizationType, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := make([]qm.QueryMod, 0)

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.OrganizationTypeWhere.ID.IN(filters.IdsIn))
		}
	}

	q := mysqlmodel.OrganizationTypes(queryMods...)
	qCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get organizations count: %v", err)
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
	q = mysqlmodel.OrganizationTypes(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get organizations: %v", err)
	}

	pagination.RowCount = len(res)
	paginated.Organizations = res
	paginated.Pagination = pagination

	return &paginated, nil
}
