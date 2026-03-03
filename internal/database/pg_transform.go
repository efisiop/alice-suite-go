package database

import (
	"regexp"
	"strings"
)

// TransformSQLiteToPostgres converts SQLite migration SQL to PostgreSQL-compatible SQL.
func TransformSQLiteToPostgres(sql string) string {
	out := sql

	// Remove PRAGMA lines (PostgreSQL doesn't use them)
	pragmaRe := regexp.MustCompile(`(?m)^\s*PRAGMA\s+.*$[\r\n]*`)
	out = pragmaRe.ReplaceAllString(out, "")

	// datetime('now') -> CURRENT_TIMESTAMP
	out = strings.ReplaceAll(out, "datetime('now')", "CURRENT_TIMESTAMP")

	// datetime('now', '-X minutes') -> CURRENT_TIMESTAMP - INTERVAL 'X minutes'
	datetimeRe := regexp.MustCompile(`datetime\('now',\s*'-(\d+)\s+(\w+)'\)`)
	out = datetimeRe.ReplaceAllString(out, "CURRENT_TIMESTAMP - INTERVAL '$1 $2'")

	// INSERT OR IGNORE -> INSERT (idempotency via DO block in migrate command)
	out = strings.ReplaceAll(out, "INSERT OR IGNORE", "INSERT")

	// INSERT OR REPLACE -> INSERT ... ON CONFLICT DO UPDATE (for migrations that use it)
	out = strings.ReplaceAll(out, "INSERT OR REPLACE", "INSERT")

	// PostgreSQL: char(10) (SQLite) -> chr(10)
	out = strings.ReplaceAll(out, "char(10)", "chr(10)")

	// PostgreSQL 9.5+: ADD COLUMN IF NOT EXISTS
	out = strings.ReplaceAll(out, "ADD COLUMN purchase_date", "ADD COLUMN IF NOT EXISTS purchase_date")

	return out
}

// WrapInsertForPostgres wraps an INSERT statement in a DO block that catches unique_violation.
// This makes INSERT idempotent (like SQLite's INSERT OR IGNORE).
func WrapInsertForPostgres(stmt string) string {
	trimmed := strings.TrimSpace(stmt)
	if !strings.HasPrefix(trimmed, "INSERT INTO") {
		return stmt
	}
	if strings.HasSuffix(trimmed, ";") {
		trimmed = strings.TrimSuffix(trimmed, ";")
	}
	return `DO $$ BEGIN ` + trimmed + `; EXCEPTION WHEN unique_violation THEN NULL; END $$`
}
