package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/models"
	"github.com/google/uuid"
)

// User Queries

// CreateUser creates a new user
func CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	query := `INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Role, user.IsVerified, time.Now(), time.Now())
	return err
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	var createdAtStr, updatedAtStr string
	query := `SELECT id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at
	          FROM users WHERE email = ?`
	err := DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.IsVerified, &createdAtStr, &updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// Parse timestamps
	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = t
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = t
		}
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	var createdAtStr, updatedAtStr string
	query := `SELECT id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at
	          FROM users WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.IsVerified, &createdAtStr, &updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// Parse timestamps
	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = t
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = t
		}
	}
	return user, nil
}

// Book Queries

// GetBookByID retrieves a book by ID
func GetBookByID(id string) (*models.Book, error) {
	book := &models.Book{}
	query := `SELECT id, title, author, description, total_pages, created_at
	          FROM books WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.Author, &book.Description, &book.TotalPages, &book.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return book, err
}

// GetAllBooks retrieves all books
func GetAllBooks() ([]*models.Book, error) {
	if DB == nil {
		return nil, sql.ErrConnDone
	}
	query := `SELECT id, title, author, description, total_pages, created_at FROM books ORDER BY created_at`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*models.Book{}
	for rows.Next() {
		book := &models.Book{}
		var createdAtStr string
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.TotalPages, &createdAtStr)
		if err != nil {
			return nil, err
		}
		// Parse created_at
		if createdAtStr != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
				book.CreatedAt = t
			}
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

// Chapter Queries

// GetChaptersByBookID retrieves all chapters for a book
func GetChaptersByBookID(bookID string) ([]*models.Chapter, error) {
	query := `SELECT id, book_id, title, number, created_at
	          FROM chapters WHERE book_id = ? ORDER BY number`
	rows, err := DB.Query(query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chapters := []*models.Chapter{}
	for rows.Next() {
		chapter := &models.Chapter{}
		var createdAtStr string
		err := rows.Scan(&chapter.ID, &chapter.BookID, &chapter.Title, &chapter.Number, &createdAtStr)
		if err != nil {
			return nil, err
		}
		// Parse created_at
		if createdAtStr != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
				chapter.CreatedAt = t
			}
		}
		chapters = append(chapters, chapter)
	}
	return chapters, rows.Err()
}

// GetChapterByID retrieves a chapter by ID
func GetChapterByID(id string) (*models.Chapter, error) {
	chapter := &models.Chapter{}
	query := `SELECT id, book_id, title, number, created_at
	          FROM chapters WHERE id = ?`

	var createdAtStr string
	err := DB.QueryRow(query, id).Scan(
		&chapter.ID, &chapter.BookID, &chapter.Title, &chapter.Number, &createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// Parse created_at
	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			chapter.CreatedAt = t
		}
	}
	return chapter, nil
}

// Page Queries

// GetPageByNumber retrieves a page with all its sections
func GetPageByNumber(bookID string, pageNumber int) (*models.Page, error) {
	// Get page
	page := &models.Page{}
	var chapterID, chapterTitle sql.NullString
	var createdAtStr string
	query := `SELECT id, book_id, page_number, chapter_id, chapter_title, content, word_count, created_at
	          FROM pages WHERE book_id = ? AND page_number = ?`

	err := DB.QueryRow(query, bookID, pageNumber).Scan(
		&page.ID, &page.BookID, &page.PageNumber, &chapterID, &chapterTitle,
		&page.Content, &page.WordCount, &createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse created_at
	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			page.CreatedAt = t
		}
	}

	if chapterID.Valid {
		page.ChapterID = &chapterID.String
	}
	if chapterTitle.Valid {
		page.ChapterTitle = &chapterTitle.String
	}

	// Get sections for this page
	sectionsQuery := `SELECT id, page_id, page_number, section_number, content, word_count, created_at
	                  FROM sections WHERE page_id = ? ORDER BY section_number`
	rows, err := DB.Query(sectionsQuery, page.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	page.Sections = []models.Section{}
	for rows.Next() {
		var section models.Section
		var sectionCreatedAtStr string
		err := rows.Scan(
			&section.ID, &section.PageID, &section.PageNumber, &section.SectionNumber,
			&section.Content, &section.WordCount, &sectionCreatedAtStr,
		)
		if err != nil {
			return nil, err
		}
		// Parse section created_at
		if sectionCreatedAtStr != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", sectionCreatedAtStr); err == nil {
				section.CreatedAt = t
			}
		}
		page.Sections = append(page.Sections, section)
	}

	return page, rows.Err()
}

// Section Queries

// GetSectionsByChapterID retrieves all sections for a chapter (legacy - for compatibility)
func GetSectionsByChapterID(chapterID string) ([]*models.Section, error) {
	// This is a legacy function - sections are now page-based
	// Return empty for now or implement if needed
	return []*models.Section{}, nil
}

// GetSectionByID retrieves a section by ID
func GetSectionByID(id string) (*models.Section, error) {
	section := &models.Section{}
	var createdAtStr string
	query := `SELECT id, page_id, page_number, section_number, content, word_count, created_at
	          FROM sections WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(
		&section.ID, &section.PageID, &section.PageNumber, &section.SectionNumber,
		&section.Content, &section.WordCount, &createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// Parse created_at
	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			section.CreatedAt = t
		}
	}
	return section, nil
}

