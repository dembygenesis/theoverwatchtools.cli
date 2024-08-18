package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/docs"
	"github.com/dembygenesis/local.tools/internal/global"
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
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
	apiV1 := a.app.Group("/api")
	v1 := apiV1.Group("/v1")

	// Category
	groupCategory := v1.Group("/category")
	groupCategory.Name("List Categories").Get("", a.ListCategories)
	groupCategory.Name("Create Category").Post("", a.CreateCategory)
	groupCategory.Name("Update Category").Patch("", a.UpdateCategory)
	groupCategory.Name("Delete Category").Delete("/:id", a.DeleteCategory)
	groupCategory.Name("Restore Category").Patch("/:id", a.RestoreCategory)

	// Organization
	groupOrganization := v1.Group("/organization")
	groupOrganization.Name("List Organizations").Get("", a.ListOrganizations)
	groupOrganization.Name("Create Organization").Post("", a.CreateOrganization)

	// Docs
	if err := a.loadStaticRoutes(); err != nil {
		return fmt.Errorf("load static routes: %w", err)
	}

	return nil
}
