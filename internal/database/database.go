package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

// DB represents the database connection
var DB *sql.DB

// InitDB initializes the SQLite database connection with WAL mode and optimal PRAGMAs
func InitDB(dbPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open database connection with WAL mode support
	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return err
	}

	// CRITICAL: Set connection pool limits for SQLite
	// SQLite works best with limited concurrent connections
	DB.SetMaxOpenConns(25)        // Limit concurrent connections
	DB.SetMaxIdleConns(5)         // Keep 5 idle connections
	DB.SetConnMaxLifetime(0)      // Reuse connections indefinitely

	// Execute PRAGMAs for optimal performance and concurrency
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",           // Enable Write-Ahead Logging (allows concurrent reads)
		"PRAGMA synchronous = NORMAL;",         // Good balance between speed and safety
		"PRAGMA foreign_keys = ON;",            // Enforce referential integrity
		"PRAGMA busy_timeout = 5000;",          // Wait 5 seconds before "database locked" error
		"PRAGMA wal_autocheckpoint = 1000;",     // Auto-checkpoint WAL every 1000 pages
		"PRAGMA cache_size = -20000;",          // 20 MB page cache (adjust based on available RAM)
		"PRAGMA temp_store = MEMORY;",          // Store temp tables in memory for speed
	}

	for _, pragma := range pragmas {
		if _, err := DB.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma %s: %w", pragma, err)
		}
	}

	// Test connection by executing a simple query (SQLite doesn't support Ping)
	_, err = DB.Exec("SELECT 1")
	if err != nil {
		return err
	}

	// Ensure consultant_prompts table exists (migration 012) so saving prompts works without running migrate
	if err := ensureConsultantPromptsTable(); err != nil {
		return fmt.Errorf("ensure consultant_prompts table: %w", err)
	}

	return nil
}

// ensureConsultantPromptsTable creates the consultant_prompts table if it doesn't exist
func ensureConsultantPromptsTable() error {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS consultant_prompts (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		book_id TEXT NOT NULL,
		page_number INTEGER NOT NULL,
		section_number INTEGER,
		prompt_text TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		dismissed_at TEXT,
		accepted_at TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}
	for _, idx := range []string{
		`CREATE INDEX IF NOT EXISTS idx_consultant_prompts_user_book ON consultant_prompts(user_id, book_id)`,
		`CREATE INDEX IF NOT EXISTS idx_consultant_prompts_page ON consultant_prompts(user_id, book_id, page_number)`,
	} {
		if _, err := DB.Exec(idx); err != nil {
			return err
		}
	}
	_, _ = DB.Exec(`ALTER TABLE consultant_prompts ADD COLUMN dismissed_at TEXT`)
	_, _ = DB.Exec(`ALTER TABLE consultant_prompts ADD COLUMN accepted_at TEXT`)
	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// RunMigrations runs SQL migration files
func RunMigrations(migrationsPath string) error {
	// This will be implemented to read and execute migration files
	// For now, return nil
	return nil
}



