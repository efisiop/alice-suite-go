package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get AI provider status
	aiStatus := aiService.GetProviderStatus()
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Alice Suite Reader API - Physical Book Companion",
		"scope":   "First 3 chapters test ground",
		"version": "1.0.0",
		"ai": map[string]interface{}{
			"active_provider": aiStatus["active_provider"],
			"configured_provider": aiStatus["configured_provider"],
		},
	})
}

// HandleStatus handles GET /api/status - provides detailed system status including AI provider info
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get detailed AI provider status
	aiStatus := aiService.GetProviderStatus()
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"services": map[string]interface{}{
			"ai": aiStatus,
		},
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

