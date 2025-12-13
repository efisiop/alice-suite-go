package models

import "time"

// User represents a user in the system (reader or consultant)
type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	PasswordHash string   `json:"-"` // Never return in JSON
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Role        string    `json:"role"` // "reader" or "consultant"
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Book represents a book in the system
type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	TotalPages  int       `json:"total_pages"`
	CreatedAt   time.Time `json:"created_at"`
}

// Chapter represents a chapter in a book
type Chapter struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	Title     string    `json:"title"`
	Number    int       `json:"number"`
	CreatedAt time.Time `json:"created_at"`
}

// Page represents a page in the physical book
type Page struct {
	ID           string     `json:"id"`
	BookID       string     `json:"book_id"`
	PageNumber   int        `json:"page_number"`
	ChapterID    *string    `json:"chapter_id,omitempty"`
	ChapterTitle *string    `json:"chapter_title,omitempty"`
	Content      string     `json:"content"`
	WordCount    int        `json:"word_count"`
	Sections     []Section  `json:"sections"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Section represents a section within a page (for reference/word clarification)
type Section struct {
	ID            string    `json:"id"`
	PageID        string    `json:"page_id"`
	PageNumber    int       `json:"page_number"`
	SectionNumber int       `json:"section_number"`
	Content       string    `json:"content"`
	WordCount     int       `json:"word_count"`
	CreatedAt     time.Time `json:"created_at"`
}

// AliceGlossary represents Alice-specific glossary terms
type AliceGlossary struct {
	ID              string    `json:"id"`
	BookID          string    `json:"book_id"`
	Term            string    `json:"term"`
	Definition      string    `json:"definition"`
	SourceSentence  string    `json:"source_sentence"`
	Example         string    `json:"example"`
	ChapterReference string   `json:"chapter_reference"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// VerificationCode represents a book verification code
type VerificationCode struct {
	Code    string    `json:"code"`
	BookID  string    `json:"book_id"`
	IsUsed  bool      `json:"is_used"`
	UsedBy  *string   `json:"used_by"`
	CreatedAt time.Time `json:"created_at"`
}

// ReadingProgress tracks user's progress in the physical book
type ReadingProgress struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	BookID    string    `json:"book_id"`
	ChapterID *string   `json:"chapter_id"`
	SectionID *string   `json:"section_id"`
	LastPage  *int      `json:"last_page"`
	LastReadAt time.Time `json:"last_read_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VocabularyLookup represents a word looked up by a user
type VocabularyLookup struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	BookID    string    `json:"book_id"`
	Word      string    `json:"word"`
	Definition string   `json:"definition"`
	ChapterID *string   `json:"chapter_id"`
	SectionID *string   `json:"section_id"`
	Context   string    `json:"context"`
	CreatedAt time.Time `json:"created_at"`
}

// AIInteraction represents an AI assistance interaction
type AIInteraction struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	BookID         string    `json:"book_id"`
	SectionID      *string   `json:"section_id"`
	InteractionType string   `json:"interaction_type"` // "explain", "quiz", "simplify", "definition", "chat"
	Question       string    `json:"question"`
	Prompt         string    `json:"prompt"`
	Response       string    `json:"response"`
	Context        string    `json:"context"`
	CreatedAt      time.Time `json:"created_at"`
}

// HelpRequest represents a help request from reader to consultant
type HelpRequest struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	BookID     string     `json:"book_id"`
	SectionID  *string    `json:"section_id"`
	Status     string     `json:"status"` // "pending", "assigned", "resolved"
	Content    string     `json:"content"`
	Context    string     `json:"context"`
	AssignedTo *string    `json:"assigned_to"`
	Response   string     `json:"response"`
	ResolvedAt *time.Time `json:"resolved_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ConsultantAssignment represents assignment of reader to consultant
type ConsultantAssignment struct {
	ID           string    `json:"id"`
	ConsultantID string    `json:"consultant_id"`
	UserID       string    `json:"user_id"`
	BookID       string    `json:"book_id"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

// ConsultantTrigger represents a prompt sent by consultant to reader
type ConsultantTrigger struct {
	ID           string     `json:"id"`
	ConsultantID *string    `json:"consultant_id"`
	UserID       string     `json:"user_id"`
	BookID       string     `json:"book_id"`
	TriggerType  string     `json:"trigger_type"`
	Message      string     `json:"message"`
	IsProcessed  bool       `json:"is_processed"`
	ProcessedAt  *time.Time `json:"processed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ReadingStats represents reading statistics for a user
type ReadingStats struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	BookID          string    `json:"book_id"`
	TotalReadingTime int      `json:"total_reading_time"` // in seconds
	PagesRead       int       `json:"pages_read"`
	VocabularyWords int       `json:"vocabulary_words"`
	LastSessionDate time.Time `json:"last_session_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}



