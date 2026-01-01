package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// SetupReaderRoutes sets up routes for the Reader app
func SetupReaderRoutes(mux *http.ServeMux) {
	// Reader app pages
	mux.HandleFunc("/reader", HandleReaderDashboard)
	mux.HandleFunc("/reader/interaction", HandleReaderInteraction)
	mux.HandleFunc("/reader/my-page", HandleReaderMyPage)
	mux.HandleFunc("/reader/book/", HandleReaderBook)
	mux.HandleFunc("/reader/statistics", HandleReaderStatistics)

	// Reader authentication pages
	mux.HandleFunc("/reader/login", HandleReaderLogin)
	mux.HandleFunc("/reader/register", HandleReaderRegister)
	// Redirect old /login and /register to /reader/login and /reader/register for consistency
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/reader/login", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/reader/register", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/verify", HandleReaderVerify)
	mux.HandleFunc("/welcome", HandleReaderWelcome)
	mux.HandleFunc("/forgot-password", HandleReaderForgotPassword)

	// Public landing page
	mux.HandleFunc("/", HandleReaderLanding)
}

// HandleReaderLanding handles GET / (landing page)
func HandleReaderLanding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "reader", "landing.html"),
	)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleReaderLogin handles GET/POST /login
// This is a PUBLIC endpoint - no authentication required
func HandleReaderLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Add headers to help with proxy issues
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		
		tmpl, err := template.ParseFiles(
			filepath.Join("internal", "templates", "base.html"),
			filepath.Join("internal", "templates", "reader", "login.html"),
		)
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// POST handled by auth handler
	HandleLogin(w, r)
}

// HandleReaderRegister handles GET/POST /register
func HandleReaderRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			filepath.Join("internal", "templates", "base.html"),
			filepath.Join("internal", "templates", "reader", "register.html"),
		)
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
		return
	}

	// POST handled by auth handler
	HandleSignUp(w, r)
}

// HandleReaderVerify handles GET/POST /verify
func HandleReaderVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			filepath.Join("internal", "templates", "base.html"),
			filepath.Join("internal", "templates", "reader", "verify.html"),
		)
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
		return
	}

	// POST handled by verification handler
	HandleVerifyBookCode(w, r)
}

// HandleReaderWelcome handles GET /welcome
func HandleReaderWelcome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "reader", "welcome.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleReaderForgotPassword handles GET/POST /forgot-password
func HandleReaderForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "reader", "forgot-password.html"))
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
		return
	}

	// POST handler - to be implemented
	http.Error(w, "Password reset - to be implemented", http.StatusNotImplemented)
}

// HandleReaderDashboard handles GET /reader
// REVERTED: Removed server-side auth check - let JavaScript handle authentication
// This restores the previous behavior that was working
func HandleReaderDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve the dashboard - authentication is handled client-side by JavaScript
	// This matches the previous working behavior
	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "reader", "dashboard.html"),
	)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleReaderInteraction handles GET /reader/interaction
func HandleReaderInteraction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "reader", "interaction.html"),
	)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleReaderBook handles GET /reader/book/:bookId
func HandleReaderBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract bookId from URL path
	bookID := r.URL.Path[len("/reader/book/"):]

	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "reader", "book.html"))
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"BookID": bookID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// HandleReaderStatistics handles GET /reader/statistics
func HandleReaderStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "reader", "statistics.html"),
	)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// HandleReaderMyPage handles GET /reader/my-page
func HandleReaderMyPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "reader", "my-page.html"),
	)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

