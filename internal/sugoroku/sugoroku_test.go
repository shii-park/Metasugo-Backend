package sugoroku

import (
	"os"
	"testing"
)

func TestGamePlay(t *testing.T) {
	const testJSON = `[
		{"id": 1, "kind": "profit", "detail": "Start", "prev_ids": [], "next_ids": [2]},
		{"id": 2, "kind": "loss", "detail": "Middle", "prev_ids": [1], "next_ids": [3]},
		{"id": 3, "kind": "quiz", "detail": "End", "prev_ids": [2], "next_ids": []}
	]`
	tmpFile := CreateTestFile(t, "test_game_play_*.json", testJSON)
	defer os.Remove(tmpFile)

	// Setup game
	game := NewGameWithTiles(tmpFile)
	player, err := game.AddPlayer("test_player")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// 1. Initial Position
	if player.position.id != InitialTileID {
		t.Errorf("Expected player to start at tile %d, but got %d", InitialTileID, player.position.id)
	}

	// 2. Move Player
	game.players["test_player"].MoveByDiceRoll(1)

	// 3. Check New Position
	expectedTileID := 2 // Assuming the next tile from InitialTileID(1) is 2
	if player.position.id != expectedTileID {
		t.Errorf("Expected player to be at tile %d after moving, but got %d", expectedTileID, player.position.id)
	}
}

func TestTileEffects_ProfitAndLoss(t *testing.T) {
	const testJSON = `[
		{"id": 1, "kind": "profit", "detail": "Start", "prev_ids": [], "next_ids": [2]},
		{"id": 2, "kind": "profit", "effect": {"type": "profit", "amount": 100}, "detail": "Profit Tile", "prev_ids": [1], "next_ids": [3]},
		{"id": 3, "kind": "loss", "effect": {"type": "loss", "amount": 50}, "detail": "Loss Tile", "prev_ids": [2], "next_ids": []}
	]`
	tmpFile := CreateTestFile(t, "test_tile_effects_*.json", testJSON)
	defer os.Remove(tmpFile)

	game := NewGameWithTiles(tmpFile)
	player, err := game.AddPlayer("test_player")
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Move to profit tile
	player.MoveByDiceRoll(1)
	err = player.position.effect.Apply(player, game)
	if err != nil {
		t.Fatalf("Error applying profit effect: %v", err)
	}
	if player.money != 100 {
		t.Errorf("Expected money to be 100 after profit tile, but got %d", player.money)
	}

	// Move to loss tile
	player.MoveByDiceRoll(1)
	err = player.position.effect.Apply(player, game)
	if err != nil {
		t.Fatalf("Error applying loss effect: %v", err)
	}
	if player.money != 50 {
		t.Errorf("Expected money to be 50 after loss tile, but got %d", player.money)
	}
}

func TestMultiplayer(t *testing.T) {
	const testJSON = `[
		{"id": 1, "kind": "profit", "detail": "Start", "prev_ids": [], "next_ids": [2]},
		{"id": 2, "kind": "loss", "detail": "Middle", "prev_ids": [1], "next_ids": [3]},
		{"id": 3, "kind": "quiz", "detail": "End", "prev_ids": [2], "next_ids": []}
	]`
	tmpFile := CreateTestFile(t, "test_multiplayer_*.json", testJSON)
	defer os.Remove(tmpFile)

	game := NewGameWithTiles(tmpFile)
	player1, err := game.AddPlayer("player1")
	if err != nil {
		t.Fatalf("Failed to add player1: %v", err)
	}

	player2, err := game.AddPlayer("player2")
	if err != nil {
		t.Fatalf("Failed to add player2: %v", err)
	}

	// Move player1
	player1.MoveByDiceRoll(1)

	// Check positions
	if player1.position.id != 2 {
		t.Errorf("Expected player1 to be at tile 2, but got %d", player1.position.id)
	}

	if player2.position.id != 1 {
		t.Errorf("Expected player2 to be at tile 1, but got %d", player2.position.id)
	}
}

