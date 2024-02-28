package helpers

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/common"
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

	logger := common.GetLogger(context.TODO())

	_, connectionParameters, cleanup := TestGetMockMariaDB(t)
	defer cleanup()

	logger.Info("Migrating...")
	ctx := context.Background()
	tables, err := Migrate(ctx, connectionParameters, migrationDir, CreateIfNotExists)
	require.NoError(t, err, "unexpected error")
	require.NotEmpty(t, tables, "unexpected empty tables")
}
