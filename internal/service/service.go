package service

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	loadOnce  sync.Once
	loadErr   error
	tilesData map[string]interface{}
)

func GetTiles() (interface{}, error) {
	loadOnce.Do(func() {
		file, err := os.Open("TILES_JSON_PATH") //TODO: パスを環境変数に設定
		if err != nil {
			loadErr = err
			return
		}
		defer file.Close()
		var m map[string]interface{}

		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&m); err != nil {
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
