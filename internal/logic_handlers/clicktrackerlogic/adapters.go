package clicktrackerlogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . persistor
type persistor interface {
	GetClickTrackers(ctx context.Context, tx persistence.TransactionHandler, filters *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error)
	CreateClickTracker(ctx context.Context, tx persistence.TransactionHandler, clickTracker *model.ClickTracker) (*model.ClickTracker, error)
	GetClickTrackerById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.ClickTracker, error)
	GetClickTrackerByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.ClickTracker, error)
	GetClickTrackerSetById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.ClickTrackerSet, error)
	UpdateClickTrackers(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateClickTracker) (*model.ClickTracker, error)
	DeleteClickTracker(ctx context.Context, tx persistence.TransactionHandler, id int) error
}
