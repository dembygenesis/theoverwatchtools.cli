package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"net/http"
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
// @Router /v1/organization [get]
func (a *Api) ListCapturePages(ctx *fiber.Ctx) error {
	filter := model.CapturePagesFilters{
		CapturePagesIsControl: true,
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
