package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// RequireAuth validates JWT token and requires authentication
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header or cookie
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Check for token in cookie as fallback
			cookie, err := r.Cookie("auth_token")
			if err != nil || cookie == nil || cookie.Value == "" {
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}
			// Safari may URL-encode cookie values, so decode if needed
			tokenValue := cookie.Value
			if decoded, err := url.QueryUnescape(tokenValue); err == nil && decoded != tokenValue {
				tokenValue = decoded
			}
			authHeader = "Bearer " + tokenValue
		}

		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// Validate token
		_, err = auth.ValidateJWT(token)
		if err != nil {
			if err == auth.ErrInvalidToken || err == auth.ErrExpiredToken {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		// Token is valid, continue
		next.ServeHTTP(w, r)
	})
}

// RequireRole requires a specific role (reader or consultant)
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header or cookie
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Check for token in cookie as fallback
				cookie, err := r.Cookie("auth_token")
				if err != nil || cookie == nil || cookie.Value == "" {
					// Try reading raw cookies as fallback (Safari sometimes doesn't parse cookies correctly)
					if cookies := r.Header.Get("Cookie"); cookies != "" {
						// Parse cookies manually
						for _, c := range parseCookies(cookies) {
							if c.Name == "auth_token" && c.Value != "" {
								cookie = c
								err = nil
								break
							}
						}
					}
					if err != nil || cookie == nil || cookie.Value == "" {
						// For consultant routes, redirect to login instead of showing error
						if requiredRole == "consultant" {
							http.Redirect(w, r, "/consultant/login", http.StatusFound)
							return
						}
						http.Error(w, "Authorization required", http.StatusUnauthorized)
						return
					}
				}
				// Safari may URL-encode cookie values, so decode if needed
				tokenValue := cookie.Value
				// Try URL decoding (Safari sometimes encodes cookies)
				if decoded, err := url.QueryUnescape(tokenValue); err == nil && decoded != tokenValue {
					tokenValue = decoded
				}
				authHeader = "Bearer " + tokenValue
			}

			token, err := auth.ExtractTokenFromHeader(authHeader)
			if err != nil {
				// For consultant routes, redirect to login instead of showing error
				if requiredRole == "consultant" {
					http.Redirect(w, r, "/consultant/login", http.StatusFound)
					return
				}
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			// Validate token and get claims
			claims, err := auth.ValidateJWT(token)
			if err != nil {
				// For consultant routes, redirect to login instead of showing error
				if requiredRole == "consultant" {
					http.Redirect(w, r, "/consultant/login", http.StatusFound)
					return
				}
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Check role
			if !auth.RequireRole(claims.Role, requiredRole) {
				// For consultant routes, redirect to login if wrong role
				if requiredRole == "consultant" {
					http.Redirect(w, r, "/consultant/login", http.StatusFound)
					return
				}
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// Role check passed, continue
			next.ServeHTTP(w, r)
		})
	}
}

// parseCookies manually parses cookie header string
func parseCookies(cookieHeader string) []*http.Cookie {
	var cookies []*http.Cookie
	for _, cookieStr := range splitCookies(cookieHeader) {
		parts := splitCookie(cookieStr, '=')
		if len(parts) == 2 {
			cookies = append(cookies, &http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}
	return cookies
}

func splitCookies(s string) []string {
	var cookies []string
	var current strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == ';' {
			if current.Len() > 0 {
				cookies = append(cookies, strings.TrimSpace(current.String()))
				current.Reset()
			}
		} else {
			current.WriteByte(s[i])
		}
	}
	if current.Len() > 0 {
		cookies = append(cookies, strings.TrimSpace(current.String()))
	}
	return cookies
}

func splitCookie(s string, sep byte) []string {
	idx := strings.IndexByte(s, sep)
	if idx == -1 {
		return []string{s, ""}
	}
	return []string{strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+1:])}
}

// RequireConsultant requires consultant role
func RequireConsultant(next http.Handler) http.Handler {
	return RequireRole("consultant")(next)
}

// RequireReader requires reader role
func RequireReader(next http.Handler) http.Handler {
	return RequireRole("reader")(next)
}

