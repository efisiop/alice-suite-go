package database

import (
	"database/sql"
	"fmt"
	"time"
)

const readingVisitGap = 30 * time.Minute

// ReadingPosition is the reader's most recently recorded position in a book.
type ReadingPosition struct {
	BookID     *string   `json:"book_id,omitempty"`
	PageNumber *int      `json:"page_number,omitempty"`
	SectionID  *string   `json:"section_id,omitempty"`
	RecordedAt time.Time `json:"recorded_at,omitempty"`
}

// ReadingVisit groups page views close together into one reading visit.
// Direction is intentionally descriptive, not a judgement about reader ability.
type ReadingVisit struct {
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	StartPage int       `json:"start_page"`
	EndPage   int       `json:"end_page"`
	PageDelta int       `json:"page_delta"`
	Direction string    `json:"direction"`
}

// ReaderJourney is the consultant-facing view of a reader's current position
// and recent reading visits.
type ReaderJourney struct {
	Position *ReadingPosition `json:"position,omitempty"`
	Visits   []ReadingVisit   `json:"visits"`
}

type pageViewEvent struct {
	bookID    *string
	page      int
	sectionID *string
	createdAt time.Time
}

// GetReaderJourney derives reading visits from PAGE_VIEW records. A visit ends
// after a 30-minute inactivity gap. It reads activity_logs only because the
// older interactions table mirrors these events for legacy dashboard support.
func GetReaderJourney(userID string, limit int) (*ReaderJourney, error) {
	if limit <= 0 {
		limit = 5
	}

	rows, err := DB.Query(Rebind(`SELECT book_id, page_number, section_id, created_at
		FROM activity_logs
		WHERE user_id = ? AND activity_type = 'PAGE_VIEW' AND page_number IS NOT NULL
		ORDER BY created_at ASC`), userID)
	if err != nil {
		return nil, fmt.Errorf("query reader page views: %w", err)
	}
	defer rows.Close()

	var events []pageViewEvent
	for rows.Next() {
		var event pageViewEvent
		var bookID, sectionID sql.NullString
		var createdAt string
		if err := rows.Scan(&bookID, &event.page, &sectionID, &createdAt); err != nil {
			return nil, fmt.Errorf("scan reader page view: %w", err)
		}
		parsed, err := parseActivityTime(createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse reader page view time: %w", err)
		}
		event.createdAt = parsed
		if bookID.Valid {
			event.bookID = &bookID.String
		}
		if sectionID.Valid {
			event.sectionID = &sectionID.String
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reader page views: %w", err)
	}

	journey := &ReaderJourney{Visits: []ReadingVisit{}}
	if len(events) == 0 {
		return journey, nil
	}

	latest := events[len(events)-1]
	journey.Position = &ReadingPosition{
		BookID: latest.bookID, PageNumber: &latest.page, SectionID: latest.sectionID, RecordedAt: latest.createdAt,
	}

	visits := []ReadingVisit{newReadingVisit(events[0])}
	for _, event := range events[1:] {
		current := &visits[len(visits)-1]
		if event.createdAt.Sub(current.EndedAt) > readingVisitGap {
			visits = append(visits, newReadingVisit(event))
			continue
		}
		current.EndedAt = event.createdAt
		current.EndPage = event.page
		current.PageDelta = current.EndPage - current.StartPage
		current.Direction = readingDirection(current.PageDelta)
	}

	start := len(visits) - limit
	if start < 0 {
		start = 0
	}
	for i := len(visits) - 1; i >= start; i-- {
		journey.Visits = append(journey.Visits, visits[i])
	}
	return journey, nil
}

func newReadingVisit(event pageViewEvent) ReadingVisit {
	return ReadingVisit{
		StartedAt: event.createdAt, EndedAt: event.createdAt, StartPage: event.page, EndPage: event.page,
		PageDelta: 0, Direction: readingDirection(0),
	}
}

func readingDirection(delta int) string {
	switch {
	case delta > 0:
		return "forward"
	case delta < 0:
		return "revisiting"
	default:
		return "unchanged"
	}
}

func parseActivityTime(value string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported timestamp %q", value)
}
