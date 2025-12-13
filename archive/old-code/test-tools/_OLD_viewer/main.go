package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "data/alice-suite.db"

func main() {
	// Open database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		fmt.Printf("❌ Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("📖 Alice Book Viewer - Command Line")
	fmt.Println("=====================================\n")

	// Get book info
	var bookTitle, bookAuthor string
	var totalPages int
	err = db.QueryRow("SELECT title, author, total_pages FROM books WHERE id = 'alice-in-wonderland'").Scan(&bookTitle, &bookAuthor, &totalPages)
	if err != nil {
		fmt.Printf("❌ Error reading book: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📚 Book: %s\n", bookTitle)
	fmt.Printf("✍️  Author: %s\n", bookAuthor)
	fmt.Printf("📄 Total Pages: %d\n\n", totalPages)

	// Get chapters
	rows, err := db.Query(`
		SELECT id, title, number 
		FROM chapters 
		WHERE book_id = 'alice-in-wonderland' 
		ORDER BY number
	`)
	if err != nil {
		fmt.Printf("❌ Error reading chapters: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	chapters := []struct {
		id     string
		title  string
		number int
	}{}

	for rows.Next() {
		var id, title string
		var number int
		if err := rows.Scan(&id, &title, &number); err != nil {
			continue
		}
		chapters = append(chapters, struct {
			id     string
			title  string
			number int
		}{id, title, number})
	}

	fmt.Println("📑 Chapters:")
	for i, ch := range chapters {
		fmt.Printf("  %d. %s\n", i+1, ch.title)
	}

	// Get sections count
	var sectionCount int
	db.QueryRow("SELECT COUNT(*) FROM sections").Scan(&sectionCount)
	fmt.Printf("\n📝 Total Sections: %d\n", sectionCount)

	// Get glossary count
	var glossaryCount int
	db.QueryRow("SELECT COUNT(*) FROM alice_glossary").Scan(&glossaryCount)
	fmt.Printf("📖 Glossary Terms: %d\n", glossaryCount)

	// Show sample section
	fmt.Println("\n📄 Sample Section (Chapter 1, Section 1):")
	var content string
	err = db.QueryRow(`
		SELECT content 
		FROM sections 
		WHERE id = 'chapter-1-section-1'
	`).Scan(&content)
	if err == nil {
		preview := content
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		fmt.Printf("\n%s\n", preview)
	}

	fmt.Println("\n✅ Database is properly loaded!")
	fmt.Println("\n💡 To view in browser:")
	fmt.Println("   1. Make sure server is running: go run cmd/reader/main.go")
	fmt.Println("   2. Open: http://localhost:8080/viewer.html")
	fmt.Println("   3. Or use DB Browser for SQLite to browse the database directly")
}

