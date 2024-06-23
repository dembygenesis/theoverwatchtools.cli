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

// ListClickTrackers fetches the click trackers
//
// @Id ListClickTrackers
// @Summary Get ClickTrackers
// @Description Returns the click trackers
// @Tags ClickTrackerService
// @Accept application/json
// @Produce application/json
// @Param filters query model.ClickTrackerFilters false "ClickTracker filters"
// @Success 200 {object} model.PaginatedClickTrackers
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/click-tracker [get]
func (a *Api) ListClickTrackers(ctx *fiber.Ctx) error {
	filter := model.ClickTrackerFilters{}
	if err := ctx.QueryParser(&filter); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	if err := filter.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	filter.SetPaginationDefaults()

	clickTrackers, err := a.cfg.ClickTrackerService.ListClickTrackers(ctx.Context(), &filter)
	return a.WriteResponse(ctx, http.StatusOK, clickTrackers, err)
}

func (a *Api) CreateClickTracker(ctx *fiber.Ctx) error {
	var body model.CreateClickTracker

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	clickTracker, err := a.cfg.ClickTrackerService.CreateClickTracker(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusCreated, clickTracker, err)
}

// UpdateClickTracker fetches the click trackers
//
// @Id UpdateClickTracker
// @Summary Update Click Tracker
// @Description Update a Click Tracker
// @Tags ClickTrackerService
// @Accept application/json
// @Produce application/json
// @Param filters body model.UpdateClickTracker false "Click Tracker body"
// @Success 200 {object} model.ClickTracker
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/clicktrackers [patch]
func (a *Api) UpdateClickTracker(ctx *fiber.Ctx) error {
	var body model.UpdateClickTracker
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}
	clicktracker, err := a.cfg.ClickTrackerService.UpdateClickTracker(ctx.Context(), &body)
	return a.WriteResponse(ctx, http.StatusOK, clicktracker, err)
}

// DeleteClickTracker deletes a click trackers by ID
//
// @Summary Delete a click trackers by ID
// @Description Deletes a click trackers by ID
// @Tags ClickTrackerService
// @Accept application/json
// @Produce application/json
// @Param filters body model.DeleteClickTracker true "Click Tracker ID to delete"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/clicktrackers/{id} [delete]
func (a *Api) DeleteClickTracker(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	clickTrackerId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	deleteParams := &model.DeleteClickTracker{ID: clickTrackerId}

	isDeleted, err := a.cfg.ClickTrackerService.GetClickTrackerByID(ctx.Context(), clickTrackerId)
	fmt.Println(strutil.GetAsJson("del cli ------------------------------------ ", isDeleted.Clicks))

	err = a.cfg.ClickTrackerService.DeleteClickTracker(ctx.Context(), deleteParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	if isDeleted.Clicks == 0 {
		return ctx.JSON("Already Deleted")
	} else {
		return ctx.JSON("del")
	}
}

// RestoreClickTracker restores a click tracker by ID
// @Summary Restore a click tracker by ID
// @Description Restores a click tracker by ID
// @Tags ClickTrackerService
// @Accept application/json
// @Produce application/json
// @Param id path int true "click tracker ID"
// @Param body model.RestoreClickTracker false "Restore parameters"
// @Success 204 "No Content"
// @Failure 400 {object} []string
// @Failure 500 {object} []string
// @Router /v1/clicktrackers/{id} [patch]
func (a *Api) RestoreClickTracker(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	clickTrackerID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(errs.ToArr(err))
	}

	restoreParams := &model.RestoreClickTracker{ID: clickTrackerID}

	isRestored, err := a.cfg.ClickTrackerService.GetClickTrackerByID(ctx.Context(), clickTrackerID)
	fmt.Println(strutil.GetAsJson("del cli ------------------------------------ ", isRestored.Clicks))

	err = a.cfg.ClickTrackerService.RestoreClickTracker(ctx.Context(), restoreParams)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errs.ToArr(err))
	}

	if isRestored.Clicks == 1 {
		return ctx.JSON("Already Restored")
	} else {
		return ctx.JSON("Res")
	}

}
