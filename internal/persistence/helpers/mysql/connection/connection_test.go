package connection

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/persistence/helpers/mysql"
	"github.com/dembygenesis/local.tools/internal/persistence/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_connection_Query_Fail(t *testing.T) {
	connectionParameters, cleanup := mock.TestMockMysqlDatabase(t)
	defer cleanup()

	settings := Settings{
		ConnectTimeout: 10 * time.Second,
		QueryTimeout:   30 * time.Second,
		ExecTimeout:    10 * time.Second,
		Parameters:     connectionParameters,
	}

	conn, err := New(&settings)
	require.NoError(t, err)

	message := "hello world"
	type queryResult struct {
		Message string `json:"message"`
	}
	var result []queryResult
	stmt := "SELECT ?x AS message"

	args := []interface{}{message}
	err = conn.Query(context.Background(), &result, stmt, args...)
	require.Error(t, err, "unexpected nil error")
}

func Test_connection_Query_Success(t *testing.T) {
	connectionParameters, cleanup := mock.TestMockMysqlDatabase(t)
	defer cleanup()

	settings := Settings{
		ConnectTimeout: 10 * time.Second,
		QueryTimeout:   30 * time.Second,
		ExecTimeout:    10 * time.Second,
		Parameters:     connectionParameters,
	}

	conn, err := New(&settings)
	require.NoError(t, err)

	message := "hello world"
	type queryResult struct {
		Message string `json:"message"`
	}
	var result []queryResult
	stmt := "SELECT ? AS message"

	args := []interface{}{message}
	err = conn.Query(context.Background(), &result, stmt, args...)
	require.NoError(t, err, "unexpected query error")

	require.Len(t, result, 1)
	require.Equal(t, result[0].Message, message)
}

func Test_connection_Exec_Success(t *testing.T) {
	connectionParameters, cleanup := mock.TestMockMysqlDatabase(t)
	defer cleanup()

	settings := Settings{
		ConnectTimeout: 10 * time.Second,
		QueryTimeout:   30 * time.Second,
		ExecTimeout:    10 * time.Second,
		Parameters:     connectionParameters,
	}

	conn, err := New(&settings)
	require.NoError(t, err)

	createTableStmt := "CREATE TABLE IF NOT EXISTS temp_test_table (id INT AUTO_INCREMENT PRIMARY KEY, message VARCHAR(255))"
	_, err = conn.Exec(context.Background(), createTableStmt)
	require.NoError(t, err, "failed to create temporary table")

	insertStmt := "INSERT INTO temp_test_table (message) VALUES (?)"
	args := []interface{}{"test message"}

	_, err = conn.Exec(context.Background(), insertStmt, args...)
	require.NoError(t, err, "unexpected error during exec")
}

func Test_connection_Exec_Fail(t *testing.T) {
	connectionParameters, cleanup := mock.TestMockMysqlDatabase(t)
	defer cleanup()

	settings := Settings{
		ConnectTimeout: 10 * time.Second,
		QueryTimeout:   30 * time.Second,
		ExecTimeout:    10 * time.Second,
		Parameters:     connectionParameters,
	}

	conn, err := New(&settings)
	require.NoError(t, err)

	stmt := "INSERT INTO non_existent_table (column_name) VALUES (?)"
	args := []interface{}{"test value"}

	_, err = conn.Exec(context.Background(), stmt, args...)
	require.Error(t, err, "expected an error but got none")
}

func TestNew(t *testing.T) {
	connectionParameters, cleanup := mock.TestMockMysqlDatabase(t)
	defer cleanup()

	type args struct {
		settings *Settings
	}
	tests := []struct {
		name       string
		args       args
		assertions func(t *testing.T, err error)
	}{
		{
			name: "nil settings",
			args: args{settings: nil},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err, "unexpected nil error")
				assert.Equal(t, err.Error(), "settings nil")
			},
		},
		{
			name: "validate required",
			args: args{settings: &Settings{}},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err, "unexpected nil error")
				assert.Contains(t, err.Error(), "required:")
			},
		},
		{
			name: "validate connection parameters",
			args: args{settings: &Settings{
				ConnectTimeout: 1 * time.Second,
				QueryTimeout:   1 * time.Second,
				ExecTimeout:    1 * time.Second,
				Parameters: &mysql.ConnectionParameters{
					Host:     connectionParameters.Host,
					User:     connectionParameters.User,
					Pass:     connectionParameters.Pass,
					Database: connectionParameters.Database,
					Port:     connectionParameters.Port,
				},
			}},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err, "unexpected nil error")
			},
		},
		{
			name: "wrong database credentials",
			args: args{settings: &Settings{
				ConnectTimeout: 1 * time.Second,
				QueryTimeout:   1 * time.Second,
				ExecTimeout:    1 * time.Second,
				Parameters: &mysql.ConnectionParameters{
					Host:     connectionParameters.Host,
					User:     connectionParameters.User + "wrong user",
					Pass:     connectionParameters.Pass + "wrong password",
					Database: connectionParameters.Database,
					Port:     connectionParameters.Port,
				},
			}},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err, "unexpected nil error")
				assert.Contains(t, err.Error(), "ping ctx:")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.settings)
			tt.assertions(t, err)
		})
	}
}
