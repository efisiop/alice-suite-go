package handlers

import (
	"net/http"

	"github.com/efisiopittau/alice-suite-go/internal/realtime"
	"github.com/efisiopittau/alice-suite-go/pkg/auth"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// In production, validate origin
		return true
	},
}

// HandleWebSocket handles WebSocket connections (optional, for bidirectional communication)
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract token
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Validate token
	user, err := auth.GetUserFromToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Register client for broadcasting
	broadcaster := realtime.GetBroadcaster()
	role := "reader"
	if user.Role == "consultant" {
		role = "consultant"
	}
	client := broadcaster.RegisterClient(user.ID, role)
	defer broadcaster.UnregisterClient(user.ID)

	// Send initial message
	conn.WriteJSON(map[string]interface{}{
		"type":    "connected",
		"user_id": user.ID,
		"role":    role,
	})

	// Handle incoming messages
	go func() {
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				break
			}
			// Handle message (echo for now)
			conn.WriteJSON(map[string]interface{}{
				"type":    "echo",
				"message": msg,
			})
		}
	}()

	// Send events from broadcaster
	for event := range client.Events {
		if err := conn.WriteJSON(event); err != nil {
			break
		}
	}
}

