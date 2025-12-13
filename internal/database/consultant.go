package database

import (
	"database/sql"
	"fmt"
	"time"
)

// ActiveReader represents a reader who is currently active
type ActiveReader struct {
	UserID      string
	Email       string
	FirstName   string
	LastName    string
	BookID      string
	CurrentPage int
	LastActiveAt time.Time
	Status      string
}

// GetActiveReaders returns readers active in the last N minutes
func GetActiveReaders(minutesThreshold int) ([]*ActiveReader, error) {
	query := `SELECT u.id, u.email, u.first_name, u.last_name, 
	                 COALESCE(rs.book_id, ''), COALESCE(rs.current_page, 0),
	                 COALESCE(rs.last_activity_at, u.last_active_at, u.created_at) as last_active,
	                 COALESCE(rs.status, 'idle') as status
	          FROM users u
	          LEFT JOIN reader_states rs ON u.id = rs.user_id
	          WHERE u.role = 'reader' 
	          AND u.is_verified = 1
	          AND (
	              rs.last_activity_at >= datetime('now', '-' || ? || ' minutes')
	              OR u.last_active_at >= datetime('now', '-' || ? || ' minutes')
	              OR (rs.last_activity_at IS NULL AND u.last_active_at IS NULL AND u.created_at >= datetime('now', '-' || ? || ' minutes'))
	          )
	          ORDER BY last_active DESC`

	rows, err := DB.Query(query, minutesThreshold, minutesThreshold, minutesThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to query active readers: %w", err)
	}
	defer rows.Close()

	var readers []*ActiveReader
	for rows.Next() {
		r := &ActiveReader{}
		var lastActiveStr string
		err := rows.Scan(
			&r.UserID, &r.Email, &r.FirstName, &r.LastName,
			&r.BookID, &r.CurrentPage, &lastActiveStr, &r.Status,
		)
		if err != nil {
			continue
		}
		timeLayout := "2006-01-02 15:04:05"
		r.LastActiveAt, _ = time.Parse(timeLayout, lastActiveStr)
		readers = append(readers, r)
	}

	return readers, rows.Err()
}

// ReaderActivitySummary represents activity summary for a reader
type ReaderActivitySummary struct {
	TotalActivities int
	ActiveDays      int
	WordLookups     int
	AIIteractions   int
	PageViews       int
}

// GetReaderActivitySummary returns activity summary for a specific reader
func GetReaderActivitySummary(userID string, hours int) (*ReaderActivitySummary, error) {
	query := `SELECT 
	          COUNT(*) as total_activities,
	          COUNT(DISTINCT DATE(created_at)) as active_days,
	          SUM(CASE WHEN activity_type = 'WORD_LOOKUP' THEN 1 ELSE 0 END) as word_lookups,
	          SUM(CASE WHEN activity_type = 'AI_INTERACTION' THEN 1 ELSE 0 END) as ai_interactions,
	          SUM(CASE WHEN activity_type = 'PAGE_VIEW' THEN 1 ELSE 0 END) as page_views
	          FROM activity_logs
	          WHERE user_id = ? 
	          AND created_at >= datetime('now', '-' || ? || ' hours')`

	summary := &ReaderActivitySummary{}
	err := DB.QueryRow(query, userID, hours).Scan(
		&summary.TotalActivities,
		&summary.ActiveDays,
		&summary.WordLookups,
		&summary.AIIteractions,
		&summary.PageViews,
	)

	if err == sql.ErrNoRows {
		// Return empty summary if no activities found
		return &ReaderActivitySummary{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reader activity summary: %w", err)
	}

	return summary, nil
}

// GetReaderState retrieves the current state of a reader
func GetReaderState(userID string) (*ReaderState, error) {
	query := `SELECT user_id, book_id, current_page, current_section_id, last_activity_type, 
	                 last_activity_at, total_pages_read, total_word_lookups, total_ai_interactions, 
	                 status, updated_at
	          FROM reader_states
	          WHERE user_id = ?`

	state := &ReaderState{}
	var bookID, sectionID, lastActivityType, lastActivityAt, status, updatedAt sql.NullString
	var currentPage, totalPagesRead, totalWordLookups, totalAIInteractions sql.NullInt64

	err := DB.QueryRow(query, userID).Scan(
		&state.UserID, &bookID, &currentPage, &sectionID, &lastActivityType,
		&lastActivityAt, &totalPagesRead, &totalWordLookups, &totalAIInteractions,
		&status, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reader state: %w", err)
	}

	if bookID.Valid {
		state.BookID = &bookID.String
	}
	if sectionID.Valid {
		state.SectionID = &sectionID.String
	}
	if currentPage.Valid {
		pageNum := int(currentPage.Int64)
		state.CurrentPage = &pageNum
	}
	if lastActivityType.Valid {
		state.LastActivityType = &lastActivityType.String
	}
	if lastActivityAt.Valid {
		timeLayout := "2006-01-02 15:04:05"
		state.LastActivityAt, _ = time.Parse(timeLayout, lastActivityAt.String)
	}
	if status.Valid {
		state.Status = status.String
	} else {
		state.Status = "idle"
	}
	if updatedAt.Valid {
		timeLayout := "2006-01-02 15:04:05"
		state.UpdatedAt, _ = time.Parse(timeLayout, updatedAt.String)
	}
	state.TotalPagesRead = int(totalPagesRead.Int64)
	state.TotalWordLookups = int(totalWordLookups.Int64)
	state.TotalAIInteractions = int(totalAIInteractions.Int64)

	return state, nil
}

// ReaderState represents the denormalized state of a reader
type ReaderState struct {
	UserID            string
	BookID            *string
	CurrentPage       *int
	SectionID         *string
	LastActivityType  *string
	LastActivityAt    time.Time
	TotalPagesRead    int
	TotalWordLookups   int
	TotalAIInteractions int
	Status            string
	UpdatedAt         time.Time
}

