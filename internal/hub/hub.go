// internal/hub/hub.go

package hub

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	mu sync.RWMutex
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) NewClient(conn *websocket.Conn, playerID string) *Client {
	return &Client{
		Hub:      h,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Receive:  make(chan []byte, 256),
		PlayerID: playerID,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.PlayerID] = client
			h.mu.Unlock()
			fmt.Printf("Client registered: %s\n", client.PlayerID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.PlayerID]; ok {
				delete(h.clients, client.PlayerID)
				close(client.Send)
				fmt.Printf("Client unregistered: %s\n", client.PlayerID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.PlayerID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Hubに新たなプレイヤーを登録する
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Hubからプレイヤーを削除する
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// 特定のプレイヤーにJSONメッセージを送信する
func (h *Hub) SendToPlayer(playerID string, message interface{}) error {
	rawMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	h.mu.RLock()
	client, ok := h.clients[playerID]
	h.mu.RUnlock()

	if !ok {
		return fmt.Errorf("client with playerID %s not found", playerID)
	}

	select {
	case client.Send <- rawMessage:
	default:
		return fmt.Errorf("client %s send channel is full, message dropped", playerID)
	}

	return nil
}
