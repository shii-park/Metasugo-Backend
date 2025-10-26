package sugoroku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateTestFileはtile_test.goに存在するため、ここでは削除しました。

func TestGame_AddPlayer(t *testing.T) {
	game := NewGameWithTilesForTest("../../tiles.json")
	player, err := game.AddPlayer("test_player")

	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, "test_player", player.GetID())

	// Test adding the same player again
	_, err = game.AddPlayer("test_player")
	assert.Error(t, err, "should return an error when adding a player with an existing ID")
}
