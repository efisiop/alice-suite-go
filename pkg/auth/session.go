package auth

import (
	"sync"
	"time"
)

// Session represents a user session
type Session struct {
	UserID    string
	Email     string
	Role      string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// SessionStore manages active sessions
type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

var (
	globalSessionStore *SessionStore
	once                sync.Once
)

// GetSessionStore returns the global session store
func GetSessionStore() *SessionStore {
	once.Do(func() {
		globalSessionStore = &SessionStore{
			sessions: make(map[string]*Session),
		}
		// Start cleanup goroutine
		go globalSessionStore.cleanupExpiredSessions()
	})
	return globalSessionStore
}

// CreateSession creates a new session
func (s *SessionStore) CreateSession(userID, email, role, token string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &Session{
		UserID:    userID,
		Email:     email,
		Role:      role,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	s.sessions[token] = session
	return session
}

// GetSession retrieves a session by token
func (s *SessionStore) GetSession(token string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return nil, false
	}

	// Check if session expired
	if time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}

// DeleteSession removes a session
func (s *SessionStore) DeleteSession(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, token)
}

// cleanupExpiredSessions periodically removes expired sessions
func (s *SessionStore) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

// IsConsultant checks if a user is a consultant
func IsConsultant(role string) bool {
	return role == "consultant"
}

// IsReader checks if a user is a reader
func IsReader(role string) bool {
	return role == "reader"
}

// RequireRole checks if a user has the required role
func RequireRole(userRole, requiredRole string) bool {
	if requiredRole == "consultant" {
		return IsConsultant(userRole)
	}
	if requiredRole == "reader" {
		return IsReader(userRole)
	}
	// Admin role can be added later
	return false
}

