package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/services"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

var (
	bookService      = services.NewBookService()
	dictionaryService = services.NewDictionaryService()
	aiService        = services.NewAIService()
	helpService      = services.NewHelpService()
)

// HealthCheck returns API health status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":  "ok",
		"message": "Alice Suite Reader API - Physical Book Companion",
		"version": "1.0.0",
		"scope":   "First 3 chapters test ground",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := auth.Login(req.Email, req.Password)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := auth.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if err == auth.ErrUserExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetBooks returns available books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	books, err := bookService.GetAllBooks()
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetChapters returns chapters for a book
func GetChapters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bookID := r.URL.Query().Get("book_id")
	if bookID == "" {
		http.Error(w, "book_id parameter required", http.StatusBadRequest)
		return
	}

	chapters, err := bookService.GetChapters(bookID)
	if err != nil {
		if err == services.ErrBookNotFound {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapters)
}

// GetSections returns sections for a chapter
func GetSections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chapterID := r.URL.Query().Get("chapter_id")
	if chapterID == "" {
		http.Error(w, "chapter_id parameter required", http.StatusBadRequest)
		return
	}

	sections, err := bookService.GetSections(chapterID)
	if err != nil {
		if err == services.ErrChapterNotFound {
			http.Error(w, "Chapter not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sections)
}

// GetPage returns a page with all its sections
func GetPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bookID := r.URL.Query().Get("book_id")
	pageNumberStr := r.URL.Query().Get("page_number")
	
	if bookID == "" || pageNumberStr == "" {
		http.Error(w, "book_id and page_number parameters required", http.StatusBadRequest)
		return
	}

	var pageNumber int
	if _, err := fmt.Sscanf(pageNumberStr, "%d", &pageNumber); err != nil {
		http.Error(w, "Invalid page_number: "+err.Error(), http.StatusBadRequest)
		return
	}

	page, err := bookService.GetPage(bookID, pageNumber)
	if err != nil {
		if err == services.ErrSectionNotFound {
			http.Error(w, "Page not found: "+err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if page == nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(page); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// LookupWord handles word definition lookup
func LookupWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookID    string  `json:"book_id"`
		Word      string  `json:"word"`
		ChapterID *string `json:"chapter_id"`
		SectionID *string `json:"section_id"`
		Context   string  `json:"context"`
		UserID    string  `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lookup word
	term, err := dictionaryService.LookupWordInContext(req.BookID, req.Word, req.ChapterID, req.SectionID)
	if err != nil && err != services.ErrTermNotFound {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Record lookup if user ID provided
	if req.UserID != "" {
		definition := ""
		if term != nil {
			definition = term.Definition
		}
		dictionaryService.RecordLookup(req.UserID, req.BookID, req.Word, definition, req.ChapterID, req.SectionID, req.Context)
	}

	w.Header().Set("Content-Type", "application/json")
	if term == nil {
		json.NewEncoder(w).Encode(map[string]string{
			"word":       req.Word,
			"definition": "Word not found in glossary",
		})
		return
	}
	json.NewEncoder(w).Encode(term)
}

// AskAI handles AI assistance requests
func AskAI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID         string  `json:"user_id"`
		BookID         string  `json:"book_id"`
		InteractionType string `json:"interaction_type"`
		Question       string  `json:"question"`
		SectionID      *string `json:"section_id"`
		Context        string  `json:"context"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	interactionType := services.InteractionType(strings.ToLower(req.InteractionType))
	if interactionType == "" {
		interactionType = services.InteractionChat
	}

	interaction, err := aiService.AskAI(req.UserID, req.BookID, interactionType, req.Question, req.SectionID, req.Context)
	if err != nil {
		if err == services.ErrAIServiceUnavailable {
			http.Error(w, "AI service unavailable", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interaction)
}

// CreateHelpRequest handles help request creation
func CreateHelpRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID    string  `json:"user_id"`
		BookID    string  `json:"book_id"`
		Content   string  `json:"content"`
		Context   string  `json:"context"`
		SectionID *string `json:"section_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	request, err := helpService.CreateHelpRequest(req.UserID, req.BookID, req.Content, req.Context, req.SectionID)
	if err != nil {
		errorMsg := fmt.Sprintf("Error creating help request: %v", err)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		fmt.Printf("[ERROR] CreateHelpRequest: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

// GetProgress returns reading progress
func GetProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	bookID := r.URL.Query().Get("book_id")

	if userID == "" || bookID == "" {
		http.Error(w, "user_id and book_id parameters required", http.StatusBadRequest)
		return
	}

	// This would use a progress service - for now, return placeholder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Progress tracking - to be fully implemented",
	})
}

// GetSectionGlossaryTerms returns glossary terms for a section
func GetSectionGlossaryTerms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract section ID from URL path: /api/dictionary/section/{section_id}/terms
	path := r.URL.Path
	sectionID := strings.TrimPrefix(path, "/api/dictionary/section/")
	sectionID = strings.TrimSuffix(sectionID, "/terms")

	if sectionID == "" {
		http.Error(w, "section_id required", http.StatusBadRequest)
		return
	}

	terms, err := dictionaryService.GetGlossaryTermsForSection(sectionID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(terms)
}

// GetUserProfile returns user profile information
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a real app, you would extract this from authentication token
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id parameter required", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Masks sensitive information for profile endpoint
	response := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetInteractions tracks AI interactions (help requests)
func GetInteractions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	bookID := r.URL.Query().Get("book_id")

	if userID == "" || bookID == "" {
		http.Error(w, "user_id and book_id parameters required", http.StatusBadRequest)
		return
	}

	// This endpoint returns AI interactions as "interactions"
	interactions, err := database.GetAIInteractions(userID, bookID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interactions)
}

// TrackEvent tracks activity events
func TrackEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID    string `json:"user_id"`
		BookID    string `json:"book_id"`
		EventType string `json:"event_type"`
		EventData string `json:"event_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For now, just return success - this is a basic tracking implementation
	// In real implementation, store event in database
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "event_tracked",
		"type":   req.EventType,
	})
}
