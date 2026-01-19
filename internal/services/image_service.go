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
	"net/url"
	"os"
	"strings"
	"time"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// containsPencilKeywords checks if prompt already contains pencil/sketch keywords
func containsPencilKeywords(prompt string) bool {
	keywords := []string{"pencil", "sketch", "charcoal", "graphite", "line art", "drawing"}
	lowerPrompt := strings.ToLower(prompt)
	for _, keyword := range keywords {
		if strings.Contains(lowerPrompt, keyword) {
			return true
		}
	}
	return false
}

var (
	ErrImageServiceUnavailable = errors.New("image generation service unavailable")
	ErrImageGenerationFailed   = errors.New("image generation failed")
	ErrImageTaskNotFound       = errors.New("image task not found")
)

// ImageProvider represents which image generation service to use
type ImageProvider string

const (
	ProviderFreepik ImageProvider = "freepik"
	ProviderDeepAI  ImageProvider = "deepai"
	ProviderReplicate ImageProvider = "replicate"
)

// ImageService handles image generation via various APIs (Freepik, DeepAI, etc.)
type ImageService struct {
	provider   ImageProvider
	apiKey     string
	client     *http.Client
	baseURL    string
}

// NewImageService creates a new image generation service
func NewImageService() *ImageService {
	// Get provider preference (default: deepai for better free tier)
	providerStr := os.Getenv("IMAGE_PROVIDER")
	if providerStr == "" {
		providerStr = "deepai" // Default to DeepAI for better free tier
	}
	provider := ImageProvider(providerStr)

	var apiKey string
	var baseURL string

	// Get API key based on provider
	switch provider {
	case ProviderFreepik:
		apiKey = os.Getenv("FREEPIK_API_KEY")
		baseURL = "https://api.freepik.com/v1"
		if apiKey == "" {
			log.Printf("Warning: FREEPIK_API_KEY not set. Image generation will be disabled.")
		}
	case ProviderDeepAI:
		apiKey = os.Getenv("DEEPAI_API_KEY")
		baseURL = "https://api.deepai.org"
		if apiKey == "" {
			log.Printf("Warning: DEEPAI_API_KEY not set. Image generation will be disabled.")
		}
	case ProviderReplicate:
		apiKey = os.Getenv("REPLICATE_API_TOKEN")
		baseURL = "https://api.replicate.com/v1"
		if apiKey == "" {
			log.Printf("Warning: REPLICATE_API_TOKEN not set. Image generation will be disabled.")
		}
	default:
		log.Printf("Warning: Unknown IMAGE_PROVIDER '%s', defaulting to deepai", providerStr)
		provider = ProviderDeepAI
		apiKey = os.Getenv("DEEPAI_API_KEY")
		baseURL = "https://api.deepai.org"
	}

	// Create HTTP client with appropriate timeout
	client := &http.Client{
		Timeout: 120 * time.Second, // 2 minutes for image generation
	}

	// For Freepik TLS certificate issues, skip verification
	if provider == ProviderFreepik && os.Getenv("FREEPIK_SKIP_TLS_VERIFY") != "false" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
			Timeout:   120 * time.Second,
		}
		log.Printf("Freepik TLS verification disabled (certificate workaround)")
	}

	log.Printf("Image service initialized with provider: %s", provider)

	return &ImageService{
		provider: provider,
		apiKey:   apiKey,
		client:   client,
		baseURL:  baseURL,
	}
}

// ImageGenerationRequest represents a request to generate an image
type ImageGenerationRequest struct {
	Prompt      string `json:"prompt"`
	AspectRatio string `json:"aspect_ratio,omitempty"` // e.g., "widescreen_16_9", "square", "portrait_4_3"
	Model       string `json:"model,omitempty"`        // e.g., "mystic", "flux", "seedream", "ideogram"
	Resolution  string `json:"resolution,omitempty"`   // e.g., "1k", "2k", "4k" (for Mystic)
}

// ImageGenerationResponse represents the response from Freepik API
type ImageGenerationResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
}

// ImageTaskStatus represents the status of an image generation task
type ImageTaskStatus struct {
	TaskID   string `json:"task_id"`
	Status   string `json:"status"` // "pending", "processing", "completed", "failed"
	ImageURL string `json:"image_url,omitempty"`
	Error    string `json:"error,omitempty"`
}

