package utils_common

import (
	"os"
	"strings"
)

// HasEnvs attempts to extract env variables
// from all keys provided, and returns
// true if all of them have a value
func HasEnvs(keys []string) bool {
	for _, key := range keys {
		key = strings.TrimSpace(key)
		envVar := strings.TrimSpace(os.Getenv(key))
		if key == "" || envVar == "" {
			return false
		}
	}
	return true
}
