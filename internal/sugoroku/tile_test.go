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

func TestInitTilesFromPath_HappyPath(t *testing.T) {
	const testJSON = `[
		{"id": 1, "kind": "profit", "detail": "Start", "prev_ids": [], "next_ids": [2]},
		{"id": 2, "kind": "loss", "detail": "Middle", "prev_ids": [1], "next_ids": [3]},
		{"id": 3, "kind": "quiz", "detail": "End", "prev_ids": [2], "next_ids": []}
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
	if len(tile1.prevs) != 0 {
		t.Errorf("Expected tile 1 to have 0 prevs, got %d", len(tile1.prevs))
	}
	if len(tile1.nexts) != 1 || tile1.nexts[0] != tile2 {
		t.Errorf("Expected tile 1 next to be [tile 2], but it was not")
	}

	if len(tile2.prevs) != 1 || tile2.prevs[0] != tile1 {
		t.Errorf("Expected tile 2 prev to be [tile 1], but it was not")
	}
	if len(tile2.nexts) != 1 || tile2.nexts[0] != tile3 {
		t.Errorf("Expected tile 2 next to be [tile 3], but it was not")
	}

	if len(tile3.prevs) != 1 || tile3.prevs[0] != tile2 {
		t.Errorf("Expected tile 3 prev to be [tile 2], but it was not")
	}
	if len(tile3.nexts) != 0 {
		t.Errorf("Expected tile 3 to have 0 nexts, got %d", len(tile3.nexts))
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
		{"id": 1, "kind": "profit", "effect": {"type": "profit", "amount": 10}, "prev_ids": [], "next_ids": [2]},
		{"id": 2, "kind": "loss", "effect": {"type": "loss", "amount": 5}, "prev_ids": [1], "next_ids": [3]},
		{"id": 3, "kind": "quiz", "effect": {"type": "quiz", "quiz_id": 1}, "prev_ids": [2], "next_ids": [4]},
		{"id": 4, "kind": "branch", "effect": {"type": "branch", "chose_id": 2}, "prev_ids": [3], "next_ids": [5, 6]},
		{"id": 5, "kind": "overall", "effect": {"type": "overall", "amount": 1}, "prev_ids": [4], "next_ids": [7]},
		{"id": 6, "kind": "neighbor", "effect": {"type": "neighbor", "amount": 2}, "prev_ids": [4], "next_ids": [7]},
		{"id": 7, "kind": "require", "effect": {"type": "require", "require_value": 3}, "prev_ids": [5, 6], "next_ids": [8]},
		{"id": 8, "kind": "gamble", "effect": {"type": "gamble"}, "prev_ids": [7], "next_ids": [9]},
		{"id": 9, "kind": "unknown", "effect": {"type": "unknown"}, "prev_ids": [8], "next_ids": [10]},
		{"id": 10, "kind": "none", "effect": null, "prev_ids": [9], "next_ids": []}
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
	if len(tileMap[1].prevs) != 0 {
		t.Error("Expected tile 1 prevs to be empty")
	}
	if len(tileMap[10].nexts) != 0 {
		t.Error("Expected tile 10 nexts to be empty")
	}

	// Check branching and merging
	if len(tileMap[4].nexts) != 2 || (tileMap[4].nexts[0] != tileMap[5] && tileMap[4].nexts[1] != tileMap[5]) || (tileMap[4].nexts[0] != tileMap[6] && tileMap[4].nexts[1] != tileMap[6]) {
		t.Errorf("Expected tile 4 to branch to tiles 5 and 6")
	}
	if len(tileMap[7].prevs) != 2 || (tileMap[7].prevs[0] != tileMap[5] && tileMap[7].prevs[1] != tileMap[5]) || (tileMap[7].prevs[0] != tileMap[6] && tileMap[7].prevs[1] != tileMap[6]) {
		t.Errorf("Expected tile 7 to merge from tiles 5 and 6")
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
