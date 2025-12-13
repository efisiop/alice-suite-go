package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrapped.statusCode,
			duration,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware validates authentication token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public routes
		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// TODO: Validate token
		// For now, just check if header exists
		next.ServeHTTP(w, r)
	})
}

// isPublicRoute checks if a route is public (doesn't require auth)
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/health",
		"/login",
		"/register",
		"/",
		"/static/",
		"/api/auth/login",
		"/api/auth/register",
		"/auth/v1/token",
		"/auth/v1/signup",
	}

	for _, route := range publicRoutes {
		if path == route || (len(route) > 0 && route[len(route)-1] == '/' && len(path) >= len(route) && path[:len(route)] == route) {
			return true
		}
	}

	return false
}

