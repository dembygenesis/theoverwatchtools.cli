package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	cli "github.com/dembygenesis/local.tools/internal/cli"
	config "github.com/dembygenesis/local.tools/internal/config"
	authlogic "github.com/dembygenesis/local.tools/internal/logic_handlers/authlogic"
	categorylogic "github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic"
	marketinglogic "github.com/dembygenesis/local.tools/internal/logic_handlers/marketinglogic"
	userlogic "github.com/dembygenesis/local.tools/internal/logic_handlers/userlogic"
	mysqltxhandler "github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltxhandler"
	mysqltxprovider "github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltxprovider"
	mysqlstore "github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	sqlx "github.com/jmoiron/sqlx"
	logrus "github.com/sirupsen/logrus"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "config_layer",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("config_layer")
				if err != nil {
					var eo *config.App
					return eo, err
				}
				b, ok := d.Build.(func() (*config.App, error))
				if !ok {
					var eo *config.App
					return eo, errors.New("could not cast build function to func() (*config.App, error)")
				}
				return b()
			},
			Unshared: false,
		},
		{
			Name:  "db_mysql",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("db_mysql")
				if err != nil {
					var eo *sqlx.DB
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *sqlx.DB
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *sqlx.DB
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				b, ok := d.Build.(func(*config.App) (*sqlx.DB, error))
				if !ok {
					var eo *sqlx.DB
					return eo, errors.New("could not cast build function to func(*config.App) (*sqlx.DB, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "logger_logrus",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logger_logrus")
				if err != nil {
					var eo *logrus.Entry
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *logrus.Entry
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *logrus.Entry
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				b, ok := d.Build.(func(*config.App) (*logrus.Entry, error))
				if !ok {
					var eo *logrus.Entry
					return eo, errors.New("could not cast build function to func(*config.App) (*logrus.Entry, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "logic_auth",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logic_auth")
				if err != nil {
					var eo *authlogic.Impl
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *authlogic.Impl
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *authlogic.Impl
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				pi1, err := ctn.SafeGet("logger_logrus")
				if err != nil {
					var eo *authlogic.Impl
					return eo, err
				}
				p1, ok := pi1.(*logrus.Entry)
				if !ok {
					var eo *authlogic.Impl
					return eo, errors.New("could not cast parameter 1 to *logrus.Entry")
				}
				pi2, err := ctn.SafeGet("tx_provider")
				if err != nil {
					var eo *authlogic.Impl
					return eo, err
				}
				p2, ok := pi2.(*mysqltxprovider.Impl)
				if !ok {
					var eo *authlogic.Impl
					return eo, errors.New("could not cast parameter 2 to *mysqltxprovider.Impl")
				}
				b, ok := d.Build.(func(*config.App, *logrus.Entry, *mysqltxprovider.Impl) (*authlogic.Impl, error))
				if !ok {
					var eo *authlogic.Impl
					return eo, errors.New("could not cast build function to func(*config.App, *logrus.Entry, *mysqltxprovider.Impl) (*authlogic.Impl, error)")
				}
				return b(p0, p1, p2)
			},
			Unshared: false,
		},
		{
			Name:  "logic_category",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logic_category")
				if err != nil {
					var eo *categorylogic.Impl
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *categorylogic.Impl
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *categorylogic.Impl
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				pi1, err := ctn.SafeGet("logger_logrus")
				if err != nil {
					var eo *categorylogic.Impl
					return eo, err
				}
				p1, ok := pi1.(*logrus.Entry)
				if !ok {
					var eo *categorylogic.Impl
					return eo, errors.New("could not cast parameter 1 to *logrus.Entry")
				}
				pi2, err := ctn.SafeGet("tx_provider")
				if err != nil {
					var eo *categorylogic.Impl
					return eo, err
				}
				p2, ok := pi2.(*mysqltxprovider.Impl)
				if !ok {
					var eo *categorylogic.Impl
					return eo, errors.New("could not cast parameter 2 to *mysqltxprovider.Impl")
				}
				pi3, err := ctn.SafeGet("persistence_mysql")
				if err != nil {
					var eo *categorylogic.Impl
					return eo, err
				}
				p3, ok := pi3.(*mysqlstore.Store)
				if !ok {
					var eo *categorylogic.Impl
					return eo, errors.New("could not cast parameter 3 to *mysqlstore.Store")
				}
				b, ok := d.Build.(func(*config.App, *logrus.Entry, *mysqltxprovider.Impl, *mysqlstore.Store) (*categorylogic.Impl, error))
				if !ok {
					var eo *categorylogic.Impl
					return eo, errors.New("could not cast build function to func(*config.App, *logrus.Entry, *mysqltxprovider.Impl, *mysqlstore.Store) (*categorylogic.Impl, error)")
				}
				return b(p0, p1, p2, p3)
			},
			Unshared: false,
		},
		{
			Name:  "logic_marketing",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logic_marketing")
				if err != nil {
					var eo *marketinglogic.Impl
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *marketinglogic.Impl
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *marketinglogic.Impl
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				pi1, err := ctn.SafeGet("logger_logrus")
				if err != nil {
					var eo *marketinglogic.Impl
					return eo, err
				}
				p1, ok := pi1.(*logrus.Entry)
				if !ok {
					var eo *marketinglogic.Impl
					return eo, errors.New("could not cast parameter 1 to *logrus.Entry")
				}
				pi2, err := ctn.SafeGet("tx_provider")
				if err != nil {
					var eo *marketinglogic.Impl
					return eo, err
				}
				p2, ok := pi2.(*mysqltxprovider.Impl)
				if !ok {
					var eo *marketinglogic.Impl
					return eo, errors.New("could not cast parameter 2 to *mysqltxprovider.Impl")
				}
				b, ok := d.Build.(func(*config.App, *logrus.Entry, *mysqltxprovider.Impl) (*marketinglogic.Impl, error))
				if !ok {
					var eo *marketinglogic.Impl
					return eo, errors.New("could not cast build function to func(*config.App, *logrus.Entry, *mysqltxprovider.Impl) (*marketinglogic.Impl, error)")
				}
				return b(p0, p1, p2)
			},
			Unshared: false,
		},
		{
			Name:  "logic_user",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logic_user")
				if err != nil {
					var eo *userlogic.Impl
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *userlogic.Impl
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *userlogic.Impl
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				pi1, err := ctn.SafeGet("logger_logrus")
				if err != nil {
					var eo *userlogic.Impl
					return eo, err
				}
				p1, ok := pi1.(*logrus.Entry)
				if !ok {
					var eo *userlogic.Impl
					return eo, errors.New("could not cast parameter 1 to *logrus.Entry")
				}
				pi2, err := ctn.SafeGet("tx_provider")
				if err != nil {
					var eo *userlogic.Impl
					return eo, err
				}
				p2, ok := pi2.(*mysqltxprovider.Impl)
				if !ok {
					var eo *userlogic.Impl
					return eo, errors.New("could not cast parameter 2 to *mysqltxprovider.Impl")
				}
				b, ok := d.Build.(func(*config.App, *logrus.Entry, *mysqltxprovider.Impl) (*userlogic.Impl, error))
				if !ok {
					var eo *userlogic.Impl
					return eo, errors.New("could not cast build function to func(*config.App, *logrus.Entry, *mysqltxprovider.Impl) (*userlogic.Impl, error)")
				}
				return b(p0, p1, p2)
			},
			Unshared: false,
		},
		{
			Name:  "persistence_mysql",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("persistence_mysql")
				if err != nil {
					var eo *mysqlstore.Store
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *mysqlstore.Store
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *mysqlstore.Store
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				pi1, err := ctn.SafeGet("logger_logrus")
				if err != nil {
					var eo *mysqlstore.Store
					return eo, err
				}
				p1, ok := pi1.(*logrus.Entry)
				if !ok {
					var eo *mysqlstore.Store
					return eo, errors.New("could not cast parameter 1 to *logrus.Entry")
				}
				pi2, err := ctn.SafeGet("db_mysql")
				if err != nil {
					var eo *mysqlstore.Store
					return eo, err
				}
				p2, ok := pi2.(*sqlx.DB)
				if !ok {
					var eo *mysqlstore.Store
					return eo, errors.New("could not cast parameter 2 to *sqlx.DB")
				}
				pi3, err := ctn.SafeGet("tx_handler")
				if err != nil {
					var eo *mysqlstore.Store
					return eo, err
				}
				p3, ok := pi3.(*mysqltxhandler.MySQLTxHandler)
				if !ok {
					var eo *mysqlstore.Store
					return eo, errors.New("could not cast parameter 3 to *mysqltxhandler.MySQLTxHandler")
				}
				b, ok := d.Build.(func(*config.App, *logrus.Entry, *sqlx.DB, *mysqltxhandler.MySQLTxHandler) (*mysqlstore.Store, error))
				if !ok {
					var eo *mysqlstore.Store
					return eo, errors.New("could not cast build function to func(*config.App, *logrus.Entry, *sqlx.DB, *mysqltxhandler.MySQLTxHandler) (*mysqlstore.Store, error)")
				}
				return b(p0, p1, p2, p3)
			},
			Unshared: false,
		},
		{
			Name:  "service_cli",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("service_cli")
				if err != nil {
					var eo *cli.Service
					return eo, err
				}
				pi0, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *cli.Service
					return eo, err
				}
				p0, ok := pi0.(*config.App)
				if !ok {
					var eo *cli.Service
					return eo, errors.New("could not cast parameter 0 to *config.App")
				}
				b, ok := d.Build.(func(*config.App) (*cli.Service, error))
				if !ok {
					var eo *cli.Service
					return eo, errors.New("could not cast build function to func(*config.App) (*cli.Service, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "tx_handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("tx_handler")
				if err != nil {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, err
				}
				pi0, err := ctn.SafeGet("logger_logrus")
				if err != nil {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, err
				}
				p0, ok := pi0.(*logrus.Entry)
				if !ok {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, errors.New("could not cast parameter 0 to *logrus.Entry")
				}
				pi1, err := ctn.SafeGet("db_mysql")
				if err != nil {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, err
				}
				p1, ok := pi1.(*sqlx.DB)
				if !ok {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, errors.New("could not cast parameter 1 to *sqlx.DB")
				}
				pi2, err := ctn.SafeGet("config_layer")
				if err != nil {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, err
				}
				p2, ok := pi2.(*config.App)
				if !ok {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, errors.New("could not cast parameter 2 to *config.App")
				}
				b, ok := d.Build.(func(*logrus.Entry, *sqlx.DB, *config.App) (*mysqltxhandler.MySQLTxHandler, error))
				if !ok {
					var eo *mysqltxhandler.MySQLTxHandler
					return eo, errors.New("could not cast build function to func(*logrus.Entry, *sqlx.DB, *config.App) (*mysqltxhandler.MySQLTxHandler, error)")
				}
				return b(p0, p1, p2)
			},
			Unshared: false,
		},
		{
			Name:  "tx_provider",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("tx_provider")
				if err != nil {
					var eo *mysqltxprovider.Impl
					return eo, err
				}
				pi0, err := ctn.SafeGet("logger_logrus")
				if err != nil {
					var eo *mysqltxprovider.Impl
					return eo, err
				}
				p0, ok := pi0.(*logrus.Entry)
				if !ok {
					var eo *mysqltxprovider.Impl
					return eo, errors.New("could not cast parameter 0 to *logrus.Entry")
				}
				pi1, err := ctn.SafeGet("db_mysql")
				if err != nil {
					var eo *mysqltxprovider.Impl
					return eo, err
				}
				p1, ok := pi1.(*sqlx.DB)
				if !ok {
					var eo *mysqltxprovider.Impl
					return eo, errors.New("could not cast parameter 1 to *sqlx.DB")
				}
				pi2, err := ctn.SafeGet("tx_handler")
				if err != nil {
					var eo *mysqltxprovider.Impl
					return eo, err
				}
				p2, ok := pi2.(*mysqltxhandler.MySQLTxHandler)
				if !ok {
					var eo *mysqltxprovider.Impl
					return eo, errors.New("could not cast parameter 2 to *mysqltxhandler.MySQLTxHandler")
				}
				b, ok := d.Build.(func(*logrus.Entry, *sqlx.DB, *mysqltxhandler.MySQLTxHandler) (*mysqltxprovider.Impl, error))
				if !ok {
					var eo *mysqltxprovider.Impl
					return eo, errors.New("could not cast build function to func(*logrus.Entry, *sqlx.DB, *mysqltxhandler.MySQLTxHandler) (*mysqltxprovider.Impl, error)")
				}
				return b(p0, p1, p2)
			},
			Unshared: false,
		},
	}
}
