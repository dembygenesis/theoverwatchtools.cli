package mysqlhelper

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_New_MariaDB_Existing(t *testing.T) {
	i := 2

	for i > 0 {
		i--

		var (
			mappedPort    = 3341
			containerName = strutil.GetUuidUnderscore()
			host          = "root"
			pass          = "secret"
		)

		cfg := MariaDBConfig{
			ContainerName:       containerName,
			Host:                host,
			Pass:                pass,
			Database:            containerName,
			ExposedInternalPort: mappedPort,
			ExposedExternalPort: mappedPort,
			RecreateContainer:   false,
		}

		mariaDbCtn, err := NewMariaDB(&cfg)
		require.NoError(t, err, "error starting new maria db container")
		require.NotNil(t, mariaDbCtn, "unexpected nil mariaDbCtn")
	}
}

func Test_New_MariaDB_Success(t *testing.T) {
	var (
		mappedPort    = 3340
		containerName = strutil.GetUuidUnderscore()
		host          = "root"
		pass          = "secret"
	)

	cfg := MariaDBConfig{
		ContainerName:       containerName,
		Host:                host,
		Pass:                pass,
		Database:            containerName,
		ExposedInternalPort: mappedPort,
		ExposedExternalPort: mappedPort,
		RecreateContainer:   true,
	}

	mariaDbCtn, err := NewMariaDB(&cfg)
	require.NoError(t, err, "error starting new maria db container")

	err = mariaDbCtn.Cleanup(context.Background())
	require.NoError(t, err, "error cleaning up database")
}

func Test_New_MariaDB_Fail_Validate(t *testing.T) {
	var (
		mappedPort    = 0
		containerName = strutil.GetUuidUnderscore()
		host          = "localhost"
		pass          = "secret"
	)

	cfg := MariaDBConfig{
		ContainerName:       containerName,
		Host:                host,
		Pass:                pass,
		Database:            containerName,
		ExposedInternalPort: mappedPort,
		RecreateContainer:   true,
	}

	_, err := NewMariaDB(&cfg)
	require.Error(t, err, "there should be a validation error")
	assert.Contains(t, err.Error(), "validate:")
}

func Test_New_MariaDB_Fail_Invalid_Setting(t *testing.T) {
	var (
		mappedPort    = 999999999999
		containerName = strutil.GetUuidUnderscore()
		host          = "localhost"
		pass          = "secret"
	)

	cfg := MariaDBConfig{
		ContainerName:       containerName,
		Host:                host,
		Pass:                pass,
		Database:            containerName,
		ExposedInternalPort: mappedPort,
		ExposedExternalPort: mappedPort,
		RecreateContainer:   true,
	}

	_, err := NewMariaDB(&cfg)
	require.Error(t, err, "there should be a validation error")
	assert.Contains(t, err.Error(), "upsert container:")
}
