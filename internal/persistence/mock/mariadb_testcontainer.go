package mock

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/dembygenesis/local.tools/internal/utils_common"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	retryAttempts = 30
	retryDelay    = 2 * time.Second
	mappedPort    = 3341
	containerName = "mock_test_database"
	host          = "localhost"
	user          = "demby"
	pass          = "secret"
)

var (
	logger = common.GetLogger(context.TODO())
)

func startContainerAndPingDB(ctx context.Context, req testcontainers.ContainerRequest, databaseName string) (*mysql.ConnectionParameters, testcontainers.Container, error) {
	var ctn testcontainers.Container
	var db *sqlx.DB
	var port nat.Port
	var err error

	for attempt := 0; attempt < retryAttempts; attempt++ {
		if attempt > 0 {
			logger.Printf("Retry attempt #%d\n", attempt)
			time.Sleep(retryDelay)
		}

		ctn, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            true,
		})
		if err != nil {
			logger.Println("Failed to start container:", err)
			continue
		}

		port, err = ctn.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", mappedPort)))
		if err != nil {
			logger.Println("Failed to get mapped port:", err)
			continue
		}

		connectionParameters := mysql.ConnectionParameters{
			Host:     host,
			User:     user,
			Pass:     pass,
			Port:     port.Int(),
			Database: databaseName,
		}
		db, err = mysql.GetClient(&mysql.ClientOptions{
			ConnString: connectionParameters.GetConnectionString(false),
			Close:      false,
		})
		if err != nil {
			logger.Println("Failed to get DB client:", err)
			continue
		}

		if err = db.Ping(); err != nil {
			logger.Println("Failed to ping DB:", err)
			continue
		}

		if err = db.Close(); err != nil {
			logger.Println("Failed to close DB:", err)
			continue
		}

		return &connectionParameters, ctn, nil
	}

	return nil, nil, fmt.Errorf("failed to start and ping DB after %d attempts", retryAttempts)
}

// TestGetMockMariaDBFromTestcontainers returns a live database instance for integration testing purposes.
// It first attempts to fetch an already running database if certain env variables,
// are set, and if it can't find them - it will spawn a new test container via the
// `testcontainer` library.
func TestGetMockMariaDBFromTestcontainers(t *testing.T) (*mysql.ConnectionParameters, func()) {
	// Old code for using an existing connection if
	// ENV vars are provided. No longer needed because
	// the function `TestGetMockMariaDB` exists.
	// But still keeping them here while continuing to refine.
	/*testDBKeys := []string{
		"TEST_DB_HOST",
		"TEST_DB_USER",
		"TEST_DB_PASS",
		"TEST_DB_PORT",
	}
	if utils_common.HasEnvs(testDBKeys) {
		return testGetExistingMysqlConnection(t)
	}*/

	return testSpawnMysqlContainer(t)
}

// testSpawnMysqlContainer returns connection parameters from env variables set.
// These env variables are going to be the running MYSQL's credentials.
func testGetExistingMysqlConnection(t *testing.T) (*mysql.ConnectionParameters, func()) {
	host := os.Getenv("TEST_DB_HOST")
	user := os.Getenv("TEST_DB_USER")
	pass := os.Getenv("TEST_DB_PASS")
	port, err := strconv.Atoi(os.Getenv("TEST_DB_PORT"))
	if err != nil {
		require.NoError(t, err, "unexpected error parsing TEST_DB_PORT")
	}

	cp := mysql.ConnectionParameters{
		Host: host,
		User: user,
		Pass: pass,
		Port: port,
	}

	err = cp.Validate(false)
	require.NoError(t, err, "unexpected test database validation error")

	db, err := mysql.GetClient(&mysql.ClientOptions{
		ConnString: cp.GetConnectionString(true),
		Close:      false,
	})
	require.NoError(t, err, "unexpected error getting client from database")

	databaseName := utils_common.GetUuidUnderscore()

	createStmt := fmt.Sprintf("CREATE DATABASE `%s`;", databaseName)
	_, err = db.Exec(createStmt)
	require.NoError(t, err, "unexpected error creating test database for localDB")
	cp.Database = databaseName

	cleanup := func() {
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE `%s`", databaseName))
		require.NoError(t, err, "unexpected error dropping test database from localDB")
	}

	return &cp, cleanup
}

