package mysqlhelper

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/global"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestMigrate(t *testing.T) {
	migrationDir := strings.TrimSpace(os.Getenv(global.OsEnvAppDir))
	if migrationDir == "" {
		t.Skip("global.OsEnvAppDir is unset, skipping...")
	}

	_, connectionParameters, cleanup := TestGetMockMariaDB(t)
	defer cleanup()

	log.Info("Migrating...")
	ctx := context.Background()
	tables, err := Migrate(ctx, connectionParameters, CreateIfNotExists)
	require.NoError(t, err, "unexpected error")
	require.NotEmpty(t, tables, "unexpected empty tables")
}
