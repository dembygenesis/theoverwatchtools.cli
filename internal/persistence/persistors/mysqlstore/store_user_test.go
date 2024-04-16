package mysqlstore

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"testing"
)

func TestUserMySQL_Create_Fail_Invalid_Fks_Created_By(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	txHandlerController, err := mysqltx.New(&mysqltx.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected non nil error")
	require.NotNil(t, txHandlerController, "unexpected nil txHandler")

	txHandler, err := txHandlerController.Db(testCtx)
	require.NoError(t, err, "unexpected non nil error")
	require.NotNil(t, txHandler, "unexpected nil txHandler")

	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}
	userMysql, err := New(cfg)
	require.NoError(t, err, "unexpected nil error")

	user := &model.User{
		Id:        0,
		Firstname: "Demby",
		Lastname:  "Abella",
		Email:     "dembygenesis@gmail.com",
		Password:  "bcrypted",
		CreatedBy: null.Int{
			Int:   3,
			Valid: true,
		},
	}

	user, err = userMysql.CreateUser(testCtx, txHandler, user)
	require.Error(t, err, "unexpected nil error creating a bogus created by FK")
	require.Nil(t, user, "unexpected non nil post user creation")
}

func TestUserMySQL_Create_Success(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	txHandlerController, err := mysqltx.New(&mysqltx.Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected non nil error")
	require.NotNil(t, txHandlerController, "unexpected nil txHandler")

	txHandler, err := txHandlerController.Db(testCtx)
	require.NoError(t, err, "unexpected non nil error")
	require.NotNil(t, txHandler, "unexpected nil txHandler")

	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}
	userMysql, err := New(cfg)
	require.NoError(t, err, "unexpected nil error")

	userTypeAdmin, err := mysqlmodel.Categories(
		qm.InnerJoin(fmt.Sprintf(
			"%s ct ON ct.id = %s.category_type_ref_id",
			mysqlmodel.TableNames.CategoryType,
			mysqlmodel.TableNames.Category,
		)),
	).One(testCtx, db)
	require.NoError(t, err, "unexpected non nil error")
	require.NotNil(t, userTypeAdmin, "unexpected nil userTypeAdmin")

	user := &model.User{
		Firstname:         "Demby",
		Lastname:          "Abella",
		Email:             "dembygenesis@gmail.com",
		Password:          "bcrypted",
		CategoryTypeRefId: userTypeAdmin.ID,
	}

	user, err = userMysql.CreateUser(testCtx, txHandler, user)
	require.NoError(t, err, "unexpected   error creating a bogus created by FK")
	require.NotNil(t, user, "unexpected non nil post user creation")
}
