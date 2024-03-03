package helpers

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/globals"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestMigrate(t *testing.T) {
	migrationDir := strings.TrimSpace(os.Getenv(globals.OsEnvMigrationDir))
	if migrationDir == "" {
		t.Skip("global.OsEnvMigrationDir is unset, skipping...")
	}

	_, connectionParameters, cleanup := TestGetMockMariaDB(t)
	defer cleanup()

	log.Info("Migrating...")
	ctx := context.Background()
	tables, err := Migrate(ctx, connectionParameters, migrationDir, CreateIfNotExists)
	require.NoError(t, err, "unexpected error")
	require.NotEmpty(t, tables, "unexpected empty tables")
}
