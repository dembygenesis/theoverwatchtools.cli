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
