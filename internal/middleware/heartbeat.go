package middleware

import (
	"net/http"
	"net/url"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// HeartbeatMiddleware updates last_active_at on every authenticated request
// This enables "who's online" queries for consultants and admin presence
func HeartbeatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header or cookie (same logic as auth middleware)
		authHeader := r.Header.Get("Authorization")
		token, err := auth.ExtractTokenFromHeader(authHeader)

		// If not in header, try cookie
		if err != nil || token == "" {
			cookie, cookieErr := r.Cookie("auth_token")
			if cookieErr != nil || cookie == nil || cookie.Value == "" {
				// Fallback: parse Cookie header manually (Safari sometimes doesn't parse correctly)
				if cookies := r.Header.Get("Cookie"); cookies != "" {
					for _, c := range parseCookies(cookies) {
						if c.Name == "auth_token" && c.Value != "" {
							cookie = c
							cookieErr = nil
							break
						}
					}
				}
			}
			if cookieErr == nil && cookie != nil && cookie.Value != "" {
				tokenValue := cookie.Value
				// URL-decode cookie value so it matches the token we hashed at login
				if decoded, decodeErr := url.QueryUnescape(tokenValue); decodeErr == nil && decoded != tokenValue {
					tokenValue = decoded
				}
				token = tokenValue
				err = nil
			}
		}

		// Update activity if token is valid (fire and forget)
		if err == nil && token != "" {
			// Get user from token
			user, err := auth.GetUserFromToken(token)
			if err == nil && user != nil {
				// Update session activity (fire and forget - don't block request)
				go func() {
					// Update session last_active_at
					database.UpdateSessionActivity(token)
					
					// Update users.last_active_at (handle column may not exist gracefully)
					database.DB.Exec(`UPDATE users SET last_active_at = datetime('now') WHERE id = ?`, user.ID)
					
					// Update reader_states.last_activity_at if reader
					if user.Role == "reader" {
						database.DB.Exec(`UPDATE reader_states SET last_activity_at = datetime('now'), status = 'active' WHERE user_id = ?`, user.ID)
					}
				}()
			}
		}

		next.ServeHTTP(w, r)
	})
}

