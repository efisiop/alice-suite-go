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

// WordCount approximates word count in text
func wordCount(text string) int {
	words := strings.Fields(text)
	return len(words)
}

// SplitIntoPages splits content into pages based on word count
// Average page: 180-200 words
// Sections: ~60-65 words each (about 1/3 of a page)
// Improved algorithm for more equalized sections
func splitIntoPagesAndSections(content string, targetWordsPerPage int, targetWordsPerSection int) ([]PageData, error) {
	// Split into sentences more carefully, preserving punctuation
	sentences := []string{}
	parts := strings.FieldsFunc(content, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})
	
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Find punctuation
		var punct string = "."
		if i < len(parts)-1 {
			idx := strings.Index(content, part)
			if idx >= 0 && idx+len(part) < len(content) {
				nextChar := content[idx+len(part)]
				if nextChar == '.' || nextChar == '!' || nextChar == '?' {
					punct = string(nextChar)
				}
			}
		}
		
		sentences = append(sentences, part+punct)
	}
	
	var pages []PageData
	currentPage := PageData{Sections: []SectionData{}}
	currentSection := SectionData{}
	
	// Target range for sections: 35-45 words (targeting 40 words average)
	lowerBound := targetWordsPerSection - 5  // 35 words
	upperBound := targetWordsPerSection + 5  // 45 words
	
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}
		
		sentenceWords := wordCount(sentence)
		currentSectionWords := wordCount(currentSection.Content)
		
		// Calculate current page word count
		currentPageWords := 0
		for _, sec := range currentPage.Sections {
			currentPageWords += sec.WordCount
		}
		currentPageWords += currentSectionWords
		
		// Check if adding this sentence would exceed page limit
		if currentPageWords+sentenceWords > targetWordsPerPage && currentPageWords > 0 {
			// Finish current section if it has content
			if currentSection.Content != "" {
				currentSection.WordCount = wordCount(currentSection.Content)
				currentPage.Sections = append(currentPage.Sections, currentSection)
			}
			
			// Finish current page
			if len(currentPage.Sections) > 0 {
				currentPage.WordCount = 0
				for _, sec := range currentPage.Sections {
					currentPage.WordCount += sec.WordCount
				}
				pages = append(pages, currentPage)
				currentPage = PageData{Sections: []SectionData{}}
			}
			
			// Start new section
			currentSection = SectionData{Content: sentence}
		} else {
			wouldBeWords := currentSectionWords + sentenceWords
			
			// If current section is at or above lower bound
			if currentSectionWords >= lowerBound {
				// If adding would exceed upper bound, start new section
				if wouldBeWords > upperBound {
					currentSection.WordCount = wordCount(currentSection.Content)
					currentPage.Sections = append(currentPage.Sections, currentSection)
					currentSection = SectionData{Content: sentence}
				} else {
					// Still within range, add to current section
					if currentSection.Content != "" {
						currentSection.Content += " "
					}
					currentSection.Content += sentence
				}
			} else {
				// Below lower bound, always add (unless it would make it way too large)
				if wouldBeWords <= upperBound+10 { // Allow slightly over for very long sentences
					if currentSection.Content != "" {
						currentSection.Content += " "
					}
					currentSection.Content += sentence
				} else {
					// Sentence is too long, split it if possible
					// For now, just add it and let it be a larger section
					if currentSection.Content != "" {
						currentSection.Content += " "
					}
					currentSection.Content += sentence
				}
			}
		}
	}
	
	// Add final section if it has content
	if currentSection.Content != "" {
		currentSection.WordCount = wordCount(currentSection.Content)
		currentPage.Sections = append(currentPage.Sections, currentSection)
	}
	
	// Add final page if it has sections
	if len(currentPage.Sections) > 0 {
		currentPage.WordCount = 0
		for _, sec := range currentPage.Sections {
			currentPage.WordCount += sec.WordCount
		}
		pages = append(pages, currentPage)
	}
	
	// Post-process: Balance sections within pages
	// Merge very small sections (< 45 words) with adjacent sections if possible
	for i := range pages {
		balancedSections := []SectionData{}
		skipNext := false
		
		for j := 0; j < len(pages[i].Sections); j++ {
			if skipNext {
				skipNext = false
				continue
			}
			
			sec := pages[i].Sections[j]
			
			// If section is too small (< 30 words)
			if sec.WordCount < 30 {
				merged := false
				
				// First try: merge with next section if combined would be <= 50 words
				if j < len(pages[i].Sections)-1 {
					nextSec := pages[i].Sections[j+1]
					combinedWords := sec.WordCount + nextSec.WordCount
					
					if combinedWords <= 50 {
						// Merge sections forward
						mergedSec := SectionData{
							Content:   sec.Content + " " + nextSec.Content,
							WordCount: combinedWords,
						}
						balancedSections = append(balancedSections, mergedSec)
						skipNext = true
						merged = true
						continue
					}
				}
				
				// Second try: merge with previous section if forward merge didn't work
				if !merged && len(balancedSections) > 0 {
					prevSec := balancedSections[len(balancedSections)-1]
					combinedWords := prevSec.WordCount + sec.WordCount
					
					if combinedWords <= 50 {
						// Merge with previous section
						balancedSections[len(balancedSections)-1] = SectionData{
							Content:   prevSec.Content + " " + sec.Content,
							WordCount: combinedWords,
						}
						merged = true
						continue
					}
				}
				
				// If couldn't merge, add as-is (might be unavoidable)
				if !merged {
					balancedSections = append(balancedSections, sec)
				}
				continue
			}
			
			// If section is too large (> 50 words), try to split it
			if sec.WordCount > 50 {
				// Split into two sections if possible
				words := strings.Fields(sec.Content)
				midPoint := len(words) / 2
				
				// Find a good split point (sentence boundary)
				firstPart := strings.Join(words[:midPoint], " ")
				secondPart := strings.Join(words[midPoint:], " ")
				
				// Try to find a period near the midpoint
				for k := midPoint; k < len(words) && k < midPoint+10; k++ {
					if strings.HasSuffix(words[k], ".") || strings.HasSuffix(words[k], "!") || strings.HasSuffix(words[k], "?") {
						firstPart = strings.Join(words[:k+1], " ")
						secondPart = strings.Join(words[k+1:], " ")
						break
					}
				}
				
				firstWords := wordCount(firstPart)
				secondWords := wordCount(secondPart)
				
				// Only split if both parts are reasonable (25-45 words)
				if firstWords >= 25 && firstWords <= 45 && secondWords >= 25 && secondWords <= 45 {
					balancedSections = append(balancedSections, SectionData{
						Content:   firstPart,
						WordCount: firstWords,
					})
					balancedSections = append(balancedSections, SectionData{
						Content:   secondPart,
						WordCount: secondWords,
					})
					continue
				}
			}
			
			// Add section as-is
			balancedSections = append(balancedSections, sec)
		}
		
		// Update sections and recalculate page word count
		pages[i].Sections = balancedSections
		pages[i].WordCount = 0
		for _, sec := range pages[i].Sections {
			pages[i].WordCount += sec.WordCount
		}
	}
	
	return pages, nil
}

