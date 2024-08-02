package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (m *Repository) DropOrganizationTable(ctx context.Context, tx persistence.TransactionHandler) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("extract context executor: %w", err)
	}

	dropStmts := []string{
		"SET FOREIGN_KEY_CHECKS = 0;",
		"DROP TABLE organization;",
		"SET FOREIGN_KEY_CHECKS = 1;",
	}

	for _, stmt := range dropStmts {
		if _, err := queries.Raw(stmt).Exec(ctxExec); err != nil {
			return fmt.Errorf("dropping organization table sql stmt: %v", err)
		}
	}

	return nil
}

func (m *Repository) UpdateOrganization(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateOrganization) (*model.Organization, error) {
	if params == nil {
		return nil, ErrOrgNil
	}
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	entry := &mysqlmodel.Organization{ID: params.Id}
	cols := []string{mysqlmodel.OrganizationColumns.ID}

	//if params.OrganizationTypeRefId.Valid {
	//    entry.OrganizationRefUsers() = params.OrganizationTypeRefId.Int
	//    cols = append(cols, mysqlmodel.OrganizationColumns.IsActive)
	//}

	if params.Name.Valid {
		entry.Name = params.Name.String
		cols = append(cols, mysqlmodel.OrganizationColumns.Name)
	}

	_, err = entry.Update(ctx, ctxExec, boil.Whitelist(cols...))
	if err != nil {
		return nil, fmt.Errorf("update failed: %v", err)
	}

	//organization, err :=

	//return oraganization, err
}

//func (m *Repository) GetOrganizationById(ctx context.Context, tx persistence.TransactionHandler) error {
//    paginated, err := m.Get
//}

func (m *Repository) GetOrganizations(ctx context.Context, tx persistence.TransactionHandler, filters *model.OrganizationFilters) (*model.PaginatedOrganization ,error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.
}

func (m *Repository) getOrganizations(ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.OrganizationFilters,) (*model.PaginatedOrganization, error) {

	var (
		paginated model.PaginatedOrganization
		pagination = model.NewPagination()
		res = make([]model.Organization, 0)
		err error()
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := []qm.QueryMod{
		qm.InnerJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				mysqlmodel.TableNames.Organization,
				mysqlmodel.CategoryTypeColumns.ID,
				mysqlmodel.Organization.ID,
				mysqlmodel.TableNames.Category,
				mysqlmodel.CategoryColumns.CategoryTypeRefID,
			),
		),
		qm.Select(
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Category,
				mysqlmodel.CategoryColumns.ID,
				mysqlmodel.CategoryColumns.ID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Category,
				mysqlmodel.CategoryColumns.Name,
				mysqlmodel.CategoryColumns.Name,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.Category,
				mysqlmodel.CategoryColumns.IsActive,
				mysqlmodel.CategoryColumns.IsActive,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CategoryType,
				mysqlmodel.CategoryTypeColumns.ID,
				mysqlmodel.CategoryColumns.CategoryTypeRefID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.CategoryType,
				mysqlmodel.CategoryTypeColumns.Name,
				"category_type",
			),
		),
	}
}