package categorymysql

import "context"

type transactor interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}
