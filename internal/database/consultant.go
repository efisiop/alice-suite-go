package database

import (
	"database/sql"
	"fmt"
	"time"
)

// sqlDateKeyYMD is SQL for calendar day YYYY-MM-DD from a timestamp/text column.
// PostgreSQL timestamps must be cast to text before SUBSTR; SQLite TEXT is unchanged.
func sqlDateKeyYMD(column string) string {
	return "SUBSTR(CAST(" + column + " AS TEXT), 1, 10)"
}

// ActiveReader represents a reader who is currently active
type ActiveReader struct {
	UserID       string
	Email        string
	FirstName    string
	LastName     string
	BookID       string
	CurrentPage  int
	LastActiveAt time.Time
	Status       string
}

// GetActiveReaders returns readers active in the last N minutes
func GetActiveReaders(minutesThreshold int) ([]*ActiveReader, error) {
	cutoff := time.Now().Add(-time.Duration(minutesThreshold) * time.Minute)
	query := `SELECT u.id, u.email, u.first_name, u.last_name, 
	                 COALESCE(rs.book_id, ''), COALESCE(rs.current_page, 0),
	                 COALESCE(rs.last_activity_at, u.last_active_at, u.created_at) as last_active,
	                 COALESCE(rs.status, 'idle') as status
	          FROM users u
	          LEFT JOIN reader_states rs ON u.id = rs.user_id
	          WHERE u.role = 'reader' 
	          AND u.is_verified = 1
	          AND (
	              rs.last_activity_at >= ?
	              OR u.last_active_at >= ?
	              OR (rs.last_activity_at IS NULL AND u.last_active_at IS NULL AND u.created_at >= ?)
	          )
	          ORDER BY last_active DESC`

	rows, err := DB.Query(Rebind(query), cutoff, cutoff, cutoff)
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
// Combines data from both interactions and activity_logs tables
func GetReaderActivitySummary(userID string, hours int) (*ReaderActivitySummary, error) {
	// Calculate the cutoff time
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	cutoffTimeStr := cutoffTime.Format("2006-01-02 15:04:05")

	// Query that combines both interactions and activity_logs tables
	// Maps event_type from interactions to match activity_type from activity_logs
	// Use COALESCE to handle NULLs properly in CASE statements
	query := `
		SELECT 
			COUNT(*) as total_activities,
			COUNT(DISTINCT ` + sqlDateKeyYMD("created_at") + `) as active_days,
			COALESCE(SUM(CASE 
				WHEN COALESCE(event_type, '') = 'DEFINITION_LOOKUP' OR COALESCE(activity_type, '') = 'WORD_LOOKUP' THEN 1
				ELSE 0 
			END), 0) as word_lookups,
			COALESCE(SUM(CASE 
				WHEN COALESCE(event_type, '') = 'AI_QUERY' OR COALESCE(activity_type, '') = 'AI_INTERACTION' THEN 1
				ELSE 0 
			END), 0) as ai_interactions,
			COALESCE(SUM(CASE 
				WHEN COALESCE(event_type, '') = 'PAGE_SYNC' OR COALESCE(activity_type, '') = 'PAGE_VIEW' THEN 1
				ELSE 0 
			END), 0) as page_views
		FROM (
			SELECT 
				created_at,
				event_type,
				CAST(NULL AS TEXT) AS activity_type
			FROM interactions
			WHERE user_id = ? 
			AND created_at >= ?
			
			UNION ALL
			
			SELECT 
				created_at,
				CAST(NULL AS TEXT) AS event_type,
				activity_type
			FROM activity_logs
			WHERE user_id = ? 
			AND created_at >= ?
		) combined_activities`

	summary := &ReaderActivitySummary{}
	err := DB.QueryRow(Rebind(query), userID, cutoffTimeStr, userID, cutoffTimeStr).Scan(
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

// ActivityDayPoint is one day of activity count for charts (consultant reader inspector).
type ActivityDayPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetReaderActivityDailySeries returns activity counts per calendar day within the last `hours` hours.
func GetReaderActivityDailySeries(userID string, hours int) ([]ActivityDayPoint, error) {
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	cutoffTimeStr := cutoffTime.Format("2006-01-02 15:04:05")

	query := `
		SELECT day, COUNT(*) as cnt FROM (
			SELECT ` + sqlDateKeyYMD("created_at") + ` AS day
			FROM interactions
			WHERE user_id = ? AND created_at >= ?
			UNION ALL
			SELECT ` + sqlDateKeyYMD("created_at") + ` AS day
			FROM activity_logs
			WHERE user_id = ? AND created_at >= ?
		) t
		GROUP BY day
		ORDER BY day`

	rows, err := DB.Query(Rebind(query), userID, cutoffTimeStr, userID, cutoffTimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get reader activity daily series: %w", err)
	}
	defer rows.Close()

	byDay := make(map[string]int)
	for rows.Next() {
		var day string
		var cnt int
		if err := rows.Scan(&day, &cnt); err != nil {
			continue
		}
		byDay[day] = cnt
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Fill every calendar day from the window start through today (zeros where no events).
	startDay := time.Date(cutoffTime.Year(), cutoffTime.Month(), cutoffTime.Day(), 0, 0, 0, 0, cutoffTime.Location())
	now := time.Now()
	endDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var out []ActivityDayPoint
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		out = append(out, ActivityDayPoint{Date: key, Count: byDay[key]})
	}

	return out, nil
}

// readerVerifiedFilter is reused for dashboard-wide aggregates (verified readers only).
const readerVerifiedFilter = `user_id IN (SELECT id FROM users WHERE role = 'reader' AND is_verified = 1)`

// DashboardReadersActivitySummary aggregates activity across all verified readers in a time window.
type DashboardReadersActivitySummary struct {
	TotalActivities int
	ActiveDays      int
	UniqueReaders   int
	WordLookups     int
	AIIteractions   int
	PageViews       int
}

// GetDashboardReadersActivitySummary returns combined stats for every verified reader in the last `hours` hours.
func GetDashboardReadersActivitySummary(hours int) (*DashboardReadersActivitySummary, error) {
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	cutoffTimeStr := cutoffTime.Format("2006-01-02 15:04:05")

	query := `
		SELECT 
			COUNT(*) as total_activities,
			COUNT(DISTINCT ` + sqlDateKeyYMD("created_at") + `) as active_days,
			COUNT(DISTINCT user_id) as unique_readers,
			COALESCE(SUM(CASE 
				WHEN COALESCE(event_type, '') = 'DEFINITION_LOOKUP' OR COALESCE(activity_type, '') = 'WORD_LOOKUP' THEN 1
				ELSE 0 
			END), 0) as word_lookups,
			COALESCE(SUM(CASE 
				WHEN COALESCE(event_type, '') = 'AI_QUERY' OR COALESCE(activity_type, '') = 'AI_INTERACTION' THEN 1
				ELSE 0 
			END), 0) as ai_interactions,
			COALESCE(SUM(CASE 
				WHEN COALESCE(event_type, '') = 'PAGE_SYNC' OR COALESCE(activity_type, '') = 'PAGE_VIEW' THEN 1
				ELSE 0 
			END), 0) as page_views
		FROM (
			SELECT user_id, created_at, event_type, CAST(NULL AS TEXT) AS activity_type
			FROM interactions
			WHERE ` + readerVerifiedFilter + `
			AND created_at >= ?
			UNION ALL
			SELECT user_id, created_at, CAST(NULL AS TEXT) AS event_type, activity_type
			FROM activity_logs
			WHERE ` + readerVerifiedFilter + `
			AND created_at >= ?
		) combined_activities`

	s := &DashboardReadersActivitySummary{}
	err := DB.QueryRow(Rebind(query), cutoffTimeStr, cutoffTimeStr).Scan(
		&s.TotalActivities,
		&s.ActiveDays,
		&s.UniqueReaders,
		&s.WordLookups,
		&s.AIIteractions,
		&s.PageViews,
	)
	if err == sql.ErrNoRows {
		return &DashboardReadersActivitySummary{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard readers activity summary: %w", err)
	}
	return s, nil
}

// GetDashboardActivityDailySeries returns total reader activity events per calendar day (all verified readers).
func GetDashboardActivityDailySeries(hours int) ([]ActivityDayPoint, error) {
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	cutoffTimeStr := cutoffTime.Format("2006-01-02 15:04:05")

	query := `
		SELECT day, COUNT(*) as cnt FROM (
			SELECT ` + sqlDateKeyYMD("created_at") + ` AS day
			FROM interactions
			WHERE ` + readerVerifiedFilter + `
			AND created_at >= ?
			UNION ALL
			SELECT ` + sqlDateKeyYMD("created_at") + ` AS day
			FROM activity_logs
			WHERE ` + readerVerifiedFilter + `
			AND created_at >= ?
		) t
		GROUP BY day
		ORDER BY day`

	rows, err := DB.Query(Rebind(query), cutoffTimeStr, cutoffTimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard activity daily series: %w", err)
	}
	defer rows.Close()

	byDay := make(map[string]int)
	for rows.Next() {
		var day string
		var cnt int
		if err := rows.Scan(&day, &cnt); err != nil {
			continue
		}
		byDay[day] = cnt
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	startDay := time.Date(cutoffTime.Year(), cutoffTime.Month(), cutoffTime.Day(), 0, 0, 0, 0, cutoffTime.Location())
	now := time.Now()
	endDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var out []ActivityDayPoint
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		out = append(out, ActivityDayPoint{Date: key, Count: byDay[key]})
	}
	return out, nil
}

// HelpDayPoint is help-request volume per day for consultant charts.
type HelpDayPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetDashboardHelpRequestsDailySeries counts new help requests per calendar day in the window.
func GetDashboardHelpRequestsDailySeries(hours int) ([]HelpDayPoint, error) {
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	cutoffTimeStr := cutoffTime.Format("2006-01-02 15:04:05")

	query := `
		SELECT ` + sqlDateKeyYMD("created_at") + ` AS day, COUNT(*) AS cnt
		FROM help_requests
		WHERE created_at >= ?
		GROUP BY ` + sqlDateKeyYMD("created_at") + `
		ORDER BY day`

	rows, err := DB.Query(Rebind(query), cutoffTimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get help requests daily series: %w", err)
	}
	defer rows.Close()

	byDay := make(map[string]int)
	for rows.Next() {
		var day string
		var cnt int
		if err := rows.Scan(&day, &cnt); err != nil {
			continue
		}
		byDay[day] = cnt
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	startDay := time.Date(cutoffTime.Year(), cutoffTime.Month(), cutoffTime.Day(), 0, 0, 0, 0, cutoffTime.Location())
	now := time.Now()
	endDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var out []HelpDayPoint
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		out = append(out, HelpDayPoint{Date: key, Count: byDay[key]})
	}
	return out, nil
}

// DashboardHelpWindow holds help-request counts inside the time window.
type DashboardHelpWindow struct {
	Total int `json:"total"`
	Open  int `json:"open"`
}

// GetDashboardHelpRequestsInWindow counts help requests created in the window; Open is pending + assigned.
func GetDashboardHelpRequestsInWindow(hours int) (*DashboardHelpWindow, error) {
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	cutoffTimeStr := cutoffTime.Format("2006-01-02 15:04:05")

	query := `
		SELECT COUNT(*),
			COALESCE(SUM(CASE WHEN status IN ('pending', 'assigned') THEN 1 ELSE 0 END), 0)
		FROM help_requests
		WHERE created_at >= ?`

	w := &DashboardHelpWindow{}
	err := DB.QueryRow(Rebind(query), cutoffTimeStr).Scan(&w.Total, &w.Open)
	if err != nil {
		return nil, fmt.Errorf("failed to get help requests in window: %w", err)
	}
	return w, nil
}

// CountVerifiedReaders returns the number of verified reader accounts.
func CountVerifiedReaders() (int, error) {
	var n int
	err := DB.QueryRow(Rebind(`SELECT COUNT(*) FROM users WHERE role = 'reader' AND is_verified = 1`)).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("failed to count readers: %w", err)
	}
	return n, nil
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

	err := DB.QueryRow(Rebind(query), userID).Scan(
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
	UserID              string     `json:"user_id"`
	BookID              *string    `json:"book_id,omitempty"`
	CurrentPage         *int       `json:"current_page,omitempty"`
	SectionID           *string    `json:"section_id,omitempty"`
	LastActivityType    *string    `json:"last_activity_type,omitempty"`
	LastActivityAt      time.Time  `json:"last_activity_at"`
	TotalPagesRead      int        `json:"total_pages_read"`
	TotalWordLookups    int        `json:"total_word_lookups"`
	TotalAIInteractions int        `json:"total_ai_interactions"`
	Status              string     `json:"status"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