// GenerateImage submits an image generation request to the configured provider
// Returns task ID that can be used to check status (for async providers) or direct image URL (for sync providers)
func (s *ImageService) GenerateImage(req ImageGenerationRequest) (*ImageGenerationResponse, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("%s API key not configured", s.provider)
	}

	// Route to appropriate provider
	switch s.provider {
	case ProviderDeepAI:
		return s.generateImageDeepAI(req)
	case ProviderFreepik:
		return s.generateImageFreepik(req)
	case ProviderReplicate:
		return s.generateImageReplicate(req)
	default:
		return nil, fmt.Errorf("unsupported image provider: %s", s.provider)
	}
}

// generateImageDeepAI generates an image using DeepAI API (synchronous)
func (s *ImageService) generateImageDeepAI(req ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// DeepAI uses form data instead of JSON
	apiURL := fmt.Sprintf("%s/api/text2img", s.baseURL)
	
	// Prepare form data
	formData := map[string]string{
		"text": req.Prompt,
	}
	
	// Add optional parameters
	if req.AspectRatio != "" {
		// DeepAI uses width/height, but we can use aspect ratio hints
		// Default to 512x512 for square, but we can adjust
		if req.AspectRatio == "square" || req.AspectRatio == "square_1_1" {
			formData["width"] = "512"
			formData["height"] = "512"
		} else {
			// Default to 512x512 for simplicity
			formData["width"] = "512"
			formData["height"] = "512"
		}
	} else {
		formData["width"] = "512"
		formData["height"] = "512"
	}

	// Add "pencil sketch" style hint to prompt for better results
	enhancedPrompt := req.Prompt
	if !containsPencilKeywords(enhancedPrompt) {
		enhancedPrompt = "minimalist pencil sketch, " + enhancedPrompt + ", graphite line art style"
	}
	formData["text"] = enhancedPrompt

	log.Printf("DeepAI API: Generating image with prompt: %s", enhancedPrompt[:min(100, len(enhancedPrompt))])

	// Create form data body (properly formatted)
	formValues := url.Values{}
	for k, v := range formData {
		formValues.Set(k, v)
	}
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formValues.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("api-key", s.apiKey)

	// DeepAI returns image synchronously
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("DeepAI API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepAI API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse DeepAI response (returns image URL directly)
	var deepAIResponse struct {
		ID       string `json:"id"`
		OutputURL string `json:"output_url"`
	}

	err = json.Unmarshal(bodyBytes, &deepAIResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DeepAI response: %w", err)
	}

	if deepAIResponse.OutputURL == "" {
		return nil, errors.New("DeepAI API returned no image URL")
	}

	// DeepAI is synchronous, so we return the image URL directly
	// Create a response that indicates completion
	return &ImageGenerationResponse{
		TaskID:  deepAIResponse.ID,
		Status:  "completed",
		Message: deepAIResponse.OutputURL, // Store URL in message for sync providers
	}, nil
}

// generateImageFreepik generates an image using Freepik API (asynchronous)
func (s *ImageService) generateImageFreepik(req ImageGenerationRequest) (*ImageGenerationResponse, error) {

	// Use Mystic model by default (good for photorealistic, pencil-style images)
	if req.Model == "" {
		req.Model = "mystic"
	}

	// Use square aspect ratio by default (good for educational illustrations)
	// Freepik uses format like "square_1_1", "widescreen_16_9", etc.
	if req.AspectRatio == "" {
		req.AspectRatio = "square_1_1"
	}
	// Convert common formats to Freepik format
	if req.AspectRatio == "square" {
		req.AspectRatio = "square_1_1"
	} else if req.AspectRatio == "portrait" {
		req.AspectRatio = "portrait_4_3"
	} else if req.AspectRatio == "landscape" || req.AspectRatio == "widescreen" {
		req.AspectRatio = "widescreen_16_9"
	}

	// Use 1K resolution by default (lighter images as requested)
	if req.Resolution == "" {
		req.Resolution = "1k"
	}

	// Build request payload
	payload := map[string]interface{}{
		"prompt": req.Prompt,
		"aspect_ratio": req.AspectRatio,
	}

	// Add model-specific parameters
	if req.Model == "mystic" {
		payload["resolution"] = req.Resolution
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/ai/%s", s.baseURL, req.Model)
	log.Printf("Freepik API: Generating image with prompt: %s (model: %s, aspect: %s, resolution: %s)", 
		req.Prompt[:min(50, len(req.Prompt))], req.Model, req.AspectRatio, req.Resolution)
	
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("x-freepik-api-key", s.apiKey)

	// Send request with context for better timeout handling
	log.Printf("Freepik API: Sending request to %s", url)
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Freepik API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("Freepik API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse struct {
		Data struct {
			TaskID string `json:"task_id"`
			Status string `json:"status"`
		} `json:"data"`
		Message string `json:"message,omitempty"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		// Try alternative response format
		var altResponse struct {
			TaskID string `json:"task_id"`
			Status string `json:"status"`
			Message string `json:"message,omitempty"`
		}
		if err2 := json.Unmarshal(body, &altResponse); err2 == nil {
			return &ImageGenerationResponse{
				TaskID:  altResponse.TaskID,
				Status:  altResponse.Status,
				Message: altResponse.Message,
			}, nil
		}
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &ImageGenerationResponse{
		TaskID:  apiResponse.Data.TaskID,
		Status:  apiResponse.Data.Status,
		Message: apiResponse.Message,
	}, nil
}

// CheckImageTaskStatus checks the status of an image generation task
func (s *ImageService) CheckImageTaskStatus(taskID, model string) (*ImageTaskStatus, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("%s API key not configured", s.provider)
	}

	// Handle different providers
	switch s.provider {
	case ProviderReplicate:
		return s.checkReplicateStatus(taskID)
	case ProviderFreepik:
		if model == "" {
			model = "mystic" // Default model
		}
		return s.checkFreepikStatus(taskID, model)
	default:
		return nil, fmt.Errorf("status checking not supported for provider: %s", s.provider)
	}
}

// checkFreepikStatus checks Freepik task status
func (s *ImageService) checkFreepikStatus(taskID, model string) (*ImageTaskStatus, error) {
	// Create request to check task status
	url := fmt.Sprintf("%s/ai/%s/tasks/%s", s.baseURL, model, taskID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("x-freepik-api-key", s.apiKey)

	// Send request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Freepik API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrImageTaskNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Freepik API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse struct {
		Data struct {
			TaskID   string `json:"task_id"`
			Status   string `json:"status"`
			ImageURL string `json:"image_url,omitempty"`
			Error    string `json:"error,omitempty"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		// Try alternative response format
		var altResponse struct {
			TaskID   string `json:"task_id"`
			Status   string `json:"status"`
			ImageURL string `json:"image_url,omitempty"`
			Error    string `json:"error,omitempty"`
		}
		if err2 := json.Unmarshal(body, &altResponse); err2 == nil {
			return &ImageTaskStatus{
				TaskID:   altResponse.TaskID,
				Status:   altResponse.Status,
				ImageURL: altResponse.ImageURL,
				Error:    altResponse.Error,
			}, nil
		}
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &ImageTaskStatus{
		TaskID:   apiResponse.Data.TaskID,
		Status:   apiResponse.Data.Status,
		ImageURL: apiResponse.Data.ImageURL,
		Error:    apiResponse.Data.Error,
	}, nil
}

// generateImageReplicate generates an image using Replicate API (asynchronous)
func (s *ImageService) generateImageReplicate(req ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// Replicate model - use stable-diffusion-3 for good quality, or you can use other models
	model := req.Model
	if model == "" {
		model = "stability-ai/stable-diffusion-3" // Default model
	}

	// Enhance prompt for pencil sketch style
	enhancedPrompt := req.Prompt
	if !containsPencilKeywords(enhancedPrompt) {
		enhancedPrompt = "minimalist pencil sketch, " + enhancedPrompt + ", graphite line art style, educational illustration"
	}

	// Build request payload
	payload := map[string]interface{}{
		"input": map[string]interface{}{
			"prompt": enhancedPrompt,
		},
	}

	// Add image dimensions (default to square 512x512 for lightweight)
	width := 512
	height := 512
	if req.AspectRatio == "portrait" || req.AspectRatio == "portrait_4_3" {
		width = 512
		height = 768
	} else if req.AspectRatio == "landscape" || req.AspectRatio == "widescreen" || req.AspectRatio == "widescreen_16_9" {
		width = 768
		height = 512
	}

	inputMap := payload["input"].(map[string]interface{})
	inputMap["width"] = width
	inputMap["height"] = height

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request - Replicate uses model path in URL
	apiURL := fmt.Sprintf("%s/models/%s/predictions", s.baseURL, model)
	log.Printf("Replicate API: Generating image with prompt: %s (model: %s)", enhancedPrompt[:min(100, len(enhancedPrompt))], model)

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Send request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Replicate API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Replicate API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse Replicate response
	var replicateResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	err = json.Unmarshal(body, &replicateResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Replicate response: %w", err)
	}

	return &ImageGenerationResponse{
		TaskID: replicateResponse.ID,
		Status: replicateResponse.Status,
	}, nil
}

// checkReplicateStatus checks the status of a Replicate prediction
func (s *ImageService) checkReplicateStatus(predictionID string) (*ImageTaskStatus, error) {
	apiURL := fmt.Sprintf("%s/predictions/%s", s.baseURL, predictionID)

	httpReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Send request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Replicate API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrImageTaskNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Replicate API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse Replicate response
	var replicateResponse struct {
		ID     string        `json:"id"`
		Status string        `json:"status"`
		Output interface{}   `json:"output"` // Can be string (URL) or array of strings
		Error  string        `json:"error,omitempty"`
	}

	err = json.Unmarshal(body, &replicateResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Replicate response: %w", err)
	}

	// Extract image URL from output (can be string or array)
	var imageURL string
	switch v := replicateResponse.Output.(type) {
	case string:
		imageURL = v
	case []interface{}:
		if len(v) > 0 {
			if urlStr, ok := v[0].(string); ok {
				imageURL = urlStr
			}
		}
	}

	// Map Replicate status to our status format
	status := replicateResponse.Status
	if status == "succeeded" {
		status = "completed"
	} else if status == "failed" || status == "canceled" {
		status = "failed"
	}

	return &ImageTaskStatus{
		TaskID:   replicateResponse.ID,
		Status:   status,
		ImageURL: imageURL,
		Error:    replicateResponse.Error,
	}, nil
}

// GenerateImageSync generates an image and waits for it to complete
// This is a convenience function that polls the task status
func (s *ImageService) GenerateImageSync(req ImageGenerationRequest) (string, error) {
	// Submit generation request
	genResp, err := s.GenerateImage(req)
	if err != nil {
		return "", err
	}

	taskID := genResp.TaskID
	model := req.Model
	if model == "" {
		// Default model depends on provider
		switch s.provider {
		case ProviderReplicate:
			model = "stability-ai/stable-diffusion-3"
		case ProviderFreepik:
			model = "mystic"
		default:
			model = "" // DeepAI doesn't need model for status check
		}
	}

	// Poll for completion (max 2 minutes, check every 2 seconds)
	maxAttempts := 60
	attempt := 0

	for attempt < maxAttempts {
		time.Sleep(2 * time.Second)
		attempt++

		status, err := s.CheckImageTaskStatus(taskID, model)
		if err != nil {
			if err == ErrImageTaskNotFound {
				return "", fmt.Errorf("image task not found: %s", taskID)
			}
			log.Printf("Error checking image task status: %v", err)
			continue
		}

		switch status.Status {
		case "completed":
			if status.ImageURL != "" {
				return status.ImageURL, nil
			}
			return "", errors.New("image generation completed but no image URL provided")
		case "failed":
			errorMsg := status.Error
			if errorMsg == "" {
				errorMsg = "image generation failed"
			}
			return "", fmt.Errorf("image generation failed: %s", errorMsg)
		case "pending", "processing":
			// Continue polling
			continue
		default:
			log.Printf("Unknown task status: %s", status.Status)
			continue
		}
	}

	return "", fmt.Errorf("image generation timed out after %d attempts", maxAttempts)
}

// IsConfigured returns whether the image service is properly configured
func (s *ImageService) IsConfigured() bool {
	return s.apiKey != ""
}
