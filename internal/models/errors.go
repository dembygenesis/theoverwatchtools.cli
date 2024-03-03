package models

import "errors"

var (
	ErrStructNil              = errors.New("struct nil")
	ErrConfigNil              = errors.New("config nil")
	ErrRootMissing            = errors.New("missing root")
	ErrContainerIdMissing     = errors.New("error, missing container id")
	ErrDatabaseNil            = errors.New("database is nil")
	ErrTimeoutNil             = errors.New("timeout is nil")
	ErrOptsNil                = errors.New("opts is nil")
	ErrInvalidCreateMode      = errors.New("invalid create mode")
	ErrEmptyDir               = errors.New("empty dir")
	ErrCommitInvalidDB        = errors.New("commit invalid for db")
	ErrRollbackInvalid        = errors.New("rollback invalid")
	ErrUnidentifiedTxType     = errors.New("unidentified tx type")
	ErrNotATransactionHandler = errors.New("not a transaction handler")
	ErrEmptyName              = errors.New("empty name")
	ErrPersistorNil           = errors.New("persistor is nil")
)
