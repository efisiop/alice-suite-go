package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/middleware"
	"github.com/efisiopittau/alice-suite-go/internal/useragent"
)

// SetupAdminRoutes sets up routes for the Administrator dashboard
func SetupAdminRoutes(mux *http.ServeMux) {
	// Admin login (public)
	mux.HandleFunc("/admin/login", HandleAdminLogin)

	// Admin API: presence (protected)
	mux.Handle("/api/admin/presence", middleware.RequireAdmin(http.HandlerFunc(HandleAdminPresence)))

	// Admin dashboard and pages (protected)
	mux.Handle("/admin", middleware.RequireAdmin(http.HandlerFunc(HandleAdminDashboard)))
	mux.Handle("/admin/", middleware.RequireAdmin(http.HandlerFunc(HandleAdminDashboard)))
}

// HandleAdminPresence returns readers online, consultants online, and a list of online readers with device (admin-only).
func HandleAdminPresence(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	readersCount, err := database.CountReadersOnline()
	if err != nil {
		log.Printf("admin presence: count readers: %v", err)
		http.Error(w, "Failed to get readers count", http.StatusInternalServerError)
		return
	}
	consultants, err := database.CountConsultantsOnline()
	if err != nil {
		log.Printf("admin presence: count consultants: %v", err)
		http.Error(w, "Failed to get consultants count", http.StatusInternalServerError)
		return
	}
	sessions, err := database.GetOnlineReaderSessions()
	if err != nil {
		log.Printf("admin presence: get reader sessions: %v", err)
		http.Error(w, "Failed to get reader sessions", http.StatusInternalServerError)
		return
	}
	readersList := make([]map[string]interface{}, 0, len(sessions))
	for _, s := range sessions {
		name := s.FirstName + " " + s.LastName
		if name == " " {
			name = s.Email
		} else {
			name = strings.TrimSpace(name)
			if name == "" {
				name = s.Email
			}
		}
		readersList = append(readersList, map[string]interface{}{
			"user_id":    s.UserID,
			"email":     s.Email,
			"name":      name,
			"device":    useragent.DeviceLabel(s.UserAgent),
			"last_seen": s.LastActiveAt.Format("15:04"),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"readers_online":     readersCount,
		"consultants_online": consultants,
		"readers":            readersList,
	})
}

// HandleAdminLogin handles GET/POST /admin/login
func HandleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			filepath.Join("internal", "templates", "base.html"),
			filepath.Join("internal", "templates", "admin", "login.html"),
		)
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
		return
	}
	// POST: use same auth as consultant/reader
	HandleLogin(w, r)
}

// HandleAdminDashboard serves the Administrator dashboard (overview of readers, consultants, and operations)
func HandleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Only /admin and /admin/ serve the dashboard; other paths could be added later
	path := r.URL.Path
	if path != "/admin" && path != "/admin/" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "base.html"),
		filepath.Join("internal", "templates", "admin", "dashboard.html"),
	)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}