func TestMultiplayerEffects_OverallAndNeighbor(t *testing.T) {
	const testJSON = `[
		{"id": 1, "kind": "normal", "detail": "Start", "prev_ids": [], "next_ids": [2]},
		{"id": 2, "kind": "overall", "effect": {"type": "overall", "profit_amount": 100}, "detail": "Overall Profit", "prev_ids": [1], "next_ids": [3]},
		{"id": 3, "kind": "neighbor", "effect": {"type": "neighbor", "loss_amount": 50}, "detail": "Neighbor Loss", "prev_ids": [2], "next_ids": []}
	]`
	tmpFile := CreateTestFile(t, "test_multiplayer_effects_*.json", testJSON)
	defer os.Remove(tmpFile)

	game := NewGameWithTiles(tmpFile)
	player1, _ := game.AddPlayer("player1")
	player2, _ := game.AddPlayer("player2")
	player3, _ := game.AddPlayer("player3")

	// Move player1 to overall profit tile
	player1.MoveByDiceRoll(1)
	err := player1.position.effect.Apply(player1, game)
	if err != nil {
		t.Fatalf("Error applying overall effect: %v", err)
	}

	// player1 gets 100, player2 and player3 lose 50 each
	if player1.money != 100 {
		t.Errorf("Player1 should have 100, but has %d", player1.money)
	}
	if player2.money != -50 {
		t.Errorf("Player2 should have -50, but has %d", player2.money)
	}
	if player3.money != -50 {
		t.Errorf("Player3 should have -50, but has %d", player3.money)
	}

	// Move player2 to neighbor loss tile (player1 is a neighbor)
	player2.MoveByDiceRoll(2)
	err = player2.position.effect.Apply(player2, game)
	if err != nil {
		t.Fatalf("Error applying neighbor effect: %v", err)
	}

	// player2 loses 50, player1 gains 50
	if player2.money != -100 {
		t.Errorf("Player2 should have -100, but has %d", player2.money)
	}
	if player1.money != 150 {
		t.Errorf("Player1 should have 150, but has %d", player1.money)
	}
	if player3.money != -50 {
		t.Errorf("Player3 should have -50, but has %d", player3.money)
	}
}

func TestNeighborEffect_WithPlayerOnSameTile(t *testing.T) {
	const testJSON = `[
		{"id": 1, "kind": "normal", "detail": "Start", "prev_ids": [], "next_ids": [2]},
		{"id": 2, "kind": "normal", "detail": "", "prev_ids": [1], "next_ids": [3]},
		{"id": 3, "kind": "neighbor", "effect": {"type": "neighbor", "profit_amount": 90}, "detail": "Neighbor Profit", "prev_ids": [2], "next_ids": [4]},
		{"id": 4, "kind": "normal", "detail": "", "prev_ids": [3], "next_ids": [5]},
		{"id": 5, "kind": "normal", "detail": "End", "prev_ids": [4], "next_ids": []}
	]`
	tmpFile := CreateTestFile(t, "test_neighbor_same_tile_*.json", testJSON)
	defer os.Remove(tmpFile)

	game := NewGameWithTiles(tmpFile)
	player1, _ := game.AddPlayer("player1") // Neighbor 1
	player2, _ := game.AddPlayer("player2") // Neighbor 2
	player3, _ := game.AddPlayer("player3") // Neighbor 3
	player4, _ := game.AddPlayer("player4") // The one to move
	player5, _ := game.AddPlayer("player5") // On the same tile

	// Position players
	player1.position = game.tileMap[2]
	player2.position = game.tileMap[4]
	player3.position = game.tileMap[2]
	player4.position = game.tileMap[1] // Start position
	player5.position = game.tileMap[3] // Same tile as player4 will land on

	// Move player4 to the neighbor effect tile
	player4.MoveByDiceRoll(2)
	err := player4.position.effect.Apply(player4, game)
	if err != nil {
		t.Fatalf("Error applying neighbor effect: %v", err)
	}

	// player4 gets 90, and the 4 neighbors (player1, player2, player3, player5) lose 22 each (90/4=22.5, truncated)
	if player4.money != 90 {
		t.Errorf("Player4 should have 90, but has %d", player4.money)
	}
	if player1.money != -22 {
		t.Errorf("Player1 should have -22, but has %d", player1.money)
	}
	if player2.money != -22 {
		t.Errorf("Player2 should have -22, but has %d", player2.money)
	}
	if player3.money != -22 {
		t.Errorf("Player3 should have -22, but has %d", player3.money)
	}
	if player5.money != -22 {
		t.Errorf("Player5 should have -22, but has %d", player5.money)
	}
}
