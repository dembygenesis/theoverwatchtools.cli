package mysqlhelper

import (
	"context"
	"errors"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"strings"
	"testing"
	"time"
)

func Test_connection_Exec_Success(t *testing.T) {
	_, connectionParameters, cleanup := TestGetMockMariaDB(t)
	defer cleanup()

	settings := Settings{
		ConnectTimeout:       10 * time.Second,
		QueryTimeout:         30 * time.Second,
		ExecTimeout:          10 * time.Second,
		ConnectionParameters: connectionParameters,
	}

	conn, err := NewConnection(&settings)
	require.NoError(t, err, "unexpected connection error")

	createTableStmt := "CREATE TABLE IF NOT EXISTS temp_test_table (id INT AUTO_INCREMENT PRIMARY KEY, message VARCHAR(255))"
	_, err = conn.Exec(context.Background(), createTableStmt)
	require.NoError(t, err, "failed to create temporary table")

	insertStmt := "INSERT INTO temp_test_table (message) VALUES (?)"
	args := []interface{}{"test message"}

	_, err = conn.Exec(context.Background(), insertStmt, args...)
	require.NoError(t, err, "unexpected error during exec")
}

func Test_connection_Exec_Fail(t *testing.T) {
	_, connectionParameters, cleanup := TestGetMockMariaDB(t)
	defer cleanup()

	settings := Settings{
		ConnectTimeout:       10 * time.Second,
		QueryTimeout:         30 * time.Second,
		ExecTimeout:          10 * time.Second,
		ConnectionParameters: connectionParameters,
	}

	conn, err := NewConnection(&settings)
	require.NoError(t, err, "unexpected connection error")

	stmt := "INSERT INTO non_existent_table (column_name) VALUES (?)"
	args := []interface{}{"test value"}

	_, err = conn.Exec(context.Background(), stmt, args...)
	require.Error(t, err, "expected an error but got none")
}

func TestNew(t *testing.T) {
	_, connectionParameters, cleanup := TestGetMockMariaDB(t)
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
			name: "validate required - missing settings",
			args: args{settings: nil},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err, "unexpected nil error")
				assert.Contains(t, err.Error(), "required:")
			},
		},
		{
			name: "validate required - missing connection parameters",
			args: args{settings: &Settings{}},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err, "unexpected nil error")
				assert.Contains(t, err.Error(), "required connection parameters:")
			},
		},
		{
			name: "validate connection parameters",
			args: args{settings: &Settings{
				ConnectTimeout: 1 * time.Second,
				QueryTimeout:   1 * time.Second,
				ExecTimeout:    1 * time.Second,
				ConnectionParameters: &mysqlutil.ConnectionSettings{
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
				ConnectionParameters: &mysqlutil.ConnectionSettings{
					Host:     connectionParameters.Host,
					User:     connectionParameters.User + "wrong user",
					Pass:     connectionParameters.Pass + "wrong password",
					Database: connectionParameters.Database,
					Port:     connectionParameters.Port,
				},
			}},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err, "unexpected nil error")
				assert.Contains(t, err.Error(), "validate:")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConnection(tt.args.settings)
			tt.assertions(t, err)
		})
	}
}

