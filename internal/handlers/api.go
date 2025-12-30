package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/middleware"
	"github.com/efisiopittau/alice-suite-go/internal/services"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// Service instances
var (
	bookService       = services.NewBookService()
	dictionaryService = services.NewDictionaryService()
	helpService       = services.NewHelpService()
	aiService         = services.NewAIService()
)

// SetupAPIRoutes sets up REST API routes (Supabase-compatible)
func SetupAPIRoutes(mux *http.ServeMux) {
	// Generic REST endpoint for all tables
	// This will catch /rest/v1/:table and /rest/v1/:table/ paths
	mux.HandleFunc("/rest/v1/", HandleRESTTable)
	// Books API
	mux.HandleFunc("/rest/v1/books", HandleBooks)
	mux.HandleFunc("/rest/v1/chapters", HandleChapters)
	mux.HandleFunc("/rest/v1/sections", HandleSections)
	mux.HandleFunc("/rest/v1/pages", HandlePages)

	// Reading progress API
	mux.HandleFunc("/rest/v1/reading_progress", HandleReadingProgress)
	mux.HandleFunc("/rest/v1/reading_stats", HandleReadingStats)

	// Dictionary/Glossary API
	mux.HandleFunc("/rest/v1/alice_glossary", HandleGlossaryTerms)
	
	// RPC functions
	mux.HandleFunc("/rest/v1/rpc/", HandleRPC)
	
	// Server-Sent Events for real-time updates
	mux.HandleFunc("/api/realtime/events", HandleSSE)
	
	// WebSocket for bidirectional communication (optional)
	mux.HandleFunc("/api/realtime/ws", HandleWebSocket)
	
	// Activity tracking
	mux.HandleFunc("/api/activity/track", HandleTrackActivity)

	// Consultant reader activity endpoints (protected with authentication middleware)
	mux.Handle("/api/consultant/reader-activities", middleware.RequireConsultant(http.HandlerFunc(HandleGetReaderActivities)))
	mux.Handle("/api/consultant/reader-activities/stream", middleware.RequireConsultant(http.HandlerFunc(HandleGetReaderActivityStream)))
	mux.Handle("/api/consultant/active-readers-count", middleware.RequireConsultant(http.HandlerFunc(HandleGetActiveReadersCount)))
	mux.Handle("/api/consultant/logged-in-readers-count", middleware.RequireConsultant(http.HandlerFunc(HandleGetLoggedInReadersCount)))
	mux.Handle("/api/consultant/logged-out-count", middleware.RequireConsultant(http.HandlerFunc(HandleGetLoggedOutCount)))
	mux.Handle("/api/consultant/todays-activity-count", middleware.RequireConsultant(http.HandlerFunc(HandleGetTodaysActivityCount)))

	// New consultant dashboard endpoints (using new activity_logs table)
	mux.Handle("/api/consultant/active-readers", middleware.RequireConsultant(http.HandlerFunc(HandleConsultantActiveReaders)))
	mux.Handle("/api/consultant/recent-activities", middleware.RequireConsultant(http.HandlerFunc(HandleConsultantRecentActivities)))
	mux.Handle("/api/consultant/reader/activity", middleware.RequireConsultant(http.HandlerFunc(HandleConsultantReaderActivity)))
	mux.Handle("/api/consultant/reader/state", middleware.RequireConsultant(http.HandlerFunc(HandleConsultantReaderState)))

	// Help requests API
	mux.HandleFunc("/rest/v1/help_requests", HandleHelpRequests)
	mux.Handle("/api/consultant/help-requests/", middleware.RequireConsultant(http.HandlerFunc(HandleGetHelpRequestByID)))

	// Interactions API
	mux.HandleFunc("/rest/v1/interactions", HandleInteractions)

	// Profiles API
	mux.HandleFunc("/rest/v1/profiles", HandleProfiles)
	mux.HandleFunc("/rest/v1/user", HandleGetUserProfile)

	// Verification codes API
	mux.HandleFunc("/rest/v1/verification_codes", HandleVerificationCodes)
	mux.HandleFunc("/rest/v1/rpc/verify-book-code", HandleVerifyBookCode)
	mux.HandleFunc("/rest/v1/rpc/check-book-verified", HandleCheckBookVerified)

	// Alternative API endpoints (for compatibility)
	mux.HandleFunc("/api/books", HandleBooks)
	mux.HandleFunc("/api/dictionary/lookup", HandleLookupWord)
	mux.HandleFunc("/api/dictionary/section/", HandleGetSectionGlossaryTerms)
	mux.HandleFunc("/api/ai/ask", HandleAskAI)
	mux.HandleFunc("/api/help", HandleCreateHelpRequest)
}

