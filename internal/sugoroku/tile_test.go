package sugoroku

import (
	"os"
	"strings"
	"testing"
)

// TestNewTile tests the NewTile function (happy path).
func TestNewTile(t *testing.T) {
	effect := ProfitEffect{Amount: 100}
	tile := NewTile(nil, nil, profit, 1, effect, "Test Tile")

	if tile == nil {
		t.Fatal("NewTile returned nil")
	}
	if tile.id != 1 {
		t.Errorf("Expected id to be 1, got %d", tile.id)
	}
	if tile.kind != profit {
		t.Errorf("Expected kind to be profit, got %s", tile.kind)
	}
}

// TestInitTilesFromPath_HappyPath tests the successful execution, including tile linking.
func TestInitTilesFromPath_HappyPath(t *testing.T) {
	const testJSON = `[
		{"id": 1, "kind": "profit", "detail": "Start", "prev_id": 0, "next_id": 2},
		{"id": 2, "kind": "loss", "detail": "Middle", "prev_id": 1, "next_id": 3},
		{"id": 3, "kind": "quiz", "detail": "End", "prev_id": 2, "next_id": 0}
	]`
	tmpFile := createTestFile(t, "happy_path_*.json", testJSON)
	defer os.Remove(tmpFile)

	tiles, err := InitTilesFromPath(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if len(tiles) != 3 {
		t.Fatalf("Expected 3 tiles, got %d", len(tiles))
	}

	// Find tiles by ID for easier assertion, as order is not guaranteed.
	tileMap := make(map[int]*Tile)
	for _, tile := range tiles {
		tileMap[tile.id] = tile
	}

	tile1, ok1 := tileMap[1]
	tile2, ok2 := tileMap[2]
	tile3, ok3 := tileMap[3]
	if !ok1 || !ok2 || !ok3 {
		t.Fatal("Could not find all tiles (1, 2, 3) in the result")
	}

	// Check links
	if tile1.prev != nil {
		t.Errorf("Expected tile 1 prev to be nil, got tile id %d", tile1.prev.id)
	}
	if tile1.next != tile2 {
		t.Errorf("Expected tile 1 next to be tile 2, but it was not")
	}

	if tile2.prev != tile1 {
		t.Errorf("Expected tile 2 prev to be tile 1, but it was not")
	}
	if tile2.next != tile3 {
		t.Errorf("Expected tile 2 next to be tile 3, but it was not")
	}

	if tile3.prev != tile2 {
		t.Errorf("Expected tile 3 prev to be tile 2, but it was not")
	}
	if tile3.next != nil {
		t.Errorf("Expected tile 3 next to be nil, got tile id %d", tile3.next.id)
	}
}

// TestInitTilesFromPath_FileNotFound tests the error case where the JSON file does not exist.
func TestInitTilesFromPath_FileNotFound(t *testing.T) {
	_, err := InitTilesFromPath("non_existent_file.json")
	if err == nil {
		t.Fatal("Expected a file open error, but got nil")
	}
	if !strings.Contains(err.Error(), "file open error") {
		t.Errorf("Expected error to contain 'file open error', but got: %v", err)
	}
}

// TestInitTilesFromPath_InvalidJSON tests the error case where the JSON file is malformed.
func TestInitTilesFromPath_InvalidJSON(t *testing.T) {
	const invalidJSON = `[{"id": 1, "kind": "profit"` // Malformed JSON
	tmpFile := createTestFile(t, "invalid_json_*.json", invalidJSON)
	defer os.Remove(tmpFile)

	_, err := InitTilesFromPath(tmpFile)
	if err == nil {
		t.Fatal("Expected a JSON decode error, but got nil")
	}
	if !strings.Contains(err.Error(), "JSON decode error") {
		t.Errorf("Expected error to contain 'JSON decode error', but got: %v", err)
	}
}

// TestInitTilesFromPath_InvalidEffect tests the error case where an effect in the JSON is malformed.
func TestInitTilesFromPath_InvalidEffect(t *testing.T) {
	// Effect is missing the 'type' field
	const invalidEffectJSON = `[{"id": 1, "kind": "profit", "effect": {"amount": 10}}]`
	tmpFile := createTestFile(t, "invalid_effect_*.json", invalidEffectJSON)
	defer os.Remove(tmpFile)

	_, err := InitTilesFromPath(tmpFile)
	if err == nil {
		t.Fatal("Expected an effect unmarshal error, but got nil")
	}
	if !strings.Contains(err.Error(), "effect type is missing") {
		t.Errorf("Expected error to contain 'effect type is missing', but got: %v", err)
	}
}

// TestInitTilesFromPath_AllCases covers all tile kinds and their linking.
func TestInitTilesFromPath_AllCases(t *testing.T) {
	const allCasesJSON = `[
		{"id": 1, "kind": "profit", "effect": {"type": "profit", "amount": 10}, "prev_id": 0, "next_id": 2},
		{"id": 2, "kind": "loss", "effect": {"type": "loss", "amount": 5}, "prev_id": 1, "next_id": 3},
		{"id": 3, "kind": "quiz", "effect": {"type": "quiz", "quiz_id": 1}, "prev_id": 2, "next_id": 4},
		{"id": 4, "kind": "branch", "effect": {"type": "branch", "chose_id": 2}, "prev_id": 3, "next_id": 5},
		{"id": 5, "kind": "overall", "effect": {"type": "overall", "amount": 1}, "prev_id": 4, "next_id": 6},
		{"id": 6, "kind": "neighbor", "effect": {"type": "neighbor", "amount": 2}, "prev_id": 5, "next_id": 7},
		{"id": 7, "kind": "require", "effect": {"type": "require", "require_value": 3}, "prev_id": 6, "next_id": 8},
		{"id": 8, "kind": "gamble", "effect": {"type": "gamble"}, "prev_id": 7, "next_id": 9},
		{"id": 9, "kind": "unknown", "effect": {"type": "unknown"}, "prev_id": 8, "next_id": 10},
		{"id": 10, "kind": "none", "effect": null, "prev_id": 9, "next_id": 0}
	]`
	tmpFile := createTestFile(t, "all_cases_*.json", allCasesJSON)
	defer os.Remove(tmpFile)

	tiles, err := InitTilesFromPath(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(tiles) != 10 {
		t.Fatalf("Expected 10 tiles, got %d", len(tiles))
	}

	// Find tiles by ID for easier assertion, as order is not guaranteed.
	tileMap := make(map[int]*Tile)
	for _, tile := range tiles {
		tileMap[tile.id] = tile
	}

	// Check kinds and effects
	if tileMap[4].kind != branch {
		t.Errorf("Expected tile 4 to be of kind 'branch', got '%s'", tileMap[4].kind)
	}
	if tileMap[8].kind != gamble {
		t.Errorf("Expected tile 8 to be of kind 'gamble', got '%s'", tileMap[8].kind)
	}
	if tileMap[9].effect != nil {
		t.Errorf("Expected tile 9 (unknown type) to have a nil effect, got %T", tileMap[9].effect)
	}
	if tileMap[10].effect != nil {
		t.Errorf("Expected tile 10 (null effect) to have a nil effect, got %T", tileMap[10].effect)
	}

	// Check links
	if tileMap[1].prev != nil {
		t.Error("Expected tile 1 prev to be nil")
	}
	if tileMap[10].next != nil {
		t.Error("Expected tile 10 next to be nil")
	}

	for i := 1; i < 10; i++ {
		currentTile := tileMap[i]
		nextTile := tileMap[i+1]
		if currentTile.next != nextTile {
			t.Errorf("Expected tile %d next to be tile %d, but it was not", i, i+1)
		}
		if nextTile.prev != currentTile {
			t.Errorf("Expected tile %d prev to be tile %d, but it was not", i+1, i)
		}
	}
}

// createTestFile is a helper function to create a temporary file with given content for testing.
func createTestFile(t *testing.T, pattern, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temporary test file: %v", err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temporary test file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary test file: %v", err)
	}
	return tmpFile.Name()
}
