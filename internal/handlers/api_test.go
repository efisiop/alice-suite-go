package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

// TestHandleBooks_Success tests successful book retrieval
func TestHandleBooks_Success(t *testing.T) {
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testDB.Close() })
	if _, err := testDB.Exec(`CREATE TABLE books (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		description TEXT,
		total_pages INTEGER NOT NULL,
		created_at TEXT
	)`); err != nil {
		t.Fatal(err)
	}
	originalDB, originalDriver := database.DB, database.DriverName
	database.DB, database.DriverName = testDB, "sqlite3"
	t.Cleanup(func() {
		database.DB, database.DriverName = originalDB, originalDriver
	})

	req, err := http.NewRequest("GET", "/rest/v1/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleBooks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want application/json", contentType)
	}
}

// TestHandleBooks_MethodNotAllowed tests wrong HTTP method
func TestHandleBooks_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("POST", "/rest/v1/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleBooks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

// TestHandleChapters_MissingBookID tests missing book_id parameter
func TestHandleChapters_MissingBookID(t *testing.T) {
	req, err := http.NewRequest("GET", "/rest/v1/chapters", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleChapters)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want application/json", contentType)
	}
}
