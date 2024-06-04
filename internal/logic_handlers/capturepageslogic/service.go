package capturepageslogic

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Config struct {
	TxProvider persistence.TransactionProvider `json:"tx_provider" validate:"required"`
	Logger     *logrus.Entry                   `json:"Logger" validate:"required"`
	Persistor  persistor                       `json:"Persistor" validate:"required"`
}

func (i *Config) Validate() error {
	return validationutils.Validate(i)
}

type Service struct {
	cfg *Config
}

func New(cfg *Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}
	return &Service{cfg}, nil
}

func (i *Service) validateCapturePageTypeId(ctx context.Context, handler persistence.TransactionHandler, id int) error {
	_, err := i.cfg.Persistor.GetCapturePageTypeById(ctx, handler, id)
	if err != nil {
		return fmt.Errorf("invalid capture_page_type_id: %v", err)
	}
	return nil
}

// CreateCapturePages creates a new capture page.
func (i *Service) CreateCapturePages(ctx context.Context, params *model.CreateCapturePage) (*model.CapturePages, error) {
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

	capture_page := params.ToCapturePage()

	_, err = i.cfg.Persistor.GetCapturePageTypeById(ctx, tx, capture_page.CapturePageSetId)
	if err != nil {
		fmt.Println("the error --- ", err)
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid capture_pages_sets_id: %v", err),
		})
	}

	capture_page, err = i.cfg.Persistor.CreateCapturePages(ctx, tx, capture_page)
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

	return capture_page, nil
}

// ListCapturePages returns paginated capture pages
func (i *Service) ListCapturePages(
	ctx context.Context,
	filter *model.CapturePagesFilters,
) (*model.PaginatedCapturePages, error) {
	db, err := i.cfg.TxProvider.Db(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}

	paginated, err := i.cfg.Persistor.GetCapturePages(ctx, db, filter)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get capture pages: %v", err),
		})
	}

	return paginated, nil
}

// UpdateCapturePages updates an existing capture page.
func (i *Service) UpdateCapturePages(ctx context.Context, params *model.UpdateCapturePages) (*model.CapturePages, error) {
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

	if params.CapturePageSetId.Valid {
		if err := i.validateCapturePageTypeId(ctx, tx, params.CapturePageSetId.Int); err != nil {
			return nil, fmt.Errorf("capture_page_set_id: %w", err)
		}
	}

	capturepages, err := i.cfg.Persistor.UpdateCapturePages(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("update capture pages: %w", err)
	}

	return capturepages, nil
}

// DeleteCapturePages deletes a capture page by ID.
func (s *Service) DeleteCapturePages(ctx context.Context, params *model.DeleteCapturePages) error {
	tx, err := s.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}
	defer tx.Rollback(ctx)

	fmt.Println("the params id ---- ", params.ID)
	err = s.cfg.Persistor.DeleteCapturePages(ctx, tx, params.ID)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("delete capture page: %v", err),
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
