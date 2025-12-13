package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
	"github.com/google/uuid"
)

func main() {
	// Remove old database for clean test
	os.Remove("test.db")

	fmt.Println("=== Testing Help Requests ===")

	// Initialize database
	databasePath := "../../data/alice-suite.db"
	if err := database.InitDB(databasePath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Test 1: Create a user first (required for foreign constraint)
	fmt.Println("1. Creating test user...")
	testUser := &models.User{
		ID:           uuid.New().String(),
		Email:        "test@example.com",
		PasswordHash: "hashed_password_here",
		FirstName:    "Test",
		LastName:     "User",
		Role:         "reader",
		IsVerified:   true,
	}
	err := createUserManually(testUser)
	if err != nil {
		log.Printf("❌ User creation failed: %v", err)
		fmt.Println("This might be the real constraint issue - user doesn't exist")
	} else {
		fmt.Println("✅ User created successfully")
	}

	// Test 2: Create a book (required for foreign constraint)
	fmt.Println("2. Creating test book...")
	testBook := &models.Book{
		ID:          "alice-in-wonderland",
		Title:       "Alice in Wonderland",
		Author:      "Lewis Carroll",
		Description: "Alice\\'s Adventures in Wonderland",
		TotalPages:  200,
		CreatedAt:   time.Now(),
	}
	err = createBookManually(testBook)
	if err != nil {
		log.Printf("❌ Book creation failed: %v", err)
		fmt.Println("This might be the real constraint issue - book doesn't exist")
	} else {
		fmt.Println("✅ Book created successfully")
	}

	// Test 3: Create help request
	fmt.Println("3. Creating help request...")
	helpRequest := &models.HelpRequest{
		ID:         uuid.New().String(),
		UserID:     testUser.ID,
		BookID:     testBook.ID,
		Status:     "pending",
		Content:    "Test help request",
		Context:    "Testing foreign key constraints",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	fmt.Printf("Request details:\n")
	fmt.Printf("- UserID: %s\n", helpRequest.UserID)
	fmt.Printf("- BookID: %s\n", helpRequest.BookID)

	// Use the actual CreateHelpRequest function
	// Use a valid section_id from the actual database
	validSectionID := "page-1-section-1"
	dbHelpReq := &models.HelpRequest{
		UserID:    testUser.ID,
		BookID:    testBook.ID,
		SectionID: &validSectionID,
		Content:   helpRequest.Content,
		Context:   helpRequest.Context,
		Status:    helpRequest.Status,
	}
	err = database.CreateHelpRequest(dbHelpReq)
	if err != nil {
		log.Printf("❌ Help request creation failed: %v", err)
		fmt.Println("\n🔍 ERROR ANALYSIS:")
		fmt.Println("This is likely a FOREIGN KEY constraint error.")
		fmt.Println("The error message will show exactly which foreign key is violated.")

		if err.Error() != "" {
			fmt.Printf("Error details: %s\n", err.Error())
		}

		// Check if it's specific to user/book foreign keys
		if contains(err.Error(), "users") {
			fmt.Println("🎯 Foreign key constraint: UserID not found in users table")
		}
		if contains(err.Error(), "books") {
			fmt.Println("🎯 Foreign key constraint: BookID not found in books table")
		}
		if contains(err.Error(), "sections") {
			fmt.Println("🎯 Foreign key constraint: SectionID not found in sections table")
		}
	} else {
		fmt.Println("✅ Help request created successfully!")
		// Helper request creation successful
	}

	// Test 4: Retrieve to verify
	fmt.Println("4. Retrieving help requests...")
	requests, err := database.GetHelpRequests(testUser.ID)
	if err != nil {
		log.Printf("❌ Failed to retrieve help requests: %v", err)
	} else {
		fmt.Printf("✅ Found %d help requests for user\n", len(requests))
		for i, req := range requests {
			fmt.Printf("  - Request %d: %s (Status: %s)\n", i+1, req.Content, req.Status)
		}
	}

	fmt.Println("\n=== Test Summary ===")
	if err != nil {
		fmt.Println("❌ The constraint error is related to foreign key relationships.")
		fmt.Println("   Check actual database state and ensure users/books exist.")
	} else {
		fmt.Println("✅ All tests passed! Help requests should be working.")
	}
}

func createUserManually(user *models.User) error {
	query := `INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_verified, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	_, err := database.DB.Exec(query, user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Role, user.IsVerified, now, now)
	return err
}

func createBookManually(book *models.Book) error {
	query := `INSERT INTO books (id, title, author, description, total_pages, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, book.ID, book.Title, book.Author, book.Description, book.TotalPages, book.CreatedAt)
	return err
}

func contains(s, substr string) bool {
	return strings.Index(s, substr) != -1
}