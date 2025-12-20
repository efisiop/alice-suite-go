package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/efisiopittau/alice-suite-go/internal/realtime"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
)

// HandleSSE handles Server-Sent Events connections
func HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Extract token from query parameter or header
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	if token == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Validate token and get user
	user, err := auth.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Determine role
	role := "reader"
	if user.Role == "consultant" {
		role = "consultant"
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Register client
	broadcaster := realtime.GetBroadcaster()
	client := broadcaster.RegisterClient(user.ID, role)
	defer broadcaster.UnregisterClient(user.ID)

	// Send initial connection event
	initialEvent := realtime.CreateEvent("connected", map[string]interface{}{
		"message": "Connected to real-time updates",
		"user_id": user.ID,
		"role":    role,
	})
	sendSSEEvent(w, initialEvent)

	// Keep connection alive with heartbeat (reduced to 15 seconds for faster detection of disconnects)
	heartbeatTicker := time.NewTicker(15 * time.Second)
	defer heartbeatTicker.Stop()

	// Flush initial headers
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// Listen for events
	for {
		select {
		case event := <-client.Events:
			if err := sendSSEEvent(w, event); err != nil {
				return // Client disconnected
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

		case <-heartbeatTicker.C:
			// Send heartbeat to keep connection alive
			heartbeat := realtime.CreateEvent("heartbeat", map[string]string{
				"timestamp": time.Now().Format(time.RFC3339),
			})
			if err := sendSSEEvent(w, heartbeat); err != nil {
				return
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

		case <-r.Context().Done():
			// Client disconnected
			return
		}
	}
}

// sendSSEEvent sends an SSE-formatted event
func sendSSEEvent(w http.ResponseWriter, event realtime.Event) error {
	// Format: "event: {type}\ndata: {json}\n\n"
	eventJSON, err := event.ToJSON()
	if err != nil {
		return err
	}

	// Convert to string and ensure it's valid JSON
	jsonStr := string(eventJSON)
	
	// Write event type line
	if _, err := fmt.Fprintf(w, "event: %s\n", event.Type); err != nil {
		return err
	}
	
	// Write data line - according to SSE spec, each line of data must be prefixed with "data: "
	// Since json.Marshal produces a single line (newlines are escaped as \n), we can write it as one line
	if _, err := fmt.Fprintf(w, "data: %s\n", jsonStr); err != nil {
		return err
	}
	
	// End event with blank line (required by SSE spec)
	_, err = fmt.Fprintf(w, "\n")
	return err
}

// BroadcastHelpRequest broadcasts a help request event to consultants
func BroadcastHelpRequest(helpRequest map[string]interface{}) {
	broadcaster := realtime.GetBroadcaster()
	event := realtime.CreateEvent(realtime.EventTypeHelpRequest, helpRequest)
	broadcaster.BroadcastToRole(event, "consultant")
}

// BroadcastHelpRequestUpdate broadcasts a help request update
func BroadcastHelpRequestUpdate(helpRequest map[string]interface{}) {
	broadcaster := realtime.GetBroadcaster()
	event := realtime.CreateEvent(realtime.EventTypeHelpRequestUpdate, helpRequest)
	broadcaster.BroadcastToRole(event, "consultant")
}

// BroadcastActivity broadcasts an activity event
func BroadcastActivity(activity map[string]interface{}) {
	broadcaster := realtime.GetBroadcaster()
	event := realtime.CreateEvent(realtime.EventTypeActivity, activity)
	broadcaster.BroadcastToRole(event, "consultant")
}

// BroadcastLogin broadcasts a login event
func BroadcastLogin(userID, email, firstName, lastName string) {
	broadcaster := realtime.GetBroadcaster()
	event := realtime.CreateEvent(realtime.EventTypeLogin, map[string]interface{}{
		"user_id":    userID,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
	})
	broadcaster.BroadcastToRole(event, "consultant")
}

// BroadcastLogout broadcasts a logout event
func BroadcastLogout(userID string) {
	broadcaster := realtime.GetBroadcaster()
	event := realtime.CreateEvent(realtime.EventTypeLogout, map[string]interface{}{
		"user_id": userID,
	})
	broadcaster.BroadcastToRole(event, "consultant")
}

