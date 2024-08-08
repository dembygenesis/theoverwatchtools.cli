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
	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (m *Repository) DeleteOrganization(
	ctx context.Context,
	tx persistence.TransactionHandler,
	id int,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("get ctx exec: %w", err)
	}

	entry := &mysqlmodel.Organization{
		ID:       id,
		IsActive: false,
	}
	if _, err = entry.Update(ctx, ctxExec, boil.Whitelist(mysqlmodel.OrganizationColumns.IsActive)); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (m *Repository) UpdateOrganization(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateOrganization) (*model.Organization, error) {
	if params == nil {
		return nil, ErrOrgNil
	}
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	entry := &mysqlmodel.Organization{ID: params.Id}
	cols := []string{mysqlmodel.OrganizationColumns.ID}

	if params.Name.Valid {
		entry.Name = params.Name.String
		cols = append(cols, mysqlmodel.OrganizationColumns.Name)
	}

	_, err = entry.Update(ctx, ctxExec, boil.Whitelist(cols...))
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	organization, err := m.GetOrganizationById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	return organization, nil
}

func (m *Repository) GetOrganizationById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.Organization, error) {
	paginated, err := m.GetOrganizations(ctx, tx, &model.OrganizationFilters{
		IdsIn: []int{id},
	})
	if err != nil {
		return nil, fmt.Errorf("organization filtered by id: %w", err)
	}

	if paginated.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, id)
	}

	return &paginated.Organizations[0], nil
}

func (m *Repository) GetOrganizationByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.Organization, error) {
	paginated, err := m.GetOrganizations(ctx, tx, &model.OrganizationFilters{
		OrganizationNameIn: []string{name},
	})
	if err != nil {
		return nil, fmt.Errorf("organization filtered by name: %w", err)
	}

	if paginated.Pagination.RowCount != 1 {
		return nil, errors.New(sysconsts.ErrExpectedExactlyOneEntry)
	}

	return &paginated.Organizations[0], nil
}

func (m *Repository) GetOrganizations(ctx context.Context, tx persistence.TransactionHandler, filters *model.OrganizationFilters) (*model.PaginatedOrganizations, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	res, err := m.getOrganizations(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read organizations: %w", err)
	}

	return res, nil
}

func (m *Repository) getOrganizations(ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.OrganizationFilters) (*model.PaginatedOrganizations, error) {

	var (
		paginated  model.PaginatedOrganizations
		pagination = model.NewPagination()
		res        = make([]model.Organization, 0)
		err        error
	)

	fmt.Println("the filters --- ", strutil.GetAsJson(filters))

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := []qm.QueryMod{
		qm.InnerJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				mysqlmodel.TableNames.User,
				mysqlmodel.TableNames.User,
				mysqlmodel.UserColumns.ID,
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.CreatedBy,
			),
		),
		qm.Select(
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.ID,
				mysqlmodel.OrganizationColumns.ID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.Name,
				mysqlmodel.OrganizationColumns.Name,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.CreatedBy,
				mysqlmodel.OrganizationColumns.CreatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.LastUpdatedBy,
				mysqlmodel.OrganizationColumns.LastUpdatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.CreatedAt,
				mysqlmodel.OrganizationColumns.CreatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.LastUpdatedAt,
				mysqlmodel.OrganizationColumns.LastUpdatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.OrganizationColumns.IsActive,
				mysqlmodel.OrganizationColumns.IsActive,
			),
		),
	}

	fmt.Println("the query mods --- ", queryMods)

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.ID.IN(filters.IdsIn))
		}

		if filters.OrganizationIsActive.Valid {
			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.IsActive.EQ(filters.OrganizationIsActive.Bool))
		}

		if len(filters.OrganizationNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.Name.IN(filters.OrganizationNameIn))
		}

		if filters.CreatedBy.Valid {
			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.CreatedBy.EQ(filters.CreatedBy))
		}

		if filters.LastUpdatedBy.Valid {
			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.LastUpdatedBy.EQ(filters.LastUpdatedBy))
		}
	}

	q := mysqlmodel.Organizations(queryMods...)
	totalCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get organizations count: %w", err)
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
	q = mysqlmodel.Organizations(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get organizations: %w", err)
	}

	pagination.RowCount = len(res)
	paginated.Organizations = res
	paginated.Pagination = pagination

	return &paginated, nil
}

func (m *Repository) CreateOrganization(ctx context.Context, tx persistence.TransactionHandler, organization *model.Organization) (*model.Organization, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	entry := mysqlmodel.Organization{
		Name: organization.Name,
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert organization: %w", err)
	}

	organization, err = m.GetOrganizationById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	return organization, nil
}

func (m *Repository) AddOrganization(ctx context.Context, tx persistence.TransactionHandler, organization *model.CreateOrganization) (*model.Organization, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	entry := &mysqlmodel.Organization{
		Name:          organization.Name,
		CreatedBy:     null.IntFrom(organization.UserId),
		LastUpdatedBy: null.IntFrom(organization.UserId),
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert organization: %w", err)
	}

	createdOrganization, err := m.GetOrganizationById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	return createdOrganization, nil
}

func (m *Repository) RestoreOrganization(
	ctx context.Context,
	tx persistence.TransactionHandler,
	id int,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("get ctx exec: %w", err)
	}

	entry := &mysqlmodel.Organization{ID: id, IsActive: true}
	if _, err = entry.Update(ctx, ctxExec, boil.Whitelist(mysqlmodel.OrganizationColumns.IsActive)); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	return nil
}
