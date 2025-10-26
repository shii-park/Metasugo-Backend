package hub

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub_Registration(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := NewClient(hub, nil, "player1")
	client2 := NewClient(hub, nil, "player2")

	// Register clients
	hub.Register(client1)
	hub.Register(client2)

	time.Sleep(50 * time.Millisecond) // Allow time for processing

	hub.mu.RLock()
	assert.Len(t, hub.clients, 2, "Should have 2 registered clients")
	assert.Contains(t, hub.clients, "player1")
	assert.Contains(t, hub.clients, "player2")
	hub.mu.RUnlock()

	// Unregister a client
	hub.Unregister(client1)

	time.Sleep(50 * time.Millisecond) // Allow time for processing

	hub.mu.RLock()
	assert.Len(t, hub.clients, 1, "Should have 1 client remaining")
	assert.NotContains(t, hub.clients, "player1")
	hub.mu.RUnlock()
}

func TestHub_SendToPlayer(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := NewClient(hub, nil, "player1")
	client2 := NewClient(hub, nil, "player2")

	hub.Register(client1)
	hub.Register(client2)
	time.Sleep(50 * time.Millisecond)

	t.Run("Send to existing player", func(t *testing.T) {
		message := map[string]string{"data": "hello player 2"}
		err := hub.SendToPlayer("player2", message)
		assert.NoError(t, err)

		// Check if the message was received by the correct client
		select {
		case msgBytes := <-client2.Send:
			var receivedMsg map[string]string
			json.Unmarshal(msgBytes, &receivedMsg)
			assert.Equal(t, "hello player 2", receivedMsg["data"])
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for message")
		}

		// Check that the other client did not receive the message
		select {
		case <-client1.Send:
			t.Fatal("Client 1 should not have received the message")
		default:
			// Correct, no message
		}
	})

	t.Run("Send to non-existent player", func(t *testing.T) {
		message := map[string]string{"data": "hello ghost"}
		err := hub.SendToPlayer("player3", message)
		assert.Error(t, err, "Should return an error when player is not found")
	})
}
