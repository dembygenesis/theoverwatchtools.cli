package dic

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	providerPkg "github.com/dembygenesis/local.tools/di/cfg"

	cli "github.com/dembygenesis/local.tools/internal/cli"
	config "github.com/dembygenesis/local.tools/internal/config"
	authlogic "github.com/dembygenesis/local.tools/internal/logic_handlers/authlogic"
	capturepageslogic "github.com/dembygenesis/local.tools/internal/logic_handlers/capturepageslogic"
	categorylogic "github.com/dembygenesis/local.tools/internal/logic_handlers/categorylogic"
	marketinglogic "github.com/dembygenesis/local.tools/internal/logic_handlers/marketinglogic"
	organizationlogic "github.com/dembygenesis/local.tools/internal/logic_handlers/organizationlogic"
	userlogic "github.com/dembygenesis/local.tools/internal/logic_handlers/userlogic"
	mysqlconn "github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	mysqltx "github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	mysqlstore "github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	sqlx "github.com/jmoiron/sqlx"
	logrus "github.com/sirupsen/logrus"
)

// C retrieves a Container from an interface.
// The function panics if the Container can not be retrieved.
//
// The interface can be :
//   - a *Container
//   - an *http.Request containing a *Container in its context.Context
//     for the dingo.ContainerKey("dingo") key.
//
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}
	r, ok := i.(*http.Request)
	if !ok {
		panic("could not get the container with dic.C()")
	}
	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok {
		panic("could not get the container from the given *http.Request in dic.C()")
	}
	return c
}

type builder struct {
	builder *di.Builder
}

// NewBuilder creates a builder that can be used to create a Container.
// You probably should use NewContainer to create the container directly.
// But using NewBuilder allows you to redefine some di services.
// This can be used for testing.
// But this behavior is not safe, so be sure to know what you are doing.
func NewBuilder(scopes ...string) (*builder, error) {
	if len(scopes) == 0 {
		scopes = []string{di.App, di.Request, di.SubRequest}
	}
	b, err := di.NewBuilder(scopes...)
	if err != nil {
		return nil, fmt.Errorf("could not create di.Builder: %v", err)
	}
	provider := &providerPkg.Provider{}
	if err := provider.Load(); err != nil {
		return nil, fmt.Errorf("could not load definitions with the Provider (Provider from github.com/dembygenesis/local.tools/di/cfg): %v", err)
	}
	for _, d := range getDiDefs(provider) {
		if err := b.Add(d); err != nil {
			return nil, fmt.Errorf("could not add di.Def in di.Builder: %v", err)
		}
	}
	return &builder{builder: b}, nil
}

// Add adds one or more definitions in the Builder.
// It returns an error if a definition can not be added.
func (b *builder) Add(defs ...di.Def) error {
	return b.builder.Add(defs...)
}

// Set is a shortcut to add a definition for an already built object.
func (b *builder) Set(name string, obj interface{}) error {
	return b.builder.Set(name, obj)
}

// Build creates a Container in the most generic scope.
func (b *builder) Build() *Container {
	return &Container{ctn: b.builder.Build()}
}

// NewContainer creates a new Container.
// If no scope is provided, di.App, di.Request and di.SubRequest are used.
// The returned Container has the most generic scope (di.App).
// The SubContainer() method should be called to get a Container in a more specific scope.
func NewContainer(scopes ...string) (*Container, error) {
	b, err := NewBuilder(scopes...)
	if err != nil {
		return nil, err
	}
	return b.Build(), nil
}

// Container represents a generated dependency injection container.
// It is a wrapper around a di.Container.
//
// A Container has a scope and may have a parent in a more generic scope
// and children in a more specific scope.
// Objects can be retrieved from the Container.
// If the requested object does not already exist in the Container,
// it is built thanks to the object definition.
// The following attempts to get this object will return the same object.
type Container struct {
	ctn di.Container
}

// Scope returns the Container scope.
func (c *Container) Scope() string {
	return c.ctn.Scope()
}

