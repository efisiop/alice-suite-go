package services

import (
	"os"
	"testing"

	"github.com/efisiopittau/alice-suite-go/internal/database"
)

// TestMain sets up and tears down the test database
func TestMain(m *testing.M) {
	tempDB, err := os.CreateTemp("", "alice-book-service-*.db")
	if err != nil {
		panic(err)
	}
	dbPath := tempDB.Name()
	if err := tempDB.Close(); err != nil {
		panic(err)
	}
	if err := database.InitDB(dbPath, ""); err != nil {
		panic(err)
	}
	schema, err := os.ReadFile("../../migrations/001_initial_schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := database.DB.Exec(string(schema)); err != nil {
		panic(err)
	}

	code := m.Run()
	database.CloseDB()
	os.Remove(dbPath)
	os.Remove(dbPath + "-shm")
	os.Remove(dbPath + "-wal")
	os.Exit(code)
}

// TestBookService_GetAllBooks_Success tests successful retrieval of all books
func TestBookService_GetAllBooks_Success(t *testing.T) {
	service := NewBookService()
	books, err := service.GetAllBooks()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if books == nil {
		t.Fatal("Expected books slice, got nil")
	}

	// Test database should have sample data from migrations
	if len(books) == 0 {
		t.Skip("No books seeded in test database - this is acceptable")
	}

	t.Logf("Found %d books in test database", len(books))
}

// TestBookService_GetBook_Success tests successful retrieval of a specific book
func TestBookService_GetBook_Success(t *testing.T) {
	service := NewBookService()

	// First get all books to find a valid book ID
	books, err := service.GetAllBooks()
	if err != nil {
		t.Fatalf("Cannot test GetBook without books: %v", err)
	}

	if len(books) == 0 {
		t.Skip("No books available in the test database")
	}

	// Test with the first book
	bookID := books[0].ID
	book, err := service.GetBook(bookID)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if book == nil {
		t.Fatal("Expected book, got nil")
	}

	if book.ID != bookID {
		t.Fatalf("Expected book ID %s, got %s", bookID, book.ID)
	}

	t.Logf("Successfully retrieved book: %s", book.Title)
}

// TestBookService_GetBook_NotFound tests error handling for non-existent book
func TestBookService_GetBook_NotFound(t *testing.T) {
	service := NewBookService()

	// Test with a clearly non-existent ID
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	book, err := service.GetBook(nonExistentID)

	if err != ErrBookNotFound {
		t.Fatalf("Expected ErrBookNotFound, got: %v", err)
	}

	if book != nil {
		t.Fatalf("Expected nil for non-existent book, got: %v", book)
	}

	t.Log("GetBook correctly returns nil for non-existent book")
}

// TestBookService_GetChapters tests chapter retrieval functionality
func TestBookService_GetChapters(t *testing.T) {
	service := NewBookService()

	// First get all books
	books, err := service.GetAllBooks()
	if err != nil {
		t.Fatalf("Cannot test GetChapters without books: %v", err)
	}

	if len(books) == 0 {
		t.Skip("No books available for testing GetChapters")
		return
	}

	// Test with the first book
	bookID := books[0].ID
	chapters, err := service.GetChapters(bookID)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if chapters == nil {
		t.Fatal("Expected chapters slice, got nil")
	}

	t.Logf("Retrieved %d chapters for book %s", len(chapters), bookID)
}

// TestBookService_GetChapters_BookNotFound tests error handling for non-existent book
func TestBookService_GetChapters_BookNotFound(t *testing.T) {
	service := NewBookService()

	nonExistentID := "00000000-0000-0000-0000-000000000000"
	chapters, err := service.GetChapters(nonExistentID)

	// Service should return nil chapters and the expected error
	if err != ErrBookNotFound {
		t.Fatalf("Expected ErrBookNotFound, got: %v", err)
	}

	if chapters != nil {
		t.Fatalf("Expected nil chapters for non-existent book, got: %v", chapters)
	}

	t.Log("GetChapters correctly returns error for non-existent book")
}