func TestConnection_QueryAsArr(t *testing.T) {
	var (
		errMockDatabase = errors.New("mock database error")
	)

	type fields struct {
		db       *sqlx.DB
		settings *Settings
	}
	type args struct {
		ctx context.Context
		f   *QueryAsArrFilter
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		getDB      func(t *testing.T) (*sqlx.DB, func(ignoreErrors ...bool))
		assertions func(t *testing.T, res [][]string, pagination *model.Pagination, err error)
	}{
		{
			name: "nil settings",
			args: args{
				ctx: context.Background(),
				f:   nil,
			},
			getDB: func(t *testing.T) (*sqlx.DB, func(ignoreErrors ...bool)) {
				return nil, func(ignoreErrors ...bool) {
					// Do nothing
				}
			},
			fields: fields{
				db:       nil,
				settings: nil,
			},
			assertions: func(t *testing.T, res [][]string, pagination *model.Pagination, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "database error",
			args: args{
				ctx: context.Background(),
				f: &QueryAsArrFilter{
					Stmt: "SELECT 5",
				},
			},
			getDB: func(t *testing.T) (*sqlx.DB, func(ignoreErrors ...bool)) {
				db, mocks, err := sqlmock.New()
				require.NoError(t, err)
				mocks.ExpectQuery(".*").WillReturnError(errMockDatabase)
				cleanup := func(ignoreErrors ...bool) {

				}
				return sqlx.NewDb(db, "sqlmock"), cleanup
			},
			fields: fields{
				settings: &Settings{
					ConnectTimeout:       1 * time.Second,
					QueryTimeout:         1 * time.Second,
					ExecTimeout:          1 * time.Second,
					ConnectionParameters: &mysqlutil.ConnectionSettings{},
				},
			},
			assertions: func(t *testing.T, res [][]string, pagination *model.Pagination, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), errMockDatabase.Error(), "not matching expect mock database error")
			},
		},
		{
			name: "database query success",
			args: args{
				ctx: context.Background(),
				f: &QueryAsArrFilter{
					Stmt: `
						SELECT 'hello' AS message
					`,
					Pagination: model.NewPagination(),
				},
			},
			getDB: func(t *testing.T) (*sqlx.DB, func(ignoreErrors ...bool)) {
				db, _, cleanup := TestGetMockMariaDB(t)
				return db, cleanup
			},
			fields: fields{
				settings: &Settings{
					ConnectTimeout:       1 * time.Second,
					QueryTimeout:         1 * time.Second,
					ExecTimeout:          1 * time.Second,
					ConnectionParameters: &mysqlutil.ConnectionSettings{},
				},
			},
			assertions: func(t *testing.T, res [][]string, pagination *model.Pagination, err error) {
				require.NoError(t, err)
				require.Len(t, res, 2, "unexpected length for database query")
				assert.Equal(t, res[1][0], "hello")
			},
		},
		{
			name: "database query success with pagination",
			args: args{
				ctx: context.Background(),
				f: func() *QueryAsArrFilter {
					entries := 100

					stmt := dummySQLStatement(entries)
					return &QueryAsArrFilter{
						Stmt: stmt,
						Pagination: &model.Pagination{
							Page:    1,
							MaxRows: 25,
						},
					}
				}(),
			},
			getDB: func(t *testing.T) (*sqlx.DB, func(ignoreErrors ...bool)) {
				db, _, cleanup := TestGetMockMariaDB(t)
				return db, cleanup
			},
			fields: fields{
				settings: &Settings{
					ConnectTimeout:       1 * time.Second,
					QueryTimeout:         1 * time.Second,
					ExecTimeout:          1 * time.Second,
					ConnectionParameters: &mysqlutil.ConnectionSettings{},
				},
			},
			assertions: func(t *testing.T, res [][]string, pagination *model.Pagination, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, pagination, "unexpected nil pagination")
				require.Len(t, res, 26, "unexpected length for database query")

				assert.Len(t, pagination.Pages, 4, "expected pages length failed")
				assert.Equal(t, 100, pagination.TotalCount, "total count not-matching")
				assert.Equal(t, 1, pagination.Page, "unexpected page count")
				assert.Equal(t, 25, pagination.MaxRows, "pagination rows not equal") // 25 + 1 (header)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := tt.getDB(t)
			defer cleanup()
			c := &Connection{
				db:       db,
				settings: tt.fields.settings,
			}

			res, pagination, err := c.QueryAsArr(tt.args.ctx, tt.args.f)
			tt.assertions(t, res, pagination, err)
		})
	}
}

func TestConnection_QueryIntoStruct(t *testing.T) {
	db, connectionParameters, cleanup := TestGetMockMariaDB(t)
	defer cleanup()

	settings := Settings{
		ConnectTimeout:       10 * time.Second,
		QueryTimeout:         30 * time.Second,
		ExecTimeout:          10 * time.Second,
		ConnectionParameters: connectionParameters,
	}

	c := Connection{
		db:       db,
		settings: &settings,
	}

	type Result struct {
		Message  string `json:"message" db:"message"`
		Message2 string `json:"message2" db:"message2"`
	}

	ctx := context.Background()
	var results []Result

	message := "hello"
	message2 := "goodbye"

	err := Paginate(ctx, &results, c.GetPaginateSettings(&PaginateSettings{
		Stmt: "SELECT ? AS message, ? AS message2",
		Args: []interface{}{message, message2},
	}))

	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, message, results[0].Message)
	require.Equal(t, message2, results[0].Message2)
}

func dummySQLStatement(entries int) string {
	stmtArr := make([]string, 0)
	i := entries

	for i > 0 {
		stmtArr = append(stmtArr, "SELECT 'hello' AS message")
		i--
	}

	return strings.Join(stmtArr, " UNION ALL ")
}