// Scopes returns the list of available scopes.
func (c *Container) Scopes() []string {
	return c.ctn.Scopes()
}

// ParentScopes returns the list of scopes wider than the Container scope.
func (c *Container) ParentScopes() []string {
	return c.ctn.ParentScopes()
}

// SubScopes returns the list of scopes that are more specific than the Container scope.
func (c *Container) SubScopes() []string {
	return c.ctn.SubScopes()
}

// Parent returns the parent Container.
func (c *Container) Parent() *Container {
	if p := c.ctn.Parent(); p != nil {
		return &Container{ctn: p}
	}
	return nil
}

// SubContainer creates a new Container in the next sub-scope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object has to belong to this scope or a more generic one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// Instead it panics.
func (c *Container) Get(name string) interface{} {
	return c.ctn.Get(name)
}

// Fill is similar to SafeGet but it does not return the object.
// Instead it fills the provided object with the value returned by SafeGet.
// The provided object must be a pointer to the value returned by SafeGet.
func (c *Container) Fill(name string, dst interface{}) error {
	return c.ctn.Fill(name, dst)
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a more specific scope.
// To do so, UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGet(name string) interface{} {
	return c.ctn.UnscopedGet(name)
}

// UnscopedFill is similar to UnscopedSafeGet but copies the object in dst instead of returning it.
func (c *Container) UnscopedFill(name string, dst interface{}) error {
	return c.ctn.UnscopedFill(name, dst)
}

