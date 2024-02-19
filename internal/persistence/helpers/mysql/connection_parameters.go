package mysql

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/utils_common"
	"strconv"
)

type ConnectionParameters struct {
	Host     string `mapstructure:"host" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Pass     string `mapstructure:"pass" validate:"required"`
	Database string `mapstructure:"database" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required,int_greater_than_zero"`
}

func (c *ConnectionParameters) Validate() error {
	err := utils_common.ValidateStruct(c)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func (c *ConnectionParameters) GetConnectionString(excludeDatabase bool) string {
	portStr := strconv.Itoa(c.Port)

	if excludeDatabase {
		return c.User + ":" +
			c.Pass + "@tcp(" +
			c.Host + ":" + portStr + ")/?charset=utf8&parseTime=true"
	}
	return c.User + ":" +
		c.Pass + "@tcp(" +
		c.Host + ":" + portStr + ")/" +
		c.Database + "?charset=utf8&parseTime=true"
}
