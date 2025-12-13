package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
)

const dbPath = "data/alice-suite.db"

// NormalizeTerm normalizes a term for matching (lowercase, remove punctuation)
func normalizeTerm(term string) string {
	term = strings.ToLower(strings.TrimSpace(term))
	// Remove common punctuation
	term = strings.Trim(term, ".,!?;:'\"()[]{}")
	return term
}

// FindTermInText finds if a term appears in text (case-insensitive, word boundaries)
func findTermInText(text, term string) bool {
	normalizedTerm := normalizeTerm(term)
	normalizedText := strings.ToLower(text)
	
	// Skip very common words that are likely false positives
	commonWords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true, "but": true,
		"not": true, "is": true, "are": true, "was": true, "were": true,
		"be": true, "been": true, "have": true, "has": true, "had": true,
		"do": true, "does": true, "did": true, "will": true, "would": true,
		"can": true, "could": true, "should": true, "may": true, "might": true,
		"it": true, "its": true, "this": true, "that": true, "these": true, "those": true,
		"i": true, "you": true, "he": true, "she": true, "we": true, "they": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "from": true, "by": true, "about": true, "into": true,
		"down": true, "up": true, "out": true, "off": true, "over": true, "under": true,
		"get": true, "got": true, "go": true, "went": true, "come": true, "came": true,
		"say": true, "said": true, "know": true, "think": true, "thought": true,
		"see": true, "saw": true, "look": true, "like": true, "well": true, "let": true,
		"must": true, "little": true, "again": true,
	}
	
	// Check if it's a common word (skip if single common word)
	if len(strings.Fields(normalizedTerm)) == 1 && commonWords[normalizedTerm] {
		return false
	}
	
	// For multi-word terms, search for exact phrase
	if strings.Contains(term, " ") || strings.Contains(term, "-") {
		return strings.Contains(normalizedText, normalizedTerm)
	}
	
	// For single words, use word boundaries (more precise)
	pattern := `\b` + regexp.QuoteMeta(normalizedTerm) + `\b`
	matched, _ := regexp.MatchString(pattern, normalizedText)
	return matched
}

func main() {
	fmt.Println("🔗 Linking Glossary Terms to Sections")
	fmt.Println("📖 Finding where glossary terms appear in book sections")
	fmt.Println("")

	// Open database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Run migration to create junction table
	fmt.Println("📄 Running migration 004...")
	migrationPath := filepath.Join("migrations", "004_link_glossary_to_sections.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("✅ Migration completed")

	// Load glossary terms from SQL file
	fmt.Println("\n📚 Loading glossary terms from alice_glossary.sql...")
	glossarySQLPath := "alice_glossary.sql"
	glossarySQL, err := os.ReadFile(glossarySQLPath)
	if err != nil {
		log.Fatalf("Failed to read glossary file: %v", err)
	}

	// Execute glossary SQL (remove BEGIN/COMMIT if present)
	glossarySQLStr := string(glossarySQL)
	glossarySQLStr = strings.ReplaceAll(glossarySQLStr, "BEGIN TRANSACTION;", "")
	glossarySQLStr = strings.ReplaceAll(glossarySQLStr, "COMMIT;", "")
	
	// Split by semicolons and execute each INSERT
	statements := strings.Split(glossarySQLStr, ";")
	insertedCount := 0
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || !strings.HasPrefix(stmt, "INSERT") {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			// Some errors are expected (duplicates), continue
			continue
		}
		insertedCount++
	}
	fmt.Printf("✅ Loaded %d glossary terms\n", insertedCount)

	// Get all glossary terms
	fmt.Println("\n🔍 Getting all glossary terms...")
	glossaryRows, err := db.Query(`
		SELECT id, term, definition
		FROM alice_glossary
		WHERE book_id = 'alice-in-wonderland'
		ORDER BY term
	`)
	if err != nil {
		log.Fatalf("Failed to query glossary: %v", err)
	}
	defer glossaryRows.Close()

	type GlossaryTerm struct {
		ID         string
		Term       string
		Definition string
	}

	terms := []GlossaryTerm{}
	for glossaryRows.Next() {
		var t GlossaryTerm
		glossaryRows.Scan(&t.ID, &t.Term, &t.Definition)
		terms = append(terms, t)
	}
	fmt.Printf("✅ Found %d glossary terms\n", len(terms))

	// Get all sections
	fmt.Println("\n📄 Getting all sections...")
	sectionRows, err := db.Query(`
		SELECT id, page_number, section_number, content
		FROM sections
		ORDER BY page_number, section_number
	`)
	if err != nil {
		log.Fatalf("Failed to query sections: %v", err)
	}
	defer sectionRows.Close()

	type Section struct {
		ID           string
		PageNumber   int
		SectionNumber int
		Content      string
	}

	sections := []Section{}
	for sectionRows.Next() {
		var s Section
		sectionRows.Scan(&s.ID, &s.PageNumber, &s.SectionNumber, &s.Content)
		sections = append(sections, s)
	}
	fmt.Printf("✅ Found %d sections\n", len(sections))

	// Find matches and create links
	fmt.Println("\n🔗 Linking terms to sections...")
	linksCreated := 0
	
	for _, term := range terms {
		for _, section := range sections {
			if findTermInText(section.Content, term.Term) {
				// Create link
				linkID := uuid.New().String()
				_, err := db.Exec(`
					INSERT OR IGNORE INTO glossary_section_links 
					(id, glossary_id, section_id, page_number, section_number, term)
					VALUES (?, ?, ?, ?, ?, ?)
				`, linkID, term.ID, section.ID, section.PageNumber, section.SectionNumber, term.Term)
				if err != nil {
					log.Printf("Warning: Failed to create link for term '%s' in section %s: %v", term.Term, section.ID, err)
				} else {
					linksCreated++
				}
			}
		}
	}

	fmt.Printf("✅ Created %d glossary-section links\n", linksCreated)

	// Show statistics
	fmt.Println("\n📊 Statistics:")
	var totalTerms, linkedTerms, totalLinks int
	db.QueryRow("SELECT COUNT(*) FROM alice_glossary WHERE book_id = 'alice-in-wonderland'").Scan(&totalTerms)
	db.QueryRow(`
		SELECT COUNT(DISTINCT glossary_id) 
		FROM glossary_section_links
	`).Scan(&linkedTerms)
	db.QueryRow("SELECT COUNT(*) FROM glossary_section_links").Scan(&totalLinks)

	fmt.Printf("  Total glossary terms: %d\n", totalTerms)
	fmt.Printf("  Terms linked to sections: %d\n", linkedTerms)
	fmt.Printf("  Total links created: %d\n", totalLinks)
	fmt.Printf("  Average links per term: %.1f\n", float64(totalLinks)/float64(linkedTerms))

	// Show sample links
	fmt.Println("\n📝 Sample links (first 10):")
	sampleRows, err := db.Query(`
		SELECT g.term, gs.page_number, gs.section_number
		FROM glossary_section_links gs
		JOIN alice_glossary g ON gs.glossary_id = g.id
		ORDER BY g.term, gs.page_number, gs.section_number
		LIMIT 10
	`)
	if err == nil {
		defer sampleRows.Close()
		for sampleRows.Next() {
			var term string
			var pageNum, sectionNum int
			sampleRows.Scan(&term, &pageNum, &sectionNum)
			fmt.Printf("  '%s' → Page %d, Section %d\n", term, pageNum, sectionNum)
		}
	}

	fmt.Println("\n✅ Glossary linking completed!")
}

