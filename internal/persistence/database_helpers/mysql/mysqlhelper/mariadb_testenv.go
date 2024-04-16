package mysqlhelper

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/global"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	m   sync.Mutex
	log = logger.New(context.TODO())
)

// TestCreateDatabase encapsulates creating a new test database and cleaning it up after tests.
func TestCreateDatabase(t *testing.T, cp *mysqlutil.ConnectionSettings) (*sqlx.DB, func(...bool)) {
	db, err := NewDbClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", cp.Database)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")

	db, err = NewDbClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	cleanup := func(ignoreFailures ...bool) {
		willIgnoreFailures := false
		if len(ignoreFailures) == 1 {
			willIgnoreFailures = ignoreFailures[0]
		}
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()

		dropStmt := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", cp.Database)
		_, err = db.ExecContext(ctx, dropStmt)
		if !willIgnoreFailures {
			assert.NoError(t, err, "unexpected error dropping test database from localDB")
		}
	}

	return db, cleanup
}

type CleanFn func(ignoreErrors ...bool)

// testSpawnMariaDB spawns a MariaDB instance from Docker.
// The main use-case for this function is to dynamically
// create an integration testing sandbox for MariaDB.
func testSpawnMariaDB(t *testing.T, cp *mysqlutil.ConnectionSettings) (*sqlx.DB, *mysqlutil.ConnectionSettings, CleanFn) {
	m.Lock()
	defer m.Unlock()

	// Set database logs on for testing
	boil.DebugMode = true
	mariaDBContainerName := strutil.GetUuidUnderscore()

	cfg := MariaDBConfig{
		RecreateContainer:   false,
		ContainerName:       mariaDBContainerName,
		Host:                cp.Host,
		Pass:                cp.Pass,
		ExposedInternalPort: cp.Port,
		ExposedExternalPort: cp.Port,
	}

	_, err := NewMariaDB(&cfg)
	require.NoError(t, err, "unexpected error spawning mariaDB")

	db, cleanup := TestCreateDatabase(t, cp)
	require.NoError(t, err, "unexpected error creating database")

	tables, err := Migrate(context.TODO(), cp, Recreate)
	require.NoError(t, err, "unexpected error during migration")
	require.NotEmpty(t, tables, "unexpected empty tables after migration")

	return db, cp, cleanup
}

func testExistingMariaDB(t *testing.T, cp *mysqlutil.ConnectionSettings) (*sqlx.DB, *mysqlutil.ConnectionSettings, CleanFn) {
	db, err := NewDbClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error creating test database for localDB")

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", cp.Database)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")

	db, err = NewDbClient(context.TODO(), &ClientOptions{
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
func TestGetMockMariaDB(t *testing.T) (*sqlx.DB, *mysqlutil.ConnectionSettings, CleanFn) {
	if os.Getenv(global.OsEnvAppDir) == "" {
		t.Error("migration dir empty")
	}

	cp := mysqlutil.ConnectionSettings{
		Host:     host,
		Pass:     pass,
		User:     "root", // mariadb's default user
		Port:     mappedPort,
		Database: strutil.GetUuidUnderscore(),
	}

	if os.Getenv("TEST_ENV_USE_EXISTING_MARIADB") == "1" {
		return testExistingMariaDB(t, &cp)
	}

	return testSpawnMariaDB(t, &cp)
}
