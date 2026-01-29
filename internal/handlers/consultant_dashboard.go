package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

// HandleUpdateBookPurchaseDate handles PUT /api/consultant/reader/purchase-date
// Updates the book purchase date for a reader
func HandleUpdateBookPurchaseDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID and book ID from request
	var req struct {
		UserID      string `json:"user_id"`
		BookID      string `json:"book_id"`
		PurchaseDate string `json:"purchase_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.BookID == "" {
		http.Error(w, "user_id and book_id are required", http.StatusBadRequest)
		return
	}

	// Update purchase date (empty string means clear the date)
	err := database.UpdateBookPurchaseDate(req.UserID, req.BookID, req.PurchaseDate)
	if err != nil {
		log.Printf("Error updating purchase date: %v", err)
		http.Error(w, fmt.Sprintf("Failed to update purchase date: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "Purchase date updated",
	})
}

// HandleGetOnlineReaders handles GET /api/consultant/online-readers
// Returns a map of reader IDs that are currently online
func HandleGetOnlineReaders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	onlineMap, err := database.GetOnlineReaderIDs()
	if err != nil {
		http.Error(w, "Failed to fetch online readers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"online_readers": onlineMap,
	})
}

// HandleConsultantAIInsight handles GET /api/consultant/ai-insight
// scope=dashboard (last N hours) or scope=reader&user_id=... (single reader)
// hours=24 (default for dashboard) or 168 for reader
func HandleConsultantAIInsight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	scope := r.URL.Query().Get("scope")
	if scope != "dashboard" && scope != "reader" {
		respondAIInsightError(w, "scope must be 'dashboard' or 'reader'", http.StatusBadRequest)
		return
	}

	hoursStr := r.URL.Query().Get("hours")
	hours := 24
	if scope == "reader" {
		hours = 168 // 7 days default for single reader
	}
	if hoursStr != "" {
		if p, err := strconv.Atoi(hoursStr); err == nil && p > 0 {
			hours = p
		}
	}

	var context strings.Builder

	if scope == "dashboard" {
		minutes := hours * 60
		active, _ := database.GetActiveReaders(minutes)
		activities, _ := database.GetRecentActivities(100)
		helpReqs, _ := database.GetRecentHelpRequests(50)

		context.WriteString(fmt.Sprintf("Time window: last %d hours\n\n", hours))
		context.WriteString(fmt.Sprintf("Active readers (last %d hours): %d\n", hours, len(active)))
		for _, r := range active {
			context.WriteString(fmt.Sprintf("  - %s %s (user_id=%s, last_active=%s)\n", r.FirstName, r.LastName, r.UserID, r.LastActiveAt.Format("2006-01-02 15:04")))
		}
		context.WriteString(fmt.Sprintf("\nRecent activities (last 100): total=%d\n", len(activities)))
		typeCount := make(map[string]int)
		for _, a := range activities {
			typeCount[a.ActivityType]++
		}
		for t, c := range typeCount {
			context.WriteString(fmt.Sprintf("  - %s: %d\n", t, c))
		}
		context.WriteString(fmt.Sprintf("\nHelp requests (last 50): total=%d\n", len(helpReqs)))
		pending := 0
		for _, h := range helpReqs {
			if h.Status == "pending" || h.Status == "assigned" {
				pending++
			}
			preview := h.Content
			if len(preview) > 120 {
				preview = preview[:120] + "..."
			}
			context.WriteString(fmt.Sprintf("  - user_id=%s status=%s at %s: %s\n", h.UserID, h.Status, h.CreatedAt.Format("2006-01-02 15:04"), preview))
		}
		context.WriteString(fmt.Sprintf("\nPending/assigned help requests: %d\n", pending))
	} else {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			respondAIInsightError(w, "user_id required when scope=reader", http.StatusBadRequest)
			return
		}
		summary, err := database.GetReaderActivitySummary(userID, hours)
		if err != nil {
			log.Printf("GetReaderActivitySummary error: %v", err)
			respondAIInsightError(w, "Failed to load reader activity", http.StatusInternalServerError)
			return
		}
		state, _ := database.GetReaderState(userID)
		helpReqs, _ := database.GetHelpRequests(userID)
		userActivities, _ := database.GetUserActivities(userID, 30)

		context.WriteString(fmt.Sprintf("Reader user_id=%s, last %d hours\n\n", userID, hours))
		context.WriteString(fmt.Sprintf("Activity summary: total=%d, active_days=%d, word_lookups=%d, ai_interactions=%d, page_views=%d\n",
			summary.TotalActivities, summary.ActiveDays, summary.WordLookups, summary.AIIteractions, summary.PageViews))
		if state != nil {
			context.WriteString(fmt.Sprintf("Current state: last_activity_type=%s, last_active=%s, status=%s\n",
				strOrNil(state.LastActivityType), state.LastActivityAt.Format("2006-01-02 15:04"), state.Status))
		}
		context.WriteString(fmt.Sprintf("\nHelp requests: %d\n", len(helpReqs)))
		for _, h := range helpReqs {
			preview := h.Content
			if len(preview) > 150 {
				preview = preview[:150] + "..."
			}
			context.WriteString(fmt.Sprintf("  - status=%s at %s: %s\n", h.Status, h.CreatedAt.Format("2006-01-02 15:04"), preview))
		}
		context.WriteString(fmt.Sprintf("\nRecent activities (last 30): %d\n", len(userActivities)))
		for i, a := range userActivities {
			if i >= 15 {
				context.WriteString("  ...\n")
				break
			}
			context.WriteString(fmt.Sprintf("  - %s type=%s\n", a.CreatedAt.Format("2006-01-02 15:04"), a.ActivityType))
		}
	}

	summary, err := aiService.ConsultantAnalyst(scope, context.String())
	if err != nil {
		log.Printf("ConsultantAnalyst error: %v", err)
		respondAIInsightError(w, "AI summary unavailable. Check that GEMINI_API_KEY or MOONSHOT_API_KEY is set.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"summary": summary})
}

func respondAIInsightError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": msg})
}

func strOrNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
