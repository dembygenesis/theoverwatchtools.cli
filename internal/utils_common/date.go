package utils_common

import (
	"time"
)

const (
	layoutDate = "2006-01-02"
	layoutTS   = "2006-01-02T15:04:05.000Z"
)

func IsValidDate(s string) bool {
	_, err := time.Parse(layoutDate, s)
	return err != nil
}
