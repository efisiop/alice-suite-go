package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// SetupAuthRoutes sets up authentication-related routes
func SetupAuthRoutes(mux *http.ServeMux) {
	// Supabase-compatible auth endpoints
	mux.HandleFunc("/auth/v1/token", HandleLogin)
	mux.HandleFunc("/auth/v1/signup", HandleSignUp)
	mux.HandleFunc("/auth/v1/user", HandleGetUser)
	mux.HandleFunc("/auth/v1/logout", HandleLogout)

	// Alternative API endpoints
	mux.HandleFunc("/api/auth/login", HandleLogin)
	mux.HandleFunc("/api/auth/register", HandleSignUp)
}

// HandleLogin handles POST /auth/v1/token (Supabase-compatible)
func HandleLogin(w http.ResponseWriter, r *http.Request) {
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
			return
		}
		// Log the actual error for debugging
		log.Printf("Login error for %s: %v", req.Email, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Create database-backed session
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()
	expiresIn := 24 * time.Hour
	_, err = database.CreateSession(user.ID, token, ipAddress, userAgent, expiresIn)
	if err != nil {
		// Log error for debugging
		log.Printf("Warning: Failed to create database session for user %s: %v", user.ID, err)
		// Don't fail login (backward compatibility)
		// Session will still work via JWT validation
	}

	// Set cookie for server-side page navigation (more reliable than client-side)
	// Cookie expires in 24 hours (same as JWT)
	expiresAt := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false, // Set to false so JavaScript can also read it if needed
		Secure:   false, // Set to false for HTTP (Render uses HTTPS but we want it to work in both)
	}
	http.SetCookie(w, cookie)

	// Track login event and broadcast for consultants
	if user.Role == "reader" {
		// Track login activity in database (both old and new tables for compatibility)
		TrackActivity(user.ID, "LOGIN", "", nil)

		// Also log to new activity_logs table
		go func() {
			activity := &database.ActivityLog{
				UserID:       user.ID,
				ActivityType: "LOGIN",
				Metadata: map[string]interface{}{
					"ip_address": ipAddress,
					"user_agent": userAgent,
				},
			}
			database.LogActivity(activity)
		}()

		// Broadcast login event with user info
		BroadcastLogin(user.ID, user.Email, user.FirstName, user.LastName)
	}

	// Supabase-compatible response format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   86400, // 24 hours in seconds
		"expires_at":   expiresAt.Unix(),
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"aud":   "authenticated",
			"role":  "authenticated",
			"user_metadata": map[string]string{
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
		},
	})
}

// HandleSignUp handles POST /auth/v1/signup
func HandleSignUp(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

// HandleGetUser handles GET /auth/v1/user (get current user from token)
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// Validate token and get user
	user, err := auth.GetUserFromToken(token)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Supabase-compatible response format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"aud":   "authenticated",
		"role":  "authenticated",
		"user_metadata": map[string]string{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

// HandleLogout handles POST /auth/v1/logout
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔓 LOGOUT API called from %s", r.RemoteAddr)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header or cookie
	authHeader := r.Header.Get("Authorization")
	var token string

	if authHeader == "" {
		// Check for token in cookie
		cookie, cookieErr := r.Cookie("auth_token")
		if cookieErr == nil && cookie.Value != "" {
			token = cookie.Value
			log.Printf("🔓 Token found in cookie")
		}
	} else {
		var err error
		token, err = auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			log.Printf("🔓 Token extraction failed: %v", err)
			token = ""
		} else {
			log.Printf("🔓 Token found in Authorization header")
		}
	}

	if token == "" {
		log.Printf("🔓 No token provided for logout")
	}

	if token != "" {
		// Get user before deleting session
		user, err := auth.GetUserFromToken(token)
		if err == nil && user != nil {
			log.Printf("🔓 Logging out user: %s %s (ID: %s, Role: %s)", user.FirstName, user.LastName, user.ID, user.Role)

			if user.Role == "reader" {
				// Track logout activity in database (both old and new tables for compatibility)
				TrackActivity(user.ID, "LOGOUT", "", nil)

				// Log to activity_logs table (synchronous to ensure it's recorded)
				activity := &database.ActivityLog{
					UserID:       user.ID,
					ActivityType: "LOGOUT",
				}
				if err := database.LogActivity(activity); err != nil {
					log.Printf("❌ Failed to log LOGOUT activity for user %s: %v", user.ID, err)
				} else {
					log.Printf("✅ Logged LOGOUT activity for user %s", user.ID)
				}

				// Broadcast logout event for consultants
				log.Printf("📡 Broadcasting logout event for user %s", user.ID)
				BroadcastLogout(user.ID)
			}

			// Delete ALL sessions for this user (complete logout across all devices)
			log.Printf("🗑️ Deleting all sessions for user %s", user.ID)
			database.DeleteAllUserSessions(user.ID)
		} else {
			log.Printf("🔓 Could not get user from token, deleting session directly")
			// Fallback: delete just this session if we couldn't get the user
			database.DeleteSession(token)
		}
	}

	// Clear auth cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)

	log.Printf("✅ Logout completed successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
