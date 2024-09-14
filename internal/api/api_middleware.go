package api

import (
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errutil"
	"github.com/gofiber/fiber/v2"
	"sync"
)

type Resource struct {
	Type  model.ResourceType
	Param string
	IsOwn bool
}

var mApi sync.Mutex

func (a *Api) LoadResources(ctx context.Context, resources []Resource) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println("==== resource loading")
		var errs errutil.List

		mLoadedResources := make(map[model.ResourceType]interface{})

		concurrency := pond.New(len(resources), len(resources))
		fn := func(r Resource) func() {
			return func() {
				var userId int
				if r.IsOwn {
					userId = 1
				}

				res, err := a.cfg.Resource.Get(ctx, r.Type, c.Params(r.Param), userId)
				if err != nil {
					errs.Add(fmt.Sprintf("load %s: %v", r.Type, err))
				}
				mApi.Lock()
				defer mApi.Unlock()
				mLoadedResources[r.Type] = res
			}
		}
		for i := range resources {
			concurrency.Submit(fn(resources[i]))
		}
		concurrency.StopAndWait()

		for i := range mLoadedResources {
			c.Locals(i, mLoadedResources[i])
		}

		if errs.HasErrors() {
			return c.Status(fiber.StatusInternalServerError).SendString(errs.Single().Error())
		}

		return c.Next()
	}
}
