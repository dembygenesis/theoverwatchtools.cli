package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

//// GetUsers attempts to fetch the users
//// entries using the given transaction layer.
//func (m *Repository) GetUsers(ctx context.Context, tx persistence.TransactionHandler, filters *model.UserFilters) (*model.PaginatedUsers, error) {
//	ctxExec, err := mysqltx.GetCtxExecutor(tx)
//	if err != nil {
//		return nil, fmt.Errorf("extract context executor: %v", err)
//	}
//
//	res, err := m.getUsers(ctx, ctxExec, filters)
//	if err != nil {
//		return nil, fmt.Errorf("read organizations: %v", err)
//	}
//
//	return res, nil
//}
//
//// getUsers performs the actual sql-queries
//// that fetches organization entries.
//func (m *Repository) getUsers(
//	ctx context.Context,
//	ctxExec boil.ContextExecutor,
//	filters *model.UserFilters,
//) (*model.PaginatedUsers, error) {
//	var (
//		paginated  model.PaginatedUsers
//		pagination = model.NewPagination()
//		res        = make([]model.User, 0)
//		err        error
//	)
//
//	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
//	defer cancel()
//
//	queryMods := []qm.QueryMod{
//		qm.InnerJoin(
//			fmt.Sprintf(
//				"%s ON %s.%s = %s.%s",
//				mysqlmodel.TableNames.OrganizationType,
//				mysqlmodel.TableNames.OrganizationType,
//				mysqlmodel.OrganizationTypeColumns.ID,
//				mysqlmodel.TableNames.Organization,
//				mysqlmodel.OrganizationColumns.OrganizationTypeRefID,
//			),
//		),
//		qm.Select(
//			fmt.Sprintf("%s.%s AS %s",
//				mysqlmodel.TableNames.Organization,
//				mysqlmodel.OrganizationColumns.ID,
//				mysqlmodel.OrganizationColumns.ID,
//			),
//			fmt.Sprintf("%s.%s AS %s",
//				mysqlmodel.TableNames.Organization,
//				mysqlmodel.OrganizationColumns.Name,
//				mysqlmodel.OrganizationColumns.Name,
//			),
//			fmt.Sprintf("%s.%s AS %s",
//				mysqlmodel.TableNames.Organization,
//				mysqlmodel.OrganizationColumns.IsActive,
//				mysqlmodel.OrganizationColumns.IsActive,
//			),
//			fmt.Sprintf("%s.%s AS %s",
//				mysqlmodel.TableNames.OrganizationType,
//				mysqlmodel.OrganizationTypeColumns.ID,
//				mysqlmodel.OrganizationColumns.OrganizationTypeRefID,
//			),
//			fmt.Sprintf("%s.%s AS %s",
//				mysqlmodel.TableNames.OrganizationType,
//				mysqlmodel.OrganizationTypeColumns.Name,
//				"organization_type",
//			),
//		),
//	}
//
//	if filters != nil {
//		if len(filters.IdsIn) > 0 {
//			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.ID.IN(filters.IdsIn))
//		}
//
//		if len(filters.UserTypeIdIn) > 0 {
//			queryMods = append(queryMods, mysqlmodel.OrganizationTypeWhere.ID.IN(filters.OrganizationTypeIdIn))
//		}
//
//		if len(filters.UserTypeNameIn) > 0 {
//			queryMods = append(queryMods, mysqlmodel.OrganizationTypeWhere.Name.IN(filters.OrganizationTypeNameIn))
//		}
//
//		if len(filters.UserIsActive) > 0 {
//			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.IsActive.IN(filters.OrganizationIsActive))
//		}
//
//		if len(filters.UserNameIn) > 0 {
//			queryMods = append(queryMods, mysqlmodel.OrganizationWhere.Name.IN(filters.OrganizationNameIn))
//		}
//	}
//
//	q := mysqlmodel.Organizations(queryMods...)
//	totalCount, err := q.Count(ctx, ctxExec)
//	if err != nil {
//		return nil, fmt.Errorf("get organizations count: %v", err)
//	}
//
//	page := pagination.Page
//	maxRows := pagination.MaxRows
//	if filters != nil {
//		if filters.Page.Valid {
//			page = filters.Page.Int
//		}
//		if filters.MaxRows.Valid {
//			maxRows = filters.MaxRows.Int
//		}
//	}
//
//	pagination.SetQueryBoundaries(page, maxRows, int(totalCount))
//
//	queryMods = append(queryMods, qm.Limit(pagination.MaxRows), qm.Offset(pagination.Offset))
//	q = mysqlmodel.Organizations(queryMods...)
//
//	if err = q.Bind(ctx, ctxExec, &res); err != nil {
//		return nil, fmt.Errorf("get organizations: %v", err)
//	}
//
//	pagination.RowCount = len(res)
//	paginated.Organizations = res
//	paginated.Pagination = pagination
//
//	return &paginated, nil
//}

