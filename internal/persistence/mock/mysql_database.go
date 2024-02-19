package mock

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/common"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	retryAttempts = 30
	retryDelay    = 2 * time.Second
	mappedPort    = 3306
	containerName = "mock_test_database"
	host          = "localhost"
	user          = "demby"
	pass          = "secret"
)

var (
	m      sync.Mutex
	logger = common.GetLogger(nil)
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
		db, err = mysql.GetClient(&connectionParameters)
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

func TestMockMysqlDatabase(t *testing.T) (*mysql.ConnectionParameters, func()) {
	ctx := context.Background()

	databaseName := strings.ReplaceAll(uuid.New().String(), "-", "_")
	m.Lock()
	defer m.Unlock()

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
