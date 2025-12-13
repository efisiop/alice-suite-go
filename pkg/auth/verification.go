package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/google/uuid"
)

var (
	ErrInvalidCode     = errors.New("invalid verification code")
	ErrCodeAlreadyUsed = errors.New("verification code already used")
	ErrUserNotVerified = errors.New("user not verified")
)

// VerifyBookCode verifies a book verification code for a user
func VerifyBookCode(code, userID string) (string, error) {
	// First, get the verification code to check if it exists
	// We'll use a simple query to get the book_id
	var bookID string
	var isUsed bool
	query := `SELECT book_id, is_used FROM verification_codes WHERE code = ?`
	err := database.DB.QueryRow(query, code).Scan(&bookID, &isUsed)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidCode
		}
		return "", err
	}

	// Check if code is already used
	if isUsed {
		return "", ErrCodeAlreadyUsed
	}

	// Mark code as used
	err = database.UseVerificationCode(code, userID)
	if err != nil {
		return "", err
	}

	// Update user's verification status (set is_verified = 1)
	updateQuery := `UPDATE users SET is_verified = 1, updated_at = ? WHERE id = ?`
	_, err = database.DB.Exec(updateQuery, time.Now(), userID)
	if err != nil {
		return "", err
	}

	return bookID, nil
}

// CheckBookVerified checks if a user has verified their book
func CheckBookVerified(userID string) (bool, error) {
	user, err := database.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, ErrUserNotFound
	}

	// For now, use is_verified field to indicate book verification
	// In the future, we might add a separate book_verified field
	return user.IsVerified, nil
}

// CreateVerificationCode creates a new verification code for a book
func CreateVerificationCode(bookID string) (string, error) {
	code := generateVerificationCode()

	query := `INSERT INTO verification_codes (code, book_id, is_used, created_at)
	          VALUES (?, ?, 0, ?)`
	_, err := database.DB.Exec(query, code, bookID, time.Now())
	if err != nil {
		return "", err
	}

	return code, nil
}

// generateVerificationCode generates a random verification code
func generateVerificationCode() string {
	// Generate a simple code (can be enhanced)
	return uuid.New().String()[:8]
}
