package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func GetClient(connectionParameters *ConnectionParameters) (*sqlx.DB, error) {
	connString :=
		connectionParameters.User + ":" +
			connectionParameters.Pass + "@tcp(" +
			connectionParameters.Host + ":" +
			strconv.Itoa(connectionParameters.Port) + ")/" +
			connectionParameters.Database + "?charset=utf8&parseTime=true"
	return sqlx.Open("mysql", connString)
}
