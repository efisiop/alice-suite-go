package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/efisiopittau/alice-suite-go/pkg/auth"
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
		t.Run(path+" rejects anonymous", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
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
	mux := http.NewServeMux()
	SetupReaderRoutes(mux)

	readerToken, err := auth.GenerateJWT("reader-1", "reader@example.test", "reader")
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
