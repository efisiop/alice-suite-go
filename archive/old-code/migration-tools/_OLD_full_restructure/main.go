package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "data/alice-suite.db"

func wordCount(text string) int {
	return len(strings.Fields(text))
}

func splitIntoPagesAndSections(content string, targetWordsPerPage int, targetWordsPerSection int) ([]PageData, error) {
	sentences := strings.Split(content, ". ")
	
	var pages []PageData
	currentPage := PageData{Sections: []SectionData{}}
	currentSection := SectionData{}
	
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}
		
		if !strings.HasSuffix(sentence, ".") && !strings.HasSuffix(sentence, "!") && !strings.HasSuffix(sentence, "?") {
			sentence += "."
		}
		
		sentenceWords := wordCount(sentence)
		currentSectionWords := wordCount(currentSection.Content)
		currentPageWords := 0
		for _, sec := range currentPage.Sections {
			currentPageWords += sec.WordCount
		}
		currentPageWords += currentSectionWords
		
		if currentPageWords+sentenceWords > targetWordsPerPage && currentPageWords > 0 {
			if currentSection.Content != "" {
				currentSection.WordCount = wordCount(currentSection.Content)
				currentPage.Sections = append(currentPage.Sections, currentSection)
			}
			
			if len(currentPage.Sections) > 0 {
				currentPage.WordCount = 0
				for _, sec := range currentPage.Sections {
					currentPage.WordCount += sec.WordCount
				}
				pages = append(pages, currentPage)
				currentPage = PageData{Sections: []SectionData{}}
			}
			
			currentSection = SectionData{Content: sentence}
		} else {
			if currentSectionWords+sentenceWords > targetWordsPerSection && currentSection.Content != "" {
				currentSection.WordCount = wordCount(currentSection.Content)
				currentPage.Sections = append(currentPage.Sections, currentSection)
				currentSection = SectionData{Content: sentence}
			} else {
				if currentSection.Content != "" {
					currentSection.Content += " "
				}
				currentSection.Content += sentence
			}
		}
	}
	
	if currentSection.Content != "" {
		currentSection.WordCount = wordCount(currentSection.Content)
		currentPage.Sections = append(currentPage.Sections, currentSection)
	}
	
	if len(currentPage.Sections) > 0 {
		currentPage.WordCount = 0
		for _, sec := range currentPage.Sections {
			currentPage.WordCount += sec.WordCount
		}
		pages = append(pages, currentPage)
	}
	
	return pages, nil
}

type PageData struct {
	PageNumber   int
	ChapterID    string
	ChapterTitle string
	Content      string
	WordCount    int
	Sections     []SectionData
}

type SectionData struct {
	SectionNumber int
	Content       string
	WordCount     int
}

