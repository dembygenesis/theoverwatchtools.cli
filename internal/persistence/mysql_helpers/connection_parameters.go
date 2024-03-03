package helpers

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/lib/validation"
	"strconv"
	"strings"
)

type ConnectionParameters struct {
	Host     string `mapstructure:"host" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Pass     string `mapstructure:"pass" validate:"required"`
	Database string `mapstructure:"database"`
	Port     int    `mapstructure:"port" validate:"required,greater_than_zero"`
}

func (c *ConnectionParameters) Validate(includeDatabase bool) error {
	err := validation.Validate(c)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if includeDatabase && strings.TrimSpace(c.Database) == "" {
		return fmt.Errorf("missing database")
	}

	return nil
}

func (c *ConnectionParameters) GetConnectionString(excludeDatabase bool) string {
	portStr := strconv.Itoa(c.Port)

	const template = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true"
	const templateExcludeDatabase = "%s:%s@tcp(%s:%s)/?charset=utf8&parseTime=true"

	if excludeDatabase {
		return fmt.Sprintf(templateExcludeDatabase,
			c.User,
			c.Pass,
			c.Host,
			portStr,
		)
	}

	return fmt.Sprintf(template,
		c.User,
		c.Pass,
		c.Host,
		portStr,
		c.Database,
	)
}
