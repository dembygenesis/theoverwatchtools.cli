package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
	"net/http"
	"strconv"
)

// ListCapturePages fetches the capture pages
//
// @Id ListCapturePages
// @Summary Get CapturePages
// @Description Returns the capture pages
// @Tags CapturePagesService
// @Accept application/json
// @Produce application/json
// @Param filters query model.CapturePagesFilters false "CapturePages filters"
// @Success 200 {object} model.PaginatedCapturePages
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/capturepage [get]
func (a *Api) ListCapturePages(ctx *fiber.Ctx) error {
	filter := model.CapturePagesFilters{
		CapturePagesIsControl: null.Bool{
			Bool:  true,
			Valid: true,
		},
	}

	if err := ctx.QueryParser(&filter); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	if err := filter.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	filter.SetPaginationDefaults()

	//fmt.Println("the filter ---- ", strutil.GetAsJson(&filter))
	organizations, err := a.cfg.CapturePagesService.ListCapturePages(ctx.Context(), &filter)
	return a.WriteResponse(ctx, http.StatusOK, organizations, err)
}

// CreateCapturePages fetches the CapturePages
//
// @Id CreateCapturePages
// @Summary Create CapturePage
// @Description Create a CapturePage
// @Tags CapturePageService
// @Accept application/json
// @Produce application/json
// @Param filters body model.CreateCapturePage false "Capture Pages filters"
// @Success 200 {object} model.CapturePage
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/capturepage [post]
func (a *Api) CreateCapturePages(ctx *fiber.Ctx) error {
	var body model.CreateCapturePage
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	capturepage, err := a.cfg.CapturePagesService.CreateCapturePages(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusCreated, capturepage, err)
}

// UpdateCapturePage fetches the capture pages
//
// @Id UpdateCapturePage
// @Summary Update Capture page
// @Description Update a capture page
// @Tags CapturePagesService
// @Accept application/json
// @Produce application/json
// @Param filters body model.UpdateCapturePages false "Category body"
// @Success 200 {object} model.CapturePages
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/capturepages [patch]
func (a *Api) UpdateCapturePages(ctx *fiber.Ctx) error {
	var body model.UpdateCapturePages
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	capturepages, err := a.cfg.CapturePagesService.UpdateCapturePages(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusOK, capturepages, err)
}

// DeleteCapturePages deletes a capture page by ID
//
// @Summary Delete a capture page by ID
// @Description Deletes a capture page by ID
// @Tags CapturePageService
// @Accept application/json
// @Produce application/json
// @Param filters body model.DeleteCapturePages true "Capture Page ID to delete"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/capturepages/{id} [delete]
func (a *Api) DeleteCapturePages(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	capturePageId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	deleteParams := &model.DeleteCapturePages{ID: capturePageId}

	err = a.cfg.CapturePagesService.DeleteCapturePages(ctx.Context(), deleteParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// RestoreCapturePages restores a capture page by ID
// @Summary Restore a capture page by ID
// @Description Restores a capture page by ID
// @Tags CapturePagesService
// @Accept application/json
// @Produce application/json
// @Param id path int true "capture page ID"
// @Param body model.RestoreCapturePages false "Restore parameters"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/capturepages/{id} [patch]
func (a *Api) RestoreCapturePages(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	capturePageID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	restoreParams := &model.RestoreCapturePages{ID: capturePageID}

	err = a.cfg.CapturePagesService.RestoreCapturePages(ctx.Context(), restoreParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
