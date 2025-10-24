package service

import (
	"encoding/json"
	"errors"
	"log"
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
		file, err := os.Open("/tiles.json")
		if err != nil {
			loadErr = err
			return
		}
		defer file.Close()
		var tilesData map[string]interface{}

		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&tilesData); err != nil {
			log.Fatal("盤面ファイルの読み込みに失敗しました: ", err)
			return
		}
	})
	if loadErr != nil {
		return nil, loadErr
	}
	if tilesData == nil {
		return nil, errors.New("盤面データを読み込めませんでした")
	}
	return tilesData, nil
}
