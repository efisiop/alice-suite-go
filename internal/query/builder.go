package query

import (
	"database/sql"
	"fmt"
	"strings"
)

// BuildSQL builds a SQL query from QueryParams
func BuildSQL(table string, params *QueryParams) (string, []interface{}, error) {
	var query strings.Builder
	var args []interface{}

	// SELECT clause
	query.WriteString("SELECT ")
	if len(params.Select) == 0 || (len(params.Select) == 1 && params.Select[0] == "*") {
		query.WriteString("*")
	} else {
		query.WriteString(strings.Join(params.Select, ", "))
	}
	query.WriteString(" FROM ")
	query.WriteString(table)

	// WHERE clause
	if len(params.Filters) > 0 {
		query.WriteString(" WHERE ")
		conditions := []string{}
		for i, filter := range params.Filters {
			condition, filterArgs := buildFilterCondition(filter, i)
			conditions = append(conditions, condition)
			args = append(args, filterArgs...)
		}
		query.WriteString(strings.Join(conditions, " AND "))
	}

	// ORDER BY clause
	if len(params.OrderBy) > 0 {
		query.WriteString(" ORDER BY ")
		orderParts := []string{}
		for _, order := range params.OrderBy {
			direction := "ASC"
			if !order.Ascending {
				direction = "DESC"
			}
			orderParts = append(orderParts, fmt.Sprintf("%s %s", order.Column, direction))
		}
		query.WriteString(strings.Join(orderParts, ", "))
	}

	// LIMIT clause
	if params.Limit != nil {
		query.WriteString(fmt.Sprintf(" LIMIT %d", *params.Limit))
	}

	// OFFSET clause
	if params.Offset != nil {
		query.WriteString(fmt.Sprintf(" OFFSET %d", *params.Offset))
	}

	return query.String(), args, nil
}

// buildFilterCondition builds a WHERE condition from a Filter
func buildFilterCondition(filter Filter, index int) (string, []interface{}) {
	placeholder := fmt.Sprintf("$%d", index+1)
	args := []interface{}{filter.Value}

	switch filter.Operator {
	case "eq":
		return fmt.Sprintf("%s = %s", filter.Column, placeholder), args
	case "neq", "not":
		return fmt.Sprintf("%s != %s", filter.Column, placeholder), args
	case "gt":
		return fmt.Sprintf("%s > %s", filter.Column, placeholder), args
	case "gte":
		return fmt.Sprintf("%s >= %s", filter.Column, placeholder), args
	case "lt":
		return fmt.Sprintf("%s < %s", filter.Column, placeholder), args
	case "lte":
		return fmt.Sprintf("%s <= %s", filter.Column, placeholder), args
	case "like":
		return fmt.Sprintf("%s LIKE %s", filter.Column, placeholder), args
	case "ilike":
		// SQLite doesn't have ILIKE, use LIKE with UPPER
		return fmt.Sprintf("UPPER(%s) LIKE UPPER(%s)", filter.Column, placeholder), args
	case "is":
		if filter.Value == nil || filter.Value == "null" {
			return fmt.Sprintf("%s IS NULL", filter.Column), []interface{}{}
		}
		return fmt.Sprintf("%s IS NOT NULL", filter.Column), []interface{}{}
	case "in":
		// Handle IN operator with multiple values
		if values, ok := filter.Value.([]interface{}); ok {
			placeholders := []string{}
			args = []interface{}{}
			for i, val := range values {
				placeholders = append(placeholders, fmt.Sprintf("$%d", index*10+i+1))
				args = append(args, val)
			}
			return fmt.Sprintf("%s IN (%s)", filter.Column, strings.Join(placeholders, ", ")), args
		}
		// Single value IN
		return fmt.Sprintf("%s IN (%s)", filter.Column, placeholder), args
	default:
		// Default to equality
		return fmt.Sprintf("%s = %s", filter.Column, placeholder), args
	}
}

// ExecuteQuery executes a query and returns rows
func ExecuteQuery(db *sql.DB, sqlQuery string, args []interface{}) (*sql.Rows, error) {
	// Convert $1, $2, etc. to ? for SQLite
	sqlQuery = convertPlaceholders(sqlQuery)
	return db.Query(sqlQuery, args...)
}

// ExecuteQueryRow executes a query and returns a single row
func ExecuteQueryRow(db *sql.DB, sqlQuery string, args []interface{}) *sql.Row {
	// Convert $1, $2, etc. to ? for SQLite
	sqlQuery = convertPlaceholders(sqlQuery)
	return db.QueryRow(sqlQuery, args...)
}

// convertPlaceholders converts PostgreSQL-style placeholders ($1, $2) to SQLite placeholders (?)
func convertPlaceholders(query string) string {
	// Simple conversion: replace $N with ?
	// This is a basic implementation - for production, use a proper SQL parser
	result := query
	for i := 100; i >= 1; i-- {
		placeholder := fmt.Sprintf("$%d", i)
		result = strings.ReplaceAll(result, placeholder, "?")
	}
	return result
}

