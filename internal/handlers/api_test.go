package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandleBooks_Success tests successful book retrieval
func TestHandleBooks_Success(t *testing.T) {
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

