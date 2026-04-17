package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

var (
	migrationsDir = "migrations"
)

// splitSQLStatements splits SQL by semicolons while respecting:
//   - single-quoted strings (including doubled '' escape, which the toggle
//     approach handles naturally: each ' flips state, pairs flip back)
//   - double-quoted identifiers
//   - -- line comments (apostrophes inside comments must NOT affect quote state)
//   - /* ... */ block comments
//
// A previous version ignored comments entirely, so text like
//   -- Alice's glossary
// toggled inSingleQuote and caused every subsequent ; to be swallowed,
// concatenating many real INSERTs into one broken statement.
func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	runes := []rune(sql)
	n := len(runes)
	i := 0
	for i < n {
		char := runes[i]

		// Skip -- line comments (only when not inside a string)
		if !inSingleQuote && !inDoubleQuote &&
			char == '-' && i+1 < n && runes[i+1] == '-' {
			// Preserve the comment text in the output so line numbers are
			// still meaningful if we ever log a statement. Stop at newline.
			for i < n && runes[i] != '\n' {
				current.WriteRune(runes[i])
				i++
			}
			continue
		}

		// Skip /* ... */ block comments (only when not inside a string)
		if !inSingleQuote && !inDoubleQuote &&
			char == '/' && i+1 < n && runes[i+1] == '*' {
			current.WriteRune('/')
			current.WriteRune('*')
			i += 2
			for i+1 < n && !(runes[i] == '*' && runes[i+1] == '/') {
				current.WriteRune(runes[i])
				i++
			}
			if i+1 < n {
				current.WriteRune('*')
				current.WriteRune('/')
				i += 2
			}
			continue
		}

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
		i++
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
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/alice-suite.db"
	}
	databaseURL := os.Getenv("DATABASE_URL")

	// Initialize database (PostgreSQL when DATABASE_URL set, else SQLite)
	if databaseURL == "" {
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}
	if err := database.InitDB(dbPath, databaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	fmt.Println("✅ Database connection established")
	if database.DriverName == "postgres" {
		fmt.Println("📦 Using PostgreSQL (persistent storage)")
	}

	// Read migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	migrationFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		path := filepath.Join(migrationsDir, filename)
		fmt.Printf("📄 Running migration: %s\n", filename)

		sqlBytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", filename, err)
		}

		sqlContent := string(sqlBytes)
		if database.DriverName == "postgres" {
			// Migration 013 uses SQLite's "create users_new, copy, drop, rename"
			// pattern to alter a CHECK constraint. That pattern doesn't work on
			// Postgres because other tables have foreign keys to users(id).
			// Use Postgres-native ALTER TABLE statements instead.
			if filename == "013_add_administrator_role.sql" {
				sqlContent = `
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('reader', 'consultant', 'administrator'));
ALTER TABLE users ADD COLUMN IF NOT EXISTS admin_level TEXT;
`
			}
			sqlContent = database.TransformSQLiteToPostgres(sqlContent)
		}

		statements := splitSQLStatements(sqlContent)
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}
			if database.DriverName == "postgres" && strings.HasPrefix(stmt, "INSERT INTO") {
				stmt = database.WrapInsertForPostgres(stmt)
			}
			_, execErr := database.DB.Exec(stmt)
			if execErr != nil {
				log.Printf("Warning: Error executing statement in %s: %v", filename, execErr)
			}
		}

		fmt.Printf("✅ Migration %s completed\n", filename)
	}

	fmt.Println("\n🎉 All migrations completed successfully!")
	if database.DriverName == "postgres" {
		fmt.Println("📊 Connected to PostgreSQL (data persists across deploys)")
	} else {
		fmt.Printf("📊 Database at: %s\n", dbPath)
	}
}
