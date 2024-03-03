package helpers

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/globals"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/utils/string_util"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	m            sync.Mutex
	log          = logger.New(context.TODO())
	cleanupsHasp = make(map[string]int)
)

// testSpawnMariaDB spawns a MariaDB instance from Docker.
// The main use-case for this function is to dynamically
// create an integration testing sandbox for MariaDB.
func testSpawnMariaDB(t *testing.T, cp *ConnectionParameters) (*sqlx.DB, *ConnectionParameters, func(ignoreErrors ...bool)) {
	m.Lock()
	defer m.Unlock()

	mariaDBContainerName := string_util.GetUuidUnderscore()

	cfg := MariaDBConfig{
		RecreateContainer:   false,
		ContainerName:       mariaDBContainerName,
		Host:                cp.Host,
		Pass:                cp.Pass,
		ExposedInternalPort: cp.Port,
		ExposedExternalPort: cp.Port,
	}

	// Upsert-spawn a mariadb instance
	_, err := NewMariaDB(&cfg)
	require.NoError(t, err, "unexpected error spawning mariaDB")

	db, err := NewClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", cp.Database)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")

	db, err = NewClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	tables, err := Migrate(context.TODO(), cp, os.Getenv(globals.OsEnvMigrationDir), Recreate)
	require.NoError(t, err, "unexpected error during migration")
	require.NotEmpty(t, tables, "unexpected empty tables after migration")

	cleanup := func(ignoreFailures ...bool) {
		willIgnoreFailures := false
		if len(ignoreFailures) == 1 {
			willIgnoreFailures = ignoreFailures[0]
		}
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()

		dropStmt := fmt.Sprintf("DROP DATABASE `%s`;", cp.Database)
		_, err = db.ExecContext(ctx, dropStmt)
		if !willIgnoreFailures {
			assert.NoError(t, err, "unexpected error dropping test database from localDB")

			// Additional info
			if err != nil {
				log.Info("cleanups hash ", cleanupsHasp)
				showDatabaseStmt := "SHOW DATABASES;"
				databases, err := queryAsStringArrays(ctx, db, showDatabaseStmt, nil)
				require.NoError(t, err, "unexpected error fetching the databases")
				log.Info("drop statement", dropStmt)
				log.Info("databases", databases)
			}

		}

		/**
		I am not dropping this container, so it can be reused
		by other tests involving a MYSQL db.

		An ideal outcome would be to only drop this after all tests complete
		(need to look more into that).
		*/
		/*err = ctn.Cleanup(context.Background())
		require.NoError(t, err, "error cleaning up container")*/
	}

	return db, cp, cleanup
}

func testExistingMariaDB(t *testing.T, cp *ConnectionParameters) (*sqlx.DB, *ConnectionParameters, func(ignoreErrors ...bool)) {
	db, err := NewClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error creating test database for localDB")

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", cp.Database)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")

	db, err = NewClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	cleanup := func(ignoreFailures ...bool) {
		willIgnoreFailures := false
		if len(ignoreFailures) == 1 {
			willIgnoreFailures = ignoreFailures[0]
		}
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE `%s`", cp.Database))
		if !willIgnoreFailures {
			require.NoError(t, err, "unexpected error dropping test database from localDB")
		}
	}

	require.NoError(t, err, "unexpected error connecting to existing maria db")
	return db, cp, cleanup
}

// TestGetMockMariaDB returns a live database instance for integration testing purposes.
// It first attempts to fetch an already running database if certain env variables,
// are set, and if it can't find them - it will spawn a new test container via the
// `testcontainer` library.
func TestGetMockMariaDB(t *testing.T) (*sqlx.DB, *ConnectionParameters, func(ignoreErrors ...bool)) {
	if os.Getenv(globals.OsEnvMigrationDir) == "" {
		t.Error("migration dir empty")
	}

	cp := ConnectionParameters{
		Host:     host,
		Pass:     pass,
		User:     "root", // mariadb's default user
		Port:     mappedPort,
		Database: string_util.GetUuidUnderscore(),
	}

	if os.Getenv("TEST_ENV_USE_EXISTING_MARIADB") == "1" {
		return testExistingMariaDB(t, &cp)
	}

	return testSpawnMariaDB(t, &cp)
}
