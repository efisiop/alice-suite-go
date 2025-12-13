package services

import (
	"errors"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
)

var (
	ErrBookNotFound   = errors.New("book not found")
	ErrChapterNotFound = errors.New("chapter not found")
	ErrSectionNotFound = errors.New("section not found")
)

// BookService handles book-related operations
type BookService struct{}

// NewBookService creates a new book service
func NewBookService() *BookService {
	return &BookService{}
}

// GetBook retrieves a book by ID
func (s *BookService) GetBook(bookID string) (*models.Book, error) {
	book, err := database.GetBookByID(bookID)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, ErrBookNotFound
	}
	return book, nil
}

// GetAllBooks retrieves all books
func (s *BookService) GetAllBooks() ([]*models.Book, error) {
	return database.GetAllBooks()
}

// GetChapters retrieves all chapters for a book
func (s *BookService) GetChapters(bookID string) ([]*models.Chapter, error) {
	// Verify book exists
	_, err := s.GetBook(bookID)
	if err != nil {
		return nil, err
	}

	return database.GetChaptersByBookID(bookID)
}

// GetChapter retrieves a chapter by ID
func (s *BookService) GetChapter(chapterID string) (*models.Chapter, error) {
	chapter, err := database.GetChapterByID(chapterID)
	if err != nil {
		return nil, err
	}
	if chapter == nil {
		return nil, ErrChapterNotFound
	}
	return chapter, nil
}

// GetSections retrieves all sections for a chapter
func (s *BookService) GetSections(chapterID string) ([]*models.Section, error) {
	// Verify chapter exists
	_, err := s.GetChapter(chapterID)
	if err != nil {
		return nil, err
	}

	return database.GetSectionsByChapterID(chapterID)
}

// GetSection retrieves a section by ID
func (s *BookService) GetSection(sectionID string) (*models.Section, error) {
	section, err := database.GetSectionByID(sectionID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrSectionNotFound
	}
	return section, nil
}

// GetSectionByPage retrieves a section by page number
func (s *BookService) GetSectionByPage(bookID string, page int) (*models.Section, error) {
	// Verify book exists
	_, err := s.GetBook(bookID)
	if err != nil {
		return nil, err
	}

	section, err := database.GetSectionByPage(bookID, page)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrSectionNotFound
	}
	return section, nil
}

// GetPage retrieves a page by page number with all its sections
func (s *BookService) GetPage(bookID string, pageNumber int) (*models.Page, error) {
	// Verify book exists
	_, err := s.GetBook(bookID)
	if err != nil {
		return nil, err
	}

	page, err := database.GetPageByNumber(bookID, pageNumber)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, ErrSectionNotFound // Reuse error for page not found
	}
	return page, nil
}

// VerifyBookAccess verifies if a user has access to a book via verification code
func (s *BookService) VerifyBookAccess(userID, bookID, code string) (bool, error) {
	vc, err := database.VerifyCode(code, bookID)
	if err != nil {
		return false, err
	}
	if vc == nil {
		return false, nil
	}
	if vc.IsUsed {
		return false, nil
	}

	// Mark code as used
	err = database.UseVerificationCode(code, userID)
	if err != nil {
		return false, err
	}

	return true, nil
}



