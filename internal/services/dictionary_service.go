package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
)

var (
	ErrTermNotFound = errors.New("term not found")
)

// DictionaryService handles dictionary and glossary operations
type DictionaryService struct {
	client *http.Client
}

// NewDictionaryService creates a new dictionary service
func NewDictionaryService() *DictionaryService {
	return &DictionaryService{
		client: &http.Client{
			Timeout: 10 * time.Second, // 10 second timeout for external API calls
		},
	}
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

// NormalizeWord normalizes a word for lookup (removes punctuation, handles plurals, etc.)
func (s *DictionaryService) NormalizeWord(word string) string {
	// Remove punctuation (keep hyphens and apostrophes for compound words)
	re := regexp.MustCompile(`[^\w'-]`)
	normalized := re.ReplaceAllString(word, "")
	
	// Convert to lowercase and trim
	normalized = strings.ToLower(strings.TrimSpace(normalized))
	
	return normalized
}

// LookupExternalDictionary looks up a word in external dictionary API (dictionaryapi.dev)
// Returns a DictionaryCache model that can be stored and reused
// This works from localhost, Docker, or any environment with internet access
func (s *DictionaryService) LookupExternalDictionary(word string) (*models.DictionaryCache, error) {
	normalizedWord := s.NormalizeWord(word)
	if normalizedWord == "" {
		return nil, ErrTermNotFound
	}

	// API endpoint: https://api.dictionaryapi.dev/api/v2/entries/en/{word}
	// This is a public API that works from any server (localhost, Docker, Render.com, etc.)
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", normalizedWord)
	
	// Make HTTP request - works from localhost as long as server has internet access
	resp, err := s.client.Get(url)
	if err != nil {
		// Network error - could be: no internet, firewall, DNS issue, or API down
		return nil, fmt.Errorf("failed to fetch from dictionary API (check internet connection): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrTermNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dictionary API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse API response - it returns an array of entries
	var entries []map[string]interface{}
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if len(entries) == 0 {
		return nil, ErrTermNotFound
	}

	// Extract first entry (most common usage)
	entry := entries[0]
	
	// Extract meanings array
	meanings, ok := entry["meanings"].([]interface{})
	if !ok || len(meanings) == 0 {
		return nil, ErrTermNotFound
	}

	// Get first meaning (most common)
	meaning, ok := meanings[0].(map[string]interface{})
	if !ok {
		return nil, ErrTermNotFound
	}

	// Extract definitions array
	definitions, ok := meaning["definitions"].([]interface{})
	if !ok || len(definitions) == 0 {
		return nil, ErrTermNotFound
	}

	// Get first definition (most common usage)
	definitionObj, ok := definitions[0].(map[string]interface{})
	if !ok {
		return nil, ErrTermNotFound
	}

	// Extract definition text
	definition, _ := definitionObj["definition"].(string)
	if definition == "" {
		return nil, ErrTermNotFound
	}

	// Collect ALL examples from multiple definitions across ALL meanings
	// The API can have multiple examples, and we want to show them all for better context
	var allExamples []string
	seenExamples := make(map[string]bool) // Track duplicates
	
	// Check all meanings (different parts of speech can have different examples)
	for _, meaningInterface := range meanings {
		meaningMap, ok := meaningInterface.(map[string]interface{})
		if !ok {
			continue
		}
		
		definitionsList, ok := meaningMap["definitions"].([]interface{})
		if !ok {
			continue
		}
		
		// Check each definition in this meaning for examples
		for _, defInterface := range definitionsList {
			def, ok := defInterface.(map[string]interface{})
			if !ok {
				continue
			}
			
			if exampleStr, ok := def["example"].(string); ok && exampleStr != "" {
				exampleStr = strings.TrimSpace(exampleStr)
				// Avoid duplicates (case-insensitive)
				exampleLower := strings.ToLower(exampleStr)
				if !seenExamples[exampleLower] && len(allExamples) < 5 { // Limit to 5 examples max
					allExamples = append(allExamples, exampleStr)
					seenExamples[exampleLower] = true
				}
			}
		}
	}
	
	// Join all examples with a special separator that we can split in the frontend
	// Using " |||| " as separator (unlikely to appear in text)
	var example string
	if len(allExamples) > 0 {
		example = strings.Join(allExamples, " |||| ")
	}

	// Extract part of speech
	var partOfSpeech string
	if pos, ok := meaning["partOfSpeech"].(string); ok {
		partOfSpeech = pos
	}

	// Extract phonetic (if available)
	var phonetic string
	if phonetics, ok := entry["phonetics"].([]interface{}); ok && len(phonetics) > 0 {
		if ph, ok := phonetics[0].(map[string]interface{}); ok {
			if text, ok := ph["text"].(string); ok {
				phonetic = text
			}
		}
	}

	// Get word from API (preserves original casing)
	apiWord, _ := entry["word"].(string)
	if apiWord == "" {
		apiWord = word // Fallback to original word
	}

	cache := &models.DictionaryCache{
		Word:         normalizedWord,
		Definition:   definition,
		Example:      example,
		Phonetic:     phonetic,
		PartOfSpeech: partOfSpeech,
		SourceAPI:    "dictionaryapi.dev",
	}

	return cache, nil
}

// LookupWordInContext looks up a word and provides context from the book
// Strategy: 1) Glossary (technical terms), 2) Cache (previously fetched), 3) External API (common words)
func (s *DictionaryService) LookupWordInContext(bookID, word string, chapterID, sectionID *string) (*models.AliceGlossary, error) {
	normalizedWord := s.NormalizeWord(word)
	
	// Step 1: Try glossary first (prioritize glossary definitions for technical terms)
	term, err := s.LookupWord(bookID, normalizedWord)
	if err == nil && term != nil {
		return term, nil
	}

	// Step 2: Check cache (previously fetched definitions)
	cached, err := database.GetCachedDefinition(normalizedWord)
	if err == nil && cached != nil {
		// Convert cache to AliceGlossary format
		glossaryTerm := &models.AliceGlossary{
			Term:       word, // Preserve original word casing
			Definition: cached.Definition,
			Example:    cached.Example,
			BookID:     bookID,
		}
		if chapterID != nil {
			glossaryTerm.ChapterReference = *chapterID
		}
		return glossaryTerm, nil
	}

	// Step 3: Fetch from external API (common words)
	cache, err := s.LookupExternalDictionary(normalizedWord)
	if err == nil && cache != nil {
		// Cache the result for future lookups
		if cacheErr := database.CacheDefinition(cache); cacheErr != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to cache definition for %s: %v\n", normalizedWord, cacheErr)
		}

		// Convert cache to AliceGlossary format
		glossaryTerm := &models.AliceGlossary{
			Term:       word, // Preserve original word casing
			Definition: cache.Definition,
			Example:    cache.Example,
			BookID:     bookID,
		}
		if chapterID != nil {
			glossaryTerm.ChapterReference = *chapterID
		}
		return glossaryTerm, nil
	}

	// Word not found in glossary, cache, or external API
	return nil, ErrTermNotFound
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



