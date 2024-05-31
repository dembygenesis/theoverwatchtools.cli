package testassets

import (
	"github.com/dembygenesis/local.tools/internal/logic_handlers/capturepageslogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/organizationlogic"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
)

type Container struct {
	CategoryService     *categorylogic.Service
	MySQLStore          *mysqlstore.Repository
	ConnProvider        *mysqlconn.Provider
	OrganizationService *organizationlogic.Service
	CapturePagesService *capturepageslogic.Service
}
