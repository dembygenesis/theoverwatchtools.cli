package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetClient(connString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("err: %v, with conn string: %s", err, connString)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}

	return db, nil
}
