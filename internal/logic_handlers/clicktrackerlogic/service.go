package clicktrackerlogic

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Config struct {
	TxProvider persistence.TransactionProvider `json:"tx_provider" validate:"required"`
	Logger     *logrus.Entry                   `json:"Logger" validate:"required"`
	Persistor  persistor                       `json:"Persistor" validate:"required"`
}

type Service struct {
	cfg *Config
}

func (i *Config) Validate() error {
	return validationutils.Validate(i)
}

func New(cfg *Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}
	return &Service{cfg}, nil
}

func (i *Service) validateClickTrackerTypeId(ctx context.Context, handler persistence.TransactionHandler, id int) error {
	_, err := i.cfg.Persistor.GetClickTrackerById(ctx, handler, id)
	if err != nil {
		return fmt.Errorf("invalid click_tracker_id: %v", err)
	}
	return nil
}

// ListClickTrackers returns paginated click trackers
func (i *Service) ListClickTrackers(
	ctx context.Context,
	filter *model.ClickTrackerFilters,
) (*model.PaginatedClickTrackers, error) {
	db, err := i.cfg.TxProvider.Db(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}

	fmt.Println("the filter at the service --- ", strutil.GetAsJson(filter))
	paginated, err := i.cfg.Persistor.GetClickTrackers(ctx, db, filter)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get click trackers: %v", err),
		})
	}

	return paginated, nil
}

// UpdateClickTracker updates an existing click tracker.
func (i *Service) UpdateClickTracker(ctx context.Context, params *model.UpdateClickTracker) (*model.ClickTracker, error) {
	tx, err := i.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if params.ClickTrackerSetId.Valid {
		if err := i.validateClickTrackerTypeId(ctx, tx, params.ClickTrackerSetId.Int); err != nil {
			return nil, fmt.Errorf("click_tracker_set_id: %w", err)
		}
	}

	clicktracker, err := i.cfg.Persistor.UpdateClickTrackers(ctx, tx, params)
	//tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("update click tracker: %w", err)
	}

	return clicktracker, nil
}

func (i *Service) CreateClickTracker(ctx context.Context, params *model.CreateClickTracker) (*model.ClickTracker, error) {
	if err := params.Validate(); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("validate: %v", err),
		})
	}

	tx, err := i.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	clickTracker := params.ToClickTracker()

	// Ensure the click tracker set exists
	_, err = i.cfg.Persistor.GetClickTrackerSetById(ctx, tx, clickTracker.ClickTrackerSetId)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid click_tracker_set_id: %v", err),
		})
	}

	//Ensure the click tracker name is unique
	exists, err := i.cfg.Persistor.GetClickTrackerByName(ctx, tx, clickTracker.Name)
	if err != nil {
		if !strings.Contains(err.Error(), sysconsts.ErrExpectedExactlyOneEntry) {
			return nil, errs.New(&errs.Cfg{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("check click tracker unique: %v", err),
			})
		}
	}
	if exists != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("click tracker already exists"),
		})
	}

	//res, err := i.cfg.Persistor.GetClickTrackers(ctx, tx, nil)
	//if err != nil {
	//	return nil, errs.New(&errs.Cfg{
	//		StatusCode: http.StatusBadRequest,
	//		Err:        fmt.Errorf("get all error in service: %v", err),
	//	})
	//}
	//
	//fmt.Println("the response from the the get ---- ", strutil.GetAsJson(res))

	clickTracker, err = i.cfg.Persistor.CreateClickTracker(ctx, tx, clickTracker)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("create: %v", err),
		})
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit: %v", err),
		})
	}

	return clickTracker, nil
}

// DeleteClickTracker deletes a click tracker by ID.
func (s *Service) DeleteClickTracker(ctx context.Context, params *model.DeleteClickTracker) error {
	tx, err := s.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	fmt.Println("the params id ---- ", params.ID)
	err = s.cfg.Persistor.DeleteClickTracker(ctx, tx, params.ID)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("delete click tracker: %v", err),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit transaction: %v", err),
		})
	}

	return nil
}
