package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Alice Suite Reader API - Physical Book Companion",
		"scope":   "First 3 chapters test ground",
		"version": "1.0.0",
	})
}

// SetupAllRoutes sets up all routes for the application
// This is a convenience function that calls all individual setup functions
func SetupAllRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/health", HealthCheck)

	// API routes (Supabase-compatible)
	SetupAPIRoutes(mux)

	// Authentication routes
	SetupAuthRoutes(mux)

	// Reader app routes
	SetupReaderRoutes(mux)

	// Consultant app routes
	SetupConsultantRoutes(mux)
}

