package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// HandleVerifyBookCode handles POST /rest/v1/rpc/verify-book-code
func HandleVerifyBookCode(w http.ResponseWriter, r *http.Request) {
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
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify book code with user_id from token
	bookID, err := auth.VerifyBookCode(req.Code, userID)
	if err != nil {
		if err == auth.ErrInvalidCode {
			http.Error(w, "Invalid verification code", http.StatusBadRequest)
			return
		}
		if err == auth.ErrCodeAlreadyUsed {
			http.Error(w, "Verification code already used", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"book_id": bookID,
	})
}

// HandleCheckBookVerified handles GET /rest/v1/rpc/check-book-verified
func HandleCheckBookVerified(w http.ResponseWriter, r *http.Request) {
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

	verified, err := auth.CheckBookVerified(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"verified": verified,
	})
}

