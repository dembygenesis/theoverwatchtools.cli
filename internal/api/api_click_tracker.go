package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"net/http"
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
