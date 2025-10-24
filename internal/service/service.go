package service

import (
	"embed"
	"encoding/json"
	"errors"
	"sync"
)

//go:embed assets/tiles.json
var tilesFS embed.FS

var (
	loadOnce  sync.Once
	loadErr   error
	tilesData map[string]interface{}
)

func GetTiles() (interface{}, error) {
	loadOnce.Do(func() {
		b, err := tilesFS.ReadFile("assets/tiles.json")
		if err != nil {
			loadErr = err
			return
		}
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			loadErr = err
			return
		}
		tilesData = m
	})

	if loadErr != nil {
		return nil, loadErr
	}
	if tilesData == nil {
		return nil, errors.New("盤面データを読み込めませんでした")
	}
	return tilesData, nil
}