// GetSectionByPage retrieves sections by page number (legacy - use GetPageByNumber instead)
func GetSectionByPage(bookID string, page int) (*models.Section, error) {
	// Legacy function - use GetPageByNumber instead
	return nil, nil
}

// Glossary Queries

// GetGlossaryTerm retrieves a glossary term
func GetGlossaryTerm(bookID, term string) (*models.AliceGlossary, error) {
	glossary := &models.AliceGlossary{}
	query := `SELECT id, book_id, term, definition, source_sentence, example, chapter_reference, created_at, updated_at
	          FROM alice_glossary WHERE book_id = ? AND term = ?`

	var sourceSentence, example, chapterRef sql.NullString
	var createdAt, updatedAt string
	
	err := DB.QueryRow(query, bookID, term).Scan(
		&glossary.ID, &glossary.BookID, &glossary.Term, &glossary.Definition,
		&sourceSentence, &example, &chapterRef,
		&createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Handle NULL values
	if sourceSentence.Valid {
		glossary.SourceSentence = sourceSentence.String
	}
	if example.Valid {
		glossary.Example = example.String
	}
	if chapterRef.Valid {
		glossary.ChapterReference = chapterRef.String
	}
	
	// Parse timestamps
	if createdAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			glossary.CreatedAt = t
		} else if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			glossary.CreatedAt = t
		}
	}
	if updatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			glossary.UpdatedAt = t
		} else if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			glossary.UpdatedAt = t
		}
	}
	
	return glossary, nil
}

// SearchGlossaryTerms searches for glossary terms
func SearchGlossaryTerms(bookID, searchTerm string) ([]*models.AliceGlossary, error) {
	query := `SELECT id, book_id, term, definition, source_sentence, example, chapter_reference, created_at, updated_at
	          FROM alice_glossary 
	          WHERE book_id = ? AND (term LIKE ? OR definition LIKE ?)
	          ORDER BY term`

	searchPattern := "%" + searchTerm + "%"
	rows, err := DB.Query(query, bookID, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	terms := []*models.AliceGlossary{}
	for rows.Next() {
		term := &models.AliceGlossary{}
		err := rows.Scan(
			&term.ID, &term.BookID, &term.Term, &term.Definition,
			&term.SourceSentence, &term.Example, &term.ChapterReference,
			&term.CreatedAt, &term.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}
	return terms, rows.Err()
}

// GetAllGlossaryTerms retrieves all glossary terms for a book
func GetAllGlossaryTerms(bookID string) ([]*models.AliceGlossary, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	
	query := `SELECT id, book_id, term, definition, source_sentence, example, chapter_reference, created_at, updated_at
	          FROM alice_glossary 
	          WHERE book_id = ?
	          ORDER BY term`

	rows, err := DB.Query(query, bookID)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	terms := []*models.AliceGlossary{}
	for rows.Next() {
		term := &models.AliceGlossary{}
		var sourceSentence, example, chapterRef sql.NullString
		var createdAt, updatedAt string
		
		err := rows.Scan(
			&term.ID, &term.BookID, &term.Term, &term.Definition,
			&sourceSentence, &example, &chapterRef,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle NULL values
		if sourceSentence.Valid {
			term.SourceSentence = sourceSentence.String
		}
		if example.Valid {
			term.Example = example.String
		}
		if chapterRef.Valid {
			term.ChapterReference = chapterRef.String
		}
		
		// Parse timestamps (SQLite stores as string)
		if createdAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
				term.CreatedAt = t
			} else if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				term.CreatedAt = t
			}
		}
		if updatedAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
				term.UpdatedAt = t
			} else if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				term.UpdatedAt = t
			}
		}
		
		terms = append(terms, term)
	}
	return terms, rows.Err()
}

