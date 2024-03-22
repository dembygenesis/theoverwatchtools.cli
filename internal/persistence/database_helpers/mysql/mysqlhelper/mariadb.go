package mysqlhelper

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/lib/dockerenv"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
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
	ExposedInternalPort int    `json:"exposed_internal_port" mapstructure:"exposed_internal_port" validate:"required,greater_than_zero"`
	ExposedExternalPort int    `json:"exposed_external_port" mapstructure:"exposed_external_port" validate:"required,greater_than_zero"`
}

func (m *MariaDBConfig) Validate() error {
	return validationutils.Validate(m)
}

func (m *MariaDBConfig) GetEnvVars() []string {
	return []string{
		"MYSQL_ROOT_HOST=%",
		fmt.Sprintf("MARIADB_ROOT_PASSWORD=%v", m.Pass),
		fmt.Sprintf("MYSQL_TCP_PORT=%v", m.ExposedInternalPort),
	}
}

// NewMariaDB spawns a new MariaDB instance.
func NewMariaDB(cfg *MariaDBConfig) (*dockerenv.DockerEnv, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %v", err)
	}

	const (
		image = "mariadb:latest"
	)

	dm, err := dockerenv.New(&dockerenv.ContainerConfig{
		Name:        cfg.ContainerName,
		Image:       image,
		Env:         cfg.GetEnvVars(),
		ExposedPort: cfg.ExposedInternalPort,
		HostPort:    cfg.ExposedInternalPort,
		WaitFor:     "mariadbd: ready for connections.",
	})
	if err != nil {
		return nil, fmt.Errorf("docker manager: %v", err)
	}

	_, err = dm.UpsertContainer(context.Background(), cfg.RecreateContainer)
	if err != nil {
		return nil, fmt.Errorf("upsert container: %v", err)
	}

	return dm, err
}
