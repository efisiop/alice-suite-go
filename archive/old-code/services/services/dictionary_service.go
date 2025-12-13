package services

import (
	"errors"
	"strings"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
)

var (
	ErrTermNotFound = errors.New("term not found")
)

// DictionaryService handles dictionary and glossary operations
type DictionaryService struct{}

// NewDictionaryService creates a new dictionary service
func NewDictionaryService() *DictionaryService {
	return &DictionaryService{}
}

// LookupWord looks up a word in the Alice glossary
// This function prioritizes glossary definitions over external dictionaries
func (s *DictionaryService) LookupWord(bookID, word string) (*models.AliceGlossary, error) {
	// Normalize word (lowercase, trim)
	normalizedWord := strings.ToLower(strings.TrimSpace(word))
	
	// Try exact match first
	term, err := database.GetGlossaryTerm(bookID, normalizedWord)
	if err != nil {
		return nil, err
	}
	if term != nil {
		return term, nil
	}
	
	// Try case-insensitive search (handles variations like "Rabbit" vs "rabbit")
	term, err = database.FindGlossaryTermInText(bookID, normalizedWord)
	if err != nil {
		return nil, err
	}
	if term != nil {
		return term, nil
	}
	
	return nil, ErrTermNotFound
}

// SearchTerms searches for glossary terms
func (s *DictionaryService) SearchTerms(bookID, searchTerm string) ([]*models.AliceGlossary, error) {
	return database.SearchGlossaryTerms(bookID, searchTerm)
}

// LookupWordInContext looks up a word and provides context from the book
// This prioritizes glossary definitions and can provide section-specific context
func (s *DictionaryService) LookupWordInContext(bookID, word string, chapterID, sectionID *string) (*models.AliceGlossary, error) {
	// First try to find in glossary (prioritize glossary definitions)
	term, err := s.LookupWord(bookID, word)
	if err == nil && term != nil {
		return term, nil
	}

	// If not found in glossary, create a basic definition
	// In production, this would call an external dictionary API
	term = &models.AliceGlossary{
		Term:       word,
		Definition: "Definition not available. This word may need to be added to the glossary.",
		BookID:     bookID,
	}
	if chapterID != nil {
		term.ChapterReference = *chapterID
	}

	return term, nil
}

// GetGlossaryTermsForSection gets all glossary terms linked to a specific section
func (s *DictionaryService) GetGlossaryTermsForSection(sectionID string) ([]*models.AliceGlossary, error) {
	return database.GetGlossaryTermBySection(sectionID)
}

// GetGlossaryTermsForPageSection gets all glossary terms for a specific page and section
func (s *DictionaryService) GetGlossaryTermsForPageSection(bookID string, pageNumber, sectionNumber int) ([]*models.AliceGlossary, error) {
	return database.GetGlossaryTermByPageAndSection(bookID, pageNumber, sectionNumber)
}

// RecordLookup records a vocabulary lookup for analytics
func (s *DictionaryService) RecordLookup(userID, bookID, word, definition string, chapterID, sectionID *string, context string) error {
	lookup := &models.VocabularyLookup{
		UserID:     userID,
		BookID:     bookID,
		Word:       word,
		Definition: definition,
		ChapterID:  chapterID,
		SectionID:  sectionID,
		Context:    context,
	}

	return database.CreateVocabularyLookup(lookup)
}

// GetUserLookups retrieves vocabulary lookups for a user
func (s *DictionaryService) GetUserLookups(userID, bookID string) ([]*models.VocabularyLookup, error) {
	return database.GetVocabularyLookups(userID, bookID)
}



