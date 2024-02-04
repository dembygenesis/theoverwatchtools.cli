package models

import "errors"

var (
	ErrConfigNil   = errors.New("config nil")
	ErrRootMissing = errors.New("missing root")
)
