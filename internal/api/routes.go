package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/docs"
	"github.com/dembygenesis/local.tools/internal/global"
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
	"github.com/gofiber/fiber/v2"
	"os"
	"regexp"
)

// loadStaticRoutes loads the static routes.
// It reads the index.html file from the docs directory,
// replaces the <redoc> tag with the correct spec-url attribute,
// and serves the file.
func (a *Api) loadStaticRoutes() error {
	dirDocs := fmt.Sprintf("%s/%s", os.Getenv(global.OsEnvAppDir), global.PublicDir)
	if a.cfg.WriteDocs {
		tmpl := docs.SwaggerTemplate
		bytesIndex, err := tmpl.ReadFile("index.html")
		if err != nil {
			return fmt.Errorf("read index.html: %w", err)
		}
		re := regexp.MustCompile(`<redoc[\s\S]*?>[\s\S]*?<\/redoc>`)
		indexStr := re.ReplaceAllString(string(bytesIndex),
			fmt.Sprintf("<redoc spec-url='%s/docs/v1/assets/swagger.yaml'></redoc>", a.cfg.BaseUrl))

		if err = fslib.CreateFileWithDirs(fmt.Sprintf("%s/index.html", dirDocs), []byte(indexStr)); err != nil {
			return fmt.Errorf("create file with dirs: %w", err)
		}
	}

	a.app.Static("/docs/v1", dirDocs)

	return nil
}

// Routes applies all routing/endpoint configurations.
func (a *Api) Routes() error {
	api := a.app.Group("/api")
	v1 := api.Group("/v1")

	/*res, err := a.cfg.Resource.Load(context.Background(), []apiresource.Resource{
		{
			Type:  apiresource.CategoryResource,
			Param: "id",
		},
	})
	if err != nil {
		return fmt.Errorf("load resources: %w", err)
	}

	fmt.Println("==== res:", res)
	*/

	// Load middleware resource test
	fnNames := func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")
		return ctx.SendString(name)
	}

	v1.Get("/test/:name", a.resourceCategory(), fnNames)
	v1.Get("/test2/:name", a.resourceOwnCategory(), fnNames)

	// Category
	gCat := v1.Group("/category")
	gCat.Name("List Categories").Get("", a.ListCategories)
	gCat.Name("Create Category").Post("", a.CreateCategory)
	// gCat.Name("Update Category").Patch("", a.UpdateCategory)
	gCat.Name("Delete Category").Delete("/:id", a.DeleteCategory)
	gCat.Name("Restore Category").Patch("/:id", a.RestoreCategory)

	// Category V2
	gCat.Name("Update Category2").Patch("/:name", a.resourceOwnCategory(), a.UpdateCategory2)

	// Docs
	if err := a.loadStaticRoutes(); err != nil {
		return fmt.Errorf("load static routes: %w", err)
	}

	return nil
}
