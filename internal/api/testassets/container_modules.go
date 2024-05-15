package testassets

import (
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
)

type Container struct {
	MySQLStore   *mysqlstore.Repository
	ConnProvider *mysqlconn.Provider
}
