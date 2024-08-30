package mysqlhelper

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/database/migration"
	"github.com/dembygenesis/local.tools/internal/global"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"strconv"
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

func testExistingMariaDB(t *testing.T, cp *mysqlutil.ConnectionSettings) (*sqlx.DB, *mysqlutil.ConnectionSettings, CleanFn) {
	db, err := NewDbClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error creating test database for localDB")

	cp.Database = strutil.GetUuidUnderscore()
	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", cp.Database)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")

	db, err = NewDbClient(context.TODO(), &ClientOptions{
		ConnString: cp.GetConnectionString(false),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	tables, err := MigrateEmbedded(context.TODO(), cp, migration.Migrations, CreateIfNotExists)
	require.NoError(t, err, "unexpected error during migration")
	require.NotEmpty(t, tables, "unexpected empty tables after migration")

	cleanup := func(ignoreFailures ...bool) {
		willIgnoreFailures := false
		if len(ignoreFailures) == 1 {
			willIgnoreFailures = ignoreFailures[0]
		}
		dropStmt := fmt.Sprintf("DROP DATABASE `%s`", cp.Database)
		_, err = db.Exec(dropStmt)
		if !willIgnoreFailures {
			require.NoError(t, err, "unexpected error dropping test database from localDB")
		}

		err = db.Close()
		assert.NoError(t, err, "unexpected error closing db connection")
	}

	require.NoError(t, err, "unexpected error connecting to existing maria db")
	return db, cp, cleanup
}

// testSpawnMariaDB spawns a MariaDB instance from Docker.
// The main use-case for this function is to dynamically
// create an integration testing sandbox for MariaDB.
func testSpawnMariaDB(t *testing.T, cp *mysqlutil.ConnectionSettings) (*sqlx.DB, *mysqlutil.ConnectionSettings, CleanFn) {
	m.Lock()
	defer m.Unlock()

	boil.DebugMode = true
	mariaDBContainerName := "test_mariadb"

	cfg := MariaDBConfig{
		RecreateContainer:   false,
		ContainerName:       mariaDBContainerName,
		Host:                cp.Host,
		Pass:                cp.Pass,
		Database:            cp.Database,
		ExposedInternalPort: cp.Port,
		ExposedExternalPort: cp.Port,
	}

	// This will create the mariadb container, and the database specified in the env
	_, err := NewMariaDB(&cfg)
	require.NoError(t, err, "unexpected error spawning mariaDB")

	cp.Database = strutil.GetUuidUnderscore()

	// This will create the new database inside it
	db, cleanup := TestCreateDatabase(t, cp)
	require.NoError(t, err, "unexpected error creating database")

	db.SetMaxIdleConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 30)

	tables, err := MigrateEmbedded(context.TODO(), cp, migration.Migrations, CreateIfNotExists)
	require.NoError(t, err, "unexpected error during migration")
	require.NotEmpty(t, tables, "unexpected empty tables after migration")

	return db, cp, cleanup
}

func applyMigration(t *testing.T, cp *mysqlutil.ConnectionSettings) (*sqlx.DB, *mysqlutil.ConnectionSettings, CleanFn) {
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

	tables, err := Migrate(context.TODO(), cp, Recreate)
	require.NoError(t, err, "unexpected error during migration")
	require.NotEmpty(t, tables, "unexpected empty tables after migration")

	cleanup := func(ignoreFailures ...bool) {
		/*willIgnoreFailures := false
		if len(ignoreFailures) == 1 {
			willIgnoreFailures = ignoreFailures[0]
		}
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE `%s`", cp.Database))
		if !willIgnoreFailures {
			require.NoError(t, err, "unexpected error dropping test database from localDB")
		}*/
	}

	require.NoError(t, err, "unexpected error connecting to existing maria db")
	return db, cp, cleanup
}

// TestGetMockMariaDB returns a live database instance for integration testing purposes.
// It first attempts to fetch an already running database if certain env variables,
// are set, and if it can't find them - it will spawn a new test container via the
// `testcontainer` library.
func TestGetMockMariaDB(t *testing.T) (*sqlx.DB, *mysqlutil.ConnectionSettings, CleanFn) {
	osTestCreds := []string{
		global.OsEnvDbTestHost,
		global.OsEnvDbTestPort,
		global.OsEnvDbTestUser,
		global.OsEnvDbTestPass,
		global.OsEnvDbTestDatabase,
	}

	// To complex, just make sure these are existing
	for i := range osTestCreds {
		if os.Getenv(osTestCreds[i]) == "" {
			t.Errorf("missing env var: %s", osTestCreds[i])
		}
	}

	port, err := strconv.Atoi(os.Getenv(global.OsEnvDbTestPort))
	require.NoError(t, err, "unexpected error converting port to int")

	cp := &mysqlutil.ConnectionSettings{
		Host:     os.Getenv(global.OsEnvDbTestHost),
		Pass:     os.Getenv(global.OsEnvDbTestPass),
		User:     os.Getenv(global.OsEnvDbTestUser),
		Database: os.Getenv(global.OsEnvDbTestDatabase),
		Port:     port,
	}

	if os.Getenv("THEOVERWATCHTOOLS_DB_USE_EXISTING_MARIADB") == "1" {
		return testExistingMariaDB(t, cp)
	}

	return testSpawnMariaDB(t, cp)
}
