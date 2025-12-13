package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// TestRequireAuth_MissingToken tests authentication middleware with missing token
func TestRequireAuth_MissingToken(t *testing.T) {
	handler := RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

// TestRequireAuth_ValidToken tests authentication middleware with valid token
func TestRequireAuth_ValidToken(t *testing.T) {
	// Generate a valid token
	token, err := auth.GenerateJWT("test-user", "test@example.com", "reader")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	handler := RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestRequireRole_Reader tests role-based access control for reader
func TestRequireRole_Reader(t *testing.T) {
	token, err := auth.GenerateJWT("test-user", "test@example.com", "reader")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	handler := RequireRole("reader")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestRequireRole_WrongRole tests role-based access control with wrong role
func TestRequireRole_WrongRole(t *testing.T) {
	token, err := auth.GenerateJWT("test-user", "test@example.com", "reader")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	handler := RequireRole("consultant")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}

