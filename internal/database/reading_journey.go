package database

import (
	"database/sql"
	"encoding/json"
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
	StartedAt time.Time     `json:"started_at"`
	EndedAt   time.Time     `json:"ended_at"`
	StartPage int           `json:"start_page"`
	EndPage   int           `json:"end_page"`
	PageDelta int           `json:"page_delta"`
	Direction string        `json:"direction"`
	Pages     []ReadingPage `json:"pages"`
}

// ReadingPage is one distinct page stop within a reading visit. Events are
// attached only while that page is the reader's most recently confirmed page.
type ReadingPage struct {
	PageNumber int            `json:"page_number"`
	StartedAt  time.Time      `json:"started_at"`
	EndedAt    time.Time      `json:"ended_at"`
	Events     []ReadingEvent `json:"events"`
}

// ReadingEvent is support activity that happened while a reader was on a page.
type ReadingEvent struct {
	Kind       string    `json:"kind"`
	Detail     string    `json:"detail"`
	OccurredAt time.Time `json:"occurred_at"`
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
	kind      string
	metadata  string
	createdAt time.Time
}

// GetReaderJourney derives reading visits from PAGE_VIEW records. A visit ends
// after a 30-minute inactivity gap. It reads activity_logs only because the
// older interactions table mirrors these events for legacy dashboard support.
func GetReaderJourney(userID string, limit int) (*ReaderJourney, error) {
	if limit <= 0 {
		limit = 5
	}

	rows, err := DB.Query(Rebind(`SELECT book_id, page_number, section_id, activity_type, metadata, created_at
		FROM activity_logs
		WHERE user_id = ?
		ORDER BY created_at ASC`), userID)
	if err != nil {
		return nil, fmt.Errorf("query reader page views: %w", err)
	}
	defer rows.Close()

	var events []pageViewEvent
	for rows.Next() {
		var event pageViewEvent
		var bookID, sectionID, metadata sql.NullString
		var createdAt string
		var pageNumber sql.NullInt64
		if err := rows.Scan(&bookID, &pageNumber, &sectionID, &event.kind, &metadata, &createdAt); err != nil {
			return nil, fmt.Errorf("scan reader page view: %w", err)
		}
		parsed, err := parseActivityTime(createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse reader page view time: %w", err)
		}
		event.createdAt = parsed
		if pageNumber.Valid {
			event.page = int(pageNumber.Int64)
		}
		if bookID.Valid {
			event.bookID = &bookID.String
		}
		if sectionID.Valid {
			event.sectionID = &sectionID.String
		}
		if metadata.Valid {
			event.metadata = metadata.String
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reader page views: %w", err)
	}

	journey := &ReaderJourney{Visits: []ReadingVisit{}}
	var pageEvents []pageViewEvent
	for _, event := range events {
		if event.kind == "PAGE_VIEW" && event.page > 0 {
			pageEvents = append(pageEvents, event)
		}
	}
	if len(pageEvents) == 0 {
		return journey, nil
	}

	latest := pageEvents[len(pageEvents)-1]
	journey.Position = &ReadingPosition{
		BookID: latest.bookID, PageNumber: &latest.page, SectionID: latest.sectionID, RecordedAt: latest.createdAt,
	}

	var visits []ReadingVisit
	for _, event := range events {
		if event.kind == "PAGE_VIEW" && event.page > 0 {
			if len(visits) == 0 {
				visits = append(visits, newReadingVisit(event))
				continue
			}
			current := &visits[len(visits)-1]
			if event.createdAt.Sub(current.EndedAt) > readingVisitGap {
				visits = append(visits, newReadingVisit(event))
				continue
			}
			current.EndedAt = event.createdAt
			current.EndPage = event.page
			current.PageDelta = current.EndPage - current.StartPage
			current.Direction = readingDirection(current.PageDelta)
			page := &current.Pages[len(current.Pages)-1]
			if page.PageNumber != event.page {
				current.Pages = append(current.Pages, ReadingPage{PageNumber: event.page, StartedAt: event.createdAt, EndedAt: event.createdAt, Events: []ReadingEvent{}})
			} else {
				page.EndedAt = event.createdAt
			}
			continue
		}

		if len(visits) == 0 {
			continue
		}
		current := &visits[len(visits)-1]
		if event.createdAt.Sub(current.EndedAt) > readingVisitGap {
			continue
		}
		page := &current.Pages[len(current.Pages)-1]
		page.Events = append(page.Events, ReadingEvent{Kind: event.kind, Detail: readingEventDetail(event.kind, event.metadata), OccurredAt: event.createdAt})
		page.EndedAt = event.createdAt
		current.EndedAt = event.createdAt
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
		Pages: []ReadingPage{{PageNumber: event.page, StartedAt: event.createdAt, EndedAt: event.createdAt, Events: []ReadingEvent{}}},
	}
}

func readingEventDetail(kind, metadata string) string {
	var values map[string]interface{}
	_ = json.Unmarshal([]byte(metadata), &values)
	switch kind {
	case "DICTIONARY_OPENED":
		return "Opened dictionary"
	case "DICTIONARY_EXAMPLES_OPENED":
		return "Opened dictionary examples"
	case "DICTIONARY_DERIVATION_OPENED":
		return "Opened word derivation"
	case "DICTIONARY_PICTURE_OPENED":
		return "Opened dictionary picture"
	case "WORD_LOOKUP":
		if word, ok := values["word"].(string); ok && word != "" {
			return "Looked up \"" + word + "\""
		}
		return "Dictionary lookup"
	case "AI_INTERACTION":
		return "Used AI help"
	case "AI_HELP_OPENED":
		return "Opened AI help"
	case "QUIZ_OPENED":
		return "Opened quiz"
	case "QUIZ_STARTED":
		return "Started quiz"
	case "QUIZ_COMPLETED":
		return "Completed quiz"
	case "AHA_MOMENTS_OPENED":
		return "Opened Ah Ah Moments"
	case "AHA_MOMENT_CREATED":
		return "Added Ah Ah Moment"
	case "CONSULTANT_HELP_OPENED":
		return "Opened consultant help"
	case "HELP_REQUEST":
		return "Asked a consultant for help"
	case "SCAN_OPENED":
		return "Opened Scan to Locate"
	case "SCAN_STARTED":
		return "Started Scan to Locate"
	case "SCAN_SUCCEEDED":
		return "Located reading position by scan"
	case "SCAN_FAILED":
		return "Scan to Locate did not find a position"
	case "SERVICE_SELECTED":
		if service, ok := values["service"].(string); ok && service != "" {
			return "Selected " + service
		}
		return "Selected a reader service"
	case "READING_CHECK_IN":
		if reaction, ok := values["reaction"].(string); ok {
			switch reaction {
			case "sad":
				return "Reading check-in: hard"
			case "bored":
				return "Reading check-in: slow"
			case "happy":
				return "Reading check-in: good"
			}
		}
		return "Reading check-in"
	case "LOGIN":
		return "Logged in"
	case "LOGOUT":
		return "Logged out"
	default:
		return "Activity"
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
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05.999999999Z07:00", "2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported timestamp %q", value)
}
