package mysqltx

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testContext = context.TODO()
	testLogger  = logger.New(testContext)
)

func Test_mySQLRepository_TxHandler_Tx_Fail(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup(true)

	cfg := &Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	}

	r, err := New(cfg)
	require.NoError(t, err, "unexpected error initialising repository")
	require.NotNil(t, r, "unexpected error")

	err = db.Close()
	require.NoError(t, err, "unexpected error closing db")

	txHandlerTx, err := r.Tx(testContext)
	require.Error(t, err, "unexpected error fetching tx from tx handler")
	require.Nil(t, txHandlerTx, "unexpected non-nil txHandlerTx")
}

func Test_mySQLRepository_TxHandler_Tx(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	}

	r, err := New(cfg)
	require.NoError(t, err, "unexpected error initialising repository")
	require.NotNil(t, r, "unexpected error")

	txHandlerTx, err := r.Tx(testContext)
	require.NoError(t, err, "unexpected error fetching tx from tx handler")
	require.NotNil(t, txHandlerTx, "unexpected nil txHandlerTx")
}

func Test_mySQLRepository_TxHandler_Db(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	}

	r, err := New(cfg)
	require.NoError(t, err, "unexpected error initialising repository")
	require.NotNil(t, r, "unexpected error")

	txHandlerDb, err := r.Db(testContext)
	require.NoError(t, err, "unexpected error fetching db from tx handler")
	require.NotNil(t, txHandlerDb, "unexpected txHandlerDb")
}

func Test_mySQLRepository_TxHandler_New(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	}

	r, err := New(cfg)
	require.NoError(t, err, "unexpected error initialising repository")
	require.NotNil(t, r, "unexpected error")
}

func Test_mySQLRepository_TxHandler_DB(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	}

	mySQLTxHandler, err := New(cfg)
	require.NoError(t, err, "unexpected error initialising repository")
	require.NotNil(t, mySQLTxHandler, "unexpected error")

	dbTypedMySQLTxHandler, err := mySQLTxHandler.Db(testContext)
	require.NoError(t, err, "unexpected error fetching mysql tx handler DB")
	require.NotNil(t, dbTypedMySQLTxHandler.db, "unexpected nil database from tx handler DB")

	createStmt := `
		CREATE TABLE person(
			id int(11) AUTO_INCREMENT,
			firstname VARCHAR(255) NOT NULL DEFAULT '',
			lastname VARCHAR(255) NOT NULL DEFAULT '',
			PRIMARY KEY (id)
		);
	`

	_, err = dbTypedMySQLTxHandler.db.Exec(createStmt)
	require.NoError(t, err, "unexpected table creation")

	showTablesStmt := "SHOW TABLES;"
	rows, err := mysqlutil.QueryAsStringArrays(
		testContext,
		dbTypedMySQLTxHandler.db,
		showTablesStmt,
		nil,
	)
	require.NoError(t, err, "unexpected error for fetching tables")
	require.True(t, len(rows) > 0, "unexpected - rows are 0")

	hasPersonTable := false
	for _, row := range rows {
		if row[0] == "person" {
			hasPersonTable = true
		}
	}
	require.True(t, hasPersonTable)
}

func Test_mySQLRepository_TxHandler_TX(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

	cfg := &Config{
		Logger:       testLogger,
		Db:           db,
		DatabaseName: cp.Database,
	}

	mySQLTxHandler, err := New(cfg)
	require.NoError(t, err, "unexpected error initializing repository")
	require.NotNil(t, mySQLTxHandler, "Handler instance should not be nil")

	dbHandler, err := mySQLTxHandler.Db(testContext)
	require.NoError(t, err, "unexpected error getting db handler")

	createTableSQL := `CREATE TABLE person (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	)`
	_, err = dbHandler.db.Exec(createTableSQL)
	require.NoError(t, err, "unexpected error creating table")

	txHandler, err := mySQLTxHandler.Tx(testContext)
	require.NoError(t, err, "unexpected error starting a transaction")

	for i := 1; i <= 5; i++ {
		insertSQL := "INSERT INTO person (name) VALUES (?)"
		_, err := txHandler.tx.Exec(insertSQL, "Test Name "+string(rune(i)))
		require.NoError(t, err, "unexpected error inserting record")
	}

	countBeforeCommit := 0
	err = dbHandler.db.Get(&countBeforeCommit, "SELECT COUNT(*) FROM person")
	require.NoError(t, err, "unexpected error querying before commit")
	require.Equal(t, 0, countBeforeCommit, "records should not be visible before commit")

	err = txHandler.Commit(context.Background())
	require.NoError(t, err, "unexpected error committing transaction")

	countAfterCommit := 0
	err = dbHandler.db.Get(&countAfterCommit, "SELECT COUNT(*) FROM person")
	require.NoError(t, err, "unexpected error querying after commit")
	require.Equal(t, 5, countAfterCommit, "records should be visible after commit")
}
