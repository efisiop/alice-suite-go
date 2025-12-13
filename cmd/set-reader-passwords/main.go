package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const dbPath = "data/alice-suite.db"

func main() {
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

	fmt.Println("🔐 Setting passwords for reader accounts...")

	// Define readers with their passwords
	readers := []struct {
		email    string
		password string
		name     string
	}{
		{"reader@example.com", "reader123", "Test Reader"},
		{"test@example.com", "test123", "Test User"},
		{"efisio@efisio.com", "efisio123", "Efisio Pittau"},
	}

	for _, reader := range readers {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reader.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", reader.email, err)
			continue
		}

		// Update password
		result, err := db.Exec(`
			UPDATE users 
			SET password_hash = ?, updated_at = datetime('now')
			WHERE email = ? AND role = 'reader'
		`, string(hashedPassword), reader.email)

		if err != nil {
			log.Printf("Failed to update password for %s: %v", reader.email, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			fmt.Printf("✅ Set password for %s (%s): %s\n", reader.email, reader.name, reader.password)
		} else {
			fmt.Printf("⚠️  User not found: %s\n", reader.email)
		}
	}

	fmt.Println("\n📋 Reader Accounts Ready for Testing:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("1. Email: reader@example.com")
	fmt.Println("   Password: reader123")
	fmt.Println("   Name: Test Reader")
	fmt.Println("")
	fmt.Println("2. Email: test@example.com")
	fmt.Println("   Password: test123")
	fmt.Println("   Name: Test User")
	fmt.Println("")
	fmt.Println("3. Email: efisio@efisio.com")
	fmt.Println("   Password: efisio123")
	fmt.Println("   Name: Efisio Pittau")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\n✅ All reader passwords set successfully!")
}

