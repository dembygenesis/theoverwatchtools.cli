package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// ListCategories fetches the categories
//
// @Id ListCategories
// @Summary Get Categories
// @Description Returns the categories
// @Tags CategoryService
// @Accept application/json
// @Produce application/json
// @Param filters query model.CategoryFilters false "Category filters"
// @Success 200 {object} model.PaginatedCategories
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/category [get]
func (a *Api) ListCategories(ctx *fiber.Ctx) error {
	filter := model.CategoryFilters{}
	if err := ctx.QueryParser(&filter); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	if err := filter.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	filter.SetPaginationDefaults()

	categories, err := a.cfg.CategoryService.ListCategories(ctx.Context(), &filter)
	return a.WriteResponse(ctx, http.StatusOK, categories, err)
}

func (a *Api) GetCategory(ctx *fiber.Ctx) error {
	return ctx.SendString("GetCategory")
}

// CreateCategory fetches the categories
//
// @Id CreateCategory
// @Summary Create Category
// @Description Create a category
// @Tags CategoryService
// @Accept application/json
// @Produce application/json
// @Param filters body model.CreateCategory false "Category filters"
// @Success 200 {object} model.Category
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/category [post]
func (a *Api) CreateCategory(ctx *fiber.Ctx) error {
	var body model.CreateCategory
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	category, err := a.cfg.CategoryService.CreateCategory(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusCreated, category, err)
}

// UpdateCategory fetches the categories
//
// @Id UpdateCategory
// @Summary Update Category
// @Description Update a category
// @Tags CategoryService
// @Accept application/json
// @Produce application/json
// @Param filters body model.UpdateCategory false "Category body"
// @Success 200 {object} model.Category
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/category [patch]
func (a *Api) UpdateCategory(ctx *fiber.Ctx) error {
	var body model.UpdateCategory
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	category, err := a.cfg.CategoryService.UpdateCategory(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusOK, category, err)
}
