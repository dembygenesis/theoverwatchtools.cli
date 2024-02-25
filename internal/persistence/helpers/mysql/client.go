package mysql

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type ClientOptions struct {
	PingTimeout time.Duration
	ConnString  string
	Close       bool
}

func GetClient(o *ClientOptions) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", o.ConnString)
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}

	if o.PingTimeout != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), o.PingTimeout)
		defer cancel()
		if err = db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("ping: %v", err)
		}
	} else {
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("ping: %v", err)
		}
	}

	if o.Close {
		if err = db.Close(); err != nil {
			return nil, fmt.Errorf("close: %v", err)
		}
	}

	return db, nil
}
