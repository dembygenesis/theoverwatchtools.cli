package mysqlstore

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
)

var (
	ErrCatNil = errors.New("category provided is nil")
	ErrOrgNil = errors.New("organization provided is nil")
)

type Config struct {
	Logger        *logrus.Entry              `json:"logger" validate:"required"`
	QueryTimeouts *persistence.QueryTimeouts `json:"query_timeouts" validate:"required"`
}

func (m *Config) Validate() error {
	return validationutils.Validate(m)
}

// Repository is the main struct that
// handles crud operations for MYSQL.
type Repository struct {
	cfg *Config
}

// New creates a new instance of a Repository.
func New(cfg *Config) (*Repository, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %v", err)
	}

	return &Repository{cfg}, nil
}
