package mock

import (
	"fmt"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/docker_testenv"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/dembygenesis/local.tools/internal/utils_common"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

var (
	m sync.Mutex
)

// testSpawnMariaDB spawns a MariaDB instance from Docker.
// The main use-case for this function is to dynamically
// create an integration testing sandbox for MariaDB.
func testSpawnMariaDB(t *testing.T, cp *mysql.ConnectionParameters) (*mysql.ConnectionParameters, func()) {
	m.Lock()
	defer m.Unlock()

	mariaDBContainerName := utils_common.GetUuidUnderscore()

	cfg := docker_testenv.MariaDBConfig{
		RecreateContainer:   false,
		ContainerName:       mariaDBContainerName,
		Host:                cp.Host,
		Pass:                cp.Pass,
		ExposedInternalPort: cp.Port,
		ExposedExternalPort: cp.Port,
	}

	// Upsert-spawn a mariadb instance
	_, err := docker_testenv.NewMariaDB(&cfg)
	require.NoError(t, err, "unexpected error spawning mariaDB")

	db, err := mysql.GetClient(&mysql.ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", cp.Database)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")

	db, err = mysql.GetClient(&mysql.ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	cleanup := func() {
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE `%s`", cp.Database))
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

	return cp, cleanup
}

func testExistingMariaDB(t *testing.T, cp *mysql.ConnectionParameters) (*mysql.ConnectionParameters, func()) {
	db, err := mysql.GetClient(&mysql.ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error creating test database for localDB")

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", cp.Database)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")

	db, err = mysql.GetClient(&mysql.ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	cleanup := func() {
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE `%s`", cp.Database))
		require.NoError(t, err, "unexpected error dropping test database from localDB")
	}

	require.NoError(t, err, "unexpected error connecting to existing maria db")
	return cp, cleanup
}

// TestGetMockMariaDB returns a live database instance for integration testing purposes.
// It first attempts to fetch an already running database if certain env variables,
// are set, and if it can't find them - it will spawn a new test container via the
// `testcontainer` library.
func TestGetMockMariaDB(t *testing.T) (*mysql.ConnectionParameters, func()) {
	cp := mysql.ConnectionParameters{
		Host:     host,
		Pass:     pass,
		User:     "root", // mariadb's default user
		Port:     mappedPort,
		Database: utils_common.GetUuidUnderscore(),
	}

	if os.Getenv("TEST_ENV_USE_EXISTING_MARIADB") == "1" {
		return testExistingMariaDB(t, &cp)
	}

	return testSpawnMariaDB(t, &cp)
}