// Clean deletes the sub-container created by UnscopedSafeGet, UnscopedGet or UnscopedFill.
func (c *Container) Clean() error {
	return c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls the Close function of their Definition on them.
// It will also call DeleteWithSubContainers on each child and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
// The sub-containers are deleted even if they are still used in other goroutines.
// It can cause errors. You may want to use the Delete method instead.
func (c *Container) DeleteWithSubContainers() error {
	return c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers if the Container does not have any child.
// But if the Container has sub-containers, it will not be deleted right away.
// The deletion only occurs when all the sub-containers have been deleted manually.
// So you have to call Delete or DeleteWithSubContainers on all the sub-containers.
func (c *Container) Delete() error {
	return c.ctn.Delete()
}

// IsClosed returns true if the Container has been deleted.
func (c *Container) IsClosed() bool {
	return c.ctn.IsClosed()
}

// SafeGetConfigLayer retrieves the "config_layer" object from the main scope.
//
// ---------------------------------------------
//
//	name: "config_layer"
//	type: *config.App
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConfigLayer() (*config.App, error) {
	i, err := c.ctn.SafeGet("config_layer")
	if err != nil {
		var eo *config.App
		return eo, err
	}
	o, ok := i.(*config.App)
	if !ok {
		return o, errors.New("could get 'config_layer' because the object could not be cast to *config.App")
	}
	return o, nil
}

// GetConfigLayer retrieves the "config_layer" object from the main scope.
//
// ---------------------------------------------
//
//	name: "config_layer"
//	type: *config.App
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConfigLayer() *config.App {
	o, err := c.SafeGetConfigLayer()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConfigLayer retrieves the "config_layer" object from the main scope.
//
// ---------------------------------------------
//
//	name: "config_layer"
//	type: *config.App
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConfigLayer() (*config.App, error) {
	i, err := c.ctn.UnscopedSafeGet("config_layer")
	if err != nil {
		var eo *config.App
		return eo, err
	}
	o, ok := i.(*config.App)
	if !ok {
		return o, errors.New("could get 'config_layer' because the object could not be cast to *config.App")
	}
	return o, nil
}

// UnscopedGetConfigLayer retrieves the "config_layer" object from the main scope.
//
// ---------------------------------------------
//
//	name: "config_layer"
//	type: *config.App
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConfigLayer() *config.App {
	o, err := c.UnscopedSafeGetConfigLayer()
	if err != nil {
		panic(err)
	}
	return o
}

// ConfigLayer retrieves the "config_layer" object from the main scope.
//
// ---------------------------------------------
//
//	name: "config_layer"
//	type: *config.App
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConfigLayer method.
// If the container can not be retrieved, it panics.
func ConfigLayer(i interface{}) *config.App {
	return C(i).GetConfigLayer()
}

// SafeGetDbMysql retrieves the "db_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "db_mysql"
//	type: *sqlx.DB
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetDbMysql() (*sqlx.DB, error) {
	i, err := c.ctn.SafeGet("db_mysql")
	if err != nil {
		var eo *sqlx.DB
		return eo, err
	}
	o, ok := i.(*sqlx.DB)
	if !ok {
		return o, errors.New("could get 'db_mysql' because the object could not be cast to *sqlx.DB")
	}
	return o, nil
}

// GetDbMysql retrieves the "db_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "db_mysql"
//	type: *sqlx.DB
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetDbMysql() *sqlx.DB {
	o, err := c.SafeGetDbMysql()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetDbMysql retrieves the "db_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "db_mysql"
//	type: *sqlx.DB
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetDbMysql() (*sqlx.DB, error) {
	i, err := c.ctn.UnscopedSafeGet("db_mysql")
	if err != nil {
		var eo *sqlx.DB
		return eo, err
	}
	o, ok := i.(*sqlx.DB)
	if !ok {
		return o, errors.New("could get 'db_mysql' because the object could not be cast to *sqlx.DB")
	}
	return o, nil
}

// UnscopedGetDbMysql retrieves the "db_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "db_mysql"
//	type: *sqlx.DB
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetDbMysql() *sqlx.DB {
	o, err := c.UnscopedSafeGetDbMysql()
	if err != nil {
		panic(err)
	}
	return o
}

// DbMysql retrieves the "db_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "db_mysql"
//	type: *sqlx.DB
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetDbMysql method.
// If the container can not be retrieved, it panics.
func DbMysql(i interface{}) *sqlx.DB {
	return C(i).GetDbMysql()
}

// SafeGetLoggerLogrus retrieves the "logger_logrus" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logger_logrus"
//	type: *logrus.Entry
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLoggerLogrus() (*logrus.Entry, error) {
	i, err := c.ctn.SafeGet("logger_logrus")
	if err != nil {
		var eo *logrus.Entry
		return eo, err
	}
	o, ok := i.(*logrus.Entry)
	if !ok {
		return o, errors.New("could get 'logger_logrus' because the object could not be cast to *logrus.Entry")
	}
	return o, nil
}

// GetLoggerLogrus retrieves the "logger_logrus" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logger_logrus"
//	type: *logrus.Entry
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLoggerLogrus() *logrus.Entry {
	o, err := c.SafeGetLoggerLogrus()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLoggerLogrus retrieves the "logger_logrus" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logger_logrus"
//	type: *logrus.Entry
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLoggerLogrus() (*logrus.Entry, error) {
	i, err := c.ctn.UnscopedSafeGet("logger_logrus")
	if err != nil {
		var eo *logrus.Entry
		return eo, err
	}
	o, ok := i.(*logrus.Entry)
	if !ok {
		return o, errors.New("could get 'logger_logrus' because the object could not be cast to *logrus.Entry")
	}
	return o, nil
}

// UnscopedGetLoggerLogrus retrieves the "logger_logrus" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logger_logrus"
//	type: *logrus.Entry
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLoggerLogrus() *logrus.Entry {
	o, err := c.UnscopedSafeGetLoggerLogrus()
	if err != nil {
		panic(err)
	}
	return o
}

// LoggerLogrus retrieves the "logger_logrus" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logger_logrus"
//	type: *logrus.Entry
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLoggerLogrus method.
// If the container can not be retrieved, it panics.
func LoggerLogrus(i interface{}) *logrus.Entry {
	return C(i).GetLoggerLogrus()
}

// SafeGetLogicAuth retrieves the "logic_auth" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_auth"
//	type: *authlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLogicAuth() (*authlogic.Impl, error) {
	i, err := c.ctn.SafeGet("logic_auth")
	if err != nil {
		var eo *authlogic.Impl
		return eo, err
	}
	o, ok := i.(*authlogic.Impl)
	if !ok {
		return o, errors.New("could get 'logic_auth' because the object could not be cast to *authlogic.Impl")
	}
	return o, nil
}

// GetLogicAuth retrieves the "logic_auth" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_auth"
//	type: *authlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLogicAuth() *authlogic.Impl {
	o, err := c.SafeGetLogicAuth()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogicAuth retrieves the "logic_auth" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_auth"
//	type: *authlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLogicAuth() (*authlogic.Impl, error) {
	i, err := c.ctn.UnscopedSafeGet("logic_auth")
	if err != nil {
		var eo *authlogic.Impl
		return eo, err
	}
	o, ok := i.(*authlogic.Impl)
	if !ok {
		return o, errors.New("could get 'logic_auth' because the object could not be cast to *authlogic.Impl")
	}
	return o, nil
}

// UnscopedGetLogicAuth retrieves the "logic_auth" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_auth"
//	type: *authlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLogicAuth() *authlogic.Impl {
	o, err := c.UnscopedSafeGetLogicAuth()
	if err != nil {
		panic(err)
	}
	return o
}

// LogicAuth retrieves the "logic_auth" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_auth"
//	type: *authlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLogicAuth method.
// If the container can not be retrieved, it panics.
func LogicAuth(i interface{}) *authlogic.Impl {
	return C(i).GetLogicAuth()
}

// SafeGetLogicCapturePages retrieves the "logic_capture_pages" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLogicCapturePages() (*capturepageslogic.Service, error) {
	i, err := c.ctn.SafeGet("logic_capture_pages")
	if err != nil {
		var eo *capturepageslogic.Service
		return eo, err
	}
	o, ok := i.(*capturepageslogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_capture_pages' because the object could not be cast to *capturepageslogic.Service")
	}
	return o, nil
}

// GetLogicCapturePages retrieves the "logic_capture_pages" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLogicCapturePages() *capturepageslogic.Service {
	o, err := c.SafeGetLogicCapturePages()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogicCapturePages retrieves the "logic_capture_pages" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLogicCapturePages() (*capturepageslogic.Service, error) {
	i, err := c.ctn.UnscopedSafeGet("logic_capture_pages")
	if err != nil {
		var eo *capturepageslogic.Service
		return eo, err
	}
	o, ok := i.(*capturepageslogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_capture_pages' because the object could not be cast to *capturepageslogic.Service")
	}
	return o, nil
}

// UnscopedGetLogicCapturePages retrieves the "logic_capture_pages" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLogicCapturePages() *capturepageslogic.Service {
	o, err := c.UnscopedSafeGetLogicCapturePages()
	if err != nil {
		panic(err)
	}
	return o
}

