package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
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
	filter := model.CategoryFilters{
		CategoryIsActive: []int{1},
	}
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

// DeleteCategory deletes a category by ID
//
// @Summary Delete a category by ID
// @Description Deletes a category by ID
// @Tags CategoryService
// @Accept application/json
// @Produce application/json
// @Param filters body model.DeleteCategory true "Category ID to delete"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/category/{id} [delete]
func (a *Api) DeleteCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	categoryId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	deleteParams := &model.DeleteCategory{ID: categoryId}

	err = a.cfg.CategoryService.DeleteCategory(ctx.Context(), deleteParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// RestoreCategory restores a category by ID
//
// @Summary Restore a category by ID
// @Description Restores a category by ID
// @Tags CategoryService
// @Accept application/json
// @Produce application/json
// @Param id path int true "Category ID"
// @Param body model.RestoreCategory false "Restore parameters"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/category/{id}/restore [patch]
func (a *Api) RestoreCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	restoreParams := &model.RestoreCategory{ID: categoryID}

	err = a.cfg.CategoryService.RestoreCategory(ctx.Context(), restoreParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
