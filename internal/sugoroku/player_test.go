package sugoroku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayer_Money(t *testing.T) {
	player := NewPlayer("test", nil)

	// Test Profit
	err := player.Profit(100)
	assert.NoError(t, err)
	assert.Equal(t, 100, player.money)

	// Test Loss
	err = player.Loss(30)
	assert.NoError(t, err)
	assert.Equal(t, 70, player.money)

	// Test invalid amounts
	err = player.Profit(-10)
	assert.Error(t, err)

	err = player.Loss(-10)
	assert.Error(t, err)
}

func TestPlayer_Getters(t *testing.T) {
	tile := &Tile{id: 1}
	player := NewPlayer("test_id", tile)

	assert.Equal(t, "test_id", player.GetID())
	assert.Equal(t, tile, player.GetPosition())
}
