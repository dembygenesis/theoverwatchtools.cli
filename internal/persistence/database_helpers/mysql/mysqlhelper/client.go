package mysqlhelper

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
	UseCache    bool
}

var cachedDb *sqlx.DB

func NewDbClient(ctx context.Context, o *ClientOptions) (*sqlx.DB, error) {
	if o.UseCache && cachedDb != nil {
		return cachedDb, nil
	}

	db, err := sqlx.Open("mysql", o.ConnString)
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)

	retries := 5
	var pingSuccess bool

	for i := 0; i < retries && !pingSuccess; i++ {
		if o.PingTimeout != 0 {
			pingCtx, cancel := context.WithTimeout(ctx, o.PingTimeout)
			defer cancel()
			err = db.PingContext(pingCtx)
		} else {
			err = db.Ping()
		}

		if err != nil {
			if o.PingTimeout == 0 {
				log.Warnf("==== ping failed, sleeping for 2 seconds, retry count: %d\n", retries-i)
				time.Sleep(2 * time.Second)
			}
		} else {
			pingSuccess = true
			break
		}
	}

	if !pingSuccess {
		return nil, fmt.Errorf("ping failed after %d retries: %v, for conn str: %s", retries, err, o.ConnString)
	}

	if o.Close {
		if err = db.Close(); err != nil {
			return nil, fmt.Errorf("close: %v", err)
		}
	}

	if o.UseCache {
		cachedDb = db
	}
	return db, nil
}
