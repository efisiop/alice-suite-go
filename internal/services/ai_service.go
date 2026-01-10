package services

import (
	"bytes"
	"crypto/tls"
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

	// Create HTTP client with custom TLS config for Moonshot (handles certificate issues)
	// Note: In production, you might want to verify certificates properly
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// For Moonshot TLS certificate issues, we can skip verification if needed
	// This is a workaround for the "x509: negative serial number" error
	// Only use this if Moonshot's certificate has issues
	if os.Getenv("MOONSHOT_SKIP_TLS_VERIFY") == "true" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}
	}

	return &AIService{
		provider:    provider,
		geminiKey:   geminiKey,
		moonshotKey: moonshotKey,
		moonshotURL: moonshotURL,
		client:      client,
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
	response, providerUsed, err := s.callAI(prompt)
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
		Provider:       string(providerUsed),
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
	basePrompt += "IMPORTANT: Please provide complete, finished answers. Do not cut off mid-sentence or leave incomplete thoughts.\n\n"

	switch interactionType {
	case InteractionExplain:
		return basePrompt + fmt.Sprintf("Please provide a complete explanation of the following passage or concept: %s\n\nContext: %s\n\nMake sure your explanation is thorough and complete.", question, context)
	case InteractionQuiz:
		return basePrompt + fmt.Sprintf("Create a complete quiz question about: %s\n\nContext: %s\n\nInclude the question, answer options, and the correct answer.", question, context)
	case InteractionSimplify:
		return basePrompt + fmt.Sprintf("Please provide a complete, simplified version or rephrasing of: %s\n\nContext: %s\n\nMake sure your simplification is complete and clear.", question, context)
	case InteractionDefinition:
		return basePrompt + fmt.Sprintf("Please provide a complete definition for: %s\n\nContext: %s\n\nMake sure your definition is comprehensive and finished.", question, context)
	case InteractionChat:
		return basePrompt + fmt.Sprintf("Question: %s\n\nContext: %s\n\nPlease provide a complete, helpful answer to the question.", question, context)
	default:
		return basePrompt + fmt.Sprintf("Question: %s\n\nPlease provide a complete, helpful answer.", question)
	}
}

// callAI calls the AI API (Gemini or Moonshot) with automatic fallback
// Returns: response, provider used, error
func (s *AIService) callAI(prompt string) (string, AIProvider, error) {
	// Determine which provider(s) to try based on configured keys
	var providers []AIProvider
	
	switch s.provider {
		case ProviderGemini:
		if s.geminiKey != "" {
			providers = []AIProvider{ProviderGemini}
		} else {
			return "", "", errors.New("GEMINI_API_KEY not set but provider is set to 'gemini'")
		}
	case ProviderMoonshot:
		if s.moonshotKey != "" {
			providers = []AIProvider{ProviderMoonshot}
		} else {
			return "", "", errors.New("MOONSHOT_API_KEY not set but provider is set to 'moonshot'")
		}
	case ProviderAuto:
		// Try Gemini first if key is available, then Moonshot if key is available
		if s.geminiKey != "" {
			providers = append(providers, ProviderGemini)
		}
		if s.moonshotKey != "" {
			providers = append(providers, ProviderMoonshot)
		}
		// If neither key is set, return error
		if len(providers) == 0 {
			return "", "", fmt.Errorf("no AI provider configured. Please set GEMINI_API_KEY or MOONSHOT_API_KEY environment variable")
		}
	default:
		// Default to auto behavior
		if s.geminiKey != "" {
			providers = append(providers, ProviderGemini)
		}
		if s.moonshotKey != "" {
			providers = append(providers, ProviderMoonshot)
		}
		if len(providers) == 0 {
			return "", "", fmt.Errorf("no AI provider configured. Please set GEMINI_API_KEY or MOONSHOT_API_KEY environment variable")
		}
	}

	// Try each provider in order (only those with keys configured)
	var lastErr error
	for _, provider := range providers {
		var response string
		var err error
		
		switch provider {
		case ProviderGemini:
			log.Printf("Trying Gemini API...")
			response, err = s.callGemini(prompt)
			if err != nil {
				log.Printf("Gemini API failed: %v", err)
			}
		case ProviderMoonshot:
			log.Printf("Trying Moonshot API...")
			response, err = s.callMoonshot(prompt)
			if err != nil {
				log.Printf("Moonshot API failed: %v", err)
			}
		default:
			continue
		}
		
		if err == nil && response != "" {
			log.Printf("AI API call successful using %s", provider)
			return response, provider, nil
		}
		lastErr = err
	}

	// All providers failed
	if lastErr != nil {
		return "", "", fmt.Errorf("all AI providers failed. Last error: %w", lastErr)
	}
	
	return "", "", fmt.Errorf("no AI provider configured. Please set GEMINI_API_KEY or MOONSHOT_API_KEY environment variable")
}

