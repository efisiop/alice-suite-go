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
	nowStr := FormatSQLDateTime(now)
	expStr := FormatSQLDateTime(expiresAt)
	_, err := DB.Exec(Rebind(query), sessionID, userID, tokenHash, ipAddress, userAgent, nowStr, nowStr, expStr)
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
	          FROM sessions WHERE token_hash = ? AND expires_at > ?`

	err := DB.QueryRow(Rebind(query), tokenHash, FormatSQLDateTime(time.Now())).Scan(
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
	_, err := DB.Exec(Rebind(`UPDATE sessions SET last_active_at = ? WHERE token_hash = ?`), FormatSQLDateTime(time.Now()), tokenHash)
	if err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}
	return nil
}

// DeleteSession removes a session
func DeleteSession(token string) error {
	tokenHash := hashToken(token)
	_, err := DB.Exec(Rebind(`DELETE FROM sessions WHERE token_hash = ?`), tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteAllUserSessions removes ALL sessions for a specific user
// This ensures complete logout across all devices/browsers
func DeleteAllUserSessions(userID string) error {
	result, err := DB.Exec(Rebind(`DELETE FROM sessions WHERE user_id = ?`), userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Deleted %d sessions for user %s\n", rowsAffected, userID)
	return nil
}

// CleanupExpiredSessions removes expired sessions (run periodically)
func CleanupExpiredSessions() error {
	result, err := DB.Exec(Rebind(`DELETE FROM sessions WHERE expires_at < ?`), FormatSQLDateTime(time.Now()))
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
	result, err := DB.Exec(Rebind(`DELETE FROM sessions WHERE last_active_at < ?`), FormatSQLDateTime(time.Now().Add(-30*time.Minute)))
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
	result, err := DB.Exec(Rebind(`DELETE FROM sessions WHERE user_id IN (SELECT id FROM users WHERE role = 'reader')`))
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
	          AND expires_at > ? 
	          AND last_active_at >= ?`
	
	err := DB.QueryRow(Rebind(query), userID, FormatSQLDateTime(time.Now()), FormatSQLDateTime(time.Now().Add(-10*time.Minute))).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user online status: %w", err)
	}
	
	return count > 0, nil
}

// OnlineReaderSession represents one online reader session (user + device/session info).
type OnlineReaderSession struct {
	UserID       string
	Email        string
	FirstName    string
	LastName     string
	UserAgent    string
	LastActiveAt time.Time
}

// GetOnlineReaderSessions returns all reader sessions active in the last 15 minutes (one row per session/device).
func GetOnlineReaderSessions() ([]OnlineReaderSession, error) {
	query := `SELECT u.id, u.email, COALESCE(u.first_name,''), COALESCE(u.last_name,''), COALESCE(s.user_agent,''), s.last_active_at
	          FROM sessions s
	          INNER JOIN users u ON s.user_id = u.id
	          WHERE u.role = 'reader'
	          AND s.expires_at > ?
	          AND s.last_active_at >= ?
	          ORDER BY s.last_active_at DESC`
	rows, err := DB.Query(Rebind(query), FormatSQLDateTime(time.Now()), FormatSQLDateTime(time.Now().Add(-15*time.Minute)))
	if err != nil {
		return nil, fmt.Errorf("failed to get online reader sessions: %w", err)
	}
	defer rows.Close()
	timeLayout := "2006-01-02 15:04:05"
	var result []OnlineReaderSession
	for rows.Next() {
		var o OnlineReaderSession
		var lastActiveStr string
		if err := rows.Scan(&o.UserID, &o.Email, &o.FirstName, &o.LastName, &o.UserAgent, &lastActiveStr); err != nil {
			continue
		}
		if t, err := time.Parse(timeLayout, lastActiveStr); err == nil {
			o.LastActiveAt = t
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

// CountReadersOnline returns the number of readers with an active session in the last 15 minutes.
func CountReadersOnline() (int, error) {
	var count int
	query := `SELECT COUNT(DISTINCT s.user_id) FROM sessions s
	          INNER JOIN users u ON s.user_id = u.id
	          WHERE u.role = 'reader'
	          AND s.expires_at > ?
	          AND s.last_active_at >= ?`
	err := DB.QueryRow(Rebind(query), FormatSQLDateTime(time.Now()), FormatSQLDateTime(time.Now().Add(-15*time.Minute))).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count readers online: %w", err)
	}
	return count, nil
}

// OnlineConsultantSession represents one online consultant session (user + device/session info).
type OnlineConsultantSession struct {
	UserID       string
	Email        string
	FirstName    string
	LastName     string
	UserAgent    string
	LastActiveAt time.Time
}

// GetOnlineConsultantSessions returns all consultant sessions active in the last 15 minutes (one row per session/device).
func GetOnlineConsultantSessions() ([]OnlineConsultantSession, error) {
	query := `SELECT u.id, u.email, COALESCE(u.first_name,''), COALESCE(u.last_name,''), COALESCE(s.user_agent,''), s.last_active_at
	          FROM sessions s
	          INNER JOIN users u ON s.user_id = u.id
	          WHERE u.role = 'consultant'
	          AND s.expires_at > ?
	          AND s.last_active_at >= ?
	          ORDER BY s.last_active_at DESC`
	rows, err := DB.Query(Rebind(query), FormatSQLDateTime(time.Now()), FormatSQLDateTime(time.Now().Add(-15*time.Minute)))
	if err != nil {
		return nil, fmt.Errorf("failed to get online consultant sessions: %w", err)
	}
	defer rows.Close()
	timeLayout := "2006-01-02 15:04:05"
	var result []OnlineConsultantSession
	for rows.Next() {
		var o OnlineConsultantSession
		var lastActiveStr string
		if err := rows.Scan(&o.UserID, &o.Email, &o.FirstName, &o.LastName, &o.UserAgent, &lastActiveStr); err != nil {
			continue
		}
		if t, err := time.Parse(timeLayout, lastActiveStr); err == nil {
			o.LastActiveAt = t
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

// CountConsultantsOnline returns the number of consultants with an active session in the last 15 minutes.
func CountConsultantsOnline() (int, error) {
	var count int
	query := `SELECT COUNT(DISTINCT s.user_id) FROM sessions s
	          INNER JOIN users u ON s.user_id = u.id
	          WHERE u.role = 'consultant'
	          AND s.expires_at > ?
	          AND s.last_active_at >= ?`
	err := DB.QueryRow(Rebind(query), FormatSQLDateTime(time.Now()), FormatSQLDateTime(time.Now().Add(-15*time.Minute))).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count consultants online: %w", err)
	}
	return count, nil
}

// GetOnlineReaderIDs returns a map of reader IDs that are currently online
func GetOnlineReaderIDs() (map[string]bool, error) {
	onlineMap := make(map[string]bool)

	// Get all readers with active sessions (not expired and active within last 10 minutes)
	query := `SELECT DISTINCT s.user_id
	          FROM sessions s
	          INNER JOIN users u ON s.user_id = u.id
	          WHERE u.role = 'reader'
	          AND s.expires_at > ?
	          AND s.last_active_at >= ?`
	
	rows, err := DB.Query(Rebind(query), FormatSQLDateTime(time.Now()), FormatSQLDateTime(time.Now().Add(-10*time.Minute)))
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
