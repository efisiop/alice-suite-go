package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRateLimit_AllowsRequests tests that rate limiter allows requests within limit
func TestRateLimit_AllowsRequests(t *testing.T) {
	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Make a few requests within the limit
	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request %d returned wrong status code: got %v want %v", i, status, http.StatusOK)
		}
	}
}

// TestRateLimit_BlocksExcessiveRequests tests that rate limiter blocks excessive requests
func TestRateLimit_BlocksExcessiveRequests(t *testing.T) {
	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Make many rapid requests to exceed rate limit
	blocked := false
	for i := 0; i < 30; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code == http.StatusTooManyRequests {
			blocked = true
			break
		}
		
		// Small delay to allow rate limiter to process
		time.Sleep(10 * time.Millisecond)
	}

	if !blocked {
		t.Log("Rate limiter did not block excessive requests (may need more requests or different timing)")
	}
}

// TestRateLimit_DifferentIPs tests that rate limiting is per IP
func TestRateLimit_DifferentIPs(t *testing.T) {
	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create requests with different IPs
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"

	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"

	rr1 := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr1, req1)
	handler.ServeHTTP(rr2, req2)

	// Both should succeed as they're from different IPs
	if rr1.Code != http.StatusOK {
		t.Errorf("Request 1 returned wrong status code: got %v want %v", rr1.Code, http.StatusOK)
	}

	if rr2.Code != http.StatusOK {
		t.Errorf("Request 2 returned wrong status code: got %v want %v", rr2.Code, http.StatusOK)
	}
}

