package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
var DB *sql.DB

// DriverName is "postgres" or "sqlite3" - used for query rebinding
var DriverName string

// InitDB initializes the database connection.
// When databaseURL is non-empty (e.g. DATABASE_URL from Render), uses PostgreSQL.
// Otherwise uses SQLite with dbPath (local development).
func InitDB(dbPath, databaseURL string) error {
	var err error

	if databaseURL != "" {
		// PostgreSQL (production: Render.com, or local with Docker)
		DriverName = "postgres"
		DB, err = sql.Open("pgx", databaseURL)
		if err != nil {
			return fmt.Errorf("open postgres: %w", err)
		}
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)
		DB.SetConnMaxLifetime(0)
	} else {
		// SQLite (local development)
		DriverName = "sqlite3"
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
		if err != nil {
			return err
		}
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)
		DB.SetConnMaxLifetime(0)

		// SQLite-only PRAGMAs
		pragmas := []string{
			"PRAGMA journal_mode = WAL;",
			"PRAGMA synchronous = NORMAL;",
			"PRAGMA foreign_keys = ON;",
			"PRAGMA busy_timeout = 5000;",
			"PRAGMA wal_autocheckpoint = 1000;",
			"PRAGMA cache_size = -20000;",
			"PRAGMA temp_store = MEMORY;",
		}
		for _, pragma := range pragmas {
			if _, err := DB.Exec(pragma); err != nil {
				return fmt.Errorf("failed to set pragma %s: %w", pragma, err)
			}
		}
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	// For PostgreSQL, migrations (012, 015) create these tables; do not run ensure* here
	// or migrate would fail on empty DB (ensure* references users which is created in 001).
	// For SQLite, ensure* is safe and keeps schema in sync if migrations were skipped.
	if databaseURL == "" {
		if err := ensureConsultantPromptsTable(); err != nil {
			return fmt.Errorf("ensure consultant_prompts table: %w", err)
		}
		if err := ensureAhAhMomentsTable(); err != nil {
			return fmt.Errorf("ensure ah_ah_moments table: %w", err)
		}
		if err := ensureReaderPreferencesTable(); err != nil {
			return fmt.Errorf("ensure reader_preferences table: %w", err)
		}
	}

	return nil
}

// Rebind converts ? placeholders to $1, $2... for PostgreSQL.
// SQLite uses ? and ignores this; call before every Exec/Query/QueryRow when using raw queries.
func Rebind(query string) string {
	if DriverName == "postgres" {
		return sqlx.Rebind(sqlx.DOLLAR, query)
	}
	return query
}

// FormatSQLDateTime formats a time for TEXT and TIMESTAMP columns. PostgreSQL's pgx
// driver cannot encode time.Time into TEXT (OID 25); use this for INSERT/UPDATE/WHERE
// on schema columns declared as TEXT in migrations.
func FormatSQLDateTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05")
}

// ensureConsultantPromptsTable creates the consultant_prompts table if it doesn't exist
func ensureConsultantPromptsTable() error {
	tsDefault := "CURRENT_TIMESTAMP"
	if DriverName == "sqlite3" {
		tsDefault = "datetime('now')"
	}
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS consultant_prompts (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		book_id TEXT NOT NULL,
		page_number INTEGER NOT NULL,
		section_number INTEGER,
		prompt_text TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (%s),
		updated_at TEXT NOT NULL DEFAULT (%s),
		dismissed_at TEXT,
		accepted_at TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
	)`, tsDefault, tsDefault)
	_, err := DB.Exec(q)
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

// ensureAhAhMomentsTable creates the ah_ah_moments table if it doesn't exist (migration 015)
func ensureAhAhMomentsTable() error {
	tsDefault := "CURRENT_TIMESTAMP"
	if DriverName == "sqlite3" {
		tsDefault = "datetime('now')"
	}
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS ah_ah_moments (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		book_id TEXT NOT NULL,
		content TEXT NOT NULL,
		page_number INTEGER,
		section_number INTEGER,
		shared INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT (%s),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
	)`, tsDefault)
	_, err := DB.Exec(q)
	if err != nil {
		return err
	}
	for _, idx := range []string{
		`CREATE INDEX IF NOT EXISTS idx_ah_ah_moments_user_book ON ah_ah_moments(user_id, book_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ah_ah_moments_book_shared ON ah_ah_moments(book_id, shared)`,
		`CREATE INDEX IF NOT EXISTS idx_ah_ah_moments_created ON ah_ah_moments(created_at DESC)`,
	} {
		if _, err := DB.Exec(idx); err != nil {
			return err
		}
	}
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