func (m *Repository) UpdateUser(
	ctx context.Context,
	tx persistence.TransactionHandler,
	user *model.User,
) (*model.User, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("get ctx executor: %v", err)
	}

	// Update
	// Question: do I just map everything?
	entry, err := mysqlmodel.FindUser(ctx, ctxExec, user.Id)
	if err != nil {
		return nil, fmt.Errorf("find user: %v", err)
	}

	// Add defaults here? lol

	// First name
	if user.Firstname != "" {
		entry.Firstname = user.Firstname
	}

	// Last Name
	if user.Lastname != "" {
		entry.Lastname = user.Lastname
	}

	// Password
	if user.Password != "" {
		entry.Password = user.Password
	}

	// Gen
	if user.Gender.Valid {
		entry.Gender = user.Gender
	}

	// Addr
	if user.Address.Valid {
		entry.Address = user.Address
	}

	// Bday
	if user.Birthday.Valid {
		entry.Birthday = user.Birthday
	}

	fmt.Println("===== here here")

	if _, err = entry.Update(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("update: %v", err)
	}

	return nil, nil
}

func (m *Repository) CreateUser(
	ctx context.Context,
	tx persistence.TransactionHandler,
	user *model.User,
) (*model.User, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("get ctx exec: %v", err)
	}

	entry := mysqlmodel.User{
		CategoryTypeRefID: user.CategoryTypeRefId,
		Firstname:         user.Firstname,
		Lastname:          user.Lastname,
		Email:             user.Email,
		Password:          user.Password,
		CreatedBy:         user.CreatedBy,
		CreatedAt:         user.CreatedAt,
		LastUpdatedBy:     user.LastUpdatedById,
		IsActive:          user.IsActive,
		ResetToken:        user.ResetToken,
		Birthday:          user.Birthday,
		Gender:            user.Gender,
	}

	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert: %v", err)
	}

	user.Id = entry.ID
	user.Firstname = entry.Firstname
	user.Lastname = entry.Lastname
	user.Email = entry.Email
	user.Password = entry.Password
	user.ResetToken = entry.ResetToken
	user.Birthday = entry.Birthday
	user.Password = entry.Password
	user.Gender = entry.Gender
	user.CreatedAt = entry.CreatedAt
	user.CreatedBy = entry.CreatedBy
	user.LastUpdatedAt = entry.LastUpdatedAt

	// Change inner join here bro, please
	// user.LastUpdatedBy = entry.LastUpdatedBy
	user.IsActive = entry.IsActive

	return user, nil
}

// DeleteUser deletes all the users,
func (m *Repository) DeleteUser(
	ctx context.Context,
	tx persistence.TransactionHandler,
	id int,
) error {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return fmt.Errorf("get ctx exec: %v", err)
	}

	entry := mysqlmodel.User{ID: id, IsActive: false}
	if _, err := entry.Update(ctx, ctxExec, boil.Infer()); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

type UserMySQL struct {
	cfg *Config
}
