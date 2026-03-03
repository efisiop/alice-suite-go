package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/config"
	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	// Initialize database (PostgreSQL when DATABASE_URL set, else SQLite)
	if err := database.InitDB(cfg.DBPath, cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	fmt.Println("🌱 Initializing test users...")

	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	// Create test reader user
	readerEmail := "reader@example.com"
	readerPassword := "reader123"
	readerID := uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(readerPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	var existingID string
	err = database.DB.QueryRow(database.Rebind("SELECT id FROM users WHERE email = ?"), readerEmail).Scan(&existingID)
	if err == nil {
		fmt.Printf("✅ Reader user already exists: %s (ID: %s)\n", readerEmail, existingID)
	} else if err == sql.ErrNoRows {
		_, err = database.DB.Exec(database.Rebind(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`), readerID, readerEmail, string(hashedPassword), "Test", "Reader", "reader", 0, now, now)
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

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(consultantPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	err = database.DB.QueryRow(database.Rebind("SELECT id FROM users WHERE email = ?"), consultantEmail).Scan(&existingID)
	if err == nil {
		fmt.Printf("✅ Consultant user already exists: %s (ID: %s)\n", consultantEmail, existingID)
	} else if err == sql.ErrNoRows {
		_, err = database.DB.Exec(database.Rebind(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`), consultantID, consultantEmail, string(hashedPassword), "Test", "Consultant", "consultant", 1, now, now)
		if err != nil {
			log.Fatalf("Failed to create consultant user: %v", err)
		}
		fmt.Printf("✅ Created consultant user: %s (Password: %s)\n", consultantEmail, consultantPassword)
	} else {
		log.Fatalf("Error checking for existing consultant: %v", err)
	}

	// Create test administrator user
	adminEmail := "admin@example.com"
	adminPassword := "admin123"
	adminID := uuid.New().String()

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	err = database.DB.QueryRow(database.Rebind("SELECT id FROM users WHERE email = ?"), adminEmail).Scan(&existingID)
	if err == nil {
		fmt.Printf("✅ Administrator user already exists: %s (ID: %s)\n", adminEmail, existingID)
	} else if err == sql.ErrNoRows {
		_, err = database.DB.Exec(database.Rebind(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`), adminID, adminEmail, string(hashedPassword), "Admin", "User", "administrator", 1, now, now)
		if err != nil {
			log.Fatalf("Failed to create administrator user: %v", err)
		}
		fmt.Printf("✅ Created administrator user: %s (Password: %s)\n", adminEmail, adminPassword)
	} else {
		log.Fatalf("Error checking for existing administrator: %v", err)
	}

	// Create verification code for reader
	verificationCode := "ALICE2024"
	if database.DriverName == "postgres" {
		_, err = database.DB.Exec(database.Rebind(`
			INSERT INTO verification_codes (code, book_id, is_used, created_at)
			VALUES (?, ?, ?, ?)
			ON CONFLICT (code) DO NOTHING
		`), verificationCode, "alice-in-wonderland", 0, now)
	} else {
		_, err = database.DB.Exec(database.Rebind(`
			INSERT OR IGNORE INTO verification_codes (code, book_id, is_used, created_at)
			VALUES (?, ?, ?, ?)
		`), verificationCode, "alice-in-wonderland", 0, now)
	}
	if err != nil {
		log.Printf("Warning: Failed to create verification code: %v", err)
	} else {
		fmt.Printf("✅ Verification code created: %s\n", verificationCode)
	}

	// Create efisio user
	efisioEmail := "efisio@efisio.com"
	efisioPassword := "efisio123"
	efisioID := uuid.New().String()

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(efisioPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	err = database.DB.QueryRow(database.Rebind("SELECT id FROM users WHERE email = ?"), efisioEmail).Scan(&existingID)
	if err == nil {
		fmt.Printf("✅ Efisio user already exists: %s (ID: %s)\n", efisioEmail, existingID)
	} else if err == sql.ErrNoRows {
		_, err = database.DB.Exec(database.Rebind(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`), efisioID, efisioEmail, string(hashedPassword), "Efisio", "Pittau", "reader", 1, now, now)
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
	fmt.Println("Administrator: admin@example.com / admin123")
	fmt.Println("Verification Code: ALICE2024")
	fmt.Println("\n✅ Test users initialized successfully!")
}
