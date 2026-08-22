package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func TestHandleSignUpRecoversMissingReaderPreferencesSchema(t *testing.T) {
	testDB, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testDB.Close() })

	originalDB, originalDriver := database.DB, database.DriverName
	database.DB, database.DriverName = testDB, "sqlite3"
	t.Cleanup(func() {
		database.DB, database.DriverName = originalDB, originalDriver
	})

	_, err = testDB.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		role TEXT NOT NULL DEFAULT 'reader',
		is_verified INTEGER NOT NULL DEFAULT 0,
		created_at TEXT,
		updated_at TEXT
	)`)
	if err != nil {
		t.Fatal(err)
	}

	body := `{"email":"schema-recovery@example.test","password":"LocalAudit123!","first_name":"Schema","last_name":"Recovery","preferred_language_code":"it"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	HandleSignUp(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("signup status = %d, want %d; body=%q", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	var language string
	if err := testDB.QueryRow(`SELECT preferred_language_code FROM reader_preferences WHERE user_id = (SELECT id FROM users WHERE email = ?)`, "schema-recovery@example.test").Scan(&language); err != nil {
		t.Fatalf("read persisted preference: %v", err)
	}
	if language != "it" {
		t.Fatalf("preferred language = %q, want %q", language, "it")
	}
}
