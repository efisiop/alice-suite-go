package database

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

// TestDB provides test database access and management
type TestDB struct {
	testing.TB
	DB  *sql.DB
	originalDB *sql.DB
}

// SetupTestDatabase creates an in-memory SQLite database for testing
// It temporarily replaces the global DB variable for test execution
func SetupTestDatabase(t testing.TB) *TestDB {
	// Save the original database connection
	originalDB := DB

	// Create in-memory SQLite database for tests
	testDB, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Set the test database as the global temporarily
	DB = testDB

	// Run migrations to set up schema
	if err := runTestMigrations(testDB); err != nil {
		testDB.Close()
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	return &TestDB{
		TB: t,
		DB: testDB,
		originalDB: originalDB,
	}
}

// runTestMigrations executes all migration files for test setup
func runTestMigrations(db *sql.DB) error {
	migrationFiles := []string{
		"/Users/efisiopittau/Project_1/alice-suite-go/migrations/001_initial_schema.sql",
		"/Users/efisiopittau/Project_1/alice-suite-go/migrations/002_seed_first_3_chapters.sql",
	}

	for _, file := range migrationFiles {
		sql, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Warning: Could not read migration %s: %v", file, err)
			continue
		}

		if _, err := db.Exec(string(sql)); err != nil {
			log.Printf("Warning: Migration %s failed: %v", file, err)
			return err
		}
	}

	return nil
}

// Cleanup restores the original database connection and closes the test database
func (td *TestDB) Cleanup() {
	// Restore original DB connection
	DB = td.originalDB

	// Close test database if it wasn't nil (safety check)
	if td.DB != nil {
		td.DB.Close()
	}
}

// WithTx provides a test-only transaction function
func (td *TestDB) WithTx(fn func(*sql.Tx) error) error {
	tx, err := td.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// ClearTable deletes all rows from a table for test isolation
func (td *TestDB) ClearTable(tableName string) error {
	if _, err := td.DB.Exec("DELETE FROM " + tableName); err != nil {
		return err
	}
	return nil
}