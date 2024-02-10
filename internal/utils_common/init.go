package utils_common

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Entry
)

// init loads all the pre-configurations required for
// utils_common functions.
func init() {
	configValidate()
	logger = common.GetLogger(context.Background())

}
