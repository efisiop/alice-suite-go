package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/handlers"
)

func main() {
	cfg := config.Load()
	// Initialize database (PostgreSQL when DATABASE_URL set, else SQLite)
	if err := database.InitDB(cfg.DBPath, cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":8080"
	}

	fmt.Println("✅ Database initialized")

	// Setup routes
	mux := http.NewServeMux()

	// Reader API routes (register BEFORE static files to avoid conflicts)
	mux.HandleFunc("/api/health", handlers.HealthCheck)
	mux.HandleFunc("/api/auth/login", handlers.Login)
	mux.HandleFunc("/api/auth/register", handlers.Register)
	mux.HandleFunc("/api/books", handlers.GetBooks)
	mux.HandleFunc("/api/chapters", handlers.GetChapters)
	mux.HandleFunc("/api/sections", handlers.GetSections)
	mux.HandleFunc("/api/pages", handlers.GetPage)
	mux.HandleFunc("/api/dictionary/lookup", handlers.LookupWord)
	mux.HandleFunc("/api/dictionary/section/", handlers.GetSectionGlossaryTerms)
	mux.HandleFunc("/api/ai/ask", handlers.AskAI)
	mux.HandleFunc("/api/help/request", handlers.CreateHelpRequest)
	mux.HandleFunc("/api/progress", handlers.GetProgress)
	mux.HandleFunc("/api/user", handlers.GetUserProfile)
	mux.HandleFunc("/api/interactions", handlers.GetInteractions)
	mux.HandleFunc("/api/track", handlers.TrackEvent)
	
	// Serve static files (register AFTER API routes)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	fmt.Printf("🚀 Alice Suite Reader API starting on port %s\n", port)
	fmt.Println("📖 Physical Book Companion App - First 3 Chapters Test Ground")
	
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

