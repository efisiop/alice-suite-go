package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
)

// bookService is shared from api.go - initialized there
// We reference it here to use the same instance

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

	pageNum := int(pageNumber)

	// First try to use the book service to get the page with sections
	// But handle database structure mismatch gracefully
	page, err := bookService.GetPage(bookID, pageNum)
	if err != nil {
		log.Printf("Error fetching page %d for book %s (likely structure mismatch): %v", pageNum, bookID, err)
		// Fall through to fallback - don't return error yet
	} else if page != nil && len(page.Sections) > 0 {
		// Success - return the page
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
		return
	}
	
	// If page was nil or had no sections, try fallback
	if page != nil && len(page.Sections) == 0 {
		log.Printf("Page %d found but has no sections, trying fallback", pageNum)
	}

	// Fallback: Query old sections structure (sections with start_page/end_page)
	// This handles databases that haven't been migrated to the new page-based structure
	log.Printf("Page %d not found in pages table, trying fallback to old sections structure", pageNum)
	
	// Check if old sections structure exists (has start_page column)
	var tableInfo string
	checkQuery := `SELECT sql FROM sqlite_master WHERE type='table' AND name='sections'`
	err = database.DB.QueryRow(checkQuery).Scan(&tableInfo)
	if err != nil || tableInfo == "" {
		log.Printf("Sections table not found or error checking structure: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Page not found - database structure not available",
		})
		return
	}

	// Try querying old structure (with start_page/end_page)
	query := `SELECT id, content FROM sections 
	          WHERE start_page <= ? AND end_page >= ?
	          ORDER BY number`
	rows, err := database.DB.Query(query, pageNum, pageNum)
	if err != nil {
		// If that fails, might be new structure - try querying by page_number
		log.Printf("Old structure query failed, trying new structure: %v", err)
		query = `SELECT id, content, page_number, section_number FROM sections 
		         WHERE page_number = ?
		         ORDER BY section_number`
		rows, err = database.DB.Query(query, pageNum)
		if err != nil {
			log.Printf("Error querying sections: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Error fetching page: %v", err),
			})
			return
		}
		defer rows.Close()

		// Build sections from new structure
		foundSections := []models.Section{}
		for rows.Next() {
			var id, content string
			var pageNum2, sectionNum int
			if err := rows.Scan(&id, &content, &pageNum2, &sectionNum); err != nil {
				continue
			}
			foundSections = append(foundSections, models.Section{
				ID:            id,
				PageID:        fmt.Sprintf("page-%d", pageNum2),
				PageNumber:    pageNum2,
				SectionNumber: sectionNum,
				Content:       content,
			})
		}

		if len(foundSections) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Page not found",
			})
			return
		}

		pageObj := &models.Page{
			ID:         fmt.Sprintf("page-%d", pageNum),
			BookID:     bookID,
			PageNumber: pageNum,
			Sections:   foundSections,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pageObj)
		return
	}
	defer rows.Close()

	// Build a page object from old sections structure
	foundSections := []models.Section{}
	for rows.Next() {
		var id, content string
		if err := rows.Scan(&id, &content); err != nil {
			continue
		}
		// Create a minimal section object
		foundSections = append(foundSections, models.Section{
			ID:            id,
			PageID:        fmt.Sprintf("page-%d", pageNum), // Synthetic page ID
			PageNumber:    pageNum,
			SectionNumber: len(foundSections) + 1,
			Content:       content,
		})
	}

	if len(foundSections) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Page not found",
		})
		return
	}

	// Return a page object with the found sections
	pageObj := &models.Page{
		ID:         fmt.Sprintf("page-%d", pageNum),
		BookID:     bookID,
		PageNumber: pageNum,
		Sections:   foundSections,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pageObj)
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
