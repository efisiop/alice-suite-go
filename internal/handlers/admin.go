package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/email"
	"github.com/efisiopittau/alice-suite-go/internal/middleware"
	"github.com/efisiopittau/alice-suite-go/internal/useragent"
)

// SetupAdminRoutes sets up routes for the Administrator dashboard
func SetupAdminRoutes(mux *http.ServeMux) {
	// Admin login (public)
	mux.HandleFunc("/admin/login", HandleAdminLogin)

	// Admin API: presence (protected)
	mux.Handle("/api/admin/presence", middleware.RequireAdmin(http.HandlerFunc(HandleAdminPresence)))
	// Admin API: settings (protected)
	mux.Handle("/api/admin/settings", middleware.RequireAdmin(http.HandlerFunc(HandleAdminSettings)))
	// Admin API: email status (for debugging why login emails might not send)
	mux.Handle("/api/admin/email-status", middleware.RequireAdmin(http.HandlerFunc(HandleAdminEmailStatus)))

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
	consultantsCount, err := database.CountConsultantsOnline()
	if err != nil {
		log.Printf("admin presence: count consultants: %v", err)
		http.Error(w, "Failed to get consultants count", http.StatusInternalServerError)
		return
	}
	readerSessions, err := database.GetOnlineReaderSessions()
	if err != nil {
		log.Printf("admin presence: get reader sessions: %v", err)
		http.Error(w, "Failed to get reader sessions", http.StatusInternalServerError)
		return
	}
	consultantSessions, err := database.GetOnlineConsultantSessions()
	if err != nil {
		log.Printf("admin presence: get consultant sessions: %v", err)
		http.Error(w, "Failed to get consultant sessions", http.StatusInternalServerError)
		return
	}
	readersList := make([]map[string]interface{}, 0, len(readerSessions))
	for _, s := range readerSessions {
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
			"user_id":   s.UserID,
			"email":    s.Email,
			"name":     name,
			"device":   useragent.DeviceLabel(s.UserAgent),
			"last_seen": s.LastActiveAt.Format("15:04"),
		})
	}
	consultantsList := make([]map[string]interface{}, 0, len(consultantSessions))
	for _, s := range consultantSessions {
		name := s.FirstName + " " + s.LastName
		if name == " " {
			name = s.Email
		} else {
			name = strings.TrimSpace(name)
			if name == "" {
				name = s.Email
			}
		}
		consultantsList = append(consultantsList, map[string]interface{}{
			"user_id":   s.UserID,
			"email":    s.Email,
			"name":     name,
			"device":   useragent.DeviceLabel(s.UserAgent),
			"last_seen": s.LastActiveAt.Format("15:04"),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"readers_online":     readersCount,
		"consultants_online": consultantsCount,
		"readers":            readersList,
		"consultants":       consultantsList,
	})
}

// HandleAdminSettings returns or updates admin settings (e.g. login email notifications). Admin-only.
func HandleAdminSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		enabled, err := database.GetLoginEmailNotificationsEnabled()
		if err != nil {
			enabled = false
		}
		recipient, err := database.GetAdminSetting("login_email_recipient")
		if err != nil || recipient == "" {
			recipient = "efisio@mylivemail.net"
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"login_email_notifications": enabled,
			"login_email_recipient":     recipient,
		})
		return
	case http.MethodPut, http.MethodPatch:
		var req struct {
			LoginEmailNotifications *bool   `json:"login_email_notifications"`
			LoginEmailRecipient     *string `json:"login_email_recipient"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.LoginEmailNotifications != nil {
			if err := database.SetLoginEmailNotificationsEnabled(*req.LoginEmailNotifications); err != nil {
				log.Printf("admin settings: set login_email_notifications: %v", err)
				http.Error(w, "Failed to update setting", http.StatusInternalServerError)
				return
			}
		}
		if req.LoginEmailRecipient != nil && strings.TrimSpace(*req.LoginEmailRecipient) != "" {
			if err := database.SetAdminSetting("login_email_recipient", strings.TrimSpace(*req.LoginEmailRecipient)); err != nil {
				log.Printf("admin settings: set login_email_recipient: %v", err)
				http.Error(w, "Failed to update recipient", http.StatusInternalServerError)
				return
			}
		}
		// Return current state
		enabled, _ := database.GetLoginEmailNotificationsEnabled()
		recipient, _ := database.GetAdminSetting("login_email_recipient")
		if recipient == "" {
			recipient = "efisio@mylivemail.net"
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"login_email_notifications": enabled,
			"login_email_recipient":     recipient,
		})
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// HandleAdminEmailStatus returns why login emails may or may not be sent (admin-only, for debugging).
func HandleAdminEmailStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	enabled, _ := database.GetLoginEmailNotificationsEnabled()
	recipient, _ := database.GetAdminSetting("login_email_recipient")
	if recipient == "" {
		recipient = "efisio@mylivemail.net"
	}
	onRender := os.Getenv("RENDER") == "true"
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpConfigured := smtpHost != "" && smtpUser != "" && os.Getenv("SMTP_PASSWORD") != ""
	willSend := enabled && email.ShouldSendLoginEmails()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"login_email_notifications_enabled": enabled,
		"login_email_recipient":              recipient,
		"on_render":                         onRender,
		"smtp_configured":                   smtpConfigured,
		"smtp_host_set":                     smtpHost != "",
		"smtp_user_set":                     smtpUser != "",
		"will_send_on_next_login":           willSend,
		"hint":                              emailStatusHint(enabled, onRender, smtpConfigured, willSend),
	})
}

func emailStatusHint(enabled, onRender, smtpConfigured, willSend bool) string {
	if willSend {
		return "Emails will be sent on reader/consultant login."
	}
	if !enabled {
		return "Turn on the toggle in Email notifications to enable."
	}
	if !onRender {
		return "Not running on Render (RENDER not true). Emails only send on Render."
	}
	if !smtpConfigured {
		return "Set SMTP_HOST, SMTP_USER, and SMTP_PASSWORD in Render Environment."
	}
	return "Check Render logs for send errors."
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
