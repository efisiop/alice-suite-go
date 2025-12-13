package handlers

import (
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// SetupConsultantRoutes sets up routes for the Consultant Dashboard app
func SetupConsultantRoutes(mux *http.ServeMux) {
	// Consultant authentication
	mux.HandleFunc("/consultant/login", HandleConsultantLogin)

	// Consultant dashboard pages
	mux.HandleFunc("/consultant", HandleConsultantDashboard)
	mux.HandleFunc("/consultant/send-prompt", HandleConsultantSendPrompt)
	mux.HandleFunc("/consultant/help-requests", HandleConsultantHelpRequests)
	mux.HandleFunc("/consultant/feedback", HandleConsultantFeedback)
	mux.HandleFunc("/consultant/readers", HandleConsultantReaders)
	mux.HandleFunc("/consultant/reports", HandleConsultantReports)
	mux.HandleFunc("/consultant/reading-insights", HandleConsultantReadingInsights)
	mux.HandleFunc("/consultant/assign-readers", HandleConsultantAssignReaders)
}

// HandleConsultantLogin handles GET/POST /consultant/login
func HandleConsultantLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			filepath.Join("internal", "templates", "base.html"),
			filepath.Join("internal", "templates", "consultant", "login.html"),
		)
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
		return
	}

	// POST handled by auth handler (with consultant role check)
	HandleLogin(w, r)
}

// HandleConsultantDashboard handles GET /consultant
func HandleConsultantDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header or cookie
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Check for token in cookie as fallback
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie == nil || cookie.Value == "" {
			// No valid token, redirect to login
			http.Redirect(w, r, "/consultant/login", http.StatusFound)
			return
		}
		// Safari may URL-encode cookie values, so decode if needed
		tokenValue := cookie.Value
		// Try URL decoding (Safari sometimes encodes cookies)
		if decoded, err := url.QueryUnescape(tokenValue); err == nil && decoded != tokenValue {
			tokenValue = decoded
		}
		authHeader = "Bearer " + tokenValue
	}

	// Extract and validate token
	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		// Invalid token format, redirect to login
		http.Redirect(w, r, "/consultant/login", http.StatusFound)
		return
	}

	// Validate token and get claims
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		// Invalid or expired token, redirect to login
		http.Redirect(w, r, "/consultant/login", http.StatusFound)
		return
	}

	// Check if user has consultant role
	if claims.Role != "consultant" {
		// Not a consultant, show forbidden message or redirect
		http.Error(w, "Access denied: Consultant privileges required", http.StatusForbidden)
		return
	}

	// Token is valid and user is a consultant, serve the dashboard
	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "consultant", "dashboard.html"),
	)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleConsultantSendPrompt handles GET /consultant/send-prompt
func HandleConsultantSendPrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Add server-side authentication check
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Check for token in cookie as fallback
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			// No valid token, redirect to login
			http.Redirect(w, r, "/consultant/login", http.StatusFound)
			return
		}
		authHeader = "Bearer " + cookie.Value
	}

	// Extract and validate token
	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Redirect(w, r, "/consultant/login", http.StatusFound)
		return
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		http.Redirect(w, r, "/consultant/login", http.StatusFound)
		return
	}

	// Check if user has consultant role
	if claims.Role != "consultant" {
		http.Error(w, "Access denied: Consultant privileges required", http.StatusForbidden)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "consultant", "send-prompt.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleConsultantHelpRequests handles GET /consultant/help-requests
func HandleConsultantHelpRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "consultant", "help-requests.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleConsultantFeedback handles GET /consultant/feedback
func HandleConsultantFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "consultant", "feedback.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleConsultantReaders handles GET /consultant/readers
func HandleConsultantReaders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "consultant", "readers.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleConsultantReports handles GET /consultant/reports
func HandleConsultantReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "consultant", "reports.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleConsultantReadingInsights handles GET /consultant/reading-insights
func HandleConsultantReadingInsights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "consultant", "reading-insights.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleConsultantAssignReaders handles GET /consultant/assign-readers
func HandleConsultantAssignReaders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "consultant", "assign-readers.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