// LogicCapturePages retrieves the "logic_capture_pages" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLogicCapturePages method.
// If the container can not be retrieved, it panics.
func LogicCapturePages(i interface{}) *capturepageslogic.Service {
	return C(i).GetLogicCapturePages()
}

// SafeGetLogicCapturePagesSets retrieves the "logic_capture_pages_sets" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages_sets"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLogicCapturePagesSets() (*capturepageslogic.Service, error) {
	i, err := c.ctn.SafeGet("logic_capture_pages_sets")
	if err != nil {
		var eo *capturepageslogic.Service
		return eo, err
	}
	o, ok := i.(*capturepageslogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_capture_pages_sets' because the object could not be cast to *capturepageslogic.Service")
	}
	return o, nil
}

// GetLogicCapturePagesSets retrieves the "logic_capture_pages_sets" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages_sets"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLogicCapturePagesSets() *capturepageslogic.Service {
	o, err := c.SafeGetLogicCapturePagesSets()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogicCapturePagesSets retrieves the "logic_capture_pages_sets" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages_sets"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLogicCapturePagesSets() (*capturepageslogic.Service, error) {
	i, err := c.ctn.UnscopedSafeGet("logic_capture_pages_sets")
	if err != nil {
		var eo *capturepageslogic.Service
		return eo, err
	}
	o, ok := i.(*capturepageslogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_capture_pages_sets' because the object could not be cast to *capturepageslogic.Service")
	}
	return o, nil
}

