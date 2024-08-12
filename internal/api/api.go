package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/global"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
	"os"
)

type Api struct {
	cfg *Config
	app *fiber.App
}

type Config struct {
	// BaseUrl is the base URL of your API.
	BaseUrl string `json:"base_url" validate:"is_url"`

	// WriteDocs is a flag to write the docs to the public folder.
	WriteDocs bool `json:"write_docs"`

	// Logger is the logger instance.
	Logger *logrus.Entry `json:"logger" validate:"required"`

	// Port is the port your API will listen to.
	Port int `json:"port" validate:"required,greater_than_zero"`

	// CategoryService is the biz function for category
	CategoryService categoryService `json:"category_manager" validate:"required"`

	// OrganizationService is the biz function for the organization
	OrganizationService organizationService `json:"organization_service" validate:"required"`

	TxHandler *mysqltx.Handler
}

func (a *Config) Validate() error {
	err := validationutils.Validate(a)
	if err != nil {
		return fmt.Errorf("required fields: %v", err)
	}
	return nil
}

// New creates a new API instance.
func New(cfg *Config) (*Api, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	docLocation := fmt.Sprintf("%s/%s", os.Getenv(global.OsEnvAppDir), "internal/docs")
	engine := html.New(docLocation, ".html")

	api := &Api{
		cfg: cfg,
		app: fiber.New(fiber.Config{
			Views:     engine,
			BodyLimit: 20971520,
		}),
	}

	api.app.Use(requestid.New())
	api.app.Use(recover.New())
	api.app.Use(cors.New())
	api.app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/New_York",
	}))

	api.app.Get("/docs", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"Title": "hello world",
		})
	})

	api.app.Get("/docs/swagger", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"Title": "hello world",
		}, "assets")
	})

	if err := api.Routes(); err != nil {
		return nil, fmt.Errorf("routes: %v", err)
	}

	return api, nil
}

// Listen makes fiber listen to the port.
func (a *Api) Listen() error {
	if err := a.Routes(); err != nil {
		return fmt.Errorf("routes: %v", err)
	}

	if err := a.app.Listen(fmt.Sprintf(":%v", a.cfg.Port)); err != nil {
		return fmt.Errorf("listen: %v", err)
	}

	return nil
}
