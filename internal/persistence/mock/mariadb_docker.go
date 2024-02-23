package mock

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/docker_testenv"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/dembygenesis/local.tools/internal/utils_common"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestGetMockMariaDBFromDocker returns a live database instance for integration testing purposes.
// It first attempts to fetch an already running database if certain env variables,
// are set, and if it can't find them - it will spawn a new test container via the
// `testcontainer` library.
func TestGetMockMariaDBFromDocker(t *testing.T) (*mysql.ConnectionParameters, func()) {
	databaseName := utils_common.GetUuidUnderscore()
	mariaDBContainerName := utils_common.GetUuidUnderscore()

	cfg := docker_testenv.MariaDBConfig{
		RecreateContainer:   false,
		ContainerName:       mariaDBContainerName,
		Host:                host,
		Pass:                pass,
		ExposedInternalPort: mappedPort,
		ExposedExternalPort: mappedPort,
	}

	// Upsert-spawn a mariadb instance
	_, err := docker_testenv.NewMariaDB(&cfg)
	require.NoError(t, err, "unexpected error spawning mariaDB")

	cp := mysql.ConnectionParameters{
		Host: cfg.Host,
		Pass: cfg.Pass,
		User: "root", // mariadb's default user
		Port: cfg.ExposedInternalPort,
	}

	db, err := mysql.GetClient(cp.GetConnectionString(true))
	require.NoError(t, err, "unexpected error getting client from database")

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", databaseName)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")
	cp.Database = databaseName

	db, err = mysql.GetClient(cp.GetConnectionString(false))
	require.NoError(t, err, "unexpected error getting client from database")

	cleanup := func() {
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE `%s`", databaseName))
		require.NoError(t, err, "unexpected error dropping test database from localDB")

		/**
		I am not dropping this container, so it can be reused
		by other tests involving a MYSQL db.

		An ideal outcome would be to only drop this after all tests complete
		(need to look more into that).
		*/
		/*err = ctn.Cleanup(context.Background())
		require.NoError(t, err, "error cleaning up container")*/
	}

	return &cp, cleanup
}