// UnscopedGetLogicCapturePagesSets retrieves the "logic_capture_pages_sets" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages_sets"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLogicCapturePagesSets() *capturepageslogic.Service {
	o, err := c.UnscopedSafeGetLogicCapturePagesSets()
	if err != nil {
		panic(err)
	}
	return o
}

// LogicCapturePagesSets retrieves the "logic_capture_pages_sets" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_capture_pages_sets"
//	type: *capturepageslogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLogicCapturePagesSets method.
// If the container can not be retrieved, it panics.
func LogicCapturePagesSets(i interface{}) *capturepageslogic.Service {
	return C(i).GetLogicCapturePagesSets()
}

// SafeGetLogicCategory retrieves the "logic_category" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_category"
//	type: *categorylogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLogicCategory() (*categorylogic.Service, error) {
	i, err := c.ctn.SafeGet("logic_category")
	if err != nil {
		var eo *categorylogic.Service
		return eo, err
	}
	o, ok := i.(*categorylogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_category' because the object could not be cast to *categorylogic.Service")
	}
	return o, nil
}

// GetLogicCategory retrieves the "logic_category" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_category"
//	type: *categorylogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLogicCategory() *categorylogic.Service {
	o, err := c.SafeGetLogicCategory()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogicCategory retrieves the "logic_category" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_category"
//	type: *categorylogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLogicCategory() (*categorylogic.Service, error) {
	i, err := c.ctn.UnscopedSafeGet("logic_category")
	if err != nil {
		var eo *categorylogic.Service
		return eo, err
	}
	o, ok := i.(*categorylogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_category' because the object could not be cast to *categorylogic.Service")
	}
	return o, nil
}

// UnscopedGetLogicCategory retrieves the "logic_category" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_category"
//	type: *categorylogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLogicCategory() *categorylogic.Service {
	o, err := c.UnscopedSafeGetLogicCategory()
	if err != nil {
		panic(err)
	}
	return o
}

// LogicCategory retrieves the "logic_category" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_category"
//	type: *categorylogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLogicCategory method.
// If the container can not be retrieved, it panics.
func LogicCategory(i interface{}) *categorylogic.Service {
	return C(i).GetLogicCategory()
}

// SafeGetLogicMarketing retrieves the "logic_marketing" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_marketing"
//	type: *marketinglogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLogicMarketing() (*marketinglogic.Impl, error) {
	i, err := c.ctn.SafeGet("logic_marketing")
	if err != nil {
		var eo *marketinglogic.Impl
		return eo, err
	}
	o, ok := i.(*marketinglogic.Impl)
	if !ok {
		return o, errors.New("could get 'logic_marketing' because the object could not be cast to *marketinglogic.Impl")
	}
	return o, nil
}

// GetLogicMarketing retrieves the "logic_marketing" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_marketing"
//	type: *marketinglogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLogicMarketing() *marketinglogic.Impl {
	o, err := c.SafeGetLogicMarketing()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogicMarketing retrieves the "logic_marketing" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_marketing"
//	type: *marketinglogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLogicMarketing() (*marketinglogic.Impl, error) {
	i, err := c.ctn.UnscopedSafeGet("logic_marketing")
	if err != nil {
		var eo *marketinglogic.Impl
		return eo, err
	}
	o, ok := i.(*marketinglogic.Impl)
	if !ok {
		return o, errors.New("could get 'logic_marketing' because the object could not be cast to *marketinglogic.Impl")
	}
	return o, nil
}

// UnscopedGetLogicMarketing retrieves the "logic_marketing" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_marketing"
//	type: *marketinglogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLogicMarketing() *marketinglogic.Impl {
	o, err := c.UnscopedSafeGetLogicMarketing()
	if err != nil {
		panic(err)
	}
	return o
}

