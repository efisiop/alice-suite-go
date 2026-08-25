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
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}
	if eventType == "" {
		return fmt.Errorf("event_type is required")
	}

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
	now := time.Now()
	createdAt := database.FormatSQLDateTime(now)

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if bookID != "" {
		query := `INSERT INTO interactions (id, user_id, event_type, book_id, section_id, page_number, content, context, created_at)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		if _, err := tx.Exec(database.Rebind(query),
			activityID,
			userID,
			eventType,
			bookID,
			sectionID,
			pageNumber,
			content,
			contextJSON,
			createdAt,
		); err != nil {
			return err
		}
	}

	// CRITICAL: Fetch user information for the broadcast - ALWAYS ensure we have user data
	var firstName, lastName, email sql.NullString
	var role string
	userQuery := `SELECT first_name, last_name, email, role FROM users WHERE id = ?`
	err = tx.QueryRow(database.Rebind(userQuery), userID).Scan(&firstName, &lastName, &email, &role)
	if err != nil {
		return fmt.Errorf("user %s not found for activity: %w", userID, err)
	}

	activityLogType := consultantActivityType(eventType)
	var activityBookID *string
	if bookID != "" {
		activityBookID = &bookID
	}
	metadata := contextJSON
	logQuery := `INSERT INTO activity_logs
	          (id, user_id, session_id, activity_type, book_id, page_number, section_id, metadata, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := tx.Exec(database.Rebind(logQuery),
		activityID,
		userID,
		nil,
		activityLogType,
		activityBookID,
		pageNumber,
		sectionID,
		metadata,
		createdAt,
	); err != nil {
		return fmt.Errorf("failed to log consultant activity: %w", err)
	}

	if role == "reader" {
		if err := upsertReaderState(tx, userID, activityBookID, pageNumber, sectionID, activityLogType, now); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
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
		"timestamp":   now.Format(time.RFC3339),
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

func consultantActivityType(eventType string) string {
	switch eventType {
	case "PAGE_SYNC":
		return "PAGE_VIEW"
	case "DEFINITION_LOOKUP":
		return "WORD_LOOKUP"
	case "AI_HELP", "AI_QUERY":
		return "AI_INTERACTION"
	default:
		return eventType
	}
}

type sqlExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func upsertReaderState(exec sqlExecutor, userID string, bookID *string, pageNumber *int, sectionID *string, activityType string, now time.Time) error {
	var exists bool
	err := exec.QueryRow(database.Rebind(`SELECT EXISTS(SELECT 1 FROM reader_states WHERE user_id = ?)`), userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check reader state: %w", err)
	}

	if !exists {
		query := `INSERT INTO reader_states
		          (user_id, book_id, current_page, current_section_id, last_activity_type, last_activity_at, status, updated_at)
		          VALUES (?, ?, ?, ?, ?, ?, 'active', ?)`
		if _, err := exec.Exec(database.Rebind(query), userID, bookID, pageNumber, sectionID, activityType, database.FormatSQLDateTime(now), database.FormatSQLDateTime(now)); err != nil {
			return fmt.Errorf("failed to create reader state: %w", err)
		}
		return nil
	}

	query := `UPDATE reader_states SET
	          book_id = COALESCE(?, book_id),
	          current_page = COALESCE(?, current_page),
	          current_section_id = COALESCE(?, current_section_id),
	          last_activity_type = ?,
	          last_activity_at = ?,
	          status = 'active',
	          updated_at = ?
	          WHERE user_id = ?`
	if _, err := exec.Exec(database.Rebind(query), bookID, pageNumber, sectionID, activityType, database.FormatSQLDateTime(now), database.FormatSQLDateTime(now), userID); err != nil {
		return fmt.Errorf("failed to update reader state: %w", err)
	}
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
		EventType  string                 `json:"event_type"`
		BookID     string                 `json:"book_id"`
		SectionID  *string                `json:"section_id"`
		PageNumber *int                   `json:"page_number"`
		Content    string                 `json:"content"`
		Context    map[string]interface{} `json:"context"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"content":     req.Content,
		"section_id":  req.SectionID,
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
