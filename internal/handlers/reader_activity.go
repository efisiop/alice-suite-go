package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
)

// ReaderActivity represents a reader interaction with user info
type ReaderActivity struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	EventType   string    `json:"event_type"`
	BookID      string    `json:"book_id"`
	SectionID   *string   `json:"section_id"`
	PageNumber  *int      `json:"page_number"`
	Content     string    `json:"content"`
	Context     string    `json:"context"`
	CreatedAt   string    `json:"created_at"`
	ParsedContext map[string]interface{} `json:"parsed_context,omitempty"`
}

// HandleGetReaderActivities handles GET /api/consultant/reader-activities
// Returns all reader interactions with user information, ordered by most recent
func HandleGetReaderActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100" // Default to 100 most recent
	}

	// Query interactions with user join
	// CRITICAL: Only show activities from actual reader users with complete user data
	// Exclude test users and ensure user_id is always present
	query := `
		SELECT 
			i.id,
			i.user_id,
			u.first_name,
			u.last_name,
			u.email,
			i.event_type,
			i.book_id,
			i.section_id,
			i.page_number,
			i.content,
			i.context,
			i.created_at
		FROM interactions i
		INNER JOIN users u ON i.user_id = u.id
		WHERE u.role = 'reader'
		AND i.user_id IS NOT NULL
		AND u.id IS NOT NULL
		AND (u.first_name IS NOT NULL OR u.last_name IS NOT NULL OR u.email IS NOT NULL)
		ORDER BY i.created_at DESC
		LIMIT ?
	`

	rows, err := database.DB.Query(query, limit)
	if err != nil {
		log.Printf("Database error in HandleGetReaderActivities: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	activities := []ReaderActivity{}
	for rows.Next() {
		var activity ReaderActivity
		var sectionID sql.NullString
		var pageNumber sql.NullInt64
		var firstName sql.NullString
		var lastName sql.NullString
		var email sql.NullString

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&firstName,
			&lastName,
			&email,
			&activity.EventType,
			&activity.BookID,
			&sectionID,
			&pageNumber,
			&activity.Content,
			&activity.Context,
			&activity.CreatedAt,
		)
		if err != nil {
			continue
		}

		// CRITICAL: Validate user data - ensure user_id is always present
		if activity.UserID == "" {
			log.Printf("WARNING: Activity %s has empty user_id - skipping", activity.ID)
			continue
		}

		if firstName.Valid {
			activity.FirstName = firstName.String
		}
		if lastName.Valid {
			activity.LastName = lastName.String
		}
		if email.Valid {
			activity.Email = email.String
		}
		
		// CRITICAL: Validate that we have at least one identifier (name or email)
		if activity.FirstName == "" && activity.LastName == "" && activity.Email == "" {
			log.Printf("WARNING: Activity %s for user %s has no name or email - skipping", activity.ID, activity.UserID)
			continue
		}
		
		if sectionID.Valid {
			activity.SectionID = &sectionID.String
		}
		if pageNumber.Valid {
			pageNum := int(pageNumber.Int64)
			activity.PageNumber = &pageNum
		}

		// Parse context JSON if available
		if activity.Context != "" {
			var parsedContext map[string]interface{}
			if err := json.Unmarshal([]byte(activity.Context), &parsedContext); err == nil {
				activity.ParsedContext = parsedContext
			}
		}

		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error reading rows in HandleGetReaderActivities: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

// HandleGetReaderActivityStream handles GET /api/consultant/reader-activities/stream
// Returns recent activities since a given timestamp
func HandleGetReaderActivityStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get since parameter (SQLite datetime format: YYYY-MM-DD HH:MM:SS)
	since := r.URL.Query().Get("since")
	if since == "" {
		// Default to last 5 minutes
		fiveMinutesAgo := time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
		since = fiveMinutesAgo
	} else {
		// Try to parse ISO format and convert to SQLite format if needed
		if parsedTime, err := time.Parse(time.RFC3339, since); err == nil {
			since = parsedTime.Format("2006-01-02 15:04:05")
		} else if parsedTime, err := time.Parse("2006-01-02T15:04:05", since); err == nil {
			since = parsedTime.Format("2006-01-02 15:04:05")
		}
		// If parsing fails, use as-is (assuming it's already in SQLite format)
	}

	query := `
		SELECT 
			i.id,
			i.user_id,
			u.first_name,
			u.last_name,
			u.email,
			i.event_type,
			i.book_id,
			i.section_id,
			i.page_number,
			i.content,
			i.context,
			i.created_at
		FROM interactions i
		INNER JOIN users u ON i.user_id = u.id
		WHERE u.role = 'reader'
		AND i.user_id IS NOT NULL
		AND u.id IS NOT NULL
		AND (u.first_name IS NOT NULL OR u.last_name IS NOT NULL OR u.email IS NOT NULL)
		AND i.created_at >= ?
		ORDER BY i.created_at DESC
		LIMIT 100
	`

	rows, err := database.DB.Query(query, since)
	if err != nil {
		log.Printf("Database error in HandleGetReaderActivityStream: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	activities := []ReaderActivity{}
	for rows.Next() {
		var activity ReaderActivity
		var sectionID sql.NullString
		var pageNumber sql.NullInt64
		var firstName sql.NullString
		var lastName sql.NullString
		var email sql.NullString

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&firstName,
			&lastName,
			&email,
			&activity.EventType,
			&activity.BookID,
			&sectionID,
			&pageNumber,
			&activity.Content,
			&activity.Context,
			&activity.CreatedAt,
		)
		if err != nil {
			continue
		}

		// CRITICAL: Validate user data - ensure user_id is always present
		if activity.UserID == "" {
			log.Printf("WARNING: Activity %s has empty user_id - skipping", activity.ID)
			continue
		}

		if firstName.Valid {
			activity.FirstName = firstName.String
		}
		if lastName.Valid {
			activity.LastName = lastName.String
		}
		if email.Valid {
			activity.Email = email.String
		}
		
		// CRITICAL: Validate that we have at least one identifier (name or email)
		if activity.FirstName == "" && activity.LastName == "" && activity.Email == "" {
			log.Printf("WARNING: Activity %s for user %s has no name or email - skipping", activity.ID, activity.UserID)
			continue
		}
		
		if sectionID.Valid {
			activity.SectionID = &sectionID.String
		}
		if pageNumber.Valid {
			pageNum := int(pageNumber.Int64)
			activity.PageNumber = &pageNum
		}

		// Parse context JSON if available
		if activity.Context != "" {
			var parsedContext map[string]interface{}
			if err := json.Unmarshal([]byte(activity.Context), &parsedContext); err == nil {
				activity.ParsedContext = parsedContext
			}
		}

		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error reading rows in HandleGetReaderActivities: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

// HandleGetActiveReadersCount handles GET /api/consultant/active-readers-count
// Returns the count of active readers (logged in but not logged out, or logged in more recently than logout)
func HandleGetActiveReadersCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query to find active readers:
	// A reader is active if:
	// 1. They have ANY activity within the last 30 minutes, OR
	// 2. Their most recent LOGIN event is within the last 30 minutes and they don't have a LOGOUT after it
	// 3. They are a reader (role = 'reader')
	thirtyMinutesAgo := time.Now().Add(-30 * time.Minute).Format("2006-01-02 15:04:05")
	
	// Find readers with any recent activity (within last 30 minutes)
	// This includes LOGIN, but also any other activity like PAGE_SYNC, DEFINITION_LOOKUP, etc.
	query := `
		SELECT DISTINCT u.id, u.first_name, u.last_name, u.email
		FROM users u
		INNER JOIN interactions i ON i.user_id = u.id
		LEFT JOIN (
			SELECT user_id, MAX(created_at) as last_logout
			FROM interactions
			WHERE event_type = 'LOGOUT'
			GROUP BY user_id
		) latest_logout ON latest_logout.user_id = u.id
		WHERE u.role = 'reader'
		AND i.created_at >= ?
		AND (latest_logout.last_logout IS NULL OR i.created_at > latest_logout.last_logout)
		ORDER BY i.created_at DESC
	`

	rows, err := database.DB.Query(query, thirtyMinutesAgo)
	if err != nil {
		log.Printf("Database error in HandleGetActiveReadersCount: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ActiveReader struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	activeReaders := []ActiveReader{}
	for rows.Next() {
		var reader ActiveReader
		var firstName sql.NullString
		var lastName sql.NullString
		var email sql.NullString

		err := rows.Scan(&reader.ID, &firstName, &lastName, &email)
		if err != nil {
			continue
		}

		if firstName.Valid {
			reader.FirstName = firstName.String
		}
		if lastName.Valid {
			reader.LastName = lastName.String
		}
		if email.Valid {
			reader.Email = email.String
		}

		activeReaders = append(activeReaders, reader)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error reading rows in HandleGetReaderActivities: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":  len(activeReaders),
		"readers": activeReaders,
	})
}

// HandleGetLoggedInReadersCount handles GET /api/consultant/logged-in-readers-count
// Returns the count of readers who are currently logged in (have active sessions)
func HandleGetLoggedInReadersCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query to find logged-in readers:
	// A reader is "logged in" if they have an active (non-expired) session 
	// AND the session was active within the last hour (to filter out stale sessions)
	query := `
		SELECT DISTINCT 
			u.id, 
			u.first_name, 
			u.last_name, 
			u.email, 
			MAX(s.last_active_at) as last_active
		FROM users u
		INNER JOIN sessions s ON s.user_id = u.id 
			AND s.expires_at > datetime('now')
			AND s.last_active_at >= datetime('now', '-1 hour')
		WHERE u.role = 'reader'
		GROUP BY u.id, u.first_name, u.last_name, u.email
		ORDER BY last_active DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Database error in HandleGetLoggedInReadersCount: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type LoggedInReader struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		LastLogin string `json:"last_login"`
	}

	loggedInReaders := []LoggedInReader{}
	for rows.Next() {
		var reader LoggedInReader
		var firstName sql.NullString
		var lastName sql.NullString
		var email sql.NullString

		err := rows.Scan(&reader.ID, &firstName, &lastName, &email, &reader.LastLogin)
		if err != nil {
			continue
		}

		if firstName.Valid {
			reader.FirstName = firstName.String
		}
		if lastName.Valid {
			reader.LastName = lastName.String
		}
		if email.Valid {
			reader.Email = email.String
		}

		loggedInReaders = append(loggedInReaders, reader)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error reading rows in HandleGetLoggedInReadersCount: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":  len(loggedInReaders),
		"readers": loggedInReaders,
	})
}

// HandleGetTodaysActivityCount handles GET /api/consultant/todays-activity-count
// Returns the count of all reader activities for today
func HandleGetTodaysActivityCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get today's date at midnight
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	todayStartStr := todayStart.Format("2006-01-02 15:04:05")

	// Count all reader activities for today (including LOGIN/LOGOUT)
	query := `
		SELECT COUNT(*) as count
		FROM interactions i
		LEFT JOIN users u ON i.user_id = u.id
		WHERE u.role = 'reader'
		AND i.created_at >= ?
	`

	var count int
	err := database.DB.QueryRow(query, todayStartStr).Scan(&count)
	if err != nil {
		log.Printf("Database error in HandleGetTodaysActivityCount: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": count,
	})
}

// HandleGetLoggedOutCount handles GET /api/consultant/logged-out-count
// Returns the count of reader logouts today
func HandleGetLoggedOutCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get today's date at midnight
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStartStr := todayStart.Format("2006-01-02 15:04:05")

	// Count LOGOUT events today from both activity_logs and interactions tables
	// (TrackActivity logs to interactions, LogActivity logs to activity_logs)
	query := `
		SELECT 
			(SELECT COUNT(*) FROM activity_logs al 
			 LEFT JOIN users u ON al.user_id = u.id 
			 WHERE u.role = 'reader' AND al.activity_type = 'LOGOUT' AND al.created_at >= ?)
			+
			(SELECT COUNT(*) FROM interactions i 
			 LEFT JOIN users u ON i.user_id = u.id 
			 WHERE u.role = 'reader' AND i.event_type = 'LOGOUT' AND i.created_at >= ?)
		as total_count
	`

	var count int
	err := database.DB.QueryRow(query, todayStartStr, todayStartStr).Scan(&count)
	if err != nil {
		log.Printf("Database error in HandleGetLoggedOutCount: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": count,
	})
}