// LogicMarketing retrieves the "logic_marketing" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_marketing"
//	type: *marketinglogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLogicMarketing method.
// If the container can not be retrieved, it panics.
func LogicMarketing(i interface{}) *marketinglogic.Impl {
	return C(i).GetLogicMarketing()
}

// SafeGetLogicOrganization retrieves the "logic_organization" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_organization"
//	type: *organizationlogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLogicOrganization() (*organizationlogic.Service, error) {
	i, err := c.ctn.SafeGet("logic_organization")
	if err != nil {
		var eo *organizationlogic.Service
		return eo, err
	}
	o, ok := i.(*organizationlogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_organization' because the object could not be cast to *organizationlogic.Service")
	}
	return o, nil
}

// GetLogicOrganization retrieves the "logic_organization" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_organization"
//	type: *organizationlogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLogicOrganization() *organizationlogic.Service {
	o, err := c.SafeGetLogicOrganization()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogicOrganization retrieves the "logic_organization" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_organization"
//	type: *organizationlogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLogicOrganization() (*organizationlogic.Service, error) {
	i, err := c.ctn.UnscopedSafeGet("logic_organization")
	if err != nil {
		var eo *organizationlogic.Service
		return eo, err
	}
	o, ok := i.(*organizationlogic.Service)
	if !ok {
		return o, errors.New("could get 'logic_organization' because the object could not be cast to *organizationlogic.Service")
	}
	return o, nil
}

// UnscopedGetLogicOrganization retrieves the "logic_organization" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_organization"
//	type: *organizationlogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLogicOrganization() *organizationlogic.Service {
	o, err := c.UnscopedSafeGetLogicOrganization()
	if err != nil {
		panic(err)
	}
	return o
}

// LogicOrganization retrieves the "logic_organization" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_organization"
//	type: *organizationlogic.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//		- "3": Service(*mysqlstore.Repository) ["persistence_mysql"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLogicOrganization method.
// If the container can not be retrieved, it panics.
func LogicOrganization(i interface{}) *organizationlogic.Service {
	return C(i).GetLogicOrganization()
}

// SafeGetLogicUser retrieves the "logic_user" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_user"
//	type: *userlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetLogicUser() (*userlogic.Impl, error) {
	i, err := c.ctn.SafeGet("logic_user")
	if err != nil {
		var eo *userlogic.Impl
		return eo, err
	}
	o, ok := i.(*userlogic.Impl)
	if !ok {
		return o, errors.New("could get 'logic_user' because the object could not be cast to *userlogic.Impl")
	}
	return o, nil
}

// GetLogicUser retrieves the "logic_user" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_user"
//	type: *userlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetLogicUser() *userlogic.Impl {
	o, err := c.SafeGetLogicUser()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogicUser retrieves the "logic_user" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_user"
//	type: *userlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetLogicUser() (*userlogic.Impl, error) {
	i, err := c.ctn.UnscopedSafeGet("logic_user")
	if err != nil {
		var eo *userlogic.Impl
		return eo, err
	}
	o, ok := i.(*userlogic.Impl)
	if !ok {
		return o, errors.New("could get 'logic_user' because the object could not be cast to *userlogic.Impl")
	}
	return o, nil
}

// UnscopedGetLogicUser retrieves the "logic_user" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_user"
//	type: *userlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetLogicUser() *userlogic.Impl {
	o, err := c.UnscopedSafeGetLogicUser()
	if err != nil {
		panic(err)
	}
	return o
}

// LogicUser retrieves the "logic_user" object from the main scope.
//
// ---------------------------------------------
//
//	name: "logic_user"
//	type: *userlogic.Impl
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*mysqlconn.Provider) ["tx_provider"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetLogicUser method.
// If the container can not be retrieved, it panics.
func LogicUser(i interface{}) *userlogic.Impl {
	return C(i).GetLogicUser()
}

// SafeGetPersistenceMysql retrieves the "persistence_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "persistence_mysql"
//	type: *mysqlstore.Repository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*sqlx.DB) ["db_mysql"]
//		- "3": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetPersistenceMysql() (*mysqlstore.Repository, error) {
	i, err := c.ctn.SafeGet("persistence_mysql")
	if err != nil {
		var eo *mysqlstore.Repository
		return eo, err
	}
	o, ok := i.(*mysqlstore.Repository)
	if !ok {
		return o, errors.New("could get 'persistence_mysql' because the object could not be cast to *mysqlstore.Repository")
	}
	return o, nil
}

