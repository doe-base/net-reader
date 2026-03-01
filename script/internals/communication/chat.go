package communication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"conceptual-lan/internals/discovery"
)

type ChatMessage struct {
	From      string `json:"from"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp"`
}

var (
	messageMu sync.RWMutex
	Messages  []ChatMessage
)

// Receive message from peer
func ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	var msg ChatMessage

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid message", http.StatusBadRequest)
		return
	}

	messageMu.Lock()
	Messages = append(Messages, msg)
	messageMu.Unlock()

	fmt.Printf("📩 Message from %s: %s\n", msg.From, msg.Body)

	w.WriteHeader(http.StatusOK)
}

// Send message to a specific peer
func SendMessageToPeer(peerIP string, msg ChatMessage) error {
	data, _ := json.Marshal(msg)

	url := fmt.Sprintf("http://%s:8080/api/message", peerIP)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Broadcast message to all peers
func BroadcastMessage(msg ChatMessage) {
	discovery.PeerMu.RLock()
	defer discovery.PeerMu.RUnlock()

	for ip := range discovery.Peers {
		go SendMessageToPeer(ip, msg)
	}
}

// Endpoint called by browser
func SendMessageFromBrowser(w http.ResponseWriter, r *http.Request) {
	var msg ChatMessage

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid message", http.StatusBadRequest)
		return
	}

	msg.Timestamp = time.Now().Unix()

	messageMu.Lock()
	Messages = append(Messages, msg)
	messageMu.Unlock()

	// Broadcast to peers
	BroadcastMessage(msg)

	w.WriteHeader(http.StatusOK)
}

// Endpoint to fetch messages for UI
func GetMessages(w http.ResponseWriter, r *http.Request) {
	messageMu.RLock()
	defer messageMu.RUnlock()

	if Messages == nil {
		Messages = []ChatMessage{} // always send empty array instead of nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Messages)
}
