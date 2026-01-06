package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.Load()

	fmt.Println("🔍 Deployment Verification Script")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println("")

	// Initialize database
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	fmt.Printf("✅ Database connected: %s\n", cfg.DBPath)
	fmt.Println("")

	// Check tables
	fmt.Println("📊 Checking database schema...")
	checkTables()

	// Check data
	fmt.Println("")
	fmt.Println("📊 Checking data...")
	checkData()

	fmt.Println("")
	fmt.Println("🎉 Verification complete!")
}

func checkTables() {
	expectedTables := []string{
		"users", "books", "chapters", "sections", "pages",
		"alice_glossary", "reading_progress", "sessions",
		"activity_logs", "reader_states", "interactions",
		"help_requests", "verification_codes", "vocabulary_lookups",
		"glossary_section_links",
	}

	for _, table := range expectedTables {
		var count int
		err := database.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)).Scan(&count)
		if err != nil {
			fmt.Printf("   ❌ Error checking table %s: %v\n", table, err)
		} else if count > 0 {
			fmt.Printf("   ✅ Table '%s' exists\n", table)
		} else {
			fmt.Printf("   ❌ Table '%s' MISSING\n", table)
		}
	}
}

func checkData() {
	checks := []struct {
		name  string
		query string
		min   int
	}{
		{"Books", "SELECT COUNT(*) FROM books", 1},
		{"Chapters", "SELECT COUNT(*) FROM chapters", 3},
		{"Sections", "SELECT COUNT(*) FROM sections", 70},
		{"Pages", "SELECT COUNT(*) FROM pages", 1},
		{"Glossary Terms", "SELECT COUNT(*) FROM alice_glossary", 1800},
		{"Users (readers)", "SELECT COUNT(*) FROM users WHERE role = 'reader'", 1},
		{"Verification Codes", "SELECT COUNT(*) FROM verification_codes", 1},
	}

	for _, check := range checks {
		var count int
		err := database.DB.QueryRow(check.query).Scan(&count)
		if err != nil {
			fmt.Printf("   ❌ Error checking %s: %v\n", check.name, err)
		} else if count >= check.min {
			fmt.Printf("   ✅ %s: %d (expected ≥%d)\n", check.name, count, check.min)
		} else {
			fmt.Printf("   ⚠️  %s: %d (expected ≥%d) - MAY NEED FIXING\n", check.name, count, check.min)
		}
	}

	// Check specific page sections
	var page1Sections int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM sections WHERE page_number = 1").Scan(&page1Sections)
	if err != nil {
		fmt.Printf("   ❌ Error checking page 1 sections: %v\n", err)
	} else if page1Sections >= 5 {
		fmt.Printf("   ✅ Page 1 sections: %d (expected ≥5)\n", page1Sections)
	} else {
		fmt.Printf("   ⚠️  Page 1 sections: %d (expected ≥5) - RUN fix-render\n", page1Sections)
	}
}

