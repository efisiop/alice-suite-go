package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/email"
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
		// Track login activity once; TrackActivity writes the consultant dashboard
		// tables and emits the real-time activity event.
		if err := TrackActivity(user.ID, "LOGIN", "", map[string]interface{}{
			"ip_address": ipAddress,
			"user_agent": userAgent,
		}); err != nil {
			log.Printf("failed to track LOGIN activity for user %s: %v", user.ID, err)
		}

		// Broadcast login event with user info
		BroadcastLogin(user.ID, user.Email, user.FirstName, user.LastName)
	}

	// Optional: email admin when a reader or consultant logs in (only on Render when enabled)
	if user.Role == "reader" || user.Role == "consultant" {
		go func() {
			enabled, err := database.GetLoginEmailNotificationsEnabled()
			if err != nil {
				log.Printf("email: login notification skipped (reader/consultant login): failed to get setting: %v", err)
				return
			}
			if !enabled {
				return // toggle off, no log needed
			}
			if !email.ShouldSendLoginEmails() {
				log.Printf("email: login notification skipped: not on Render or SMTP not configured (RENDER=%q, SMTP_HOST set=%v)", os.Getenv("RENDER"), os.Getenv("SMTP_HOST") != "")
				return
			}
			to, _ := database.GetAdminSetting("login_email_recipient")
			if to == "" {
				to = "efisio@mylivemail.net"
			}
			name := strings.TrimSpace(user.FirstName + " " + user.LastName)
			if name == "" {
				name = user.Email
			}
			log.Printf("email: sending login notification to %s for %s %s", to, user.Role, user.Email)
			email.SendLoginNotification(to, user.Role, name, user.Email)
		}()
	}

	// Supabase-compatible response format
	preferredLanguageCode := database.DefaultReaderLanguageCode
	if user.Role == "reader" {
		if pref, err := database.GetReaderPreference(user.ID); err == nil && pref != nil {
			preferredLanguageCode = pref.PreferredLanguageCode
		}
	}
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
			"role":  user.Role,
			"user_metadata": map[string]string{
				"first_name":              user.FirstName,
				"last_name":               user.LastName,
				"preferred_language_code": preferredLanguageCode,
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
		Email                 string `json:"email"`
		Password              string `json:"password"`
		FirstName             string `json:"first_name"`
		LastName              string `json:"last_name"`
		PreferredLanguageCode string `json:"preferred_language_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := auth.Register(req.Email, req.Password, req.FirstName, req.LastName, req.PreferredLanguageCode)
	if err != nil {
		if err == auth.ErrUserExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		log.Printf("Registration error for %s: %v", req.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Registration continues directly to book verification, so establish the
	// same reader session that a normal login would create.
	token, err := auth.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	if _, err := database.CreateSession(user.ID, token, r.RemoteAddr, r.UserAgent(), 24*time.Hour); err != nil {
		log.Printf("Warning: Failed to create registration session for user %s: %v", user.ID, err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
		Secure:   false,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   86400,
		"expires_at":   expiresAt.Unix(),
		"user": map[string]interface{}{
			"id":                      user.ID,
			"email":                   user.Email,
			"first_name":              user.FirstName,
			"last_name":               user.LastName,
			"preferred_language_code": database.NormalizeReaderLanguageCode(req.PreferredLanguageCode),
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
	preferredLanguageCode := database.DefaultReaderLanguageCode
	if user.Role == "reader" {
		if pref, err := database.GetReaderPreference(user.ID); err == nil && pref != nil {
			preferredLanguageCode = pref.PreferredLanguageCode
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"aud":   "authenticated",
		"role":  user.Role,
		"user_metadata": map[string]string{
			"first_name":              user.FirstName,
			"last_name":               user.LastName,
			"preferred_language_code": preferredLanguageCode,
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
				// Track logout activity once; TrackActivity writes the consultant dashboard
				// tables and emits the real-time activity event.
				if err := TrackActivity(user.ID, "LOGOUT", "", nil); err != nil {
					log.Printf("❌ Failed to track LOGOUT activity for user %s: %v", user.ID, err)
				} else {
					log.Printf("✅ Tracked LOGOUT activity for user %s", user.ID)
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