// GetPersistenceMysql retrieves the "persistence_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "persistence_mysql"
//	type: *mysqlstore.Repository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*sqlx.DB) ["db_mysql"]
//		- "3": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetPersistenceMysql() *mysqlstore.Repository {
	o, err := c.SafeGetPersistenceMysql()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetPersistenceMysql retrieves the "persistence_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "persistence_mysql"
//	type: *mysqlstore.Repository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*sqlx.DB) ["db_mysql"]
//		- "3": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetPersistenceMysql() (*mysqlstore.Repository, error) {
	i, err := c.ctn.UnscopedSafeGet("persistence_mysql")
	if err != nil {
		var eo *mysqlstore.Repository
		return eo, err
	}
	o, ok := i.(*mysqlstore.Repository)
	if !ok {
		return o, errors.New("could get 'persistence_mysql' because the object could not be cast to *mysqlstore.Repository")
	}
	return o, nil
}

// UnscopedGetPersistenceMysql retrieves the "persistence_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "persistence_mysql"
//	type: *mysqlstore.Repository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*sqlx.DB) ["db_mysql"]
//		- "3": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetPersistenceMysql() *mysqlstore.Repository {
	o, err := c.UnscopedSafeGetPersistenceMysql()
	if err != nil {
		panic(err)
	}
	return o
}

// PersistenceMysql retrieves the "persistence_mysql" object from the main scope.
//
// ---------------------------------------------
//
//	name: "persistence_mysql"
//	type: *mysqlstore.Repository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//		- "1": Service(*logrus.Entry) ["logger_logrus"]
//		- "2": Service(*sqlx.DB) ["db_mysql"]
//		- "3": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetPersistenceMysql method.
// If the container can not be retrieved, it panics.
func PersistenceMysql(i interface{}) *mysqlstore.Repository {
	return C(i).GetPersistenceMysql()
}

// SafeGetServiceCli retrieves the "service_cli" object from the main scope.
//
// ---------------------------------------------
//
//	name: "service_cli"
//	type: *cli.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetServiceCli() (*cli.Service, error) {
	i, err := c.ctn.SafeGet("service_cli")
	if err != nil {
		var eo *cli.Service
		return eo, err
	}
	o, ok := i.(*cli.Service)
	if !ok {
		return o, errors.New("could get 'service_cli' because the object could not be cast to *cli.Service")
	}
	return o, nil
}

// GetServiceCli retrieves the "service_cli" object from the main scope.
//
// ---------------------------------------------
//
//	name: "service_cli"
//	type: *cli.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetServiceCli() *cli.Service {
	o, err := c.SafeGetServiceCli()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetServiceCli retrieves the "service_cli" object from the main scope.
//
// ---------------------------------------------
//
//	name: "service_cli"
//	type: *cli.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetServiceCli() (*cli.Service, error) {
	i, err := c.ctn.UnscopedSafeGet("service_cli")
	if err != nil {
		var eo *cli.Service
		return eo, err
	}
	o, ok := i.(*cli.Service)
	if !ok {
		return o, errors.New("could get 'service_cli' because the object could not be cast to *cli.Service")
	}
	return o, nil
}

// UnscopedGetServiceCli retrieves the "service_cli" object from the main scope.
//
// ---------------------------------------------
//
//	name: "service_cli"
//	type: *cli.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetServiceCli() *cli.Service {
	o, err := c.UnscopedSafeGetServiceCli()
	if err != nil {
		panic(err)
	}
	return o
}

