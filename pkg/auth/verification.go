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
	tx, err := database.DB.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Look up and claim the code in the same transaction so a later failure does
	// not consume a valid one without also verifying the reader.
	var bookID string
	var isUsed bool
	query := `SELECT book_id, is_used FROM verification_codes WHERE code = ?`
	err = tx.QueryRow(database.Rebind(query), code).Scan(&bookID, &isUsed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCode
		}
		return "", err
	}

	// Check if code is already used
	if isUsed {
		return "", ErrCodeAlreadyUsed
	}

	result, err := tx.Exec(
		database.Rebind(`UPDATE verification_codes SET is_used = 1, used_by = ? WHERE code = ? AND is_used = 0`),
		userID,
		code,
	)
	if err != nil {
		return "", err
	}
	claimed, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if claimed != 1 {
		return "", ErrCodeAlreadyUsed
	}

	// Update user's verification status (set is_verified = 1)
	updateQuery := `UPDATE users SET is_verified = 1, updated_at = ? WHERE id = ?`
	userUpdate, err := tx.Exec(database.Rebind(updateQuery), database.FormatSQLDateTime(time.Now()), userID)
	if err != nil {
		return "", err
	}
	updatedUsers, err := userUpdate.RowsAffected()
	if err != nil {
		return "", err
	}
	if updatedUsers != 1 {
		return "", ErrUserNotFound
	}

	if err := tx.Commit(); err != nil {
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
	_, err := database.DB.Exec(database.Rebind(query), code, bookID, database.FormatSQLDateTime(time.Now()))
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
