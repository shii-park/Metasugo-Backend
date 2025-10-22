package hub

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHub(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		// In a real application, you would create a client here.
		// For this test, we are creating the client outside the server
		// to have more control over it.
		_ = NewClient(hub, conn, "test-user")
	}))
	defer server.Close()

	// Convert the http:// server URL to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to websocket: %v", err)
	}
	defer conn.Close()

	client := NewClient(hub, conn, "test-user")

	// Register the client
	hub.register <- client

	// Allow some time for registration to complete
	time.Sleep(100 * time.Millisecond)

	if len(hub.clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(hub.clients))
	}

	// Broadcast a message
	message := []byte("hello")
	hub.broadcast <- message

	// Check if the client received the message
	select {
	case msg := <-client.send:
		if string(msg) != "hello" {
			t.Errorf("Expected 'hello', got '%s'", string(msg))
		}
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for message")
	}

	// Unregister the client
	hub.unregister <- client

	// Allow some time for unregistration to complete
	time.Sleep(100 * time.Millisecond)

	if len(hub.clients) != 0 {
		t.Errorf("Expected 0 clients, got %d", len(hub.clients))
	}
}