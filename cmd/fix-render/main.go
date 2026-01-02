package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load configuration
	cfg := config.Load()

	fmt.Println("🔧 Render.com Sections Fix Script")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("")

	// Initialize database
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	fmt.Println("✅ Database connected:", cfg.DBPath)
	fmt.Println("")

	// Step 1: Diagnose current state
	fmt.Println("📊 Step 1: Diagnosing current database state...")
	fmt.Println("-" + strings.Repeat("-", 60))

	var tableSQL string
	err := database.DB.QueryRow(`
		SELECT sql FROM sqlite_master 
		WHERE type='table' AND name='sections'
	`).Scan(&tableSQL)

	if err == sql.ErrNoRows {
		log.Fatal("❌ ERROR: sections table does not exist! Run migrations first.")
	} else if err != nil {
		log.Fatalf("❌ Error checking table: %v", err)
	}

	// Check structure
	isNewStructure := strings.Contains(tableSQL, "page_number") && strings.Contains(tableSQL, "section_number")
	if !isNewStructure {
		log.Fatal("❌ ERROR: Old database structure detected. Please run migrations first.")
	}

	// Count existing sections
	var totalSections int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections").Scan(&totalSections)
	fmt.Printf("   Current sections in database: %d\n", totalSections)

	var page1Count int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections WHERE page_number = 1").Scan(&page1Count)
	fmt.Printf("   Sections for page 1: %d\n", page1Count)

	if page1Count >= 5 && totalSections >= 70 {
		fmt.Println("")
		fmt.Println("✅ Database already has correct data!")
		fmt.Printf("   Page 1 has %d sections (expected 5+)\n", page1Count)
		fmt.Printf("   Total sections: %d (expected 70+)\n", totalSections)
		fmt.Println("")
		fmt.Println("🎉 No fix needed. Exiting.")
		os.Exit(0)
	}

	fmt.Println("")
	fmt.Printf("⚠️  Issue detected: Page 1 has only %d section(s) (expected 5+)\n", page1Count)
	fmt.Println("")

	// Step 2: Check if pages table has data
	fmt.Println("📊 Step 2: Checking pages table...")
	fmt.Println("-" + strings.Repeat("-", 60))

	var pageCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM pages").Scan(&pageCount)
	fmt.Printf("   Pages in database: %d\n", pageCount)

	if pageCount == 0 {
		fmt.Println("")
		fmt.Println("⚠️  WARNING: No pages found in database!")
		fmt.Println("   Sections require pages to exist. Please run migrations/seeds first.")
		fmt.Println("")
	}

	// Step 3: Import sections data
	fmt.Println("")
	fmt.Println("📥 Step 3: Preparing to import sections data...")
	fmt.Println("-" + strings.Repeat("-", 60))

	// Check if we should clear existing sections
	clearExisting := false
	if totalSections > 0 && totalSections < 70 {
		fmt.Printf("   Found %d sections (expected 70+), will replace with correct data\n", totalSections)
		clearExisting = true
	}

	// Read the embedded sections data
	sectionsData := getSectionsData()

	if sectionsData == "" {
		log.Fatal("❌ ERROR: Could not load sections data")
	}

	// Parse and count INSERT statements
	insertCount := strings.Count(sectionsData, "INSERT INTO sections")
	fmt.Printf("   Found %d sections to import\n", insertCount)

	// Step 4: Import the data
	fmt.Println("")
	fmt.Println("💾 Step 4: Importing sections data...")
	fmt.Println("-" + strings.Repeat("-", 60))

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatalf("❌ Failed to begin transaction: %v", err)
	}

	// Clear existing sections if needed
	if clearExisting {
		fmt.Println("   Clearing existing sections...")
		_, err = tx.Exec("DELETE FROM sections")
		if err != nil {
			tx.Rollback()
			log.Fatalf("❌ Failed to clear existing sections: %v", err)
		}
		fmt.Printf("   ✓ Cleared %d existing sections\n", totalSections)
	}

	// Split into individual INSERT statements and execute
	// Each line in the file is one INSERT statement
	lines := strings.Split(sectionsData, "\n")
	importedCount := 0
	failedCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "INSERT INTO sections") {
			continue
		}

		// Ensure statement ends with semicolon (for safety, though each line should have it)
		stmt := line
		if !strings.HasSuffix(stmt, ";") {
			stmt += ";"
		}

		_, err = tx.Exec(stmt)
		if err != nil {
			// Check if it's a duplicate key error (OK if we're re-running)
			if strings.Contains(err.Error(), "UNIQUE constraint") || strings.Contains(err.Error(), "constraint failed") {
				// Section already exists, skip it
				continue
			}
			// Log error but continue (might be other recoverable issues)
			failedCount++
			fmt.Printf("   ⚠️  Warning: Failed to import one section (continuing...): %v\n", err)
			// Don't fail completely, just skip this one
			continue
		}
		importedCount++
	}

	if failedCount > 0 {
		fmt.Printf("   ⚠️  %d sections failed to import (may already exist)\n", failedCount)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Fatalf("❌ Failed to commit transaction: %v", err)
	}

	fmt.Printf("   ✓ Successfully imported %d sections\n", importedCount)
	fmt.Println("")

	// Step 5: Verify the fix
	fmt.Println("✅ Step 5: Verifying import...")
	fmt.Println("-" + strings.Repeat("-", 60))

	var newTotal int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections").Scan(&newTotal)
	fmt.Printf("   Total sections after import: %d\n", newTotal)

	var newPage1Count int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections WHERE page_number = 1").Scan(&newPage1Count)
	fmt.Printf("   Sections for page 1: %d\n", newPage1Count)

	// Show sample sections for page 1
	fmt.Println("")
	fmt.Println("   Sample sections for page 1:")
	rows, err := database.DB.Query(`
		SELECT section_number, SUBSTR(content, 1, 60) as preview
		FROM sections 
		WHERE page_number = 1 
		ORDER BY section_number
		LIMIT 5
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var sectionNum int
			var preview string
			rows.Scan(&sectionNum, &preview)
			fmt.Printf("     Section %d: %s...\n", sectionNum, preview)
		}
	}

	fmt.Println("")
	if newPage1Count >= 5 {
		fmt.Println("🎉 SUCCESS! Fix completed successfully!")
		fmt.Printf("   Page 1 now has %d sections (expected 5+)\n", newPage1Count)
		fmt.Println("   You can now test the Render.com reader app.")
	} else {
		fmt.Println("⚠️  WARNING: Page 1 still has less than 5 sections")
		fmt.Println("   Please check the import and verify pages table has data.")
	}
	fmt.Println("")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getSectionsData returns the embedded sections SQL data
func getSectionsData() string {
	// Try multiple possible paths (works in different environments)
	possiblePaths := []string{
		"scripts/sections-data.sql",
		"../scripts/sections-data.sql",
		"../../scripts/sections-data.sql",
		"./scripts/sections-data.sql",
		// For Render.com build environment
		"/opt/render/project/src/scripts/sections-data.sql",
	}

	for _, path := range possiblePaths {
		if data, err := os.ReadFile(path); err == nil {
			fmt.Printf("   ✓ Found sections data at: %s\n", path)
			return string(data)
		}
	}

	// If file not found, try embedded data (fallback)
	return embeddedSectionsData
}

// embeddedSectionsData contains the actual SQL INSERT statements
// This is a fallback if the file can't be read
var embeddedSectionsData = getEmbeddedSectionsData()

func getEmbeddedSectionsData() string {
	// Try to read from the file using all possible paths
	possiblePaths := []string{
		"scripts/sections-data.sql",
		"../scripts/sections-data.sql",
		"../../scripts/sections-data.sql",
	}

	for _, path := range possiblePaths {
		if data, err := os.ReadFile(path); err == nil {
			return string(data)
		}
	}
	return ""
}

