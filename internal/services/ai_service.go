package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
)

var (
	ErrAIServiceUnavailable = errors.New("AI service unavailable")
	ErrInvalidInteractionType = errors.New("invalid interaction type")
)

// AIProvider represents which AI service to use
type AIProvider string

const (
	ProviderGemini   AIProvider = "gemini"
	ProviderMoonshot AIProvider = "moonshot"
	ProviderAuto     AIProvider = "auto" // Try Gemini first, fallback to Moonshot
)

// AIService handles AI interactions
type AIService struct {
	provider     AIProvider
	geminiKey    string
	moonshotKey  string
	moonshotURL  string
	client       *http.Client
}

// NewAIService creates a new AI service with support for multiple providers
func NewAIService() *AIService {
	// Get provider preference (default: auto = try Gemini first, fallback to Moonshot)
	providerStr := os.Getenv("AI_PROVIDER")
	if providerStr == "" {
		providerStr = "auto"
	}
	provider := AIProvider(providerStr)
	
	// Get API keys
	geminiKey := os.Getenv("GEMINI_API_KEY")
	moonshotKey := os.Getenv("MOONSHOT_API_KEY")
	if moonshotKey == "" {
		moonshotKey = os.Getenv("ANTHROPIC_AUTH_TOKEN") // Fallback to old env var name
	}

	moonshotURL := os.Getenv("MOONSHOT_BASE_URL") // Use MOONSHOT_BASE_URL instead of ANTHROPIC_BASE_URL
	if moonshotURL == "" {
		moonshotURL = os.Getenv("ANTHROPIC_BASE_URL") // Fallback to old env var name
	}
	// Validate and fix incorrect Moonshot URLs
	if moonshotURL != "" {
		// Fix common incorrect URLs
		if strings.Contains(moonshotURL, "moonshot.ai") || strings.Contains(moonshotURL, "/anthropic") {
			log.Printf("Warning: ANTHROPIC_BASE_URL or MOONSHOT_BASE_URL is set to incorrect value: %s. Using default Moonshot API URL instead.", moonshotURL)
			moonshotURL = "https://api.moonshot.cn/v1"
		}
	}
	if moonshotURL == "" {
		moonshotURL = "https://api.moonshot.cn/v1" // Default Moonshot API (correct URL)
	}

	return &AIService{
		provider:    provider,
		geminiKey:   geminiKey,
		moonshotKey: moonshotKey,
		moonshotURL: moonshotURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// InteractionType represents the type of AI interaction
type InteractionType string

const (
	InteractionExplain   InteractionType = "explain"
	InteractionQuiz      InteractionType = "quiz"
	InteractionSimplify  InteractionType = "simplify"
	InteractionDefinition InteractionType = "definition"
	InteractionChat      InteractionType = "chat"
)

// AskAI sends a question to the AI and returns a response
func (s *AIService) AskAI(userID, bookID string, interactionType InteractionType, question string, sectionID *string, context string) (*models.AIInteraction, error) {
	// Validate interaction type
	validTypes := map[InteractionType]bool{
		InteractionExplain:    true,
		InteractionQuiz:      true,
		InteractionSimplify:  true,
		InteractionDefinition: true,
		InteractionChat:      true,
	}
	if !validTypes[interactionType] {
		return nil, ErrInvalidInteractionType
	}

	// Build prompt based on interaction type
	prompt := s.buildPrompt(interactionType, question, context)

	// Call AI API (with automatic fallback if using "auto" provider)
	response, err := s.callAI(prompt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAIServiceUnavailable, err)
	}

	// Create interaction record
	interaction := &models.AIInteraction{
		UserID:         userID,
		BookID:         bookID,
		SectionID:      sectionID,
		InteractionType: string(interactionType),
		Question:       question,
		Prompt:         prompt,
		Response:       response,
		Context:        context,
	}

	// Save to database
	err = database.CreateAIInteraction(interaction)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to save AI interaction: %v\n", err)
	}

	return interaction, nil
}

