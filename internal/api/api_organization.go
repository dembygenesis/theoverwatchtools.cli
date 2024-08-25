package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
	"net/http"
	"strconv"
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
	return a.WriteResponse(ctx, http.StatusCreated, organization, err)
}

func (a *Api) DeleteOrganization(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	organizationId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	deleteParams := &model.DeleteOrganization{
		ID: organizationId,
	}

	err = a.cfg.OrganizationService.DeleteOrganization(ctx.Context(), deleteParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (a *Api) UpdateOrganization(ctx *fiber.Ctx) error {
	var body model.UpdateOrganization
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	organization, err := a.cfg.OrganizationService.UpdateOrganization(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusOK, organization, err)
}

func (a *Api) RestoreOrganization(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	organizationId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid organization ID: %v", err)
	}

	restoreParams := &model.RestoreOrganization{ID: organizationId}

	err = a.cfg.OrganizationService.RestoreOrganization(ctx.Context(), restoreParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
