package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errutil"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// ListCategories fetches the categories
//
// @Id ListCategories
// @Summary Get Categories
// @Description Returns the categories
// @Tags Category
// @Accept application/json
// @Produce application/json
// @Param you query string false "name search by q" Format(email)
// @Param fuck query string false "Category name"
// @Success 200 {object} model.Category
// @Failure 400 {object} model.Category
// @Failure 500 {object} model.Category
// @Router /v1/category [get]
func (a *Api) ListCategories(ctx *fiber.Ctx) error {
	filter := model.CategoryFilters{}
	if err := ctx.QueryParser(&filter); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errutil.ToArr(err))
	}

	if err := filter.ValidatePagination(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errutil.ToArr(err))
	}
	filter.SetPaginationDefaults()

	categories, err := a.cfg.CategoryManager.GetCategories(ctx.Context(), &filter)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errutil.ToArr(err))
	}

	return ctx.Status(http.StatusOK).JSON(categories)
}

func (a *Api) GetCategory(ctx *fiber.Ctx) error {
	return ctx.SendString("GetCategory")
}

func (a *Api) CreateCategory(ctx *fiber.Ctx) error {
	return ctx.SendString("CreateCategory")
}

func (a *Api) UpdateCategory(ctx *fiber.Ctx) error {
	return ctx.SendString("UpdateCategory")
}

func (a *Api) DeleteCategory(ctx *fiber.Ctx) error {
	return ctx.SendString("DeleteCategory")
}