// buildPrompt builds a prompt based on interaction type
func (s *AIService) buildPrompt(interactionType InteractionType, question, context string) string {
	basePrompt := "You are a helpful reading assistant for Alice's Adventures in Wonderland. "
	basePrompt += "This is a physical book companion app - users read from their physical book and use this app for assistance.\n\n"

	switch interactionType {
	case InteractionExplain:
		return basePrompt + fmt.Sprintf("Please explain the following passage or concept: %s\n\nContext: %s", question, context)
	case InteractionQuiz:
		return basePrompt + fmt.Sprintf("Create a quiz question about: %s\n\nContext: %s", question, context)
	case InteractionSimplify:
		return basePrompt + fmt.Sprintf("Please simplify or rephrase the following: %s\n\nContext: %s", question, context)
	case InteractionDefinition:
		return basePrompt + fmt.Sprintf("Please provide a definition for: %s\n\nContext: %s", question, context)
	case InteractionChat:
		return basePrompt + fmt.Sprintf("Question: %s\n\nContext: %s", question, context)
	default:
		return basePrompt + question
	}
}

// callAI calls the AI API (Gemini or Moonshot) with automatic fallback
func (s *AIService) callAI(prompt string) (string, error) {
	// Determine which provider(s) to try
	var providers []AIProvider
	
	switch s.provider {
	case ProviderGemini:
		providers = []AIProvider{ProviderGemini}
	case ProviderMoonshot:
		providers = []AIProvider{ProviderMoonshot}
	case ProviderAuto:
		// Try Gemini first, then Moonshot
		providers = []AIProvider{ProviderGemini, ProviderMoonshot}
	default:
		providers = []AIProvider{ProviderGemini, ProviderMoonshot} // Default to auto behavior
	}

	// Try each provider in order
	var lastErr error
	for _, provider := range providers {
		var response string
		var err error
		
		switch provider {
		case ProviderGemini:
			response, err = s.callGemini(prompt)
		case ProviderMoonshot:
			response, err = s.callMoonshot(prompt)
		default:
			continue
		}
		
		if err == nil && response != "" {
			return response, nil
		}
		lastErr = err
	}

	// All providers failed
	if lastErr != nil {
		return "", lastErr
	}
	
	return "", fmt.Errorf("no AI provider configured. Please set GEMINI_API_KEY or MOONSHOT_API_KEY environment variable")
}

// callGemini calls the Google Gemini API
func (s *AIService) callGemini(prompt string) (string, error) {
	if s.geminiKey == "" {
		return "", errors.New("GEMINI_API_KEY not set")
	}

	// Gemini API endpoint (using gemini-1.5-flash for faster responses and free tier)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", s.geminiKey)

	// Gemini uses a different request format
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.7,
			"maxOutputTokens": 1000,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Gemini request: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gemini API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Gemini response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse Gemini response (different structure)
	var geminiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &geminiResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	if geminiResponse.Error != nil {
		return "", fmt.Errorf("Gemini API error: %s", geminiResponse.Error.Message)
	}

	if len(geminiResponse.Candidates) == 0 || len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response from Gemini API")
	}

	return geminiResponse.Candidates[0].Content.Parts[0].Text, nil
}

// callMoonshot calls the Moonshot/Kimi API
func (s *AIService) callMoonshot(prompt string) (string, error) {
	if s.moonshotKey == "" {
		return "", errors.New("MOONSHOT_API_KEY not set")
	}

	// Moonshot API endpoint
	url := s.moonshotURL + "/chat/completions"

	// Moonshot uses OpenAI-compatible format
	payload := map[string]interface{}{
		"model": "moonshot-v1-8k", // or "moonshot-v1-32k" for longer context
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Moonshot request: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create Moonshot request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.moonshotKey)

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Moonshot API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Moonshot response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Moonshot API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse Moonshot response (OpenAI-compatible format)
	var moonshotResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &moonshotResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse Moonshot response: %w", err)
	}

	if moonshotResponse.Error != nil {
		return "", fmt.Errorf("Moonshot API error: %s", moonshotResponse.Error.Message)
	}

	if len(moonshotResponse.Choices) == 0 {
		return "", errors.New("no response from Moonshot API")
	}

	return moonshotResponse.Choices[0].Message.Content, nil
}

// GetUserInteractions retrieves AI interactions for a user
func (s *AIService) GetUserInteractions(userID, bookID string) ([]*models.AIInteraction, error) {
	return database.GetAIInteractions(userID, bookID)
}