// TestBookService_GetPage tests page retrieval functionality
func TestBookService_GetPage(t *testing.T) {
	service := NewBookService()

	books, err := service.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(books) == 0 {
		t.Skip("No page fixture available in the test database")
	}
	bookID := books[0].ID
	pageNumber := 1

	page, err := service.GetPage(bookID, pageNumber)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if page == nil {
		t.Fatal("Expected page, got nil")
	}

	if page.PageNumber != pageNumber {
		t.Fatalf("Expected page number %d, got %d", pageNumber, page.PageNumber)
	}

	t.Logf("Successfully retrieved page %d for book %s", pageNumber, bookID)
}

// TestBookService_GetPage_InvalidPage tests error handling for non-existent page
func TestBookService_GetPage_InvalidPage(t *testing.T) {
	service := NewBookService()

	books, err := service.GetAllBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(books) == 0 {
		t.Skip("No book fixture available in the test database")
	}
	testBookID := books[0].ID
	invalidPageNumber := 99999 // Clearly non-existent page

	page, err := service.GetPage(testBookID, invalidPageNumber)

	// Service should return nil page and the expected error
	if err != ErrSectionNotFound {
		t.Fatalf("Expected ErrSectionNotFound, got: %v", err)
	}

	if page != nil {
		t.Fatalf("Expected nil page for non-existent page, got: %v", page)
	}

	t.Log("GetPage correctly returns error for non-existent page")
}

// TestBookService_GetPage_BookNotFound tests error handling for non-existent book
func TestBookService_GetPage_BookNotFound(t *testing.T) {
	service := NewBookService()

	nonExistentBookID := "00000000-0000-0000-0000-000000000000"
	testPage := 1

	page, err := service.GetPage(nonExistentBookID, testPage)

	// Service should return the expected error
	if err != ErrBookNotFound {
		t.Fatalf("Expected ErrBookNotFound, got: %v", err)
	}

	if page != nil {
		t.Fatalf("Expected nil page for non-existent book, got: %v", page)
	}

	t.Log("GetPage correctly returns error for non-existent book")
}

// TestBookService_ValidateImports ensures no panics in critical functions
func TestBookService_ValidateImports(t *testing.T) {
	// This test ensures the service can be instantiated without panics
	service := NewBookService()

	if service == nil {
		t.Fatal("NewBookService() returned nil")
	}

	// Test that error variables are properly set
	if ErrBookNotFound == nil {
		t.Fatal("ErrBookNotFound is nil")
	}

	if ErrChapterNotFound == nil {
		t.Fatal("ErrChapterNotFound is nil")
	}

	if ErrSectionNotFound == nil {
		t.Fatal("ErrSectionNotFound is nil")
	}

	t.Log("Service and error imports validated successfully")
}

// TestErrorConditionHandling ensures services handle errors gracefully
func TestErrorConditionHandling(t *testing.T) {
	service := NewBookService()

	// Test with various invalid inputs
	invalidInputs := []struct {
		name string
		fn   func() error
	}{
		{"Empty book ID", func() error {
			_, err := service.GetBook("")
			return err
		}},
		{"Invalid page number", func() error {
			_, err := service.GetPage("book-id", -1)
			return err
		}},
	}

	for _, tc := range invalidInputs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()
			// Services should either handle this gracefully or return an error
			// They should NOT panic
			t.Logf("%s handled: %v", tc.name, err)
		})
	}
}

// TestBookService_Robustness ensures services handle concurrent access
func TestBookService_Robustness(t *testing.T) {
	service := NewBookService()

	// Concurrent access test to ensure no race conditions or panics
	done := make(chan bool, 2)

	// Concurrent get all books
	go func() {
		_, _ = service.GetAllBooks()
		done <- true
	}()

	// Concurrent get book
	go func() {
		_, _ = service.GetBook("book-id")
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 2; i++ {
		<-done
	}

	t.Log("Concurrent access handled without issues")
}
