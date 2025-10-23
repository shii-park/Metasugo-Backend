package service

import (
	"encoding/json"
	"log"
	"os"
)

type ResponceData struct {
	Id     uint   "json:'id'"
	Kind   string "json:'kind'"
	Detail string "json:'detail'"
	PrevId int    "json:'prev_id'"
	NextId int    "json:'next_id'"
}

func GetData(params map[string]interface{}) (interface{}, error) {
	file, err := os.Open("../../tiles.json")
	if err != nil {
		log.Fatal("盤面ファイルを開けませんでした: ", err)
		return nil, err
	}
	defer file.Close()
	var tilesData map[string]interface{}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&tilesData); err != nil {
		log.Fatal("盤面ファイルの読み込みに失敗しました: ", err)
		return nil, err
	}
	return tilesData, nil
}
