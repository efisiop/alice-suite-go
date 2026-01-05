package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Session represents a database-backed session
type Session struct {
	ID           string
	UserID       string
	TokenHash    string
	IPAddress    string
	UserAgent    string
	CreatedAt    time.Time
	LastActiveAt time.Time
	ExpiresAt    time.Time
}

// CreateSession creates a new session in the database
func CreateSession(userID, token, ipAddress, userAgent string, expiresIn time.Duration) (*Session, error) {
	sessionID := uuid.New().String()
	tokenHash := hashToken(token)
	expiresAt := time.Now().Add(expiresIn)

	query := `INSERT INTO sessions (id, user_id, token_hash, ip_address, user_agent, created_at, last_active_at, expires_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	_, err := DB.Exec(query, sessionID, userID, tokenHash, ipAddress, userAgent, now, now, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Session{
		ID:           sessionID,
		UserID:       userID,
		TokenHash:    tokenHash,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    now,
		LastActiveAt: now,
		ExpiresAt:    expiresAt,
	}, nil
}

// GetSessionByToken retrieves a session by token hash
func GetSessionByToken(token string) (*Session, error) {
	tokenHash := hashToken(token)

	var s Session
	var createdAtStr, lastActiveStr, expiresStr string

	query := `SELECT id, user_id, token_hash, ip_address, user_agent, created_at, last_active_at, expires_at
	          FROM sessions WHERE token_hash = ? AND expires_at > datetime('now')`

	err := DB.QueryRow(query, tokenHash).Scan(
		&s.ID, &s.UserID, &s.TokenHash, &s.IPAddress, &s.UserAgent,
		&createdAtStr, &lastActiveStr, &expiresStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Parse timestamps
	timeLayout := "2006-01-02 15:04:05"
	s.CreatedAt, err = time.Parse(timeLayout, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	s.LastActiveAt, err = time.Parse(timeLayout, lastActiveStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last_active_at: %w", err)
	}

	s.ExpiresAt, err = time.Parse(timeLayout, expiresStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expires_at: %w", err)
	}

	return &s, nil
}

// UpdateSessionActivity updates last_active_at for a session
func UpdateSessionActivity(token string) error {
	tokenHash := hashToken(token)
	_, err := DB.Exec(`UPDATE sessions SET last_active_at = datetime('now') WHERE token_hash = ?`, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}
	return nil
}

// DeleteSession removes a session
func DeleteSession(token string) error {
	tokenHash := hashToken(token)
	_, err := DB.Exec(`DELETE FROM sessions WHERE token_hash = ?`, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteAllUserSessions removes ALL sessions for a specific user
// This ensures complete logout across all devices/browsers
func DeleteAllUserSessions(userID string) error {
	result, err := DB.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Deleted %d sessions for user %s\n", rowsAffected, userID)
	return nil
}

// CleanupExpiredSessions removes expired sessions (run periodically)
func CleanupExpiredSessions() error {
	result, err := DB.Exec(`DELETE FROM sessions WHERE expires_at < datetime('now')`)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("🧹 Cleaned up %d expired sessions\n", rowsAffected)
	}
	return nil
}

// CleanupStaleSessions removes sessions that haven't been active for more than 30 minutes
// This handles cases where users close the browser without logging out
func CleanupStaleSessions() error {
	result, err := DB.Exec(`DELETE FROM sessions WHERE last_active_at < datetime('now', '-30 minutes')`)
	if err != nil {
		return fmt.Errorf("failed to cleanup stale sessions: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("🧹 Cleaned up %d stale sessions (inactive for 30+ minutes)\n", rowsAffected)
	}
	return nil
}

// CleanupAllReaderSessions removes all reader sessions (for fresh start)
func CleanupAllReaderSessions() error {
	result, err := DB.Exec(`DELETE FROM sessions WHERE user_id IN (SELECT id FROM users WHERE role = 'reader')`)
	if err != nil {
		return fmt.Errorf("failed to cleanup reader sessions: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("🧹 Cleaned up %d reader sessions\n", rowsAffected)
	return nil
}

// IsUserOnline checks if a user has an active session (online)
// A user is considered online if they have an active session that hasn't expired
// and has been active within the last 10 minutes
func IsUserOnline(userID string) (bool, error) {
	var count int
	// Check if user has an active session (not expired and active within last 10 minutes)
	query := `SELECT COUNT(*) FROM sessions 
	          WHERE user_id = ? 
	          AND expires_at > datetime('now') 
	          AND last_active_at >= datetime('now', '-10 minutes')`
	
	err := DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user online status: %w", err)
	}
	
	return count > 0, nil
}

// GetOnlineReaderIDs returns a map of reader IDs that are currently online
func GetOnlineReaderIDs() (map[string]bool, error) {
	onlineMap := make(map[string]bool)
	
	// Get all readers with active sessions (not expired and active within last 10 minutes)
	query := `SELECT DISTINCT s.user_id 
	          FROM sessions s
	          INNER JOIN users u ON s.user_id = u.id
	          WHERE u.role = 'reader'
	          AND s.expires_at > datetime('now') 
	          AND s.last_active_at >= datetime('now', '-10 minutes')`
	
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get online readers: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		onlineMap[userID] = true
	}
	
	return onlineMap, rows.Err()
}

// hashToken creates a SHA-256 hash of the token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
