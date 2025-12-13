package query

import (
	"fmt"
	"net/url"
	"strings"
)

// QueryParams represents parsed Supabase query parameters
type QueryParams struct {
	Select   []string            // Columns to select
	Filters  []Filter            // WHERE conditions
	OrderBy  []OrderBy           // ORDER BY clauses
	Limit    *int                // LIMIT
	Offset   *int                // OFFSET
	Joins    []Join              // JOIN specifications
}

// Filter represents a WHERE condition
type Filter struct {
	Column   string
	Operator string // eq, neq, gt, gte, lt, lte, like, ilike, is, in
	Value    interface{}
}

// OrderBy represents an ORDER BY clause
type OrderBy struct {
	Column    string
	Ascending bool
}

// Join represents a join specification
type Join struct {
	Table      string   // e.g., "profiles"
	ForeignKey string   // e.g., "user_id"
	Columns    []string // Columns to select from joined table
}

// ParseQuery parses Supabase-style query parameters
func ParseQuery(query url.Values) (*QueryParams, error) {
	params := &QueryParams{
		Select:  []string{},
		Filters: []Filter{},
		OrderBy: []OrderBy{},
		Joins:   []Join{},
	}

	// Parse select parameter
	if selectParam := query.Get("select"); selectParam != "" {
		params.Select = parseSelect(selectParam)
		// Extract joins from select
		params.Joins = parseJoins(selectParam)
	} else {
		// Default to all columns
		params.Select = []string{"*"}
	}

	// Parse filters (eq, neq, gt, gte, lt, lte, like, ilike, is, in)
	for key, values := range query {
		if len(values) == 0 {
			continue
		}
		value := values[0]

		// Skip special parameters
		if key == "select" || key == "order" || key == "limit" || key == "offset" {
			continue
		}

		// Parse filter operators
		filter := parseFilter(key, value)
		if filter != nil {
			params.Filters = append(params.Filters, *filter)
		}
	}

	// Parse order parameter
	if orderParam := query.Get("order"); orderParam != "" {
		params.OrderBy = parseOrder(orderParam)
	}

	// Parse limit
	if limitParam := query.Get("limit"); limitParam != "" {
		var limit int
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err == nil {
			params.Limit = &limit
		}
	}

	// Parse offset
	if offsetParam := query.Get("offset"); offsetParam != "" {
		var offset int
		if _, err := fmt.Sscanf(offsetParam, "%d", &offset); err == nil {
			params.Offset = &offset
		}
	}

	return params, nil
}

// parseSelect parses the select parameter
func parseSelect(selectParam string) []string {
	// Remove join syntax for now (handled separately)
	selectParam = removeJoinSyntax(selectParam)
	
	// Split by comma
	parts := strings.Split(selectParam, ",")
	columns := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			columns = append(columns, part)
		}
	}
	
	if len(columns) == 0 {
		return []string{"*"}
	}
	
	return columns
}

// parseJoins parses join syntax like "profiles:user_id(first_name,last_name)"
func parseJoins(selectParam string) []Join {
	joins := []Join{}
	
	// Find patterns like "table:foreign_key(columns)"
	// Regex would be better, but we'll use string parsing
	parts := strings.Split(selectParam, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, ":") && strings.Contains(part, "(") {
			// Extract table:foreign_key(columns)
			colonIdx := strings.Index(part, ":")
			parenIdx := strings.Index(part, "(")
			
			if colonIdx > 0 && parenIdx > colonIdx {
				table := strings.TrimSpace(part[:colonIdx])
				foreignKey := strings.TrimSpace(part[colonIdx+1 : parenIdx])
				columnsStr := strings.TrimSpace(part[parenIdx+1 : len(part)-1])
				
				columns := []string{}
				if columnsStr != "" {
					for _, col := range strings.Split(columnsStr, ",") {
						col = strings.TrimSpace(col)
						if col != "" {
							columns = append(columns, col)
						}
					}
				}
				
				joins = append(joins, Join{
					Table:      table,
					ForeignKey: foreignKey,
					Columns:    columns,
				})
			}
		}
	}
	
	return joins
}

// removeJoinSyntax removes join syntax from select string
func removeJoinSyntax(selectParam string) string {
	// Remove patterns like "profiles:user_id(...)"
	result := selectParam
	for {
		start := strings.Index(result, ":")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ")")
		if end == -1 {
			break
		}
		// Find the matching opening parenthesis
		parenStart := strings.LastIndex(result[:start+end], "(")
		if parenStart > start {
			// Remove from : to )
			result = result[:start] + result[start+end+1:]
		} else {
			break
		}
	}
	return result
}

// parseFilter parses a filter from query parameter
func parseFilter(key, value string) *Filter {
	// Check if key contains operator (format: column.operator=value)
	if strings.Contains(key, ".") {
		parts := strings.Split(key, ".")
		if len(parts) == 2 {
			column := parts[0]
			operator := parts[1]
			
			validOperators := map[string]bool{
				"eq": true, "neq": true, "not": true,
				"gt": true, "gte": true, "lt": true, "lte": true,
				"like": true, "ilike": true, "is": true, "in": true,
			}
			
			if validOperators[operator] {
				return &Filter{
					Column:   column,
					Operator: operator,
					Value:    value,
				}
			}
		}
	}
	
	// Check if value contains operator (format: column=operator.value)
	if strings.Contains(value, ".") {
		parts := strings.SplitN(value, ".", 2)
		if len(parts) == 2 {
			operator := parts[0]
			actualValue := parts[1]
			
			validOperators := map[string]bool{
				"eq": true, "neq": true, "not": true,
				"gt": true, "gte": true, "lt": true, "lte": true,
				"like": true, "ilike": true, "is": true, "in": true,
			}
			
			if validOperators[operator] {
				return &Filter{
					Column:   key,
					Operator: operator,
					Value:    actualValue,
				}
			}
		}
	}
	
	// Default to equality
	return &Filter{
		Column:   key,
		Operator: "eq",
		Value:    value,
	}
}

// parseOrder parses the order parameter
func parseOrder(orderParam string) []OrderBy {
	orders := []OrderBy{}
	
	// Format: "column.asc" or "column.desc" or "column"
	parts := strings.Split(orderParam, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		order := OrderBy{Ascending: true} // Default ascending
		
		if strings.Contains(part, ".") {
			orderParts := strings.Split(part, ".")
			if len(orderParts) == 2 {
				order.Column = strings.TrimSpace(orderParts[0])
				direction := strings.ToLower(strings.TrimSpace(orderParts[1]))
				order.Ascending = direction != "desc"
			}
		} else {
			order.Column = part
		}
		
		if order.Column != "" {
			orders = append(orders, order)
		}
	}
	
	return orders
}

