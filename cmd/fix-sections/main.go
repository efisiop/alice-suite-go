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
	// Load configuration
	cfg := config.Load()

	// Initialize database (PostgreSQL when DATABASE_URL set, else SQLite)
	if err := database.InitDB(cfg.DBPath, cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	fmt.Println("🔍 Diagnosing sections table structure...")
	fmt.Println("=" + string(make([]byte, 60)))

	// Check current sections table structure
	var tableSQL string
	err := database.DB.QueryRow(`
		SELECT sql FROM sqlite_master 
		WHERE type='table' AND name='sections'
	`).Scan(&tableSQL)

	if err == sql.ErrNoRows {
		fmt.Println("❌ ERROR: sections table does not exist!")
		os.Exit(1)
	} else if err != nil {
		log.Fatalf("Error checking table: %v", err)
	}

	fmt.Println("\n📋 Current sections table structure:")
	fmt.Println(tableSQL)

	// Check if it's old structure (has chapter_id, start_page, end_page)
	isOldStructure := false
	if contains(tableSQL, "chapter_id") || contains(tableSQL, "start_page") {
		isOldStructure = true
		fmt.Println("\n⚠️  Detected OLD structure (chapter-based)")
	} else if contains(tableSQL, "page_id") && contains(tableSQL, "page_number") && contains(tableSQL, "section_number") {
		fmt.Println("\n✅ Detected NEW structure (page-based)")
	} else {
		fmt.Println("\n❓ Unknown structure")
	}

	// Check sections_new table
	var sectionsNewSQL string
	err = database.DB.QueryRow(`
		SELECT sql FROM sqlite_master 
		WHERE type='table' AND name='sections_new'
	`).Scan(&sectionsNewSQL)

	hasSectionsNew := err == nil
	if hasSectionsNew {
		fmt.Println("\n📋 sections_new table exists:")
		fmt.Println(sectionsNewSQL)
	}

	// Count sections in current table
	var sectionCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections").Scan(&sectionCount)
	fmt.Printf("\n📊 Sections in 'sections' table: %d\n", sectionCount)

	// Check sections for page 1
	var page1Count int
	if !isOldStructure {
		database.DB.QueryRow("SELECT COUNT(*) FROM sections WHERE page_number = 1").Scan(&page1Count)
		fmt.Printf("📄 Sections for page 1: %d\n", page1Count)
	} else {
		fmt.Println("📄 Cannot check page 1 (old structure)")
	}

	// Count sections_new if it exists
	if hasSectionsNew {
		var sectionsNewCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM sections_new").Scan(&sectionsNewCount)
		fmt.Printf("📊 Sections in 'sections_new' table: %d\n", sectionsNewCount)
	}

	fmt.Println("\n" + string(make([]byte, 61)))

	// Determine what needs to be done
	if isOldStructure && hasSectionsNew {
		fmt.Println("\n🔧 FIX NEEDED: Old structure detected, sections_new exists but may be empty")
		fmt.Println("   Action: Need to migrate data from old to new structure")
		fmt.Println("   OR: Drop old sections, rename sections_new to sections, and populate data")
	} else if isOldStructure && !hasSectionsNew {
		fmt.Println("\n🔧 FIX NEEDED: Old structure detected, no sections_new table")
		fmt.Println("   Action: Migration 003 didn't complete properly")
	} else if !isOldStructure && page1Count < 2 {
		fmt.Println("\n🔧 FIX NEEDED: New structure detected but page 1 has only", page1Count, "section(s)")
		fmt.Println("   Action: Data needs to be populated/seeded")
	} else {
		fmt.Println("\n✅ Structure looks correct!")
		if page1Count >= 2 {
			fmt.Println("   Page 1 has", page1Count, "sections (expected 5+)")
		}
	}

	// Show sample data if new structure
	if !isOldStructure && sectionCount > 0 {
		fmt.Println("\n📝 Sample sections (first 5):")
		rows, err := database.DB.Query(`
			SELECT page_number, section_number, 
			       SUBSTR(content, 1, 50) as preview
			FROM sections 
			ORDER BY page_number, section_number 
			LIMIT 5
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var pageNum, sectionNum int
				var preview string
				rows.Scan(&pageNum, &sectionNum, &preview)
				fmt.Printf("   Page %d, Section %d: %s...\n", pageNum, sectionNum, preview)
			}
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
