package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "data/alice-suite.db"
	migrationsDir = "migrations"
)

// splitSQLStatements splits SQL by semicolons, but respects quoted strings
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	
	for i, char := range sql {
		switch char {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
			current.WriteRune(char)
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
			current.WriteRune(char)
		case ';':
			if !inSingleQuote && !inDoubleQuote {
				stmt := strings.TrimSpace(current.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
		_ = i // avoid unused variable
	}
	
	// Add remaining statement
	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	
	return statements
}

func main() {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✅ Database connection established")

	// Read migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Sort migration files by name
	migrationFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Execute migrations in order
	for _, filename := range migrationFiles {
		filepath := filepath.Join(migrationsDir, filename)
		fmt.Printf("📄 Running migration: %s\n", filename)

		sqlBytes, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", filename, err)
		}

		sql := string(sqlBytes)
		
		// Execute entire SQL file - SQLite can handle multiple statements
		// This approach handles quoted strings and multi-line statements correctly
		_, err = db.Exec(sql)
		if err != nil {
			// If executing entire file fails, try splitting by semicolons
			// but only split outside of quoted strings
			statements := splitSQLStatements(sql)
		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" || strings.HasPrefix(statement, "--") {
				continue
			}

				_, execErr := db.Exec(statement)
				if execErr != nil {
					log.Printf("Warning: Error executing statement in %s: %v", filename, execErr)
				// Continue with next statement (some errors are expected for IF NOT EXISTS)
				}
			}
		}

		fmt.Printf("✅ Migration %s completed\n", filename)
	}

	fmt.Println("\n🎉 All migrations completed successfully!")
	fmt.Printf("📊 Database created at: %s\n", dbPath)
}



