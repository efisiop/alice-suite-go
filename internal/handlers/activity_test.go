package handlers

import (
	"database/sql"
	"testing"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
)

func setupActivityTestDB(t *testing.T) func() {
	t.Helper()

	previousDB := database.DB
	previousDriver := database.DriverName

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)

	database.DB = db
	database.DriverName = "sqlite3"

	schema := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			email TEXT,
			role TEXT NOT NULL
		)`,
		`CREATE TABLE books (
			id TEXT PRIMARY KEY
		)`,
		`CREATE TABLE interactions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			book_id TEXT NOT NULL,
			section_id TEXT,
			page_number INTEGER,
			content TEXT,
			context TEXT,
			created_at TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (book_id) REFERENCES books(id)
		)`,
		`CREATE TABLE activity_logs (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			session_id TEXT,
			activity_type TEXT NOT NULL,
			book_id TEXT,
			page_number INTEGER,
			section_id TEXT,
			metadata TEXT,
			created_at TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (book_id) REFERENCES books(id)
		)`,
		`CREATE TABLE reader_states (
			user_id TEXT PRIMARY KEY,
			book_id TEXT,
			current_page INTEGER,
			current_section_id TEXT,
			last_activity_type TEXT,
			last_activity_at TEXT,
			status TEXT,
			updated_at TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (book_id) REFERENCES books(id)
		)`,
		`INSERT INTO users (id, first_name, last_name, email, role)
		 VALUES ('reader-1', 'Alice', 'Reader', 'alice@example.com', 'reader')`,
		`INSERT INTO books (id) VALUES ('alice-in-wonderland')`,
	}

	for _, query := range schema {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("setup query failed: %v", err)
		}
	}

	return func() {
		db.Close()
		database.DB = previousDB
		database.DriverName = previousDriver
	}
}

func TestGetReaderJourneyGroupsVisitsAndDescribesDirection(t *testing.T) {
	cleanup := setupActivityTestDB(t)
	defer cleanup()

	insert := func(page int, at string) {
		t.Helper()
		if _, err := database.DB.Exec(`INSERT INTO activity_logs (id, user_id, activity_type, book_id, page_number, created_at)
			VALUES (?, 'reader-1', 'PAGE_VIEW', 'alice-in-wonderland', ?, ?)`, at, page, at); err != nil {
			t.Fatalf("insert page view: %v", err)
		}
	}
	insert(10, "2026-09-15 09:00:00+02:00")
	insert(14, "2026-09-15 09:15:00+02:00")
	insert(12, "2026-09-15 11:00:00+02:00")
	insert(11, "2026-09-15 11:10:00+02:00")
	if _, err := database.DB.Exec(`INSERT INTO activity_logs (id, user_id, activity_type, book_id, metadata, created_at)
		VALUES ('lookup-1', 'reader-1', 'WORD_LOOKUP', 'alice-in-wonderland', '{"word":"rabbit"}', '2026-09-15 11:15:00+02:00')`); err != nil {
		t.Fatalf("insert word lookup: %v", err)
	}

	journey, err := database.GetReaderJourney("reader-1", 5)
	if err != nil {
		t.Fatalf("GetReaderJourney returned error: %v", err)
	}
	if journey.Position == nil || journey.Position.PageNumber == nil || *journey.Position.PageNumber != 11 {
		t.Fatalf("latest position = %#v, want page 11", journey.Position)
	}
	if len(journey.Visits) != 2 {
		t.Fatalf("visit count = %d, want 2", len(journey.Visits))
	}
	if got := journey.Visits[0]; got.StartPage != 12 || got.EndPage != 11 || got.Direction != "revisiting" {
		t.Fatalf("latest visit = %#v, want pages 12 to 11 revisiting", got)
	}
	if got := journey.Visits[1]; got.StartPage != 10 || got.EndPage != 14 || got.Direction != "forward" {
		t.Fatalf("older visit = %#v, want pages 10 to 14 forward", got)
	}
	if got := journey.Visits[0].EndedAt.Sub(journey.Visits[0].StartedAt); got != 15*time.Minute {
		t.Fatalf("latest visit duration = %s, want 15m", got)
	}
	latestPage := journey.Visits[0].Pages[1]
	if len(latestPage.Events) != 1 || latestPage.Events[0].Detail != `Looked up "rabbit"` {
		t.Fatalf("page events = %#v, want rabbit lookup", latestPage.Events)
	}
}

func TestGetReaderJourneyDescribesReaderToolActions(t *testing.T) {
	cleanup := setupActivityTestDB(t)
	defer cleanup()

	entries := []struct {
		id, kind, metadata, createdAt string
	}{
		{"page", "PAGE_VIEW", "{}", "2026-09-15 09:00:00+02:00"},
		{"dictionary", "DICTIONARY_OPENED", "{}", "2026-09-15 09:01:00+02:00"},
		{"quiz", "QUIZ_STARTED", `{"scope":"section"}`, "2026-09-15 09:02:00+02:00"},
		{"scan", "SCAN_SUCCEEDED", "{}", "2026-09-15 09:03:00+02:00"},
	}
	for _, entry := range entries {
		pageNumber := interface{}(nil)
		if entry.kind == "PAGE_VIEW" {
			pageNumber = 12
		}
		if _, err := database.DB.Exec(`INSERT INTO activity_logs (id, user_id, activity_type, book_id, page_number, metadata, created_at)
			VALUES (?, 'reader-1', ?, 'alice-in-wonderland', ?, ?, ?)`, entry.id, entry.kind, pageNumber, entry.metadata, entry.createdAt); err != nil {
			t.Fatalf("insert %s: %v", entry.kind, err)
		}
	}

	journey, err := database.GetReaderJourney("reader-1", 5)
	if err != nil {
		t.Fatalf("GetReaderJourney returned error: %v", err)
	}
	events := journey.Visits[0].Pages[0].Events
	want := []string{"Opened dictionary", "Started quiz", "Located reading position by scan"}
	if len(events) != len(want) {
		t.Fatalf("event count = %d, want %d: %#v", len(events), len(want), events)
	}
	for i, detail := range want {
		if events[i].Detail != detail {
			t.Errorf("event %d detail = %q, want %q", i, events[i].Detail, detail)
		}
	}
}

func TestTrackActivityWritesConsultantDashboardPath(t *testing.T) {
	cleanup := setupActivityTestDB(t)
	defer cleanup()

	err := TrackActivity("reader-1", "PAGE_SYNC", "alice-in-wonderland", map[string]interface{}{
		"page_number":   float64(12),
		"section_id":    "section-3",
		"section_index": float64(2),
	})
	if err != nil {
		t.Fatalf("TrackActivity returned error: %v", err)
	}

	var interactionType string
	if err := database.DB.QueryRow(`SELECT event_type FROM interactions WHERE user_id = 'reader-1'`).Scan(&interactionType); err != nil {
		t.Fatalf("interaction row missing: %v", err)
	}
	if interactionType != "PAGE_SYNC" {
		t.Fatalf("interaction event_type = %q, want PAGE_SYNC", interactionType)
	}

	var activityType string
	var pageNumber int
	if err := database.DB.QueryRow(`SELECT activity_type, page_number FROM activity_logs WHERE user_id = 'reader-1'`).Scan(&activityType, &pageNumber); err != nil {
		t.Fatalf("activity_logs row missing: %v", err)
	}
	if activityType != "PAGE_VIEW" {
		t.Fatalf("activity_logs activity_type = %q, want PAGE_VIEW", activityType)
	}
	if pageNumber != 12 {
		t.Fatalf("activity_logs page_number = %d, want 12", pageNumber)
	}

	var stateType string
	if err := database.DB.QueryRow(`SELECT last_activity_type FROM reader_states WHERE user_id = 'reader-1'`).Scan(&stateType); err != nil {
		t.Fatalf("reader_states row missing: %v", err)
	}
	if stateType != "PAGE_VIEW" {
		t.Fatalf("reader_states last_activity_type = %q, want PAGE_VIEW", stateType)
	}
}

func TestTrackActivityWithoutBookStillUpdatesConsultantPath(t *testing.T) {
	cleanup := setupActivityTestDB(t)
	defer cleanup()

	if err := TrackActivity("reader-1", "LOGIN", "", nil); err != nil {
		t.Fatalf("TrackActivity returned error: %v", err)
	}

	var interactionsCount int
	if err := database.DB.QueryRow(`SELECT COUNT(*) FROM interactions`).Scan(&interactionsCount); err != nil {
		t.Fatalf("count interactions: %v", err)
	}
	if interactionsCount != 0 {
		t.Fatalf("interactions count = %d, want 0 for bookless login", interactionsCount)
	}

	var activityType string
	if err := database.DB.QueryRow(`SELECT activity_type FROM activity_logs WHERE user_id = 'reader-1'`).Scan(&activityType); err != nil {
		t.Fatalf("activity_logs row missing: %v", err)
	}
	if activityType != "LOGIN" {
		t.Fatalf("activity_logs activity_type = %q, want LOGIN", activityType)
	}

	var stateType string
	if err := database.DB.QueryRow(`SELECT last_activity_type FROM reader_states WHERE user_id = 'reader-1'`).Scan(&stateType); err != nil {
		t.Fatalf("reader_states row missing: %v", err)
	}
	if stateType != "LOGIN" {
		t.Fatalf("reader_states last_activity_type = %q, want LOGIN", stateType)
	}
}
