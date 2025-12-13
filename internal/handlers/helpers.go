package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/efisiopittau/alice-suite-go/internal/errors"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// Login wraps the auth Login function to match API standards
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
		errors.HandleError(w, errors.InternalError("login failed", err))
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		errors.HandleError(w, errors.InternalError("failed to generate token", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": token,
		"message":      "Login successful",
	})
}

// Register wraps the auth Register function to match API standards
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
		errors.HandleError(w, errors.InternalError("registration failed", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetBooks provides book listing for reader API
func GetBooks(w http.ResponseWriter, r *http.Request) {
	HandleBooks(w, r) // Reuse existing implementation
}

// GetChapters provides chapter listing for reader API
func GetChapters(w http.ResponseWriter, r *http.Request) {
	HandleChapters(w, r) // Reuse existing implementation
}

// GetSections provides section listing for reader API
func GetSections(w http.ResponseWriter, r *http.Request) {
	HandleSections(w, r) // Reuse existing implementation
}

// GetPage provides page content for reader API
func GetPage(w http.ResponseWriter, r *http.Request) {
	HandlePages(w, r) // Reuse existing implementation
}

// LookupWord provides dictionary lookup for reader API
func LookupWord(w http.ResponseWriter, r *http.Request) {
	HandleLookupWord(w, r) // Reuse existing implementation
}

// GetSectionGlossaryTerms provides glossary for reader API
func GetSectionGlossaryTerms(w http.ResponseWriter, r *http.Request) {
	HandleGetSectionGlossaryTerms(w, r) // Reuse existing implementation
}

// AskAI provides AI service for reader API
func AskAI(w http.ResponseWriter, r *http.Request) {
	HandleAskAI(w, r) // Reuse existing implementation
}

// CreateHelpRequest provides help system for reader API
func CreateHelpRequest(w http.ResponseWriter, r *http.Request) {
	HandleCreateHelpRequest(w, r) // Reuse existing implementation
}

// GetProgress provides reading progress for reader API
func GetProgress(w http.ResponseWriter, r *http.Request) {
	HandleReadingProgress(w, r) // Reuse existing implementation
}

// GetUserProfile provides user profile for reader API
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	HandleGetUserProfile(w, r) // Reuse existing implementation
}

// GetInteractions provides interactions for reader API
func GetInteractions(w http.ResponseWriter, r *http.Request) {
	HandleInteractions(w, r) // Reuse existing implementation
}

// TrackEvent provides activity tracking for reader API
func TrackEvent(w http.ResponseWriter, r *http.Request) {
	HandleTrackActivity(w, r) // Reuse existing implementation
}

