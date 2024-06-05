package mysqlstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// DropCategoryTable drops the category table (for testing purposes).
func (m *Repository) DropCategoryTable(
	ctx context.Context,
	tx persistence.TransactionHandler,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("extract context executor: %v", err)
	}

	dropStmts := []string{
		"SET FOREIGN_KEY_CHECKS = 0;",
		"DROP TABLE category;",
		"SET FOREIGN_KEY_CHECKS = 1;",
	}

	for _, stmt := range dropStmts {
		if _, err := queries.Raw(stmt).Exec(ctxExec); err != nil {
			return fmt.Errorf("dropping category table sql stmt: %v", err)
		}
	}

	return nil
}

func (m *Repository) UpdateCategory(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateCategory) (*model.Category, error) {
	if params == nil {
		return nil, ErrCatNil
	}
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	entry := &mysqlmodel.Category{ID: params.Id}
	cols := []string{mysqlmodel.CategoryColumns.ID}

	fmt.Println("====== entry UpdateCategory:", strutil.GetAsJson(entry))
	fmt.Println("====== params UpdateCategory:", strutil.GetAsJson(params))

	if params.CategoryTypeRefId.Valid {
		entry.CategoryTypeRefID = params.CategoryTypeRefId.Int
		cols = append(cols, mysqlmodel.CategoryColumns.CategoryTypeRefID)
	}
	if params.Name.Valid {
		entry.Name = params.Name.String
		cols = append(cols, mysqlmodel.CategoryColumns.Name)
	}

	_, err = entry.Update(ctx, ctxExec, boil.Whitelist(cols...))
	tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("update failed: %v", err)
	}

	category, err := m.GetCategoryById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get category by id: %v", err)
	}

	return category, nil
}

// AddCategory attempts to add a new category
func (m *Repository) AddCategory(ctx context.Context, tx persistence.TransactionHandler, category *model.Category) (*model.Category, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	entry := &mysqlmodel.Category{
		Name:              category.Name,
		CategoryTypeRefID: category.CategoryTypeRefId,
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert category: %v", err)
	}

	category, err = m.GetCategoryById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get category by id: %v", err)
	}

	return category, nil
}

func (m *Repository) CreateCategory(ctx context.Context, tx persistence.TransactionHandler, category *model.Category) (*model.Category, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	entry := mysqlmodel.Category{
		Name:              category.Name,
		CategoryTypeRefID: category.CategoryTypeRefId,
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert category: %v", err)
	}

	category, err = m.GetCategoryById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get category by id: %v", err)
	}

	return category, nil
}

// GetCategoryById attempts to fetch the category.
func (m *Repository) GetCategoryById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.Category, error) {
	paginated, err := m.GetCategories(ctx, tx, &model.CategoryFilters{
		IdsIn: []int{id},
	})
	if err != nil {
		return nil, fmt.Errorf("category filtered by id: %v", err)
	}

	if paginated.Pagination.RowCount != 1 {

		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, id)
	}

	return &paginated.Categories[0], nil
}

// GetCategoryByName attempts to fetch the category.
func (m *Repository) GetCategoryByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.Category, error) {
	paginated, err := m.GetCategories(ctx, tx, &model.CategoryFilters{
		CategoryNameIn: []string{name},
	})
	if err != nil {
		return nil, fmt.Errorf("category filtered by name: %v", err)
	}

	if paginated.Pagination.RowCount != 1 {
		return nil, errors.New(sysconsts.ErrExpectedExactlyOneEntry)
	}

	return &paginated.Categories[0], nil
}

// GetCategories attempts to fetch the category
// entries using the given transaction layer.
func (m *Repository) GetCategories(ctx context.Context, tx persistence.TransactionHandler, filters *model.CategoryFilters) (*model.PaginatedCategories, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	fmt.Println("the filters cats --- ", strutil.GetAsJson(filters))

	res, err := m.getCategories(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read categories: %v", err)
	}

	return res, nil
}

// getCategories performs the actual sql-queries
// that fetches category entries.
func (m *Repository) getCategories(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CategoryFilters,
) (*model.PaginatedCategories, error) {
	var (
		paginated  model.PaginatedCategories
		pagination = model.NewPagination()
		res        = make([]model.Category, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := []qm.QueryMod{
		qm.InnerJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				mysqlmodel.TableNames.CategoryType,
				mysqlmodel.TableNames.CategoryType,
				mysqlmodel.CategoryTypeColumns.ID,
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

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryWhere.ID.IN(filters.IdsIn))
		}

		if len(filters.CategoryTypeIdIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryTypeWhere.ID.IN(filters.CategoryTypeIdIn))
			fmt.Println("the cat query mods dude --- ", queryMods)
		}

		if len(filters.CategoryTypeNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryTypeWhere.Name.IN(filters.CategoryTypeNameIn))
		}

		if len(filters.CategoryIsActive) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryWhere.IsActive.IN(filters.CategoryIsActive))
		}

		if len(filters.CategoryNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CategoryWhere.Name.IN(filters.CategoryNameIn))
		}
	}

	fmt.Println("the queryMods ----- ", queryMods)

	q := mysqlmodel.Categories(queryMods...)
	totalCount, err := q.Count(ctx, ctxExec)
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

	pagination.SetQueryBoundaries(page, maxRows, int(totalCount))

	queryMods = append(queryMods, qm.Limit(pagination.MaxRows), qm.Offset(pagination.Offset))
	q = mysqlmodel.Categories(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get categories: %v", err)
	}

	pagination.RowCount = len(res)
	paginated.Categories = res
	paginated.Pagination = pagination

	return &paginated, nil
}

// DeleteCategory delete the category.
func (m *Repository) DeleteCategory(
	ctx context.Context,
	tx persistence.TransactionHandler,
	id int,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("get ctx exec: %v", err)
	}

	entry := &mysqlmodel.Category{ID: id, IsActive: 0}
	if _, err = entry.Update(ctx, ctxExec, boil.Whitelist("is_active")); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// RestoreCategory restores a category.
func (m *Repository) RestoreCategory(
	ctx context.Context,
	tx persistence.TransactionHandler,
	id int,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("get ctx exec: %v", err)
	}

	entry := &mysqlmodel.Category{ID: id, IsActive: 1}
	if _, err = entry.Update(ctx, ctxExec, boil.Whitelist("is_active")); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	return nil
}
