package capturepageslogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . persistor
type persistor interface {
	GetCapturePages(ctx context.Context, tx persistence.TransactionHandler, filters *model.CapturePagesFilters) (*model.PaginatedCapturePages, error)
	CreateCapturePages(ctx context.Context, tx persistence.TransactionHandler, capturePage *model.CapturePages) (*model.CapturePages, error)
	GetCapturePageTypeById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.CapturePageType, error)
}
