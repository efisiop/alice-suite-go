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