// callGemini calls the Google Gemini API
func (s *AIService) callGemini(prompt string) (string, error) {
	if s.geminiKey == "" {
		return "", errors.New("GEMINI_API_KEY not set")
	}

	// First, try to get available models from the API
	availableModels, err := s.listGeminiModels()
	if err != nil {
		log.Printf("Warning: Could not list available Gemini models: %v. Using fallback model list.", err)
	}

	// Try multiple model names in order (fallback if one doesn't work)
	// Start with models from the API if available, then fallback to known models
	modelNames := []string{}
	
	// Add available models first
	for _, model := range availableModels {
		modelNames = append(modelNames, model)
	}
	
	// Add fallback models
	fallbackModels := []string{
		"gemini-1.5-flash-001",
		"gemini-1.5-pro-002",
		"gemini-1.5-pro-001",
		"gemini-1.5-flash",
		"gemini-1.5-pro",
	}
	for _, model := range fallbackModels {
		// Only add if not already in the list
		found := false
		for _, existing := range modelNames {
			if existing == model {
				found = true
				break
			}
		}
		if !found {
			modelNames = append(modelNames, model)
		}
	}

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
			"maxOutputTokens": 4096, // Increased from 1000 to allow complete responses
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Gemini request: %w", err)
	}

	// Try each model name until one works
	var lastErr error
	for _, modelName := range modelNames {
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/%s:generateContent?key=%s", modelName, s.geminiKey)

		// Create request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("failed to create Gemini request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("Gemini API request failed: %w", err)
			continue
		}

		// Read response
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read Gemini response: %w", err)
			continue
		}

		// If successful, parse and return
		if resp.StatusCode == http.StatusOK {
			return s.parseGeminiResponse(body)
		}

		// If 404, try next model; otherwise return error
		if resp.StatusCode != http.StatusNotFound {
			return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
		}

		lastErr = fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
		log.Printf("Model %s not available, trying next model...", modelName)
	}

	// All models failed
	return "", fmt.Errorf("all Gemini models failed. Last error: %w", lastErr)
}

// listGeminiModels queries the Gemini API to get a list of available models
func (s *AIService) listGeminiModels() ([]string, error) {
	if s.geminiKey == "" {
		return nil, errors.New("GEMINI_API_KEY not set")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models?key=%s", s.geminiKey)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models (status %d): %s", resp.StatusCode, string(body))
	}

	var modelsResponse struct {
		Models []struct {
			Name         string   `json:"name"`
			SupportedMethods []string `json:"supportedGenerationMethods"`
		} `json:"models"`
	}

	err = json.Unmarshal(body, &modelsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse models response: %w", err)
	}

	var availableModels []string
	for _, model := range modelsResponse.Models {
		// Check if model supports generateContent
		supportsGenerateContent := false
		for _, method := range model.SupportedMethods {
			if method == "generateContent" {
				supportsGenerateContent = true
				break
			}
		}
		
		if supportsGenerateContent {
			// Extract just the model name (format is "models/gemini-1.5-pro")
			parts := strings.Split(model.Name, "/")
			if len(parts) > 0 {
				modelName := parts[len(parts)-1]
				availableModels = append(availableModels, modelName)
			}
		}
	}

	return availableModels, nil
}

// parseGeminiResponse parses the Gemini API response
func (s *AIService) parseGeminiResponse(body []byte) (string, error) {
	// Parse Gemini response (different structure)
	var geminiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"` // "STOP", "MAX_TOKENS", "SAFETY", etc.
		} `json:"candidates"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	err := json.Unmarshal(body, &geminiResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	if geminiResponse.Error != nil {
		return "", fmt.Errorf("Gemini API error: %s", geminiResponse.Error.Message)
	}

	if len(geminiResponse.Candidates) == 0 || len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response from Gemini API")
	}

	responseText := geminiResponse.Candidates[0].Content.Parts[0].Text
	
	// Check if response was truncated due to token limit
	if geminiResponse.Candidates[0].FinishReason == "MAX_TOKENS" {
		log.Printf("Warning: Gemini response may be incomplete (MAX_TOKENS finish reason)")
		// Optionally append a note or try to extend the response
		// For now, just log it - the response should still be useful
	}
	
	return responseText, nil
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
		"max_tokens":  4096, // Increased from 1000 to allow complete responses
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
			FinishReason string `json:"finish_reason"` // "stop", "length", "content_filter", etc.
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

	responseText := moonshotResponse.Choices[0].Message.Content
	
	// Check if response was truncated due to token limit
	if moonshotResponse.Choices[0].FinishReason == "length" {
		log.Printf("Warning: Moonshot response may be incomplete (length finish reason)")
		// The response should still be useful, just log it
	}
	
	return responseText, nil
}

// GetUserInteractions retrieves AI interactions for a user
func (s *AIService) GetUserInteractions(userID, bookID string) ([]*models.AIInteraction, error) {
	return database.GetAIInteractions(userID, bookID)
}

// GetProviderStatus returns information about the current AI provider configuration
func (s *AIService) GetProviderStatus() map[string]interface{} {
	status := map[string]interface{}{
		"configured_provider": string(s.provider),
		"gemini_configured":   s.geminiKey != "",
		"moonshot_configured": s.moonshotKey != "",
	}

	// Determine which providers are available
	var availableProviders []string
	if s.geminiKey != "" {
		availableProviders = append(availableProviders, "gemini")
	}
	if s.moonshotKey != "" {
		availableProviders = append(availableProviders, "moonshot")
	}
	status["available_providers"] = availableProviders

	// Determine which provider will be used based on configuration
	var activeProvider string
	switch s.provider {
	case ProviderGemini:
		if s.geminiKey != "" {
			activeProvider = "gemini"
		} else {
			activeProvider = "none (gemini key not set)"
		}
	case ProviderMoonshot:
		if s.moonshotKey != "" {
			activeProvider = "moonshot"
		} else {
			activeProvider = "none (moonshot key not set)"
		}
	case ProviderAuto:
		if s.geminiKey != "" {
			activeProvider = "gemini (auto mode - will try moonshot if gemini fails)"
		} else if s.moonshotKey != "" {
			activeProvider = "moonshot (auto mode - gemini not configured)"
		} else {
			activeProvider = "none (no API keys configured)"
		}
	default:
		// Default to auto behavior
		if s.geminiKey != "" {
			activeProvider = "gemini (auto mode - will try moonshot if gemini fails)"
		} else if s.moonshotKey != "" {
			activeProvider = "moonshot (auto mode - gemini not configured)"
		} else {
			activeProvider = "none (no API keys configured)"
		}
	}
	status["active_provider"] = activeProvider

	return status
}



