package realtime

import (
	"encoding/json"
	"sync"
	"time"
)

// Event represents a real-time event
type Event struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Client represents an SSE client connection
type Client struct {
	ID     string
	Role   string // "reader" or "consultant"
	Events chan Event
}

// Broadcaster manages SSE connections and broadcasts events
type Broadcaster struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var globalBroadcaster *Broadcaster
var once sync.Once

// GetBroadcaster returns the global broadcaster instance
func GetBroadcaster() *Broadcaster {
	once.Do(func() {
		globalBroadcaster = &Broadcaster{
			clients: make(map[string]*Client),
		}
	})
	return globalBroadcaster
}

// RegisterClient registers a new SSE client
func (b *Broadcaster) RegisterClient(id, role string) *Client {
	b.mu.Lock()
	defer b.mu.Unlock()

	client := &Client{
		ID:     id,
		Role:   role,
		Events: make(chan Event, 10), // Buffered channel
	}

	b.clients[id] = client
	return client
}

// UnregisterClient removes a client
func (b *Broadcaster) UnregisterClient(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if client, exists := b.clients[id]; exists {
		close(client.Events)
		delete(b.clients, id)
	}
}

// Broadcast sends an event to all clients matching the filter
func (b *Broadcaster) Broadcast(event Event, filter func(*Client) bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, client := range b.clients {
		if filter == nil || filter(client) {
			select {
			case client.Events <- event:
			default:
				// Channel full, skip (non-blocking)
			}
		}
	}
}

// BroadcastToRole sends an event to all clients with a specific role
func (b *Broadcaster) BroadcastToRole(event Event, role string) {
	b.Broadcast(event, func(client *Client) bool {
		return client.Role == role
	})
}

// BroadcastToClient sends an event to a specific client
func (b *Broadcaster) BroadcastToClient(event Event, clientID string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if client, exists := b.clients[clientID]; exists {
		select {
		case client.Events <- event:
		default:
			// Channel full, skip
		}
	}
}

// GetClientCount returns the number of connected clients
func (b *Broadcaster) GetClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}

// GetClientsByRole returns clients with a specific role
func (b *Broadcaster) GetClientsByRole(role string) []*Client {
	b.mu.RLock()
	defer b.mu.RUnlock()

	clients := []*Client{}
	for _, client := range b.clients {
		if client.Role == role {
			clients = append(clients, client)
		}
	}
	return clients
}

// Event types
const (
	EventTypeLogin           = "login"
	EventTypeLogout          = "logout"
	EventTypeHelpRequest     = "help_request"
	EventTypeHelpRequestUpdate = "help_request_update"
	EventTypeActivity        = "activity"
	EventTypeReadingProgress = "reading_progress"
	EventTypeOnlineUsers     = "online_users"
)

// CreateEvent creates a new event
func CreateEvent(eventType string, data interface{}) Event {
	return Event{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// ToJSON converts event to JSON
func (e Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

