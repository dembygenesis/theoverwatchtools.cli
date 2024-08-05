package organizationlogic

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
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

func (i *Service) validateOrganizationTypeId(ctx context.Context, handler persistence.TransactionHandler, id int) error {
	_, err := i.cfg.Persistor.GetOrganizationById(ctx, handler, id)
	if err != nil {
		return fmt.Errorf("invalid organization_id: %v", err)
	}
	return nil
}
