package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
	"github.com/google/uuid"
)

// TrackActivity tracks a user activity event
func TrackActivity(userID, eventType, bookID string, data map[string]interface{}) error {
	// Insert into interactions table
	query := `INSERT INTO interactions (id, user_id, event_type, book_id, section_id, page_number, content, context, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	content := ""
	contextJSON := "{}"
	var sectionID *string
	var pageNumber *int
	
	if data != nil {
		if c, ok := data["content"].(string); ok {
			content = c
		}
		if sid, ok := data["section_id"].(*string); ok && sid != nil {
			sectionID = sid
		} else if sidStr, ok := data["section_id"].(string); ok && sidStr != "" {
			sectionID = &sidStr
		}
		if pn, ok := data["page_number"].(*int); ok && pn != nil {
			pageNumber = pn
		} else if pnFloat, ok := data["page_number"].(float64); ok {
			pnInt := int(pnFloat)
			pageNumber = &pnInt
		}
		if ctxJSON, err := json.Marshal(data); err == nil {
			contextJSON = string(ctxJSON)
		}
	}

	activityID := uuid.New().String()
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	
	_, err := database.DB.Exec(query,
		activityID,
		userID,
		eventType,
		bookID,
		sectionID,
		pageNumber,
		content,
		contextJSON,
		createdAt,
	)

	if err != nil {
		return err
	}

	// CRITICAL: Fetch user information for the broadcast - ALWAYS ensure we have user data
	var firstName, lastName, email sql.NullString
	userQuery := `SELECT first_name, last_name, email FROM users WHERE id = ?`
	err = database.DB.QueryRow(userQuery, userID).Scan(&firstName, &lastName, &email)
	if err != nil {
		// If user not found, log error and skip broadcast (don't send incomplete data)
		log.Printf("ERROR: User %s not found for activity broadcast - skipping broadcast", userID)
		// Still save the activity, but don't broadcast incomplete data
		return nil
	}
	
	// Validate we have at least email or name
	if !firstName.Valid && !lastName.Valid && !email.Valid {
		log.Printf("ERROR: User %s has no name or email - skipping broadcast", userID)
		return nil
	}

	// Broadcast activity to consultants with full user info
	// CRITICAL: Always include user_id and ensure all user fields are properly set
	activityData := map[string]interface{}{
		"id":          activityID,
		"user_id":     userID, // CRITICAL: Always include user_id for identification
		"first_name":  "",
		"last_name":   "",
		"email":       "",
		"event_type":  eventType,
		"book_id":     bookID,
		"section_id":  sectionID,
		"page_number": pageNumber,
		"content":     content,
		"context":     contextJSON,
		"created_at":  createdAt,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	
	// Set user fields only if valid
	if firstName.Valid {
		activityData["first_name"] = firstName.String
	}
	if lastName.Valid {
		activityData["last_name"] = lastName.String
	}
	if email.Valid {
		activityData["email"] = email.String
	}
	
	// CRITICAL: Validate user_id is present before broadcasting
	if userID == "" {
		log.Printf("ERROR: Cannot broadcast activity - user_id is empty")
		return nil
	}
	
	// Parse context if it's JSON
	var parsedContext map[string]interface{}
	if contextJSON != "" && contextJSON != "{}" {
		json.Unmarshal([]byte(contextJSON), &parsedContext)
		activityData["parsed_context"] = parsedContext
	}
	
	BroadcastActivity(activityData)

	return nil
}

// HandleTrackActivity handles POST /api/activity/track
func HandleTrackActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract and validate token to get user_id (SECURITY: Never trust user_id from request body)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Extract user_id from token (not from request body)
	userID := claims.UserID

	var req struct {
		EventType string                 `json:"event_type"`
		BookID    string                 `json:"book_id"`
		SectionID *string                `json:"section_id"`
		PageNumber *int                  `json:"page_number"`
		Content   string                 `json:"content"`
		Context   map[string]interface{} `json:"context"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"content":    req.Content,
		"section_id": req.SectionID,
		"page_number": req.PageNumber,
	}
	if req.Context != nil {
		for k, v := range req.Context {
			data[k] = v
		}
	}

	err = TrackActivity(userID, req.EventType, req.BookID, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error tracking activity: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "tracked",
	})
}