// HandleBooks handles GET /rest/v1/books
func HandleBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	books, err := bookService.GetAllBooks()
	if err != nil {
		log.Printf("Internal error in HandleBooks: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// HandleChapters handles GET /rest/v1/chapters
func HandleChapters(w http.ResponseWriter, r *http.Request) {
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

// HandleSections handles GET /rest/v1/sections
func HandleSections(w http.ResponseWriter, r *http.Request) {
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

// HandlePages handles GET /rest/v1/pages
func HandlePages(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid page_number", http.StatusBadRequest)
		return
	}

	page, err := bookService.GetPage(bookID, pageNumber)
	if err != nil {
		if err == services.ErrSectionNotFound {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		log.Printf("Internal error in HandlePages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if page == nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

// HandleReadingProgress handles GET/POST /rest/v1/reading_progress
func HandleReadingProgress(w http.ResponseWriter, r *http.Request) {
	// Extract and validate token to get user_id (SECURITY: Never trust user_id from query/body)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Extract user_id from token (not from query parameter or request body)
	userID := claims.UserID

	switch r.Method {
	case http.MethodGet:
		bookID := r.URL.Query().Get("book_id")
		if bookID == "" {
			http.Error(w, "book_id parameter required", http.StatusBadRequest)
			return
		}
		// TODO: Implement reading progress retrieval
		// Use userID from token, bookID from query
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Reading progress - to be implemented",
			"user_id": userID,
			"book_id": bookID,
		})

	case http.MethodPost, http.MethodPut:
		var req struct {
			BookID    string `json:"book_id"`
			SectionID string `json:"section_id"`
			Position  string `json:"last_position"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// TODO: Implement reading progress save
		// Use userID from token, req.BookID from body
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "saved"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleReadingStats handles GET /rest/v1/reading_stats
func HandleReadingStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract and validate token to get user_id (SECURITY: Never trust user_id from query parameter)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Extract user_id from token (not from query parameter)
	userID := claims.UserID

	bookID := r.URL.Query().Get("book_id")
	if bookID == "" {
		http.Error(w, "book_id parameter required", http.StatusBadRequest)
		return
	}

	// TODO: Implement reading stats retrieval
	// Use userID from token, bookID from query
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reading stats - to be implemented",
		"user_id": userID,
		"book_id": bookID,
	})
}

// HandleGlossaryTerms handles GET /rest/v1/alice_glossary
func HandleGlossaryTerms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sectionID := r.URL.Query().Get("section_id")

	if sectionID != "" {
		terms, err := dictionaryService.GetGlossaryTermsForSection(sectionID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(terms)
		return
	}

	// TODO: Implement book-level glossary terms
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// HandleGetDefinitionWithContext handles POST /rest/v1/rpc/get_definition_with_context
func HandleGetDefinitionWithContext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Term      string  `json:"term"`
		BookID    string  `json:"book_id"`
		SectionID *string `json:"section_id"`
		Context   string  `json:"context"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	term, err := dictionaryService.LookupWordInContext(req.BookID, req.Term, nil, req.SectionID)
	if err != nil && err != services.ErrTermNotFound {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if term == nil {
		json.NewEncoder(w).Encode(map[string]string{
			"term":       req.Term,
			"definition": "Word not found in glossary",
		})
		return
	}
	json.NewEncoder(w).Encode(term)
}

// HandleGetSectionsForPage handles POST /rest/v1/rpc/get_sections_for_page
func HandleGetSectionsForPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookID     string `json:"book_id"`
		PageNumber int    `json:"page_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	page, err := bookService.GetPage(req.BookID, req.PageNumber)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

// HandleHelpRequests handles GET/POST /rest/v1/help_requests
// PATCH, PUT, DELETE are forwarded to the generic REST handler
func HandleHelpRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// For GET requests, use the generic REST handler to support query parameters
		// This allows filtering like ?user_id=eq.{id}&order=created_at.desc
		HandleRESTTable(w, r)
		return

	case http.MethodPost:
		// Extract and validate token to get user_id (SECURITY: Never trust user_id from request body)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		// Extract user_id from token (not from request body)
		userID := claims.UserID

		var req struct {
			BookID    string  `json:"book_id"`
			Content   string  `json:"content"`
			Context   string  `json:"context"`
			SectionID *string `json:"section_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Invalid request body in HandleHelpRequests: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		request, err := helpService.CreateHelpRequest(userID, req.BookID, req.Content, req.Context, req.SectionID)
		if err != nil {
			log.Printf("Error creating help request: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(request)

	case http.MethodPatch, http.MethodPut, http.MethodDelete:
		// Forward PATCH, PUT, DELETE to generic REST handler
		HandleRESTTable(w, r)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleGetHelpRequestByID handles GET /api/consultant/help-requests/:id
// Returns a single help request by ID (consultant-only)
func HandleGetHelpRequestByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	// Path will be /api/consultant/help-requests/:id
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	requestID := pathParts[3]

	if requestID == "" {
		http.Error(w, "Help request ID required", http.StatusBadRequest)
		return
	}

	// Get help request from database
	request, err := helpService.GetHelpRequestByID(requestID)
	if err != nil {
		log.Printf("Error fetching help request %s: %v", requestID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if request == nil {
		http.Error(w, "Help request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(request)
}

// HandleInteractions handles GET/POST /rest/v1/interactions
func HandleInteractions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Implement interaction retrieval
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})

	case http.MethodPost:
		// TODO: Implement interaction tracking
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "tracked"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProfiles handles GET /rest/v1/profiles
func HandleProfiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement profile retrieval
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// HandleGetUserProfile handles GET /rest/v1/user (legacy endpoint)
func HandleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Redirect to auth handler
	HandleGetUser(w, r)
}

// HandleVerificationCodes handles GET /rest/v1/verification_codes
func HandleVerificationCodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement verification code lookup
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// HandleVerifyBookCode is now in verification.go

// Legacy handlers from handlers.go (kept for compatibility)
func HandleLookupWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract and validate token to get user_id (SECURITY: Never trust user_id from request body)
	authHeader := r.Header.Get("Authorization")
	var userID string
	if authHeader != "" {
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err == nil {
			claims, err := auth.ValidateJWT(token)
			if err == nil {
				userID = claims.UserID
			}
		}
	}

	var req struct {
		BookID    string  `json:"book_id"`
		Word      string  `json:"word"`
		ChapterID *string `json:"chapter_id"`
		SectionID *string `json:"section_id"`
		Context   string  `json:"context"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	term, err := dictionaryService.LookupWordInContext(req.BookID, req.Word, req.ChapterID, req.SectionID)
	if err != nil && err != services.ErrTermNotFound {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Record lookup if user is authenticated (user_id from token, not request body)
	if userID != "" {
		definition := ""
		if term != nil {
			definition = term.Definition
		}
		dictionaryService.RecordLookup(userID, req.BookID, req.Word, definition, req.ChapterID, req.SectionID, req.Context)
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

func HandleAskAI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract and validate token to get user_id (SECURITY: Never trust user_id from request body)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Extract user_id from token (not from request body)
	userID := claims.UserID

	var req struct {
		BookID          string  `json:"book_id"`
		InteractionType string  `json:"interaction_type"`
		Question        string  `json:"question"`
		SectionID       *string `json:"section_id"`
		Context         string  `json:"context"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	interactionType := services.InteractionType(strings.ToLower(req.InteractionType))
	if interactionType == "" {
		interactionType = services.InteractionChat
	}

	interaction, err := aiService.AskAI(userID, req.BookID, interactionType, req.Question, req.SectionID, req.Context)
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

func HandleCreateHelpRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract and validate token to get user_id (SECURITY: Never trust user_id from request body)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Extract user_id from token (not from request body)
	userID := claims.UserID

	var req struct {
		BookID    string  `json:"book_id"`
		Content   string  `json:"content"`
		Context   string  `json:"context"`
		SectionID *string `json:"section_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body in HandleCreateHelpRequest: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	request, err := helpService.CreateHelpRequest(userID, req.BookID, req.Content, req.Context, req.SectionID)
	if err != nil {
		log.Printf("Error creating help request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

func HandleGetSectionGlossaryTerms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
