package models

import "errors"

var (
	ErrConfigNil          = errors.New("config nil")
	ErrRootMissing        = errors.New("missing root")
	ErrContainerIdMissing = errors.New("error, missing container id")
	ErrDatabaseNil        = errors.New("database is nil")
	ErrTimeoutNil         = errors.New("timeout is nil")
)
