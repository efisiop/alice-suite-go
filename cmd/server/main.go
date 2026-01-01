package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/handlers"
	"github.com/efisiopittau/alice-suite-go/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()
	cfg.Validate()

	// Initialize database
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	log.Println("Database initialized successfully")

	// Clean up stale sessions on startup (fresh start)
	log.Println("🧹 Cleaning up stale sessions on startup...")
	if err := database.CleanupExpiredSessions(); err != nil {
		log.Printf("Warning: Failed to cleanup expired sessions: %v", err)
	}
	if err := database.CleanupStaleSessions(); err != nil {
		log.Printf("Warning: Failed to cleanup stale sessions: %v", err)
	}

	// Start periodic session cleanup (every 5 minutes)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("🧹 Running periodic session cleanup...")
			database.CleanupExpiredSessions()
			database.CleanupStaleSessions()
		}
	}()

	// Setup routes
	mux := http.NewServeMux()

	// Static files (CSS, JS, images)
	staticDir := filepath.Join("internal", "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	log.Println("Static files directory configured:", staticDir)

	// Health check endpoint (no rate limiting)
	mux.HandleFunc("/health", handlers.HealthCheck)

	// Setup API routes (these will include protected consultant endpoints)
	handlers.SetupAPIRoutes(mux)

	// Authentication routes
	handlers.SetupAuthRoutes(mux)

	// Reader app routes
	handlers.SetupReaderRoutes(mux)

	// Consultant app routes with authentication middleware
	// Note: Go's ServeMux matches exact paths, so we need to handle /consultant and /consultant/
	// Create a sub-router for consultant routes with authentication
	consultantHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the path after /consultant
		path := strings.TrimPrefix(r.URL.Path, "/consultant")

		// Check if path matches /readers/:id pattern before other cases
		if strings.HasPrefix(path, "/readers/") {
			handlers.HandleConsultantReaderInspector(w, r)
			return
		}
		
		switch path {
		case "", "/":
			handlers.HandleConsultantDashboard(w, r)
		case "/help-requests":
			handlers.HandleConsultantHelpRequests(w, r)
		case "/feedback":
			handlers.HandleConsultantFeedback(w, r)
		case "/readers":
			handlers.HandleConsultantReaders(w, r)
		case "/reports":
			handlers.HandleConsultantReports(w, r)
		case "/reading-insights":
			handlers.HandleConsultantReadingInsights(w, r)
		case "/assign-readers":
			handlers.HandleConsultantAssignReaders(w, r)
		case "/send-prompt":
			handlers.HandleConsultantSendPrompt(w, r)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	// Wrap consultant handler with authentication middleware
	consultantAuthHandler := middleware.RequireConsultant(consultantHandler)
	mux.Handle("/consultant/", consultantAuthHandler)

	// Also handle /consultant without trailing slash
	mux.Handle("/consultant", consultantAuthHandler)

	// Leave login page without authentication (public access)
	mux.HandleFunc("/consultant/login", handlers.HandleConsultantLogin)

	// Wrap entire mux with heartbeat middleware (updates last_active_at on every request)
	// Then wrap with rate limiting middleware
	handler := middleware.RateLimit(middleware.HeartbeatMiddleware(mux))

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Health check: http://localhost:%s/health", cfg.Port)
	log.Printf("Reader app: http://localhost:%s/reader", cfg.Port)
	log.Printf("Consultant dashboard: http://localhost:%s/consultant", cfg.Port)
	log.Fatal(http.ListenAndServe(""+":"+cfg.Port, handler))
}
