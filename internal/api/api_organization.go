package api

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
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
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	restoreParams := &model.RestoreOrganization{ID: organizationId}

	fmt.Println("the restore params --- ", strutil.GetAsJson(restoreParams))

	err = a.cfg.OrganizationService.RestoreOrganization(ctx.Context(), restoreParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
