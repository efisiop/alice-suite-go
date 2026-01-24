package main

import (
	"fmt"
	"log"
	"os"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.Load()

	fmt.Println("🔍 Sections Diagnostic Tool")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println("")

	// Initialize database
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	fmt.Printf("✅ Database connected: %s\n", cfg.DBPath)
	fmt.Println("")

	// Check sections table structure
	fmt.Println("📊 Checking sections table structure...")
	var tableSQL string
	err := database.DB.QueryRow(`
		SELECT sql FROM sqlite_master 
		WHERE type='table' AND name='sections'
	`).Scan(&tableSQL)
	if err != nil {
		log.Fatalf("❌ Sections table not found: %v", err)
	}
	fmt.Println("✅ Sections table exists")
	fmt.Println("")

	// Count total sections
	var totalSections int
	database.DB.QueryRow("SELECT COUNT(*) FROM sections").Scan(&totalSections)
	fmt.Printf("📊 Total sections in database: %d\n", totalSections)

	// Count sections per page
	fmt.Println("")
	fmt.Println("📄 Sections per page (first 20 pages):")
	rows, err := database.DB.Query(`
		SELECT page_number, COUNT(*) as count 
		FROM sections 
		GROUP BY page_number 
		ORDER BY page_number 
		LIMIT 20
	`)
	if err != nil {
		log.Fatalf("Error querying sections: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pageNum, count int
		if err := rows.Scan(&pageNum, &count); err != nil {
			log.Printf("Error scanning: %v", err)
			continue
		}
		fmt.Printf("   Page %d: %d sections\n", pageNum, count)
	}

	// Check page 1 specifically
	fmt.Println("")
	fmt.Println("📄 Page 1 sections detail:")
	rows2, err := database.DB.Query(`
		SELECT id, section_number, LENGTH(content) as content_length, 
		       SUBSTR(content, 1, 50) as preview
		FROM sections 
		WHERE page_number = 1 
		ORDER BY section_number
	`)
	if err != nil {
		log.Fatalf("Error querying page 1: %v", err)
	}
	defer rows2.Close()

	page1Count := 0
	for rows2.Next() {
		var id, preview string
		var sectionNum, contentLen int
		if err := rows2.Scan(&id, &sectionNum, &contentLen, &preview); err != nil {
			log.Printf("Error scanning: %v", err)
			continue
		}
		page1Count++
		fmt.Printf("   Section %d: %s (content: %d chars)\n", sectionNum, id, contentLen)
		fmt.Printf("      Preview: %s...\n", preview)
	}
	fmt.Printf("\n✅ Page 1 has %d sections\n", page1Count)

	// Test the query used by the API
	fmt.Println("")
	fmt.Println("🔍 Testing API query for page 1...")
	query := `SELECT id, content, page_number, section_number FROM sections 
	         WHERE page_number = ? ORDER BY section_number`
	rows3, err := database.DB.Query(query, 1)
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	defer rows3.Close()

	apiSectionCount := 0
	for rows3.Next() {
		var id, content string
		var pageNum, sectionNum int
		if err := rows3.Scan(&id, &content, &pageNum, &sectionNum); err != nil {
			log.Printf("Error scanning: %v", err)
			continue
		}
		apiSectionCount++
	}
	fmt.Printf("✅ API query returned %d sections for page 1\n", apiSectionCount)

	if page1Count < 5 {
		fmt.Println("")
		fmt.Println("⚠️  WARNING: Page 1 has less than 5 sections!")
		fmt.Println("   Run: ./bin/fix-render")
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("🎉 All checks passed!")
}