// testSpawnMysqlContainer spawns a new mysql container from test container
func testSpawnMysqlContainer(t *testing.T) (*mysql.ConnectionParameters, func()) {
	ctx := context.Background()

	databaseName := utils_common.GetUuidUnderscore()

	req := testcontainers.ContainerRequest{
		Image:        "mariadb:latest",
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", mappedPort)},
		WaitingFor:   wait.ForAll(wait.ForExposedPort()),
		Env: map[string]string{
			"MARIADB_ALLOW_EMPTY_ROOT_PASSWORD": "true",
			"MYSQL_DATABASE":                    databaseName,
			"MYSQL_TCP_PORT":                    strconv.Itoa(mappedPort),
			"MYSQL_USER_HOST":                   "%s",
			"MYSQL_USER":                        user,
			"MYSQL_PASSWORD":                    pass,
		},
		Name: containerName,
	}

	connectionParameters, ctn, err := startContainerAndPingDB(ctx, req, databaseName)
	require.NoError(t, err, "failed to start and ping mysql test container")

	cleanup := func() {
		if err := ctn.Terminate(ctx); err != nil {
			logger.Errorf("Error terminating container: %v", err)
		}
	}

	return connectionParameters, cleanup
}

func testGetMysqlContainer(t *testing.T) (*mysql.ConnectionParameters, func()) {
	ctx := context.Background()

	databaseName := utils_common.GetUuidUnderscore()

	req := testcontainers.ContainerRequest{
		Image:        "mariadb:latest",
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", mappedPort)},
		WaitingFor:   wait.ForAll(wait.ForExposedPort()),
		Env: map[string]string{
			"MARIADB_ALLOW_EMPTY_ROOT_PASSWORD": "true",
			"MYSQL_TCP_PORT":                    strconv.Itoa(mappedPort),
			"MYSQL_USER_HOST":                   "%s",
			"MYSQL_USER":                        user,
			"MYSQL_PASSWORD":                    pass,
		},
		Name: containerName,
	}

	connectionParameters, ctn, err := startContainerAndPingDBV2(ctx, req, databaseName)
	require.NoError(t, err, "failed to start and ping mysql test container")

	cleanup := func() {
		// Not really needed, because Ryuk automatically terminates the connector.
		// What's needed here is the database cleanup functionality
		if err := ctn.Terminate(ctx); err != nil {
			logger.Errorf("Error terminating container: %v", err)
		}
	}

	return connectionParameters, cleanup
}

func startContainerAndPingDBV2(ctx context.Context, req testcontainers.ContainerRequest, databaseName string) (*mysql.ConnectionParameters, testcontainers.Container, error) {
	var ctn testcontainers.Container
	var db *sqlx.DB
	var port nat.Port
	var err error

	for attempt := 0; attempt < retryAttempts; attempt++ {
		if attempt > 0 {
			logger.Printf("Retry attempt #%d\n", attempt)
			time.Sleep(retryDelay)
		}

		ctn, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            true,
		})
		if err != nil {
			logger.Println("Failed to start container:", err)
			continue
		}

		port, err = ctn.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", mappedPort)))
		if err != nil {
			logger.Println("Failed to get mapped port:", err)
			continue
		}

		connectionParameters := mysql.ConnectionParameters{
			Host: host,
			User: user,
			Pass: pass,
			Port: port.Int(),
		}

		db, err = mysql.GetClient(&mysql.ClientOptions{
			ConnString: connectionParameters.GetConnectionString(false),
			Close:      false,
		})
		if err != nil {
			logger.Println("Failed to get DB client:", err)
			continue
		}

		if err = db.Ping(); err != nil {
			logger.Println("Failed to ping DB:", err)
			continue
		}

		if err = db.Close(); err != nil {
			logger.Println("Failed to close DB:", err)
			continue
		}

		return &connectionParameters, ctn, nil
	}

	return nil, nil, fmt.Errorf("failed to start and ping DB after %d attempts", retryAttempts)
}
