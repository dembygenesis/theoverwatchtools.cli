package mysqlhelper

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/lib/dockerenv"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

// MariaDBConfig is our setting struct for creating
// maria db instances.
//
// It purposefully doesn't have "User" as an attribute,
// because it is hardcoded as "root".
//
// The primary, and only purpose of this container
// is for testing, hence we want full root access.
//
// This could "change" in the future when we want to test
// granular user privileges.
type MariaDBConfig struct {
	RecreateContainer   bool   `json:"recreate_container" mapstructure:"recreate_container"`
	ContainerName       string `json:"container_name" mapstructure:"container_name" validate:"required"`
	Host                string `json:"host" mapstructure:"host" validate:"required"`
	Pass                string `json:"pass" mapstructure:"pass" validate:"required"`
	Database            string `json:"database" mapstructure:"database" validate:"required"`
	ExposedInternalPort int    `json:"exposed_internal_port" mapstructure:"exposed_internal_port" validate:"required,greater_than_zero"`
	ExposedExternalPort int    `json:"exposed_external_port" mapstructure:"exposed_external_port" validate:"required,greater_than_zero"`
}

func (m *MariaDBConfig) Validate() error {
	return validationutils.Validate(m)
}

func (m *MariaDBConfig) GetEnvVars() []string {
	return []string{
		"MYSQL_ROOT_HOST=%",
		fmt.Sprintf("MYSQL_ROOT_PASSWORD=%v", m.Pass),
		fmt.Sprintf("MYSQL_TCP_PORT=%v", m.ExposedInternalPort),
		fmt.Sprintf("MYSQL_DATABASE=%v", m.Database),
	}
}

// func NewMariaDBTestContainer(t *testing.T, cfg *MariaDBConfig) (*dockerenv.DockerEnv, error) {

func NewMariaDBTestContainer(t *testing.T, cfg *mysqlutil.ConnectionSettings) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:           "mysql:8.0",
		AlwaysPullImage: true,
		Name:            "test_container_mysql",
		ExposedPorts:    []string{"3306/tcp"},
		WaitingFor:      wait.ForLog("X Plugin ready for connections. Bind-address:"),
		Env: map[string]string{
			"MYSQL_ROOT_HOST":     "%",
			"MYSQL_ROOT_PASSWORD": cfg.Pass,
		},
		Cmd: []string{
			"--default-authentication-plugin=mysql_native_password",
			"--sql_mode=STRICT_TRANS_TABLES",
		},
	}

	// Lets try this a couple of times
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		t.Fatalf("Could not start mysql: %s", err)
	}

	port, err := container.MappedPort(ctx, "3306/tcp")
	require.NoError(t, err, "unexpected error getting mapped port")

	defer func() {
		/*if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Could not stop redis: %s", err)
		}*/
	}()
	cfg.Port = port.Int()
}

// NewMariaDB spawns a new MariaDB instance.
func NewMariaDB(cfg *MariaDBConfig) (*dockerenv.DockerEnv, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	fmt.Println("=== hello cfg.ContainerName:", cfg.ContainerName)

	const (
		image = "mysql:8.0"
	)

	dm, err := dockerenv.New(&dockerenv.ContainerConfig{
		Name:        cfg.ContainerName,
		Image:       image,
		Env:         cfg.GetEnvVars(),
		ExposedPort: cfg.ExposedInternalPort,
		HostPort:    cfg.ExposedInternalPort,
		WaitFor:     fmt.Sprintf("port: %v", cfg.ExposedInternalPort),
		Cmd: []string{
			"--default-authentication-plugin=mysql_native_password",
			"--sql-mode=STRICT_TRANS_TABLES",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("docker manager: %v", err)
	}

	// Ah, I think it's returning the already existing container
	_, err = dm.UpsertContainer(context.Background(), cfg.RecreateContainer)
	if err != nil {
		return nil, fmt.Errorf("upsert container: %v", err)
	}

	// I should have a check here tho if the database exists... Else,
	// manually recreate it... But why would it disppear on the first place?
	// Shouldn't happen...
	// But i think it's nice because... It will handle some weird edge case scenarios

	return dm, err
}
