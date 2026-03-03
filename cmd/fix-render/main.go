package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sections-data.sql
var embeddedSectionsData string

func main() {
	// Load configuration
	cfg := config.Load()

	fmt.Println("🔧 Render.com Sections Fix Script")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("")

	// Initialize database (PostgreSQL when DATABASE_URL set, else SQLite)
	if err := database.InitDB(cfg.DBPath, cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	if cfg.DatabaseURL != "" {
		fmt.Println("✅ Database connected: PostgreSQL (persistent)")
	} else {
		fmt.Println("✅ Database connected:", cfg.DBPath)
	}
	fmt.Println("")

	// Step 1: Diagnose current state
	fmt.Println("📊 Step 1: Diagnosing current database state...")
	fmt.Println("-" + strings.Repeat("-", 60))

	hasSectionsTable := false
	var err error

	if database.DriverName == "postgres" {
		// PostgreSQL: use information_schema (migrations create correct schema from the start)
		var n int
		err = database.DB.QueryRow(`
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = 'sections'
		`).Scan(&n)
		hasSectionsTable = (err == nil)
		if !hasSectionsTable {
			log.Fatal("❌ ERROR: No sections table found. Please run migrations first.")
		}
		fmt.Printf("   PostgreSQL: sections table exists (migrations define schema)\n")
	} else {
		// SQLite: check sqlite_master for structure and sections_new migration
		var tableSQL string
		err = database.DB.QueryRow(`
			SELECT sql FROM sqlite_master 
			WHERE type='table' AND name='sections'
		`).Scan(&tableSQL)
		hasSectionsTable = (err == nil)

		if hasSectionsTable {
			isNewStructure := strings.Contains(tableSQL, "page_number") && strings.Contains(tableSQL, "section_number")
			if !isNewStructure {
				fmt.Printf("   Detected OLD structure in 'sections' table\n")
				var sectionsNewSQL string
				err2 := database.DB.QueryRow(`
					SELECT sql FROM sqlite_master 
					WHERE type='table' AND name='sections_new'
				`).Scan(&sectionsNewSQL)
				if err2 == nil {
					fmt.Println("   Migrating: Dropping old 'sections', renaming 'sections_new' to 'sections'...")
					_, _ = database.DB.Exec("DROP TABLE IF EXISTS sections")
					_, err = database.DB.Exec("ALTER TABLE sections_new RENAME TO sections")
					if err != nil {
						log.Fatalf("❌ Failed to rename sections_new: %v", err)
					}
					fmt.Println("   ✓ Migration completed")
				} else {
					log.Fatal("❌ ERROR: Old database structure detected. Please run migrations first.")
				}
			} else {
				fmt.Printf("   Detected NEW structure in 'sections' table\n")
			}
		} else {
			var sectionsNewSQL string
			err2 := database.DB.QueryRow(`
				SELECT sql FROM sqlite_master WHERE type='table' AND name='sections_new'
			`).Scan(&sectionsNewSQL)
			if err2 == nil {
				_, _ = database.DB.Exec("ALTER TABLE sections_new RENAME TO sections")
				hasSectionsTable = true
				fmt.Println("   ✓ Renamed 'sections_new' to 'sections'")
			} else {
				log.Fatal("❌ ERROR: No sections table found. Please run migrations first.")
			}
		}
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

	// Step 3: Create pages if they don't exist
	fmt.Println("")
	fmt.Println("📄 Step 3: Creating pages if needed...")
	fmt.Println("-" + strings.Repeat("-", 60))

	// Read sections data to extract page numbers
	sectionsData := getSectionsData()
	if sectionsData == "" {
		log.Fatal("❌ ERROR: Could not load sections data")
	}

	// Extract unique page numbers from sections data (page-1, page-2, etc.)
	pageMap := make(map[int]bool)
	lines := strings.Split(sectionsData, "\n")
	for _, line := range lines {
		// Extract page_id from INSERT statements like: INSERT INTO sections VALUES('page-1-section-1','page-1',1,1,...
		if strings.Contains(line, "'page-") {
			// Find the page number in the line (e.g., 'page-1' -> 1)
			parts := strings.Split(line, ",")
			if len(parts) >= 3 {
				// parts[1] should be the page_id like 'page-1'
				pageID := strings.Trim(parts[1], " '")
				if strings.HasPrefix(pageID, "page-") {
					var pageNum int
					_, err := fmt.Sscanf(pageID, "page-%d", &pageNum)
					if err == nil && pageNum > 0 {
						pageMap[pageNum] = true
					}
				}
			}
		}
	}

	if len(pageMap) == 0 {
		log.Fatal("❌ ERROR: Could not extract page numbers from sections data")
	}

	fmt.Printf("   Found %d unique pages needed (page 1-%d)\n", len(pageMap), len(pageMap))

	// Create pages that don't exist
	createdPages := 0
	txPages, err := database.DB.Begin()
	if err != nil {
		log.Fatalf("❌ Failed to begin transaction: %v", err)
	}

	// Check if book exists
	var bookExists int
	database.DB.QueryRow("SELECT COUNT(*) FROM books WHERE id = 'alice-in-wonderland'").Scan(&bookExists)
	if bookExists == 0 {
		// Create the book
		_, err = txPages.Exec(`
			INSERT INTO books (id, title, author, description, total_pages)
			VALUES ('alice-in-wonderland', 'Alice''s Adventures in Wonderland', 'Lewis Carroll',
			        'The classic tale of a girl who falls through a rabbit hole into a fantasy world.', 100)
		`)
		if err != nil {
			txPages.Rollback()
			log.Printf("⚠️  Warning: Could not create book (may already exist): %v", err)
		} else {
			fmt.Println("   ✓ Created book 'alice-in-wonderland'")
		}
	}

	// Create each page (idempotent: INSERT OR IGNORE for SQLite, ON CONFLICT for Postgres)
	pageInsertQuery := `INSERT OR IGNORE INTO pages (id, book_id, page_number) VALUES (?, 'alice-in-wonderland', ?)`
	if database.DriverName == "postgres" {
		pageInsertQuery = `INSERT INTO pages (id, book_id, page_number) VALUES ($1, 'alice-in-wonderland', $2) ON CONFLICT (id) DO NOTHING`
	}
	for pageNum := range pageMap {
		pageID := fmt.Sprintf("page-%d", pageNum)
		_, err = txPages.Exec(database.Rebind(pageInsertQuery), pageID, pageNum)
		if err != nil {
			txPages.Rollback()
			log.Fatalf("❌ Failed to create page %d: %v", pageNum, err)
		}
		createdPages++
	}

	if err = txPages.Commit(); err != nil {
		log.Fatalf("❌ Failed to commit pages: %v", err)
	}

	fmt.Printf("   ✓ Created/verified %d pages\n", createdPages)

	// Step 4: Import sections data
	fmt.Println("")
	fmt.Println("📥 Step 4: Preparing to import sections data...")
	fmt.Println("-" + strings.Repeat("-", 60))

	// Check if we should clear existing sections
	clearExisting := false
	if totalSections > 0 && totalSections < 70 {
		fmt.Printf("   Found %d sections (expected 70+), will replace with correct data\n", totalSections)
		clearExisting = true
	}

	// Parse and count INSERT statements
	insertCount := strings.Count(sectionsData, "INSERT INTO sections")
	fmt.Printf("   Found %d sections to import\n", insertCount)

	// Step 5: Import the data
	fmt.Println("")
	fmt.Println("💾 Step 5: Importing sections data...")
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

	// Transform sections data for PostgreSQL (char(10) -> chr(10), etc.)
	dataToImport := sectionsData
	if database.DriverName == "postgres" {
		dataToImport = database.TransformSQLiteToPostgres(sectionsData)
	}

	// Split into individual INSERT statements and execute
	lines2 := strings.Split(dataToImport, "\n")
	importedCount := 0
	failedCount := 0

	for _, line := range lines2 {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "INSERT INTO sections") {
			continue
		}

		stmt := line
		if !strings.HasSuffix(stmt, ";") {
			stmt += ";"
		}

		_, err = tx.Exec(stmt)
		if err != nil {
			// Duplicate key errors (OK when re-running)
			errStr := err.Error()
			if strings.Contains(errStr, "UNIQUE constraint") ||
				strings.Contains(errStr, "constraint failed") ||
				strings.Contains(errStr, "duplicate key") ||
				strings.Contains(errStr, "unique_violation") {
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

	// Step 6: Verify the fix
	fmt.Println("✅ Step 6: Verifying import...")
	fmt.Println("-" + strings.Repeat("-", 60))

	var newTotal int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections").Scan(&newTotal)
	fmt.Printf("   Total sections after import: %d\n", newTotal)

	var newPage1Count int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections WHERE page_number = 1").Scan(&newPage1Count)
	fmt.Printf("   Sections for page 1: %d\n", newPage1Count)

	// Show sample sections for page 1 (SUBSTR works in both SQLite and Postgres)
	previewQuery := database.Rebind(`SELECT section_number, SUBSTR(content, 1, 60) as preview
		FROM sections WHERE page_number = 1 ORDER BY section_number LIMIT 5`)
	fmt.Println("")
	fmt.Println("   Sample sections for page 1:")
	rows, err := database.DB.Query(previewQuery)
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
	// First, try to use embedded data (most reliable - always available)
	if embeddedSectionsData != "" && len(embeddedSectionsData) > 100 {
		fmt.Printf("   ✓ Using embedded sections data (%d bytes)\n", len(embeddedSectionsData))
		return embeddedSectionsData
	}

	// Fallback: Try multiple possible paths (for development/debugging)
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
			fmt.Printf("   ✓ Found sections data at: %s (%d bytes)\n", path, len(data))
			return string(data)
		}
	}

	// If nothing found, return empty (should not happen with embedded data)
	fmt.Printf("   ⚠️  Warning: No sections data found (embedded=%d bytes)\n", len(embeddedSectionsData))
	return embeddedSectionsData
}
