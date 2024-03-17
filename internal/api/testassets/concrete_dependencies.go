package testassets

import (
	"github.com/dembygenesis/local.tools/di/ctn/dic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/stretchr/testify/require"
	"testing"
)

type Factory struct {
	LogicHandlers LogicHandlers
}

type LogicHandlers struct {
	Category *categorylogic.Impl
}

func TestGetConcreteLogicHandlers(t *testing.T) (*LogicHandlers, mysqlhelper.CleanFn) {
	db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	ctn, err := dic.NewContainer()
	require.NoError(t, err, "unexpected error instantiating a new container")

	mySQLTxHandler, err := ctn.SafeGetTxHandler()
	require.NoError(t, err, "unexpected error: SafeGetTxHandler")

	err = mySQLTxHandler.Set(db)
	require.Nil(t, err, "unexpected error for setting mysqlTxHandler.Set")

	category, err := ctn.SafeGetLogicCategory()
	require.NoError(t, err, "unexpected error: SafeGetLogicCategory")

	return &LogicHandlers{Category: category}, cleanup
}

func TestGetMockLogicHandlers(t *testing.T) {

}
