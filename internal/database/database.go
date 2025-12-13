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



