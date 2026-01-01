package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func getDBPath() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/alice-suite.db"
	}
	return dbPath
}

func main() {
	dbPath := getDBPath()

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Database not found at %s. Please run migrations first.", dbPath)
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

	fmt.Println("🌱 Initializing test users...")

	// Create test reader user
	readerEmail := "reader@example.com"
	readerPassword := "reader123"
	readerID := uuid.New().String()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(readerPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Check if user already exists
	var existingID string
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", readerEmail).Scan(&existingID)
	if err == nil {
		fmt.Printf("✅ Reader user already exists: %s (ID: %s)\n", readerEmail, existingID)
	} else if err == sql.ErrNoRows {
		// Insert reader user
		_, err = db.Exec(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
		`, readerID, readerEmail, string(hashedPassword), "Test", "Reader", "reader", 0)
		if err != nil {
			log.Fatalf("Failed to create reader user: %v", err)
		}
		fmt.Printf("✅ Created reader user: %s (Password: %s)\n", readerEmail, readerPassword)
	} else {
		log.Fatalf("Error checking for existing user: %v", err)
	}

	// Create test consultant user
	consultantEmail := "consultant@example.com"
	consultantPassword := "consultant123"
	consultantID := uuid.New().String()

	// Hash password
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(consultantPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Check if consultant already exists
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", consultantEmail).Scan(&existingID)
	if err == nil {
		fmt.Printf("✅ Consultant user already exists: %s (ID: %s)\n", consultantEmail, existingID)
	} else if err == sql.ErrNoRows {
		// Insert consultant user
		_, err = db.Exec(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
		`, consultantID, consultantEmail, string(hashedPassword), "Test", "Consultant", "consultant", 1)
		if err != nil {
			log.Fatalf("Failed to create consultant user: %v", err)
		}
		fmt.Printf("✅ Created consultant user: %s (Password: %s)\n", consultantEmail, consultantPassword)
	} else {
		log.Fatalf("Error checking for existing consultant: %v", err)
	}

	// Create verification code for reader
	verificationCode := "ALICE2024"
	_, err = db.Exec(`
		INSERT OR IGNORE INTO verification_codes (code, book_id, is_used, created_at)
		VALUES (?, ?, ?, datetime('now'))
	`, verificationCode, "alice-in-wonderland", 0)
	if err != nil {
		log.Printf("Warning: Failed to create verification code: %v", err)
	} else {
		fmt.Printf("✅ Verification code created: %s\n", verificationCode)
	}

	// Create efisio user
	efisioEmail := "efisio@efisio.com"
	efisioPassword := "efisio123"
	efisioID := uuid.New().String()

	// Hash password
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(efisioPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Check if efisio already exists
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", efisioEmail).Scan(&existingID)
	if err == nil {
		fmt.Printf("✅ Efisio user already exists: %s (ID: %s)\n", efisioEmail, existingID)
	} else if err == sql.ErrNoRows {
		// Insert efisio user
		_, err = db.Exec(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
		`, efisioID, efisioEmail, string(hashedPassword), "Efisio", "Pittau", "reader", 1)
		if err != nil {
			log.Fatalf("Failed to create efisio user: %v", err)
		}
		fmt.Printf("✅ Created efisio user: %s (Password: %s)\n", efisioEmail, efisioPassword)
	} else {
		log.Fatalf("Error checking for existing efisio user: %v", err)
	}

	fmt.Println("\n📋 Test Users Created:")
	fmt.Println("Reader: reader@example.com / reader123")
	fmt.Println("Efisio: efisio@efisio.com / efisio123")
	fmt.Println("Consultant: consultant@example.com / consultant123")
	fmt.Println("Verification Code: ALICE2024")
	fmt.Println("\n✅ Test users initialized successfully!")
}