// GetGlossaryTermBySection gets glossary terms linked to a specific section
func GetGlossaryTermBySection(sectionID string) ([]*models.AliceGlossary, error) {
	query := `SELECT g.id, g.book_id, g.term, g.definition, g.source_sentence, g.example, g.chapter_reference, g.created_at, g.updated_at
	          FROM alice_glossary g
	          JOIN glossary_section_links gs ON g.id = gs.glossary_id
	          WHERE gs.section_id = ?
	          ORDER BY g.term`

	rows, err := DB.Query(query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	terms := []*models.AliceGlossary{}
	for rows.Next() {
		term := &models.AliceGlossary{}
		var sourceSentence, example, chapterRef sql.NullString
		var createdAt, updatedAt string
		
		err := rows.Scan(
			&term.ID, &term.BookID, &term.Term, &term.Definition,
			&sourceSentence, &example, &chapterRef,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle NULL values
		if sourceSentence.Valid {
			term.SourceSentence = sourceSentence.String
		}
		if example.Valid {
			term.Example = example.String
		}
		if chapterRef.Valid {
			term.ChapterReference = chapterRef.String
		}
		
		// Parse timestamps
		if createdAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
				term.CreatedAt = t
			} else if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				term.CreatedAt = t
			}
		}
		if updatedAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
				term.UpdatedAt = t
			} else if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				term.UpdatedAt = t
			}
		}
		
		terms = append(terms, term)
	}
	return terms, rows.Err()
}

// GetGlossaryTermByPageAndSection gets glossary terms for a specific page and section
func GetGlossaryTermByPageAndSection(bookID string, pageNumber, sectionNumber int) ([]*models.AliceGlossary, error) {
	query := `SELECT DISTINCT g.id, g.book_id, g.term, g.definition, g.source_sentence, g.example, g.chapter_reference, g.created_at, g.updated_at
	          FROM alice_glossary g
	          JOIN glossary_section_links gs ON g.id = gs.glossary_id
	          JOIN sections s ON gs.section_id = s.id
	          WHERE g.book_id = ? AND gs.page_number = ? AND gs.section_number = ?
	          ORDER BY g.term`

	rows, err := DB.Query(query, bookID, pageNumber, sectionNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	terms := []*models.AliceGlossary{}
	for rows.Next() {
		term := &models.AliceGlossary{}
		var sourceSentence, example, chapterRef sql.NullString
		var createdAt, updatedAt string
		
		err := rows.Scan(
			&term.ID, &term.BookID, &term.Term, &term.Definition,
			&sourceSentence, &example, &chapterRef,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle NULL values
		if sourceSentence.Valid {
			term.SourceSentence = sourceSentence.String
		}
		if example.Valid {
			term.Example = example.String
		}
		if chapterRef.Valid {
			term.ChapterReference = chapterRef.String
		}
		
		// Parse timestamps
		if createdAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
				term.CreatedAt = t
			} else if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				term.CreatedAt = t
			}
		}
		if updatedAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
				term.UpdatedAt = t
			} else if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				term.UpdatedAt = t
			}
		}
		
		terms = append(terms, term)
	}
	return terms, rows.Err()
}

// FindGlossaryTermInText finds if a word appears in glossary and returns the term
func FindGlossaryTermInText(bookID, word string) (*models.AliceGlossary, error) {
	// Try exact match first (case-insensitive)
	term, err := GetGlossaryTerm(bookID, strings.ToLower(word))
	if err == nil && term != nil {
		return term, nil
	}

	// Try case-insensitive search
	query := `SELECT id, book_id, term, definition, source_sentence, example, chapter_reference, created_at, updated_at
	          FROM alice_glossary 
	          WHERE book_id = ? AND LOWER(term) = LOWER(?)
	          LIMIT 1`

	term = &models.AliceGlossary{}
	var sourceSentence, example, chapterRef sql.NullString
	var createdAt, updatedAt string
	
	err = DB.QueryRow(query, bookID, word).Scan(
		&term.ID, &term.BookID, &term.Term, &term.Definition,
		&sourceSentence, &example, &chapterRef,
		&createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Handle NULL values
	if sourceSentence.Valid {
		term.SourceSentence = sourceSentence.String
	}
	if example.Valid {
		term.Example = example.String
	}
	if chapterRef.Valid {
		term.ChapterReference = chapterRef.String
	}
	
	// Parse timestamps
	if createdAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			term.CreatedAt = t
		} else if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			term.CreatedAt = t
		}
	}
	if updatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			term.UpdatedAt = t
		} else if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			term.UpdatedAt = t
		}
	}
	
	return term, nil
}

