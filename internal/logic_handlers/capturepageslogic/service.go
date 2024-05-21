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
