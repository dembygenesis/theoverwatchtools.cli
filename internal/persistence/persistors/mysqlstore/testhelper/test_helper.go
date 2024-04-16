package testhelper

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

// DropTable drops a table from the database,
// for testing purposes.
func DropTable(t *testing.T, db *sqlx.DB, table string) {
	t.Helper()
	stmts := []string{
		"SET FOREIGN_KEY_CHECKS = 0;",
		"DROP TABLE " + table + ";",
		"SET FOREIGN_KEY_CHECKS = 1;",
	}

	for _, stmt := range stmts {
		_, err := db.Exec(stmt)
		require.NoErrorf(t, err, "stmt %s: %v", table, err)
	}
}
