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

// TestInitTilesFromPath_HappyPath tests the successful execution of InitTilesFromPath.
func TestInitTilesFromPath_HappyPath(t *testing.T) {
	const testJSON = `[{"id": 1, "kind": "profit", "detail": "Start", "effect": {"type": "profit", "amount": 10}, "prev_id": 0, "next_id": 0}]`
	tmpFile := createTestFile(t, "happy_path_*.json", testJSON)
	defer os.Remove(tmpFile)

	tiles, err := InitTilesFromPath(tmpFile)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if len(tiles) != 1 {
		t.Fatalf("Expected 1 tile, got %d", len(tiles))
	}
	if tiles[0].id != 1 {
		t.Errorf("Expected tile ID to be 1, got %d", tiles[0].id)
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

// TestInitTilesFromPath_AllCases covers all tile kinds.
func TestInitTilesFromPath_AllCases(t *testing.T) {
	const allCasesJSON = `[
		{"id": 1, "kind": "profit", "effect": {"type": "profit", "amount": 10}},
		{"id": 2, "kind": "loss", "effect": {"type": "loss", "amount": 5}},
		{"id": 3, "kind": "quiz", "effect": {"type": "quiz", "quiz_id": 1}},
		{"id": 4, "kind": "branch", "effect": {"type": "branch", "chose_id": 2}},
		{"id": 5, "kind": "overall", "effect": {"type": "overall", "amount": 1}},
		{"id": 6, "kind": "neighbor", "effect": {"type": "neighbor", "amount": 2}},
		{"id": 7, "kind": "require", "effect": {"type": "require", "require_value": 3}},
		{"id": 8, "kind": "gamble", "effect": {"type": "gamble"}},
		{"id": 9, "kind": "unknown", "effect": {"type": "unknown"}},
		{"id": 10, "kind": "none", "effect": null}
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

	// Check a few tile kinds to ensure they were parsed correctly
	if tiles[3].kind != branch {
		t.Errorf("Expected tile 4 to be of kind 'branch', got '%s'", tiles[3].kind)
	}
	if tiles[7].kind != gamble {
		t.Errorf("Expected tile 8 to be of kind 'gamble', got '%s'", tiles[7].kind)
	}
	if tiles[8].effect != nil {
		t.Errorf("Expected tile 9 (unknown type) to have a nil effect, got %T", tiles[8].effect)
	}
	if tiles[9].effect != nil {
		t.Errorf("Expected tile 10 (null effect) to have a nil effect, got %T", tiles[9].effect)
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
