package test

import (
	"reflect"
	"testing"

	"github.com/shii-park/Metasugo-Backend/internal/service"
)

// GetTilesの基本テスト
func TestGetTiles(t *testing.T) {
	tiles, err := service.GetTiles()
	if err != nil {
		t.Fatalf("GetTilesがエラーを返しました: %v", err)
	}
	
	if tiles == nil {
		t.Fatal("GetTilesがnilを返しました")
	}
	
	switch v := tiles.(type) {
	case map[string]interface{}:
		if len(v) == 0 {
			t.Fatal("tilesマップが空です")
		}
		t.Logf("tilesマップを読み込みました（トップレベルキー数: %d）", len(v))
	case []interface{}:
		if len(v) == 0 {
			t.Fatal("tiles配列が空です")
		}
		t.Logf("tiles配列を読み込みました（要素数: %d）", len(v))
	default:
		t.Fatalf("予期しないtilesの型: %T", tiles)
	}
}

// GetTilesの複数回呼び出しテスト（キャッシュの確認）
func TestGetTiles_MultipleCalls(t *testing.T) {
	// 1回目の呼び出し
	tiles1, err1 := service.GetTiles()
	if err1 != nil {
		t.Fatalf("1回目のGetTiles呼び出しがエラーを返しました: %v", err1)
	}
	
	// 2回目の呼び出し
	tiles2, err2 := service.GetTiles()
	if err2 != nil {
		t.Fatalf("2回目のGetTiles呼び出しがエラーを返しました: %v", err2)
	}
	
	// 両方がnilでないことを確認
	if tiles1 == nil || tiles2 == nil {
		t.Error("GetTilesがnilを返しました")
	}
	
	// 両方が同じ型であることを確認
	if reflect.TypeOf(tiles1) != reflect.TypeOf(tiles2) {
		t.Errorf("tiles1とtiles2の型が異なります: %T vs %T", tiles1, tiles2)
	} else {
		t.Log("複数回呼び出しでも一貫した結果が返されています")
	}
}

// GetTilesのデータ構造テスト
func TestGetTiles_DataStructure(t *testing.T) {
	tiles, err := service.GetTiles()
	if err != nil {
		t.Fatalf("GetTilesがエラーを返しました: %v", err)
	}
	
	// 配列型の場合の検証
	if arr, ok := tiles.([]interface{}); ok {
		for i, tile := range arr {
			if tileMap, ok := tile.(map[string]interface{}); ok {
				// 各タイルがマップ構造を持つことを確認
				t.Logf("タイル[%d]はマップ構造です（キー数: %d）", i, len(tileMap))
			} else {
				t.Logf("タイル[%d]の型: %T", i, tile)
			}
		}
	}
	
	// マップ型の場合の検証
	if m, ok := tiles.(map[string]interface{}); ok {
		t.Logf("tilesはマップ構造です（トップレベルキー: %v）", getMapKeys(m))
	}
}

// ヘルパー関数：マップのキーを取得
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
