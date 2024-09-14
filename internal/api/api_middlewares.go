package api

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/api/resource"
	"github.com/gofiber/fiber/v2"
)

func (a *Api) resourceCategory() func(c *fiber.Ctx) error {
	return a.LoadResources(context.Background(), []Resource{
		{
			Type:  resource.CategoryResource,
			Param: "name",
		},
	})
}

func (a *Api) resourceOwnCategory() func(c *fiber.Ctx) error {
	return a.LoadResources(context.Background(), []Resource{
		{
			Type:  resource.CategoryResource,
			Param: "name",
			IsOwn: true,
		},
	})
}

// auth
func (a *Api) auth() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		// Need a desconstruction .. I need to check.
		return c.Next()
	}

	/*return a.LoadResources(context.Background(), []Resource{
		{
			Type:  resource.CategoryResource,
			Param: "name",
			IsOwn: true,
		},
	})*/
}
