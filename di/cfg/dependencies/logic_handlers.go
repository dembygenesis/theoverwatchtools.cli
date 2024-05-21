package dependencies

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/authlogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/capturepageslogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/marketinglogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/organizationlogic"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/userlogic"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

const (
<<<<<<< HEAD
	logicCategory        = "logic_category"
	logicUser            = "logic_user"
	logicAuth            = "logic_auth"
	logicMarketing       = "logic_marketing"
	logicOrganization    = "logic_organization"
	logicCapturePages    = "logic_capture_pages"
	logicCapturePageSets = "logic_capture_pages_sets"
=======
	logicCategory     = "logic_category"
	logicUser         = "logic_user"
	logicAuth         = "logic_auth"
	logicMarketing    = "logic_marketing"
	logicOrganization = "logic_organization"
>>>>>>> upstream/main
)

func GetLogicHandlers() []dingo.Def {
	return []dingo.Def{
		{
<<<<<<< HEAD
			Name: logicCapturePageSets,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
				store *mysqlstore.Repository,
			) (*capturepageslogic.Service, error) {
				logic, err := capturepageslogic.New(&capturepageslogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
					Persistor:  store,
				})
				if err != nil {
					return nil, fmt.Errorf("logiccapturepagesset: %v", err)
				}
				return logic, nil
			},
		},
		{
			Name: logicCapturePages,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
				store *mysqlstore.Repository,
			) (*capturepageslogic.Service, error) {
				logic, err := capturepageslogic.New(&capturepageslogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
					Persistor:  store,
				})
				if err != nil {
					return nil, fmt.Errorf("logiccapturepages: %v", err)
				}
				return logic, nil
			},
		},
		{
=======
>>>>>>> upstream/main
			Name: logicOrganization,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
				store *mysqlstore.Repository,
			) (*organizationlogic.Service, error) {
				logic, err := organizationlogic.New(&organizationlogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
					Persistor:  store,
				})
				if err != nil {
					return nil, fmt.Errorf("logicorganization: %v", err)
				}
				return logic, nil
			},
		},
		{
			Name: logicCategory,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
				store *mysqlstore.Repository,
			) (*categorylogic.Service, error) {
				logic, err := categorylogic.New(&categorylogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
					Persistor:  store,
				})
				if err != nil {
					return nil, fmt.Errorf("logicategory: %v", err)
				}
				return logic, nil
			},
		},
		{
			Name: logicUser,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
			) (*userlogic.Impl, error) {
				logic, err := userlogic.New(&userlogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
				})
				if err != nil {
					return nil, fmt.Errorf("logicuser: %v", err)
				}
				return logic, nil
			},
		},
		{
			Name: logicAuth,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
			) (*authlogic.Impl, error) {
				logic, err := authlogic.New(&authlogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
				})
				if err != nil {
					return nil, fmt.Errorf("logicauth: %v", err)
				}
				return logic, nil
			},
		},
		{
			Name: logicMarketing,
			Build: func(
				cfg *config.App,
				logger *logrus.Entry,
				txProvider *mysqlconn.Provider,
			) (*marketinglogic.Impl, error) {
				logic, err := marketinglogic.New(&marketinglogic.Config{
					TxProvider: txProvider,
					Logger:     logger,
				})
				if err != nil {
					return nil, fmt.Errorf("logicmarketing: %v", err)
				}
				return logic, nil
			},
		},
	}
}