func main() {
	fmt.Println("🔄 Complete Database Restructuring")
	fmt.Println("📖 Physical book structure: Pages -> Sections (60-65 words each)")
	fmt.Println("")
	
	// Backup old database
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Println("📦 Backing up existing database...")
		os.Rename(dbPath, dbPath+".backup")
	}
	
	// Create fresh database
	fmt.Println("🗄️  Creating fresh database...")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// Run initial schema
	fmt.Println("📄 Running initial schema migration...")
	schemaPath := filepath.Join("migrations", "001_initial_schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read schema: %v", err)
	}
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		log.Fatalf("Schema migration failed: %v", err)
	}
	
	// Load seed data using the seed tool approach (Go code)
	fmt.Println("🌱 Loading seed data...")
	go run cmd/seed/main.go
	
	// Wait a moment for seed to complete, then run page restructuring
	fmt.Println("📄 Running page restructuring migration...")
	restructurePath := filepath.Join("migrations", "003_restructure_pages_and_sections.sql")
	restructureSQL, err := os.ReadFile(restructurePath)
	if err != nil {
		log.Fatalf("Failed to read restructure migration: %v", err)
	}
	_, err = db.Exec(string(restructureSQL))
	if err != nil {
		log.Fatalf("Restructure migration failed: %v", err)
	}
	
	// Now migrate data
	fmt.Println("\n📚 Reading chapter data...")
	chapterRows, err := db.Query(`
		SELECT c.id, c.title, c.number, c.book_id
		FROM chapters c
		WHERE c.book_id = 'alice-in-wonderland'
		ORDER BY c.number
	`)
	if err != nil {
		log.Fatalf("Failed to query chapters: %v", err)
	}
	defer chapterRows.Close()
	
	type ChapterInfo struct {
		ID     string
		Title  string
		Number int
		BookID string
	}
	
	chapters := []ChapterInfo{}
	for chapterRows.Next() {
		var ch ChapterInfo
		chapterRows.Scan(&ch.ID, &ch.Title, &ch.Number, &ch.BookID)
		chapters = append(chapters, ch)
	}
	
	if len(chapters) == 0 {
		log.Fatalf("No chapters found! Make sure seed data is loaded.")
	}
	
	pageNumber := 1
	for _, chapter := range chapters {
		fmt.Printf("  Processing Chapter %d: %s\n", chapter.Number, chapter.Title)
		
		sectionRows, err := db.Query(`
			SELECT content
			FROM sections
			WHERE chapter_id = ?
			ORDER BY number
		`, chapter.ID)
		if err != nil {
			log.Printf("Warning: Failed to query sections: %v", err)
			continue
		}
		
		var fullContent strings.Builder
		for sectionRows.Next() {
			var content string
			sectionRows.Scan(&content)
			if fullContent.Len() > 0 {
				fullContent.WriteString(" ")
			}
			fullContent.WriteString(content)
		}
		sectionRows.Close()
		
		pages, err := splitIntoPagesAndSections(fullContent.String(), 190, 63)
		if err != nil {
			log.Printf("Warning: Failed to split: %v", err)
			continue
		}
		
		firstPageOfChapter := true
		for _, pageData := range pages {
			pageID := fmt.Sprintf("page-%d", pageNumber)
			
			var chapterTitle *string
			if firstPageOfChapter {
				title := chapter.Title
				chapterTitle = &title
			}
			
			_, err = db.Exec(`
				INSERT INTO pages (id, book_id, page_number, chapter_id, chapter_title, content, word_count)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, pageID, chapter.BookID, pageNumber, chapter.ID, chapterTitle, pageData.Content, pageData.WordCount)
			if err != nil {
				log.Printf("Warning: Failed to insert page %d: %v", pageNumber, err)
				pageNumber++
				continue
			}
			
			for i, sectionData := range pageData.Sections {
				sectionID := fmt.Sprintf("page-%d-section-%d", pageNumber, i+1)
				_, err = db.Exec(`
					INSERT INTO sections_new (id, page_id, page_number, section_number, content, word_count)
					VALUES (?, ?, ?, ?, ?, ?)
				`, sectionID, pageID, pageNumber, i+1, sectionData.Content, sectionData.WordCount)
				if err != nil {
					log.Printf("Warning: Failed to insert section: %v", err)
				}
			}
			
			fmt.Printf("    ✅ Page %d: %d sections, %d words\n", pageNumber, len(pageData.Sections), pageData.WordCount)
			pageNumber++
			firstPageOfChapter = false
		}
	}
	
	// Replace old sections table
	fmt.Println("\n🔄 Replacing sections table...")
	_, err = db.Exec("DROP TABLE IF EXISTS sections")
	if err != nil {
		log.Printf("Warning: %v", err)
	}
	
	_, err = db.Exec("ALTER TABLE sections_new RENAME TO sections")
	if err != nil {
		log.Fatalf("Failed to rename: %v", err)
	}
	
	fmt.Println("✅ Database restructured successfully!")
	
	var pageCount, sectionCount int
	db.QueryRow("SELECT COUNT(*) FROM pages").Scan(&pageCount)
	db.QueryRow("SELECT COUNT(*) FROM sections").Scan(&sectionCount)
	
	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("  Pages: %d\n", pageCount)
	fmt.Printf("  Sections: %d\n", sectionCount)
	if pageCount > 0 {
		fmt.Printf("  Average sections per page: %.1f\n", float64(sectionCount)/float64(pageCount))
	}
}

