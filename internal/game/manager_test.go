package game

import (
	"testing"

	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
	"github.com/stretchr/testify/assert"
)

func TestNewGameManager(t *testing.T) {
	g := sugoroku.NewGameWithTilesForTest("../../tiles.json")
	h := hub.NewHub()
	gm := NewGameManager(g, h)

	assert.NotNil(t, gm)
	assert.Equal(t, g, gm.game)
	assert.Equal(t, h, gm.hub)
	assert.NotNil(t, gm.playerClients)
}

func TestRegisterPlayerClient(t *testing.T) {
	g := sugoroku.NewGameWithTilesForTest("../../tiles.json")
	h := hub.NewHub()
	gm := NewGameManager(g, h)

	userID := "testUser"
	// hub.Client is an empty struct, so we can create it directly.
	client := &hub.Client{}

	gm.RegisterPlayerClient(userID, client)

	// Check if the player is added to the game
	players := g.GetAllPlayers()
	foundPlayer := false
	for _, p := range players {
		// Assuming Player has a method ID() that returns the player's ID.
		// This might need adjustment based on the actual implementation of the Player struct.
		if p.GetID() == userID {
			foundPlayer = true
			break
		}
	}
	assert.True(t, foundPlayer, "Player should be added to the game")

	// Check if the client is registered
	registeredClient, ok := gm.playerClients[userID]
	assert.True(t, ok)
	assert.Equal(t, client, registeredClient)
}
