package api

import (
	"context"
	"errors"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
)

func TestApi_WriteResponse(t *testing.T) {
	invalidStatusCode := -100

	type args struct {
		statusCode int
		data       interface{}
		respErr    error
	}
	tests := []struct {
		name       string
		args       args
		assertions func(t *testing.T, err error)
	}{
		{
			name: "invalid_status_code",
			args: args{
				statusCode: invalidStatusCode,
				data:       nil,
				respErr:    nil,
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected not nil error")
			},
		},
		{
			name: "success_only_status",
			args: args{
				statusCode: 200,
				data:       nil,
				respErr:    nil,
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected not nil error")
			},
		},
		{
			name: "success_status_and_data",
			args: args{
				statusCode: 200,
				data:       []int{1, 2, 3, 4, 5},
				respErr:    nil,
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected not nil error")
			},
		},
		{
			name: "fail_status_and_data",
			args: args{
				statusCode: 200,
				data:       func() {},
				respErr:    nil,
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected not nil error")
			},
		},
		{
			name: "fail_not_errs.Util",
			args: args{
				statusCode: http.StatusInternalServerError,
				data:       func() {},
				respErr:    errors.New("mock error"),
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected not nil error")
			},
		},
		{
			name: "fail_errs.Util_invalid_status_code",
			args: args{
				data: func() {},
				respErr: &errs.Util{
					StatusCode: invalidStatusCode,
				},
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected not nil error")
			},
		},
		{
			name: "fail_success",
			args: args{
				data: func() {},
				respErr: &errs.Util{
					StatusCode: 100,
				},
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected not nil error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Api{
				cfg: &Config{
					Logger: logger.New(context.TODO()),
				},
				app: fiber.New(),
			}

			err := a.WriteResponse(
				fiber.New().AcquireCtx(&fasthttp.RequestCtx{}),
				tt.args.statusCode,
				tt.args.data,
				tt.args.respErr,
			)

			tt.assertions(t, err)
		})
	}
}
