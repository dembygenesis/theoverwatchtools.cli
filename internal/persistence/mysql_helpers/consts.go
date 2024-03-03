package helpers

import "time"

const (
	retryAttempts = 30
	retryDelay    = 2 * time.Second
	mappedPort    = 3341
	containerName = "mock_test_database"
	host          = "localhost"
	user          = "demby"
	pass          = "secret"
)