type PageData struct {
	PageNumber  int
	ChapterID   string
	ChapterTitle string
	Content     string
	WordCount   int
	Sections    []SectionData
}

type SectionData struct {
	SectionNumber int
	Content       string
	WordCount     int
}

func main() {
	fmt.Println("🔄 Restructuring database to page-based system...")
	fmt.Println("📖 Physical book structure: Pages -> Sections (40 words average)")
	fmt.Println("")
	
	// Open database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// Read migration file
	migrationPath := filepath.Join("migrations", "003_restructure_pages_and_sections.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}
	
	// Execute migration
	fmt.Println("📄 Running migration 003...")
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("✅ Migration completed")
	
	// Get all chapters and their sections
	fmt.Println("\n📚 Reading existing chapter data...")
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
	
	// Get all sections for each chapter
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
			log.Printf("Warning: Failed to query sections for chapter %s: %v", chapter.ID, err)
			continue
		}
		
		// Combine all sections into one content string
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
		
		// Split into pages and sections
		// Target: 180-200 words per page, 40 words per section (average)
		pages, err := splitIntoPagesAndSections(fullContent.String(), 190, 40)
		if err != nil {
			log.Printf("Warning: Failed to split chapter %s: %v", chapter.ID, err)
			continue
		}
		
		// Insert pages and sections
		firstPageOfChapter := true
		for _, pageData := range pages {
			pageID := fmt.Sprintf("page-%d", pageNumber)
			
			// Determine chapter title (only on first page of chapter)
			var chapterTitle *string
			if firstPageOfChapter {
				title := chapter.Title
				chapterTitle = &title
			}
			
			// Insert page
			_, err = db.Exec(`
				INSERT INTO pages (id, book_id, page_number, chapter_id, chapter_title, content, word_count)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, pageID, chapter.BookID, pageNumber, chapter.ID, chapterTitle, pageData.Content, pageData.WordCount)
			if err != nil {
				log.Printf("Warning: Failed to insert page %d: %v", pageNumber, err)
				pageNumber++
				continue
			}
			
			// Insert sections for this page
			for i, sectionData := range pageData.Sections {
				sectionID := fmt.Sprintf("page-%d-section-%d", pageNumber, i+1)
				_, err = db.Exec(`
					INSERT INTO sections_new (id, page_id, page_number, section_number, content, word_count)
					VALUES (?, ?, ?, ?, ?, ?)
				`, sectionID, pageID, pageNumber, i+1, sectionData.Content, sectionData.WordCount)
				if err != nil {
					log.Printf("Warning: Failed to insert section %d for page %d: %v", i+1, pageNumber, err)
				}
			}
			
			fmt.Printf("    ✅ Page %d: %d sections, %d words\n", pageNumber, len(pageData.Sections), pageData.WordCount)
			pageNumber++
			firstPageOfChapter = false
		}
	}
	
	// Drop old sections table and rename new one
	fmt.Println("\n🔄 Replacing old sections table...")
	_, err = db.Exec("DROP TABLE IF EXISTS sections")
	if err != nil {
		log.Printf("Warning: Failed to drop old sections table: %v", err)
	}
	
	_, err = db.Exec("ALTER TABLE sections_new RENAME TO sections")
	if err != nil {
		log.Fatalf("Failed to rename sections table: %v", err)
	}
	
	fmt.Println("✅ Database restructured successfully!")
	fmt.Printf("\n📊 Summary:\n")
	
	var pageCount, sectionCount int
	db.QueryRow("SELECT COUNT(*) FROM pages").Scan(&pageCount)
	db.QueryRow("SELECT COUNT(*) FROM sections").Scan(&sectionCount)
	
	fmt.Printf("  Pages: %d\n", pageCount)
	fmt.Printf("  Sections: %d\n", sectionCount)
	fmt.Printf("  Average sections per page: %.1f\n", float64(sectionCount)/float64(pageCount))
}

