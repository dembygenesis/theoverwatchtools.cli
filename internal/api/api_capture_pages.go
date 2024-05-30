package api

import (
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
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
// @Router /v1/capturepage [get]
func (a *Api) ListCapturePages(ctx *fiber.Ctx) error {
	filter := model.CapturePagesFilters{
		CapturePagesIsControl: null.Bool{
			Bool:  true,
			Valid: false,
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
