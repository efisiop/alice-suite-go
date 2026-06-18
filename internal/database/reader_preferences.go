package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/models"
)

const DefaultReaderLanguageCode = "en"

var supportedReaderLanguageNames = map[string]string{
	"en": "English",
	"it": "Italian",
	"da": "Danish",
	"es": "Spanish",
	"fr": "French",
	"de": "German",
	"pt": "Portuguese",
}

// NormalizeReaderLanguageCode returns a supported BCP-47-like language code.
func NormalizeReaderLanguageCode(languageCode string) string {
	code := strings.ToLower(strings.TrimSpace(languageCode))
	if len(code) > 2 {
		code = strings.Split(code, "-")[0]
	}
	if _, ok := supportedReaderLanguageNames[code]; ok {
		return code
	}
	return DefaultReaderLanguageCode
}

// ReaderLanguageName returns the display language name for prompting/UI.
func ReaderLanguageName(languageCode string) string {
	code := NormalizeReaderLanguageCode(languageCode)
	return supportedReaderLanguageNames[code]
}

// GetReaderPreference retrieves preferences for a reader, creating defaults when missing.
func GetReaderPreference(userID string) (*models.ReaderPreference, error) {
	pref := &models.ReaderPreference{}
	var createdAtStr, updatedAtStr string
	query := `SELECT user_id, preferred_language_code, created_at, updated_at
	          FROM reader_preferences WHERE user_id = ?`
	err := DB.QueryRow(Rebind(query), userID).Scan(
		&pref.UserID, &pref.PreferredLanguageCode, &createdAtStr, &updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return UpsertReaderPreference(userID, DefaultReaderLanguageCode)
	}
	if err != nil {
		return nil, err
	}
	pref.PreferredLanguageCode = NormalizeReaderLanguageCode(pref.PreferredLanguageCode)
	if createdAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			pref.CreatedAt = t
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			pref.UpdatedAt = t
		}
	}
	return pref, nil
}

// UpsertReaderPreference creates or updates preferences for a reader.
func UpsertReaderPreference(userID, languageCode string) (*models.ReaderPreference, error) {
	code := NormalizeReaderLanguageCode(languageCode)
	now := FormatSQLDateTime(time.Now())
	if DriverName == "postgres" {
		query := `INSERT INTO reader_preferences (user_id, preferred_language_code, created_at, updated_at)
		          VALUES (?, ?, ?, ?)
		          ON CONFLICT (user_id) DO UPDATE SET preferred_language_code = EXCLUDED.preferred_language_code, updated_at = EXCLUDED.updated_at`
		if _, err := DB.Exec(Rebind(query), userID, code, now, now); err != nil {
			return nil, err
		}
	} else {
		query := `INSERT INTO reader_preferences (user_id, preferred_language_code, created_at, updated_at)
		          VALUES (?, ?, ?, ?)
		          ON CONFLICT(user_id) DO UPDATE SET preferred_language_code = excluded.preferred_language_code, updated_at = excluded.updated_at`
		if _, err := DB.Exec(query, userID, code, now, now); err != nil {
			return nil, err
		}
	}
	return &models.ReaderPreference{
		UserID:                userID,
		PreferredLanguageCode: code,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}, nil
}

func ensureReaderPreferencesTable() error {
	tsDefault := "CURRENT_TIMESTAMP"
	if DriverName == "sqlite3" {
		tsDefault = "datetime('now')"
	}
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS reader_preferences (
		user_id TEXT PRIMARY KEY,
		preferred_language_code TEXT NOT NULL DEFAULT 'en',
		created_at TEXT NOT NULL DEFAULT (%s),
		updated_at TEXT NOT NULL DEFAULT (%s),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`, tsDefault, tsDefault)
	if _, err := DB.Exec(q); err != nil {
		return err
	}
	_, err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_reader_preferences_language ON reader_preferences(preferred_language_code)`)
	return err
}
