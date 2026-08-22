package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
	_ "github.com/mattn/go-sqlite3"
)

func TestGenericUsersEndpointRequiresConsultantAndOmitsPasswordHash(t *testing.T) {
	testDB, err := sql.Open("sqlite3", ":memory:")
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
		email TEXT,
		password_hash TEXT,
		first_name TEXT,
		last_name TEXT,
		role TEXT,
		is_verified INTEGER,
		created_at TEXT,
		updated_at TEXT,
		admin_level TEXT
	)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = testDB.Exec(`INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"reader-1", "reader@example.test", "sensitive-hash", "Test", "Reader", "reader", 1)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	SetupAPIRoutes(mux)

	t.Run("anonymous", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/rest/v1/users?select=*", nil))
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusUnauthorized, recorder.Body.String())
		}
	})

	t.Run("reader", func(t *testing.T) {
		token, err := auth.GenerateJWT("reader-1", "reader@example.test", "reader")
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/rest/v1/users?select=*", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusForbidden, recorder.Body.String())
		}
	})

	t.Run("consultant", func(t *testing.T) {
		token, err := auth.GenerateJWT("consultant-1", "consultant@example.test", "consultant")
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/rest/v1/users?select=*", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusOK, recorder.Body.String())
		}
		if strings.Contains(recorder.Body.String(), "password_hash") || strings.Contains(recorder.Body.String(), "sensitive-hash") {
			t.Fatalf("response exposed password hash: %s", recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), "reader@example.test") {
			t.Fatalf("response omitted expected safe user data: %s", recorder.Body.String())
		}
	})
}

func TestGenericPersonalDataEndpointRestrictsReaderToOwnRows(t *testing.T) {
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testDB.Close() })

	originalDB, originalDriver := database.DB, database.DriverName
	database.DB, database.DriverName = testDB, "sqlite3"
	t.Cleanup(func() {
		database.DB, database.DriverName = originalDB, originalDriver
	})

	_, err = testDB.Exec(`CREATE TABLE activity_logs (id TEXT PRIMARY KEY, user_id TEXT, activity_type TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = testDB.Exec(`INSERT INTO activity_logs (id, user_id, activity_type) VALUES
		('own', 'reader-1', 'OWN_EVENT'),
		('other', 'reader-2', 'OTHER_EVENT')`)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	SetupAPIRoutes(mux)
	token, err := auth.GenerateJWT("reader-1", "reader@example.test", "reader")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("own rows", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/rest/v1/activity_logs?user_id=eq.reader-1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusOK, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), "OWN_EVENT") || strings.Contains(recorder.Body.String(), "OTHER_EVENT") {
			t.Fatalf("unexpected response rows: %s", recorder.Body.String())
		}
	})

	for name, path := range map[string]string{
		"other reader": "/rest/v1/activity_logs?user_id=eq.reader-2",
		"no owner filter": "/rest/v1/activity_logs",
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)
			if recorder.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want %d; body=%q", recorder.Code, http.StatusForbidden, recorder.Body.String())
			}
		})
	}
}
