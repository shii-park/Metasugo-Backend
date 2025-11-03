package test

import (
	"testing"

	"github.com/shii-park/Metasugo-Backend/internal/service"
)

func TestGetTiles(t *testing.T) {
	tiles, err := service.GetTiles()
	if err != nil {
		t.Fatalf("GetTiles returned error: %v", err)
	}
	switch v := tiles.(type) {
	case map[string]any:
		if len(v) == 0 {
			t.Fatalf("tiles map is empty")
		}
		t.Logf("loaded %d top-level keys", len(v))
	case []any:
		if len(v) == 0 {
			t.Fatalf("tiles array is empty")
		}
		t.Logf("loaded tiles array with %d elements", len(v))
	default:
		t.Fatalf("unexpected tiles type: %T", tiles)
	}
}