// Verification Code Queries

// VerifyCode checks if a verification code is valid and unused
func VerifyCode(code, bookID string) (*models.VerificationCode, error) {
	vc := &models.VerificationCode{}
	var usedBy sql.NullString
	query := `SELECT code, book_id, is_used, used_by, created_at
	          FROM verification_codes WHERE code = ? AND book_id = ?`

	err := DB.QueryRow(query, code, bookID).Scan(
		&vc.Code, &vc.BookID, &vc.IsUsed, &usedBy, &vc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if usedBy.Valid {
		vc.UsedBy = &usedBy.String
	}
	return vc, nil
}

// UseVerificationCode marks a verification code as used
func UseVerificationCode(code, userID string) error {
	query := `UPDATE verification_codes SET is_used = 1, used_by = ?, used_at = datetime('now') WHERE code = ?`
	_, err := DB.Exec(query, userID, code)
	return err
}

// Reading Progress Queries

// GetReadingProgress retrieves reading progress for a user and book
func GetReadingProgress(userID, bookID string) (*models.ReadingProgress, error) {
	progress := &models.ReadingProgress{}
	var chapterID, sectionID, purchaseDate sql.NullString
	var lastPage sql.NullInt64
	query := `SELECT id, user_id, book_id, chapter_id, section_id, last_page, last_read_at, purchase_date, created_at, updated_at
	          FROM reading_progress WHERE user_id = ? AND book_id = ?`

	err := DB.QueryRow(query, userID, bookID).Scan(
		&progress.ID, &progress.UserID, &progress.BookID, &chapterID, &sectionID,
		&lastPage, &progress.LastReadAt, &purchaseDate, &progress.CreatedAt, &progress.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if chapterID.Valid {
		progress.ChapterID = &chapterID.String
	}
	if sectionID.Valid {
		progress.SectionID = &sectionID.String
	}
	if lastPage.Valid {
		pageNum := int(lastPage.Int64)
		progress.LastPage = &pageNum
	}
	if purchaseDate.Valid {
		progress.PurchaseDate = &purchaseDate.String
	}
	return progress, nil
}

// UpdateReadingProgress updates reading progress
// IMPORTANT: If PurchaseDate is nil, it preserves the existing purchase_date in the database
// Only updates purchase_date if it's explicitly provided
func UpdateReadingProgress(progress *models.ReadingProgress) error {
	if progress.ID == "" {
		progress.ID = uuid.New().String()
		var purchaseDate interface{}
		if progress.PurchaseDate != nil {
			purchaseDate = *progress.PurchaseDate
		}
		query := `INSERT INTO reading_progress (id, user_id, book_id, chapter_id, section_id, last_page, last_read_at, purchase_date, created_at, updated_at)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := DB.Exec(query, progress.ID, progress.UserID, progress.BookID, progress.ChapterID, progress.SectionID,
			progress.LastPage, time.Now(), purchaseDate, time.Now(), time.Now())
		return err
	}

	// For updates, we need to preserve purchase_date if it's not explicitly provided
	// First, get the existing purchase_date from the database
	var existingPurchaseDate sql.NullString
	err := DB.QueryRow(`SELECT purchase_date FROM reading_progress WHERE id = ?`, progress.ID).Scan(&existingPurchaseDate)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get existing purchase_date: %w", err)
	}

	// Determine what purchase_date value to use
	var purchaseDate interface{}
	if progress.PurchaseDate != nil {
		// Explicitly provided - use it (even if empty string to clear it)
		if *progress.PurchaseDate == "" {
			purchaseDate = nil
		} else {
			purchaseDate = *progress.PurchaseDate
		}
	} else {
		// Not provided - preserve existing value
		if existingPurchaseDate.Valid {
			purchaseDate = existingPurchaseDate.String
		} else {
			purchaseDate = nil
		}
	}

	query := `UPDATE reading_progress SET chapter_id = ?, section_id = ?, last_page = ?, last_read_at = ?, purchase_date = ?, updated_at = ?
	          WHERE id = ?`
	resolvedAt := time.Now()
	_, err = DB.Exec(query, progress.ChapterID, progress.SectionID, progress.LastPage, resolvedAt, purchaseDate, resolvedAt, progress.ID)
	return err
}

// UpdateBookPurchaseDate updates the purchase date for a reader's book
func UpdateBookPurchaseDate(userID, bookID, purchaseDate string) error {
	// Handle empty string as NULL
	var purchaseDateValue interface{}
	if purchaseDate == "" {
		purchaseDateValue = nil
	} else {
		purchaseDateValue = purchaseDate
	}

	// First, check if reading_progress record exists
	var existingID string
	err := DB.QueryRow(`SELECT id FROM reading_progress WHERE user_id = ? AND book_id = ?`, userID, bookID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Create new reading_progress record if it doesn't exist
		progressID := uuid.New().String()
		query := `INSERT INTO reading_progress (id, user_id, book_id, purchase_date, created_at, updated_at)
		          VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))`
		_, err = DB.Exec(query, progressID, userID, bookID, purchaseDateValue)
		if err != nil {
			// Check if error is due to missing column
			if strings.Contains(err.Error(), "no such column: purchase_date") {
				return fmt.Errorf("purchase_date column does not exist. Please run migration 007_add_book_purchase_date.sql: %w", err)
			}
			return fmt.Errorf("failed to insert reading progress: %w", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to check reading progress: %w", err)
	}

	// Update existing record
	query := `UPDATE reading_progress SET purchase_date = ?, updated_at = datetime('now')
	          WHERE user_id = ? AND book_id = ?`
	_, err = DB.Exec(query, purchaseDateValue, userID, bookID)
	if err != nil {
		// Check if error is due to missing column
		if strings.Contains(err.Error(), "no such column: purchase_date") {
			return fmt.Errorf("purchase_date column does not exist. Please run migration 007_add_book_purchase_date.sql: %w", err)
		}
		return fmt.Errorf("failed to update purchase date: %w", err)
	}
	return nil
}

// Vocabulary Lookup Queries

// CreateVocabularyLookup creates a vocabulary lookup record
func CreateVocabularyLookup(lookup *models.VocabularyLookup) error {
	lookup.ID = uuid.New().String()
	query := `INSERT INTO vocabulary_lookups (id, user_id, book_id, word, definition, chapter_id, section_id, context, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, lookup.ID, lookup.UserID, lookup.BookID, lookup.Word, lookup.Definition,
		lookup.ChapterID, lookup.SectionID, lookup.Context, time.Now())
	return err
}

// GetVocabularyLookups retrieves vocabulary lookups for a user
func GetVocabularyLookups(userID, bookID string) ([]*models.VocabularyLookup, error) {
	query := `SELECT id, user_id, book_id, word, definition, chapter_id, section_id, context, created_at
	          FROM vocabulary_lookups WHERE user_id = ? AND book_id = ? ORDER BY created_at DESC`
	rows, err := DB.Query(query, userID, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lookups := []*models.VocabularyLookup{}
	for rows.Next() {
		lookup := &models.VocabularyLookup{}
		var chapterID, sectionID sql.NullString
		err := rows.Scan(
			&lookup.ID, &lookup.UserID, &lookup.BookID, &lookup.Word, &lookup.Definition,
			&chapterID, &sectionID, &lookup.Context, &lookup.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if chapterID.Valid {
			lookup.ChapterID = &chapterID.String
		}
		if sectionID.Valid {
			lookup.SectionID = &sectionID.String
		}
		lookups = append(lookups, lookup)
	}
	return lookups, rows.Err()
}

// AI Interaction Queries

// CreateAIInteraction creates an AI interaction record
func CreateAIInteraction(interaction *models.AIInteraction) error {
	interaction.ID = uuid.New().String()
	// Try with provider field first, fallback to old schema if column doesn't exist
	query := `INSERT INTO ai_interactions (id, user_id, book_id, section_id, interaction_type, question, prompt, response, context, provider, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, interaction.ID, interaction.UserID, interaction.BookID, interaction.SectionID,
		interaction.InteractionType, interaction.Question, interaction.Prompt, interaction.Response,
		interaction.Context, interaction.Provider, time.Now())
	if err != nil {
		// Fallback to old schema without provider field
		queryOld := `INSERT INTO ai_interactions (id, user_id, book_id, section_id, interaction_type, question, prompt, response, context, created_at)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = DB.Exec(queryOld, interaction.ID, interaction.UserID, interaction.BookID, interaction.SectionID,
			interaction.InteractionType, interaction.Question, interaction.Prompt, interaction.Response,
			interaction.Context, time.Now())
	}
	return err
}

// GetAIInteractions retrieves AI interactions for a user
func GetAIInteractions(userID, bookID string) ([]*models.AIInteraction, error) {
	// Try to select with provider field first
	query := `SELECT id, user_id, book_id, section_id, interaction_type, question, prompt, response, context, provider, created_at
	          FROM ai_interactions WHERE user_id = ? AND book_id = ? ORDER BY created_at DESC`
	rows, err := DB.Query(query, userID, bookID)
	hasProviderColumn := true
	
	if err != nil {
		// If that fails, try without provider column (old schema)
		queryOld := `SELECT id, user_id, book_id, section_id, interaction_type, question, prompt, response, context, created_at
		          FROM ai_interactions WHERE user_id = ? AND book_id = ? ORDER BY created_at DESC`
		rows, err = DB.Query(queryOld, userID, bookID)
		hasProviderColumn = false
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	interactions := []*models.AIInteraction{}
	for rows.Next() {
		interaction := &models.AIInteraction{}
		var sectionID sql.NullString
		var provider sql.NullString
		
		if hasProviderColumn {
			// Scan with provider field
			err := rows.Scan(
				&interaction.ID, &interaction.UserID, &interaction.BookID, &sectionID,
				&interaction.InteractionType, &interaction.Question, &interaction.Prompt,
				&interaction.Response, &interaction.Context, &provider, &interaction.CreatedAt,
			)
			if err != nil {
				return nil, err
			}
			if provider.Valid {
				interaction.Provider = provider.String
			} else {
				interaction.Provider = ""
			}
		} else {
			// Scan without provider field (old schema)
			err := rows.Scan(
				&interaction.ID, &interaction.UserID, &interaction.BookID, &sectionID,
				&interaction.InteractionType, &interaction.Question, &interaction.Prompt,
				&interaction.Response, &interaction.Context, &interaction.CreatedAt,
			)
			if err != nil {
				return nil, err
			}
			interaction.Provider = "" // Empty for old records
		}
		if sectionID.Valid {
			interaction.SectionID = &sectionID.String
		}
		interactions = append(interactions, interaction)
	}
	return interactions, rows.Err()
}

// Help Request Queries

// CreateHelpRequest creates a help request
func CreateHelpRequest(request *models.HelpRequest) error {
	request.ID = uuid.New().String()
	query := `INSERT INTO help_requests (id, user_id, book_id, section_id, status, content, context, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, request.ID, request.UserID, request.BookID, request.SectionID,
		request.Status, request.Content, request.Context, time.Now(), time.Now())
	return err
}

// GetHelpRequests retrieves help requests for a user
func GetHelpRequests(userID string) ([]*models.HelpRequest, error) {
	query := `SELECT id, user_id, book_id, section_id, status, content, context, assigned_to, response, resolved_at, created_at, updated_at
	          FROM help_requests WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []*models.HelpRequest{}
	for rows.Next() {
		request := &models.HelpRequest{}
		var sectionID, assignedTo, response sql.NullString
		var resolvedAt sql.NullTime
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&request.ID, &request.UserID, &request.BookID, &sectionID, &request.Status,
			&request.Content, &request.Context, &assignedTo, &response, &resolvedAt,
			&createdAtStr, &updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps from SQLite string format
		if createdAtStr != "" {
			request.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				return nil, err
			}
		}
		if updatedAtStr != "" {
			request.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
			if err != nil {
				return nil, err
			}
		}

		if sectionID.Valid {
			request.SectionID = &sectionID.String
		}
		if assignedTo.Valid {
			request.AssignedTo = &assignedTo.String
		}
		if response.Valid {
			request.Response = response.String
		}
		if resolvedAt.Valid {
			request.ResolvedAt = &resolvedAt.Time
		}
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

// GetHelpRequestByID retrieves a help request by ID
func GetHelpRequestByID(id string) (*models.HelpRequest, error) {
	query := `SELECT id, user_id, book_id, section_id, status, content, context, assigned_to, response, resolved_at, created_at, updated_at
	          FROM help_requests WHERE id = ?`

	request := &models.HelpRequest{}
	var sectionID, assignedTo, response sql.NullString
	var resolvedAt sql.NullTime

	err := DB.QueryRow(query, id).Scan(
		&request.ID, &request.UserID, &request.BookID, &sectionID, &request.Status,
		&request.Content, &request.Context, &assignedTo, &response, &resolvedAt,
		&request.CreatedAt, &request.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if sectionID.Valid {
		request.SectionID = &sectionID.String
	}
	if assignedTo.Valid {
		request.AssignedTo = &assignedTo.String
	}
	if response.Valid {
		request.Response = response.String
	}
	if resolvedAt.Valid {
		request.ResolvedAt = &resolvedAt.Time
	}

	return request, nil
}

// GetHelpRequestsByConsultant retrieves help requests assigned to a consultant
func GetHelpRequestsByConsultant(consultantID string) ([]*models.HelpRequest, error) {
	query := `SELECT id, user_id, book_id, section_id, status, content, context, assigned_to, response, resolved_at, created_at, updated_at
	          FROM help_requests WHERE assigned_to = ? ORDER BY created_at DESC`
	rows, err := DB.Query(query, consultantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []*models.HelpRequest{}
	for rows.Next() {
		request := &models.HelpRequest{}
		var sectionID, assignedTo, response sql.NullString
		var resolvedAt sql.NullTime
		err := rows.Scan(
			&request.ID, &request.UserID, &request.BookID, &sectionID, &request.Status,
			&request.Content, &request.Context, &assignedTo, &response, &resolvedAt,
			&request.CreatedAt, &request.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if sectionID.Valid {
			request.SectionID = &sectionID.String
		}
		if assignedTo.Valid {
			request.AssignedTo = &assignedTo.String
		}
		if response.Valid {
			request.Response = response.String
		}
		if resolvedAt.Valid {
			request.ResolvedAt = &resolvedAt.Time
		}
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

// UpdateHelpRequest updates a help request
func UpdateHelpRequest(request *models.HelpRequest) error {
	query := `UPDATE help_requests SET status = ?, assigned_to = ?, response = ?, resolved_at = ?, updated_at = ?
	          WHERE id = ?`
	var resolvedAt interface{}
	if request.ResolvedAt != nil {
		resolvedAt = request.ResolvedAt
	}
	_, err := DB.Exec(query, request.Status, request.AssignedTo, request.Response, resolvedAt, time.Now(), request.ID)
	return err
}

// Dictionary Cache Functions

// GetCachedDefinition retrieves a cached definition from dictionary_cache table
func GetCachedDefinition(word string) (*models.DictionaryCache, error) {
	normalizedWord := strings.ToLower(strings.TrimSpace(word))
	query := `SELECT id, word, definition, example, phonetic, part_of_speech, source_api, created_at, updated_at
	          FROM dictionary_cache WHERE word = ? LIMIT 1`

	cache := &models.DictionaryCache{}
	err := DB.QueryRow(query, normalizedWord).Scan(
		&cache.ID, &cache.Word, &cache.Definition, &cache.Example,
		&cache.Phonetic, &cache.PartOfSpeech, &cache.SourceAPI,
		&cache.CreatedAt, &cache.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// CacheDefinition stores a definition in the dictionary_cache table
func CacheDefinition(cache *models.DictionaryCache) error {
	normalizedWord := strings.ToLower(strings.TrimSpace(cache.Word))
	if cache.ID == "" {
		cache.ID = fmt.Sprintf("dict-cache-%s-%d", normalizedWord, time.Now().UnixNano())
	}

	query := `INSERT OR REPLACE INTO dictionary_cache 
	          (id, word, definition, example, phonetic, part_of_speech, source_api, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`

	_, err := DB.Exec(query, cache.ID, normalizedWord, cache.Definition, cache.Example,
		cache.Phonetic, cache.PartOfSpeech, cache.SourceAPI)
	return err
}
