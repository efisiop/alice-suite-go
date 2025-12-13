package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/database"
)

// HandleRPC handles POST /rest/v1/rpc/:function
func HandleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract function name from path
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/rest/v1/rpc/")
	function := strings.TrimSuffix(path, "/")

	if function == "" {
		http.Error(w, "RPC function name required", http.StatusBadRequest)
		return
	}

	// Parse request body for parameters
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		// If no body, use empty params
		params = make(map[string]interface{})
	}

	// Route to specific RPC function
	switch function {
	case "get_definition_with_context":
		handleGetDefinitionWithContext(w, r, params)
	case "get_sections_for_page":
		handleGetSectionsForPage(w, r, params)
	case "verify-book-code":
		// This is handled by HandleVerifyBookCode in verification.go
		HandleVerifyBookCode(w, r)
	case "check-book-verified":
		HandleCheckBookVerified(w, r)
	case "check_table_exists":
		handleCheckTableExists(w, r, params)
	default:
		http.Error(w, fmt.Sprintf("Unknown RPC function: %s", function), http.StatusNotFound)
	}
}

// handleGetDefinitionWithContext handles get_definition_with_context RPC
func handleGetDefinitionWithContext(w http.ResponseWriter, r *http.Request, params map[string]interface{}) {
	term, _ := params["term"].(string)
	bookID, _ := params["book_id"].(string)
	_ = params["section_id"] // sectionID not used yet but may be needed for context

	if term == "" {
		http.Error(w, "term parameter required", http.StatusBadRequest)
		return
	}

	// Query glossary
	query := `SELECT term, definition, example FROM alice_glossary WHERE term = ? AND book_id = ? LIMIT 1`
	var definition, example string
	err := database.DB.QueryRow(query, term, bookID).Scan(&term, &definition, &example)
	if err != nil {
		// Term not found
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"term":       term,
			"definition": "Word not found in glossary",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"term":       term,
		"definition": definition,
		"example":    example,
	})
}

// handleGetSectionsForPage handles get_sections_for_page RPC
func handleGetSectionsForPage(w http.ResponseWriter, r *http.Request, params map[string]interface{}) {
	bookID, _ := params["book_id"].(string)
	pageNumber, ok := params["page_number"].(float64)
	if !ok {
		http.Error(w, "book_id and page_number parameters required", http.StatusBadRequest)
		return
	}

	// Query sections for the page
	query := `SELECT id, page_id, page_number, section_number, content, word_count 
	          FROM sections WHERE page_id IN (
	              SELECT id FROM pages WHERE book_id = ? AND page_number = ?
	          ) ORDER BY section_number`
	rows, err := database.DB.Query(query, bookID, int(pageNumber))
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	sections := []map[string]interface{}{}
	for rows.Next() {
		var id, pageID, content string
		var pageNum, sectionNum, wordCount int
		if err := rows.Scan(&id, &pageID, &pageNum, &sectionNum, &content, &wordCount); err != nil {
			continue
		}
		sections = append(sections, map[string]interface{}{
			"id":             id,
			"page_id":        pageID,
			"page_number":    pageNum,
			"section_number": sectionNum,
			"content":        content,
			"word_count":     wordCount,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sections)
}

// handleCheckTableExists handles check_table_exists RPC
func handleCheckTableExists(w http.ResponseWriter, r *http.Request, params map[string]interface{}) {
	tableName, _ := params["table_name"].(string)
	if tableName == "" {
		http.Error(w, "table_name parameter required", http.StatusBadRequest)
		return
	}

	// Check if table exists
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
	var name string
	err := database.DB.QueryRow(query, tableName).Scan(&name)
	exists := err == nil && name == tableName

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"exists": exists,
	})
}

