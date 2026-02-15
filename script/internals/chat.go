package internals

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ChatMessage defines the data structure for communication
type ChatMessage struct {
	From      string `json:"from"`
	To        string `json:"to"` // "all" or ClientID
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// Client represents a single connected user
type Client struct {
	ID   string // Use "Daniel-PC" or a Username here
	Conn *websocket.Conn
	Send chan ChatMessage
}

// Hub manages all active clients and message routing
type Hub struct {
	clients    map[string]*Client
	broadcast  chan ChatMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan ChatMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Inside your Hub.Run() case client := <-h.register:
			h.mu.Lock()
			// Instead of RemoteAddr, we use the name/IP from the discovery service
			h.clients[client.ID] = client
			h.mu.Unlock()
			fmt.Printf("Client %s joined\n", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			fmt.Printf("Client %s left\n", client.ID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			if msg.To == "all" {
				// Broadcast to everyone
				for _, client := range h.clients {
					select {
					case client.Send <- msg:
					default: // If buffer is full, drop message to prevent blocking
						close(client.Send)
						delete(h.clients, client.ID)
					}
				}
			} else {
				// Direct Message
				if target, ok := h.clients[msg.To]; ok {
					target.Send <- msg
				}
			}
			h.mu.RUnlock()
		}
	}
}

// writePump handles sending messages TO the browser
func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) ChatWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Using RemoteAddr as ID (In production, use a UserID or UUID)
	client := &Client{
		ID:   r.RemoteAddr,
		Conn: conn,
		Send: make(chan ChatMessage, 256),
	}

	h.register <- client

	// Start the writer goroutine
	go client.writePump()

	// Main Read loop: handles incoming messages FROM browser
	defer func() {
		h.unregister <- client
		conn.Close()
	}()

	for {
		var msg ChatMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		if msg.Type == "identify" {
			client.ID = msg.From // Now 'Daniel-PC' is the key, not the random port!
			h.register <- client
			continue
		}
	}
}
