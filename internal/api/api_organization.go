package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
	"net/http"
)

func (a *Api) ListOrganizations(ctx *fiber.Ctx) error {
	filter := model.OrganizationFilters{
		OrganizationIsActive: null.BoolFrom(true),
	}
	if err := ctx.QueryParser(&filter); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	if err := filter.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	filter.SetPaginationDefaults()

	organizations, err := a.cfg.OrganizationService.ListOrganizations(ctx.Context(), &filter)
	return a.WriteResponse(ctx, http.StatusOK, organizations, err)
}

// CreateOrganization fetches the categories
//
// @Id CreateOrganization
// @Summary Create Organization
// @Description Create an organization
// @Tags OrganizationService
// @Accept application/json
// @Produce application/json
// @Param filters body model.CreateOrganization false "Organization filters"
// @Success 200 {object} model.Organization
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/organization [post]
func (a *Api) CreateOrganization(ctx *fiber.Ctx) error {
	var body model.CreateOrganization
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	organization, err := a.cfg.OrganizationService.CreateOrganization(ctx.Context(), &body)
	fmt.Println("=========================================================================================++>", organization)
	return a.WriteResponse(ctx, http.StatusCreated, organization, err)
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
//func (a *Api) UpdateCategory(ctx *fiber.Ctx) error {
//	var body model.UpdateCategory
//	if err := ctx.BodyParser(&body); err != nil {
//		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
//	}
//	category, err := a.cfg.CategoryService.UpdateCategory(ctx.Context(), &body)
//	return a.WriteResponse(ctx, http.StatusOK, category, err)
//}

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
//func (a *Api) DeleteCategory(ctx *fiber.Ctx) error {
//	id := ctx.Params("id")
//	categoryId, err := strconv.Atoi(id)
//	if err != nil {
//		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
//	}
//
//	deleteParams := &model.DeleteCategory{ID: categoryId}
//
//	err = a.cfg.CategoryService.DeleteCategory(ctx.Context(), deleteParams)
//	if err != nil {
//		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
//	}
//
//	return ctx.SendStatus(http.StatusNoContent)
//}

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
//func (a *Api) RestoreCategory(ctx *fiber.Ctx) error {
//	id := ctx.Params("id")
//	categoryID, err := strconv.Atoi(id)
//	if err != nil {
//		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
//	}
//
//	restoreParams := &model.RestoreCategory{ID: categoryID}
//
//	err = a.cfg.CategoryService.RestoreCategory(ctx.Context(), restoreParams)
//	if err != nil {
//		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
//	}
//
//	return ctx.SendStatus(http.StatusNoContent)
//}
