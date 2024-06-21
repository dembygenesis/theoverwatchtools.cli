package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// ListOrganizations fetches the organizations
//
// @Id ListOrganizations
// @Summary Get Organizations
// @Description Returns the organizations
// @Tags OrganizationService
// @Accept application/json
// @Produce application/json
// @Param filters query model.OrganizationFilters false "Organization filters"
// @Success 200 {object} model.PaginatedOrganizations
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/organization [get]
func (a *Api) ListOrganizations(ctx *fiber.Ctx) error {
	filter := model.OrganizationFilters{
		OrganizationIsActive: []int{1},
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

// CreateOrganization fetches the organizations
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
	category, err := a.cfg.OrganizationService.CreateOrganization(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusCreated, category, err)
}

// DeleteOrganization deletes an organization by ID
//
// @Summary Delete an organization by ID
// @Description Deletes an organization by ID
// @Tags OrganizationService
// @Accept application/json
// @Produce application/json
// @Param filters body model.DeleteOrganization true "Organization ID to delete"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/organization/{id} [delete]
func (a *Api) DeleteOrganization(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	organizationId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	deleteParams := &model.DeleteOrganization{ID: organizationId}

	isDeleted, err := a.cfg.OrganizationService.GetOrganizationByID(ctx.Context(), organizationId)
	fmt.Println(strutil.GetAsJson("del cli ------------------------------------ ", isDeleted.IsActive))

	err = a.cfg.OrganizationService.DeleteOrganization(ctx.Context(), deleteParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	if isDeleted.IsActive == 0 {
		return ctx.JSON("Already Deleted")
	} else {
		return ctx.JSON("del")
	}
}

// RestoreOrganization restores an organization by ID
//
// @Summary Restore an organization by ID
// @Description Restores an organization by ID
// @Tags OrganizationService
// @Accept application/json
// @Produce application/json
// @Param id path int true "Organization ID"
// @Param body model.RestoreOrganization false "Restore parameters"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/organization/{id}/restore [patch]
func (a *Api) RestoreOrganization(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	organizationID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	restoreParams := &model.RestoreOrganization{ID: organizationID}

	isRestored, err := a.cfg.OrganizationService.GetOrganizationByID(ctx.Context(), organizationID)
	fmt.Println(strutil.GetAsJson("del cli ------------------------------------ ", isRestored.IsActive))

	err = a.cfg.OrganizationService.RestoreOrganization(ctx.Context(), restoreParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	if isRestored.IsActive == 1 {
		return ctx.JSON("Already Restored")
	} else {
		return ctx.JSON("Res")
	}
}

// UpdateOrganization fetches the organizations
//
// @Id UpdateOrganization
// @Summary Update Organization
// @Description Update an organization
// @Tags OrganizationService
// @Accept application/json
// @Produce application/json
// @Param filters body model.UpdateOrganization false "Organization body"
// @Success 200 {object} model.Organization
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/organization [patch]
func (a *Api) UpdateOrganization(ctx *fiber.Ctx) error {
	var body model.UpdateOrganization
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	category, err := a.cfg.OrganizationService.UpdateOrganization(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusOK, category, err)
}
