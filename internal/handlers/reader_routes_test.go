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

func TestReaderPagesRequireReaderRole(t *testing.T) {
	mux := http.NewServeMux()
	SetupReaderRoutes(mux)

	protectedPaths := []string{
		"/reader",
		"/reader/interaction",
		"/reader/my-page",
		"/reader/book/alice-in-wonderland",
		"/reader/statistics",
	}

	consultantToken, err := auth.GenerateJWT("consultant-1", "consultant@example.test", "consultant")
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range protectedPaths {
		t.Run(path+" redirects anonymous to login", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
			if recorder.Code != http.StatusFound {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusFound)
			}
			if location := recorder.Header().Get("Location"); location != "/reader/login" {
				t.Fatalf("location = %q, want %q", location, "/reader/login")
			}
		})

		t.Run(path+" rejects consultant", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", "Bearer "+consultantToken)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)
			if recorder.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
			}
		})
	}
}

func TestReaderOnboardingExplainsPhysicalBookAndUsesOneSessionStore(t *testing.T) {
	t.Chdir("../..")
	withReaderVerificationDB(t)
	mux := http.NewServeMux()
	SetupReaderRoutes(mux)

	landing := httptest.NewRecorder()
	mux.ServeHTTP(landing, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(strings.ToLower(landing.Body.String()), "physical copy") {
		t.Fatalf("landing page does not clearly explain that a physical book is required")
	}

	anonymousVerify := httptest.NewRecorder()
	mux.ServeHTTP(anonymousVerify, httptest.NewRequest(http.MethodGet, "/verify", nil))
	if anonymousVerify.Code != http.StatusFound || anonymousVerify.Header().Get("Location") != "/reader/login" {
		t.Fatalf("anonymous verification status/location = %d %q, want 302 /reader/login", anonymousVerify.Code, anonymousVerify.Header().Get("Location"))
	}

	token, err := auth.GenerateJWT("reader-unverified", "unverified@example.test", "reader")
	if err != nil {
		t.Fatal(err)
	}
	verify := httptest.NewRecorder()
	verifyRequest := httptest.NewRequest(http.MethodGet, "/verify", nil)
	verifyRequest.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(verify, verifyRequest)
	if verify.Code != http.StatusOK {
		t.Fatalf("verification page status = %d, want %d", verify.Code, http.StatusOK)
	}
	if !strings.Contains(verify.Body.String(), "sessionStorage.getItem('auth_token')") {
		t.Fatal("verification page does not use the reader login session")
	}
	if strings.Contains(verify.Body.String(), "localStorage.getItem('auth_token')") {
		t.Fatal("verification page still reads the obsolete login store")
	}
	if strings.Contains(verify.Body.String(), "fetch('/auth/v1/user'") || strings.Contains(verify.Body.String(), "user_id: user.id") {
		t.Fatal("verification page still asks the browser to identify the reader")
	}
}

func TestReaderAuthPagesRemainPublic(t *testing.T) {
	t.Chdir("../..")
	mux := http.NewServeMux()
	SetupReaderRoutes(mux)

	for _, path := range []string{"/reader/login", "/reader/register"} {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", path, recorder.Code, http.StatusOK)
		}
	}
}

func TestReaderRoleCanAccessPages(t *testing.T) {
	t.Chdir("../..")
	withReaderVerificationDB(t)
	mux := http.NewServeMux()
	SetupReaderRoutes(mux)

	readerToken, err := auth.GenerateJWT("reader-verified", "verified@example.test", "reader")
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/reader", "/reader/interaction", "/reader/my-page", "/reader/statistics"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer "+readerToken)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", path, recorder.Code, http.StatusOK)
		}
	}
}

func TestUnverifiedReaderMustVerifyBeforeReading(t *testing.T) {
	withReaderVerificationDB(t)
	mux := http.NewServeMux()
	SetupReaderRoutes(mux)

	token, err := auth.GenerateJWT("reader-unverified", "unverified@example.test", "reader")
	if err != nil {
		t.Fatal(err)
	}

	pageRequest := httptest.NewRequest(http.MethodGet, "/reader", nil)
	pageRequest.Header.Set("Authorization", "Bearer "+token)
	pageRecorder := httptest.NewRecorder()
	mux.ServeHTTP(pageRecorder, pageRequest)
	if pageRecorder.Code != http.StatusFound || pageRecorder.Header().Get("Location") != "/verify" {
		t.Fatalf("reader page status/location = %d %q, want 302 /verify", pageRecorder.Code, pageRecorder.Header().Get("Location"))
	}

	apiMux := http.NewServeMux()
	SetupAPIRoutes(apiMux)
	apiRequest := httptest.NewRequest(http.MethodGet, "/api/reader/preferences", nil)
	apiRequest.Header.Set("Authorization", "Bearer "+token)
	apiRecorder := httptest.NewRecorder()
	apiMux.ServeHTTP(apiRecorder, apiRequest)
	if apiRecorder.Code != http.StatusForbidden {
		t.Fatalf("reader API status = %d, want %d", apiRecorder.Code, http.StatusForbidden)
	}
}

func TestReaderAPIsRejectConsultantRole(t *testing.T) {
	mux := http.NewServeMux()
	SetupAPIRoutes(mux)

	consultantToken, err := auth.GenerateJWT("consultant-1", "consultant@example.test", "consultant")
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		"/api/reader/prompts",
		"/api/reader/quiz",
		"/api/reader/ah-ah-moments",
		"/api/reader/preferences",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer "+consultantToken)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("%s status = %d, want %d", path, recorder.Code, http.StatusForbidden)
		}
	}
}

func withReaderVerificationDB(t *testing.T) {
	t.Helper()
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
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		role TEXT NOT NULL,
		is_verified INTEGER NOT NULL,
		created_at TEXT,
		updated_at TEXT
	);
	INSERT INTO users VALUES
		('reader-verified', 'verified@example.test', '', 'Verified', 'Reader', 'reader', 1, '2026-08-22 00:00:00', '2026-08-22 00:00:00'),
		('reader-unverified', 'unverified@example.test', '', 'Unverified', 'Reader', 'reader', 0, '2026-08-22 00:00:00', '2026-08-22 00:00:00')`)
	if err != nil {
		t.Fatal(err)
	}
}