// ServiceCli retrieves the "service_cli" object from the main scope.
//
// ---------------------------------------------
//
//	name: "service_cli"
//	type: *cli.Service
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetServiceCli method.
// If the container can not be retrieved, it panics.
func ServiceCli(i interface{}) *cli.Service {
	return C(i).GetServiceCli()
}

// SafeGetTxHandler retrieves the "tx_handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_handler"
//	type: *mysqltx.Handler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetTxHandler() (*mysqltx.Handler, error) {
	i, err := c.ctn.SafeGet("tx_handler")
	if err != nil {
		var eo *mysqltx.Handler
		return eo, err
	}
	o, ok := i.(*mysqltx.Handler)
	if !ok {
		return o, errors.New("could get 'tx_handler' because the object could not be cast to *mysqltx.Handler")
	}
	return o, nil
}

// GetTxHandler retrieves the "tx_handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_handler"
//	type: *mysqltx.Handler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetTxHandler() *mysqltx.Handler {
	o, err := c.SafeGetTxHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetTxHandler retrieves the "tx_handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_handler"
//	type: *mysqltx.Handler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetTxHandler() (*mysqltx.Handler, error) {
	i, err := c.ctn.UnscopedSafeGet("tx_handler")
	if err != nil {
		var eo *mysqltx.Handler
		return eo, err
	}
	o, ok := i.(*mysqltx.Handler)
	if !ok {
		return o, errors.New("could get 'tx_handler' because the object could not be cast to *mysqltx.Handler")
	}
	return o, nil
}

// UnscopedGetTxHandler retrieves the "tx_handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_handler"
//	type: *mysqltx.Handler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetTxHandler() *mysqltx.Handler {
	o, err := c.UnscopedSafeGetTxHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// TxHandler retrieves the "tx_handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_handler"
//	type: *mysqltx.Handler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*config.App) ["config_layer"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetTxHandler method.
// If the container can not be retrieved, it panics.
func TxHandler(i interface{}) *mysqltx.Handler {
	return C(i).GetTxHandler()
}

// SafeGetTxProvider retrieves the "tx_provider" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_provider"
//	type: *mysqlconn.Provider
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetTxProvider() (*mysqlconn.Provider, error) {
	i, err := c.ctn.SafeGet("tx_provider")
	if err != nil {
		var eo *mysqlconn.Provider
		return eo, err
	}
	o, ok := i.(*mysqlconn.Provider)
	if !ok {
		return o, errors.New("could get 'tx_provider' because the object could not be cast to *mysqlconn.Provider")
	}
	return o, nil
}

// GetTxProvider retrieves the "tx_provider" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_provider"
//	type: *mysqlconn.Provider
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetTxProvider() *mysqlconn.Provider {
	o, err := c.SafeGetTxProvider()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetTxProvider retrieves the "tx_provider" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_provider"
//	type: *mysqlconn.Provider
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetTxProvider() (*mysqlconn.Provider, error) {
	i, err := c.ctn.UnscopedSafeGet("tx_provider")
	if err != nil {
		var eo *mysqlconn.Provider
		return eo, err
	}
	o, ok := i.(*mysqlconn.Provider)
	if !ok {
		return o, errors.New("could get 'tx_provider' because the object could not be cast to *mysqlconn.Provider")
	}
	return o, nil
}

// UnscopedGetTxProvider retrieves the "tx_provider" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_provider"
//	type: *mysqlconn.Provider
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetTxProvider() *mysqlconn.Provider {
	o, err := c.UnscopedSafeGetTxProvider()
	if err != nil {
		panic(err)
	}
	return o
}

// TxProvider retrieves the "tx_provider" object from the main scope.
//
// ---------------------------------------------
//
//	name: "tx_provider"
//	type: *mysqlconn.Provider
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(*logrus.Entry) ["logger_logrus"]
//		- "1": Service(*sqlx.DB) ["db_mysql"]
//		- "2": Service(*mysqltx.Handler) ["tx_handler"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetTxProvider method.
// If the container can not be retrieved, it panics.
func TxProvider(i interface{}) *mysqlconn.Provider {
	return C(i).GetTxProvider()
}
