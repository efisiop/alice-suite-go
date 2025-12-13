package database

import (
	"database/sql"
	"time"
)

// VerificationCode represents a verification code in the database
type VerificationCode struct {
	Code      string
	BookID    string
	IsUsed    bool
	UsedBy    *string
	CreatedAt time.Time
}

// GetVerificationCode retrieves a verification code by code string
func GetVerificationCode(code string) (*VerificationCode, error) {
	vc := &VerificationCode{}
	var usedBy sql.NullString
	var createdAtStr string

	query := `SELECT code, book_id, is_used, used_by, created_at
	          FROM verification_codes WHERE code = ?`
	err := DB.QueryRow(query, code).Scan(
		&vc.Code, &vc.BookID, &vc.IsUsed, &usedBy, &createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if usedBy.Valid {
		vc.UsedBy = &usedBy.String
	}

	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			vc.CreatedAt = t
		}
	}

	return vc, nil
}

// MarkVerificationCodeUsed marks a verification code as used
func MarkVerificationCodeUsed(code, userID string) error {
	query := `UPDATE verification_codes SET is_used = 1, used_by = ? WHERE code = ?`
	_, err := DB.Exec(query, userID, code)
	return err
}

// CreateVerificationCode creates a new verification code
func CreateVerificationCode(vc *VerificationCode) error {
	query := `INSERT INTO verification_codes (code, book_id, is_used, used_by, created_at)
	          VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, vc.Code, vc.BookID, vc.IsUsed, vc.UsedBy, vc.CreatedAt)
	return err
}

// UpdateUserVerification updates a user's verification status
func UpdateUserVerification(userID string, verified bool) error {
	query := `UPDATE users SET is_verified = ?, updated_at = ? WHERE id = ?`
	verifiedInt := 0
	if verified {
		verifiedInt = 1
	}
	_, err := DB.Exec(query, verifiedInt, time.Now(), userID)
	return err
}

