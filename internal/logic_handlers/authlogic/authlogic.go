package authlogic

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
)

type Config struct {
	TxProvider persistence.TransactionProvider `json:"tx_provider" validate:"required"`
	Logger     *logrus.Entry                   `json:"logger" validate:"required"`
}

func (i *Config) Validate() error {
	return validationutils.Validate(i)
}

type Impl struct {
	cfg *Config
}

func New(cfg *Config) (*Impl, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}
	return &Impl{cfg}, nil
}

func (i *Impl) GetCategories() {

}
