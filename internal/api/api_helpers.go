package api

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/errutil"
	"github.com/dembygenesis/local.tools/internal/utilities/resputil"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	requestIdKey = "requestid"
)

// validStatusCode checks if the status code is invalid, else
// fails with http.StatusInternalServerError.
func (a *Api) validStatusCode(ctx *fiber.Ctx, statusCode int) bool {
	isValid := true
	if !resputil.IsValidHTTPStatusCode(statusCode) {
		a.cfg.Logger.Error(logrus.Fields{
			"err":            fmt.Errorf(sysconsts.ErrInvalidStatusCode, statusCode),
			"correlation_id": ctx.Get(requestIdKey),
		})
		errSend := ctx.SendStatus(http.StatusInternalServerError)
		if errSend != nil {
			a.cfg.Logger.Error(logrus.Fields{
				"err":            fmt.Errorf(sysconsts.ErrSendResp, errSend),
				"correlation_id": ctx.Get(requestIdKey),
			})
		}

		isValid = false
	}
	return isValid
}

// WriteResponse writes the response to the client.
func (a *Api) WriteResponse(ctx *fiber.Ctx, statusCode int, data interface{}, respErr error) error {
	if respErr == nil {
		if valid := a.validStatusCode(ctx, statusCode); !valid {
			return nil
		}
		if data == nil {
			return ctx.SendStatus(statusCode)
		}
		err := ctx.Status(statusCode).JSON(data)
		if err != nil {
			a.cfg.Logger.Error(logrus.Fields{
				"err":            fmt.Errorf(sysconsts.ErrSendResp, err),
				"correlation_id": ctx.Get(requestIdKey),
			})
		}
		return nil
	}

	errUtil, ok := errutil.ErrAsUtil(respErr)
	if !ok {
		a.cfg.Logger.Error(logrus.Fields{
			"err":            errors.New(sysconsts.ErrNotUtilErr),
			"correlation_id": ctx.Get(requestIdKey),
		})
		return ctx.Status(http.StatusInternalServerError).JSON([]string{respErr.Error()})
	}
	if valid := a.validStatusCode(ctx, errUtil.StatusCode); !valid {
		return nil
	}
	return ctx.Status(errUtil.StatusCode).JSON(errUtil.List.ErrsAsStrArr())
}
