package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ActivityLog represents an activity log entry
type ActivityLog struct {
	ID          string
	UserID      string
	SessionID   *string
	ActivityType string
	BookID      *string
	PageNumber  *int
	SectionID   *string
	Metadata    map[string]interface{}
	CreatedAt   time.Time
}

// LogActivity records an activity in the database
func LogActivity(activity *ActivityLog) error {
	activity.ID = uuid.New().String()

	var metadataJSON string
	if activity.Metadata != nil {
		jsonBytes, err := json.Marshal(activity.Metadata)
		if err == nil {
			metadataJSON = string(jsonBytes)
		}
	}

	query := `INSERT INTO activity_logs 
	          (id, user_id, session_id, activity_type, book_id, page_number, section_id, metadata, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(query,
		activity.ID, activity.UserID, activity.SessionID, activity.ActivityType,
		activity.BookID, activity.PageNumber, activity.SectionID, metadataJSON, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	// Update reader_states table (denormalized)
	return updateReaderState(activity)
}

// updateReaderState updates the reader_states table
func updateReaderState(activity *ActivityLog) error {
	// Only update reader_states for readers
	// First check if user is a reader
	var userRole string
	err := DB.QueryRow(`SELECT role FROM users WHERE id = ?`, activity.UserID).Scan(&userRole)
	if err != nil {
		// If user not found, skip state update
		return nil
	}
	if userRole != "reader" {
		// Only update states for readers
		return nil
	}

	// Check if reader_state exists
	var exists bool
	err = DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM reader_states WHERE user_id = ?)`, activity.UserID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check reader_state existence: %w", err)
	}

	if !exists {
		// Create new reader_state
		query := `INSERT INTO reader_states 
		          (user_id, book_id, current_page, current_section_id, last_activity_type, last_activity_at, status, updated_at)
		          VALUES (?, ?, ?, ?, ?, datetime('now'), 'active', datetime('now'))`
		_, err = DB.Exec(query, activity.UserID, activity.BookID, activity.PageNumber, activity.SectionID, activity.ActivityType)
		if err != nil {
			return fmt.Errorf("failed to create reader_state: %w", err)
		}
		return nil
	}

	// Update existing reader_state
	query := `UPDATE reader_states SET
	          book_id = COALESCE(?, book_id),
	          current_page = COALESCE(?, current_page),
	          current_section_id = COALESCE(?, current_section_id),
	          last_activity_type = ?,
	          last_activity_at = datetime('now'),
	          status = 'active',
	          updated_at = datetime('now')
	          WHERE user_id = ?`

	_, err = DB.Exec(query, activity.BookID, activity.PageNumber, activity.SectionID, activity.ActivityType, activity.UserID)
	if err != nil {
		return fmt.Errorf("failed to update reader_state: %w", err)
	}
	return nil
}

// GetRecentActivities retrieves recent activities (for consultant dashboard)
func GetRecentActivities(limit int) ([]*ActivityLog, error) {
	query := `SELECT id, user_id, session_id, activity_type, book_id, page_number, section_id, metadata, created_at
	          FROM activity_logs
	          ORDER BY created_at DESC
	          LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent activities: %w", err)
	}
	defer rows.Close()

	var activities []*ActivityLog
	for rows.Next() {
		activity := &ActivityLog{}
		var sessionID, bookID, sectionID sql.NullString
		var pageNumber sql.NullInt64
		var metadataJSON sql.NullString
		var createdAtStr string

		err := rows.Scan(
			&activity.ID, &activity.UserID, &sessionID, &activity.ActivityType,
			&bookID, &pageNumber, &sectionID, &metadataJSON, &createdAtStr,
		)
		if err != nil {
			continue
		}

		if sessionID.Valid {
			activity.SessionID = &sessionID.String
		}
		if bookID.Valid {
			activity.BookID = &bookID.String
		}
		if sectionID.Valid {
			activity.SectionID = &sectionID.String
		}
		if pageNumber.Valid {
			pageNum := int(pageNumber.Int64)
			activity.PageNumber = &pageNum
		}
		if metadataJSON.Valid && metadataJSON.String != "" {
			json.Unmarshal([]byte(metadataJSON.String), &activity.Metadata)
		}
		timeLayout := "2006-01-02 15:04:05"
		activity.CreatedAt, _ = time.Parse(timeLayout, createdAtStr)

		activities = append(activities, activity)
	}

	return activities, rows.Err()
}

// GetUserActivities retrieves activities for a specific user
func GetUserActivities(userID string, limit int) ([]*ActivityLog, error) {
	query := `SELECT id, user_id, session_id, activity_type, book_id, page_number, section_id, metadata, created_at
	          FROM activity_logs
	          WHERE user_id = ?
	          ORDER BY created_at DESC
	          LIMIT ?`

	rows, err := DB.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query user activities: %w", err)
	}
	defer rows.Close()

	var activities []*ActivityLog
	for rows.Next() {
		activity := &ActivityLog{}
		var sessionID, bookID, sectionID sql.NullString
		var pageNumber sql.NullInt64
		var metadataJSON sql.NullString
		var createdAtStr string

		err := rows.Scan(
			&activity.ID, &activity.UserID, &sessionID, &activity.ActivityType,
			&bookID, &pageNumber, &sectionID, &metadataJSON, &createdAtStr,
		)
		if err != nil {
			continue
		}

		if sessionID.Valid {
			activity.SessionID = &sessionID.String
		}
		if bookID.Valid {
			activity.BookID = &bookID.String
		}
		if sectionID.Valid {
			activity.SectionID = &sectionID.String
		}
		if pageNumber.Valid {
			pageNum := int(pageNumber.Int64)
			activity.PageNumber = &pageNum
		}
		if metadataJSON.Valid && metadataJSON.String != "" {
			json.Unmarshal([]byte(metadataJSON.String), &activity.Metadata)
		}
		timeLayout := "2006-01-02 15:04:05"
		activity.CreatedAt, _ = time.Parse(timeLayout, createdAtStr)

		activities = append(activities, activity)
	}

	return activities, rows.Err()
}

