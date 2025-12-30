package services

import (
	"errors"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
)

var (
	ErrHelpRequestNotFound = errors.New("help request not found")
	ErrUnauthorized       = errors.New("unauthorized")
)

// HelpService handles help request operations
type HelpService struct{}

// NewHelpService creates a new help service
func NewHelpService() *HelpService {
	return &HelpService{}
}

// CreateHelpRequest creates a new help request
func (s *HelpService) CreateHelpRequest(userID, bookID, content, context string, sectionID *string) (*models.HelpRequest, error) {
	request := &models.HelpRequest{
		UserID:    userID,
		BookID:    bookID,
		SectionID: sectionID,
		Status:    "pending",
		Content:   content,
		Context:   context,
	}

	err := database.CreateHelpRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// GetHelpRequests retrieves help requests for a user
func (s *HelpService) GetHelpRequests(userID string) ([]*models.HelpRequest, error) {
	return database.GetHelpRequests(userID)
}

// GetHelpRequestsByConsultant retrieves help requests assigned to a consultant
func (s *HelpService) GetHelpRequestsByConsultant(consultantID string) ([]*models.HelpRequest, error) {
	return database.GetHelpRequestsByConsultant(consultantID)
}

// GetHelpRequestByID retrieves a help request by ID
func (s *HelpService) GetHelpRequestByID(requestID string) (*models.HelpRequest, error) {
	return database.GetHelpRequestByID(requestID)
}

// AssignHelpRequest assigns a help request to a consultant
func (s *HelpService) AssignHelpRequest(requestID, consultantID string) (*models.HelpRequest, error) {
	// Get request
	request, err := database.GetHelpRequestByID(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrHelpRequestNotFound
	}

	// Update request
	request.Status = "assigned"
	request.AssignedTo = &consultantID

	err = database.UpdateHelpRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// ResolveHelpRequest resolves a help request with a response
func (s *HelpService) ResolveHelpRequest(requestID, consultantID, response string) (*models.HelpRequest, error) {
	// Get request
	request, err := database.GetHelpRequestByID(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrHelpRequestNotFound
	}

	// Verify consultant is assigned
	if request.AssignedTo == nil || *request.AssignedTo != consultantID {
		return nil, ErrUnauthorized
	}

	// Update request
	request.Status = "resolved"
	request.Response = response
	now := time.Now()
	request.ResolvedAt = &now

	err = database.UpdateHelpRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

