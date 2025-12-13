package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/query"
	"github.com/google/uuid"
)

// HandleRESTTable handles GET/POST/PATCH/DELETE /rest/v1/:table
func HandleRESTTable(w http.ResponseWriter, r *http.Request) {
	// Extract table name from path
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/rest/v1/")
	table := strings.Split(path, "?")[0] // Remove query string
	table = strings.TrimSuffix(table, "/")

	if table == "" {
		http.Error(w, "Table name required", http.StatusBadRequest)
		return
	}

	// Validate table name (prevent SQL injection)
	if !isValidTableName(table) {
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGETTable(w, r, table)
	case http.MethodPost:
		handlePOSTTable(w, r, table)
	case http.MethodPatch:
		handlePATCHTable(w, r, table)
	case http.MethodDelete:
		handleDELETETable(w, r, table)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGETTable handles GET /rest/v1/:table
func handleGETTable(w http.ResponseWriter, r *http.Request, table string) {
	// Parse query parameters
	queryParams, err := query.ParseQuery(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid query parameters: %v", err), http.StatusBadRequest)
		return
	}

	// Build SQL query
	sqlQuery, args, err := query.BuildSQL(table, queryParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error building query: %v", err), http.StatusInternalServerError)
		return
	}

	// Execute query
	rows, err := query.ExecuteQuery(database.DB, sqlQuery, args)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting columns: %v", err), http.StatusInternalServerError)
		return
	}

	// Scan rows into maps
	results := []map[string]interface{}{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val != nil {
				// Convert SQLite INTEGER booleans to actual booleans
				if b, ok := convertBooleanValue(table, col, val); ok {
					row[col] = b
				} else {
					row[col] = val
				}
			} else {
				row[col] = nil
			}
		}

		results = append(results, row)
	}

	// Handle joins (post-process results)
	if len(queryParams.Joins) > 0 {
		results = applyJoins(results, queryParams.Joins, table)
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// handlePOSTTable handles POST /rest/v1/:table
func handlePOSTTable(w http.ResponseWriter, r *http.Request, table string) {
	// Parse request body
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Generate UUID for id if missing
	if _, exists := data["id"]; !exists {
		data["id"] = uuid.New().String()
	}

	// Add timestamps if missing
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, exists := data["created_at"]; !exists {
		data["created_at"] = now
	}
	if _, exists := data["updated_at"]; !exists {
		data["updated_at"] = now
	}

	// Validate foreign keys
	if err := validateForeignKeys(table, data); err != nil {
		http.Error(w, fmt.Sprintf("Foreign key validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// Build INSERT query
	columns := []string{}
	placeholders := []string{}
	values := []interface{}{}
	index := 1

	for col, val := range data {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		values = append(values, val)
		index++
	}

	sqlQuery := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// Execute INSERT
	_, err := database.DB.Exec(sqlQuery, values...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	// Get inserted row ID
	insertedID := data["id"].(string)

	// Check if select parameter is present
	selectParam := r.URL.Query().Get("select")
	if selectParam != "" {
		// Return the inserted row
		queryParams, err := query.ParseQuery(r.URL.Query())
		if err == nil {
			queryParams.Filters = []query.Filter{
				{Column: "id", Operator: "eq", Value: insertedID},
			}
			sqlQuery, args, err := query.BuildSQL(table, queryParams)
			if err == nil {
				row := database.DB.QueryRow(sqlQuery, args...)

				columns, err := getTableColumns(table)
				if err == nil {
					values := make([]interface{}, len(columns))
					valuePtrs := make([]interface{}, len(columns))
					for i := range values {
						valuePtrs[i] = &values[i]
					}

					if err := row.Scan(valuePtrs...); err == nil {
						resultRow := make(map[string]interface{})
						for i, col := range columns {
							if values[i] != nil {
								if b, ok := convertBooleanValue(table, col, values[i]); ok {
									resultRow[col] = b
								} else {
									resultRow[col] = values[i]
								}
							} else {
								resultRow[col] = nil
							}
						}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode([]map[string]interface{}{resultRow})
		
		// Broadcast help request to consultants
		BroadcastHelpRequest(resultRow)
		return
	}
				}
			}
		}
	}

	// Return success with ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      insertedID,
		"success": true,
		"rows":    1,
	})
}

// handlePATCHTable handles PATCH /rest/v1/:table
func handlePATCHTable(w http.ResponseWriter, r *http.Request, table string) {
	// Parse query parameters for WHERE clause
	queryParams, err := query.ParseQuery(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid query parameters: %v", err), http.StatusBadRequest)
		return
	}

	// Parse request body
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Add updated_at timestamp
	data["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

	// Build UPDATE query
	setParts := []string{}
	values := []interface{}{}
	for col, val := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", col))
		values = append(values, val)
	}

	sqlQuery := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(setParts, ", "))

	// Add WHERE clause
	if len(queryParams.Filters) > 0 {
		whereParts := []string{}
		for _, filter := range queryParams.Filters {
			whereParts = append(whereParts, fmt.Sprintf("%s = ?", filter.Column))
			values = append(values, filter.Value)
		}
		sqlQuery += " WHERE " + strings.Join(whereParts, " AND ")
	}

	// Execute UPDATE
	result, err := database.DB.Exec(sqlQuery, values...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"rows":    rowsAffected,
	})
}

// handleDELETETable handles DELETE /rest/v1/:table
func handleDELETETable(w http.ResponseWriter, r *http.Request, table string) {
	// Parse query parameters for WHERE clause
	queryParams, err := query.ParseQuery(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid query parameters: %v", err), http.StatusBadRequest)
		return
	}

	// Build DELETE query
	sqlQuery := fmt.Sprintf("DELETE FROM %s", table)

	// Add WHERE clause
	values := []interface{}{}
	if len(queryParams.Filters) > 0 {
		whereParts := []string{}
		for _, filter := range queryParams.Filters {
			whereParts = append(whereParts, fmt.Sprintf("%s = ?", filter.Column))
			values = append(values, filter.Value)
		}
		sqlQuery += " WHERE " + strings.Join(whereParts, " AND ")
	}

	// Execute DELETE
	result, err := database.DB.Exec(sqlQuery, values...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"rows":    rowsAffected,
	})
}

// Helper functions

// isValidTableName validates table name to prevent SQL injection
func isValidTableName(table string) bool {
	// Only allow alphanumeric and underscore
	for _, char := range table {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return len(table) > 0 && len(table) < 100
}

// convertBooleanValue converts SQLite INTEGER booleans to actual booleans
func convertBooleanValue(table, column string, value interface{}) (bool, bool) {
	// Tables/columns that should be converted to boolean
	booleanColumns := map[string]map[string]bool{
		"users": {
			"is_verified": true,
		},
		"verification_codes": {
			"is_used": true,
		},
		"help_requests": {
			"is_public": false, // Add if exists
		},
	}

	if cols, ok := booleanColumns[table]; ok {
		if cols[column] {
			switch v := value.(type) {
			case int64:
				return v == 1, true
			case int:
				return v == 1, true
			case string:
				return v == "1" || v == "true", true
			}
		}
	}

	return false, false
}

// validateForeignKeys validates foreign key constraints
func validateForeignKeys(table string, data map[string]interface{}) error {
	// Foreign key mappings
	fkMappings := map[string]map[string]string{
		"help_requests": {
			"user_id": "users",
			"book_id": "books",
		},
		"interactions": {
			"user_id": "users",
			"book_id": "books",
		},
		"reading_progress": {
			"user_id": "users",
			"book_id": "books",
		},
	}

	if fks, ok := fkMappings[table]; ok {
		for fkColumn, refTable := range fks {
			if fkValue, exists := data[fkColumn]; exists && fkValue != nil {
				// Check if referenced record exists
				query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", refTable)
				var count int
				err := database.DB.QueryRow(query, fkValue).Scan(&count)
				if err != nil {
					return fmt.Errorf("error checking foreign key %s: %v", fkColumn, err)
				}
				if count == 0 {
					return fmt.Errorf("foreign key constraint failed: %s.%s references non-existent %s.id", table, fkColumn, refTable)
				}
			}
		}
	}

	return nil
}

// applyJoins applies join logic to results (post-processing)
func applyJoins(results []map[string]interface{}, joins []query.Join, mainTable string) []map[string]interface{} {
	// This is a simplified join implementation
	// For production, use proper SQL JOINs
	for _, join := range joins {
		for _, result := range results {
			if fkValue, exists := result[join.ForeignKey]; exists {
				// Fetch joined data
				query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?",
					strings.Join(join.Columns, ", "),
					join.Table)
				row := database.DB.QueryRow(query, fkValue)

				values := make([]interface{}, len(join.Columns))
				valuePtrs := make([]interface{}, len(join.Columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}

				if err := row.Scan(valuePtrs...); err == nil {
					joinedData := make(map[string]interface{})
					for i, col := range join.Columns {
						joinedData[col] = values[i]
					}
					result[join.Table] = joinedData
				}
			}
		}
	}
	return results
}

// getTableColumns gets column names for a table
func getTableColumns(table string) ([]string, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", table)
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := []string{}
	for rows.Next() {
		var cid int
		var name string
		var dataType string
		var notNull int
		var defaultValue interface{}
		var pk int

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
			continue
		}
		columns = append(columns, name)
	}

	return columns, nil
}

