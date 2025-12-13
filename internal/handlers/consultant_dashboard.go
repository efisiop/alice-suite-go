package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/efisiopittau/alice-suite-go/internal/database"
)

// HandleConsultantActiveReaders handles GET /api/consultant/active-readers
// Returns list of readers active in the last N minutes
func HandleConsultantActiveReaders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get minutes threshold from query parameter (default: 30 minutes)
	minutesStr := r.URL.Query().Get("minutes")
	minutes := 30 // default
	if minutesStr != "" {
		if parsed, err := strconv.Atoi(minutesStr); err == nil && parsed > 0 {
			minutes = parsed
		}
	}

	readers, err := database.GetActiveReaders(minutes)
	if err != nil {
		http.Error(w, "Failed to fetch active readers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":   len(readers),
		"readers": readers,
	})
}

// HandleConsultantReaderActivity handles GET /api/consultant/reader/:id/activity
// Returns activity summary for a specific reader
func HandleConsultantReaderActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	// Assuming URL pattern: /api/consultant/reader/:id/activity
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		// Try to extract from path
		// This is a simplified version - you may need to adjust based on your routing
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	// Get hours threshold from query parameter (default: 24 hours)
	hoursStr := r.URL.Query().Get("hours")
	hours := 24 // default
	if hoursStr != "" {
		if parsed, err := strconv.Atoi(hoursStr); err == nil && parsed > 0 {
			hours = parsed
		}
	}

	summary, err := database.GetReaderActivitySummary(userID, hours)
	if err != nil {
		http.Error(w, "Failed to fetch reader activity", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// HandleConsultantRecentActivities handles GET /api/consultant/recent-activities
// Returns recent activity feed for all readers
func HandleConsultantRecentActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get limit from query parameter (default: 100)
	limitStr := r.URL.Query().Get("limit")
	limit := 100 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	activities, err := database.GetRecentActivities(limit)
	if err != nil {
		http.Error(w, "Failed to fetch recent activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":      len(activities),
		"activities": activities,
	})
}

// HandleConsultantReaderState handles GET /api/consultant/reader/:id/state
// Returns current state of a specific reader
func HandleConsultantReaderState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from query parameter
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	state, err := database.GetReaderState(userID)
	if err != nil {
		http.Error(w, "Failed to fetch reader state", http.StatusInternalServerError)
		return
	}

	if state == nil {
		http.Error(w, "Reader state not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

