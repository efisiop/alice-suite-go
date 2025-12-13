package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
)

var (
	ErrAIServiceUnavailable = errors.New("AI service unavailable")
	ErrInvalidInteractionType = errors.New("invalid interaction type")
)

// AIService handles AI interactions
type AIService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewAIService creates a new AI service
func NewAIService() *AIService {
	apiKey := os.Getenv("MOONSHOT_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_AUTH_TOKEN") // Fallback to Anthropic token
	}

	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.moonshot.cn/v1" // Default Moonshot API
	}

	return &AIService{
		apiKey:  apiKey,
		baseURL: baseURL,
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

	// Call AI API
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

// callAI calls the AI API (Moonshot/Kimi or Anthropic)
func (s *AIService) callAI(prompt string) (string, error) {
	if s.apiKey == "" {
		return "AI service not configured. Please set MOONSHOT_API_KEY or ANTHROPIC_AUTH_TOKEN environment variable.", nil
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"model": "moonshot-v1-8k", // or "kimi-v1" depending on your setup
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
		return "", err
	}

	// Create request
	req, err := http.NewRequest("POST", s.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API error: %s", string(body))
	}

	// Parse response
	var aiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(body, &aiResponse)
	if err != nil {
		return "", err
	}

	if len(aiResponse.Choices) == 0 {
		return "", errors.New("no response from AI")
	}

	return aiResponse.Choices[0].Message.Content, nil
}

// GetUserInteractions retrieves AI interactions for a user
func (s *AIService) GetUserInteractions(userID, bookID string) ([]*models.AIInteraction, error) {
	return database.GetAIInteractions(userID, bookID)
}



