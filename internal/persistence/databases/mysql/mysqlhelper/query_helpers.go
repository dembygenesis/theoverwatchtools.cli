package mysqlhelper

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/jmoiron/sqlx"
	"math"
	"strconv"
	"time"
)

// preparePagination prepares the SQL query with pagination and calculates pagination details.
func preparePagination(ctx context.Context, db *sqlx.DB, query string, pagination *model.Pagination, args ...interface{}) (string, error) {
	totalCount, err := getQueryCount(ctx, db, query, args...)
	if err != nil {
		return "", fmt.Errorf("error retrieving query count: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.Rows)))
	if pagination.Page > totalPages {
		pagination.Page = totalPages
	}
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	offset := (pagination.Page - 1) * pagination.Rows

	paginatedQuery := query + " LIMIT " + strconv.Itoa(pagination.Rows) + " OFFSET " + strconv.Itoa(offset)

	pages := make([]int, totalPages)
	for i := 1; i <= totalPages; i++ {
		pages[i-1] = i
	}

	pagination.SetData(
		pages,
		pagination.Rows,
		pagination.Page,
		totalCount,
	)

	if err = pagination.ValidatePagination(); err != nil {
		return "", fmt.Errorf("validate: %v", err)
	}

	return paginatedQuery, nil
}

// queryAsStringArrays performs a paginated query and converts the results to a slice of string arrays.
func queryAsStringArrays(ctx context.Context, db *sqlx.DB, query string, pagination *model.Pagination) ([][]string, error) {
	var (
		paginatedQuery string
		err            error
	)

	if pagination == nil {
		paginatedQuery = query
	} else {
		paginatedQuery, err = preparePagination(ctx, db, query, pagination)
		if err != nil {
			return nil, err
		}
	}

	rows, err := db.QueryxContext(ctx, paginatedQuery)
	if err != nil {
		return nil, fmt.Errorf("queryx context: %v", err)
	}
	defer rows.Close()

	var results [][]string
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("columns: %v", err)
	}
	results = append(results, columns)

	for rows.Next() {
		rowValues, err := getRowValues(rows, len(columns))
		if err != nil {
			return nil, fmt.Errorf("get row values: %v", err)
		}
		results = append(results, rowValues)
	}

	return results, nil
}

// getQueryCount retrieves the total number of records for a given SQL query.
func getQueryCount(
	ctx context.Context,
	db *sqlx.DB,
	query string,
	args ...interface{},
) (int, error) {
	wrappedQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_wrapper"
	var count int
	if err := db.GetContext(ctx, &count, wrappedQuery, args...); err != nil {
		return 0, fmt.Errorf("error running count query: %v", err)
	}
	return count, nil
}

// getRowValues reads the current row's values and converts them into a string slice
func getRowValues(row *sqlx.Rows, colCount int) ([]string, error) {
	values := make([]interface{}, colCount)
	valuePtrs := make([]interface{}, colCount)
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := row.Scan(valuePtrs...); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	strValues := make([]string, colCount)
	for i, val := range values {
		strValues[i] = fmt.Sprintf("%v", convertSQLValueToString(val))
	}

	return strValues, nil
}

func convertSQLValueToString(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return ""
	case []byte:
		return string(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case bool, int, int32, int64, float64, string:
		return fmt.Sprintf("%v", v)
	default:
		log.Printf("Unhandled type for %v: %T\n", val, val)
		return fmt.Sprintf("%v", val)
	}
}

// paginate executes a query, and stores
// the results in the `dest`.
func paginate[T any](
	ctx context.Context,
	db *sqlx.DB,
	dest *[]T,
	stmt string,
	pagination *model.Pagination,
	args ...interface{},
) error {
	paginatedQuery, err := preparePagination(ctx, db, stmt, pagination, args...)
	if err != nil {
		return fmt.Errorf("preparePagination: %v", err)
	}

	if err = db.SelectContext(ctx, dest, paginatedQuery, args...); err != nil {
		return fmt.Errorf("select: %v", err)
	}

	return nil
}

func Paginate[T any](
	ctx context.Context,
	dest *[]T,
	settings *PaginateSettings,
) error {
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("validate: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, settings.timeoutDuration)
	defer cancel()

	return paginate(ctx, settings.db, dest, settings.Stmt, settings.Pagination, settings.Args...)
}
