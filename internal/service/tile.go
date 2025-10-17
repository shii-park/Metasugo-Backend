package sugoroku

import (
	"encoding/json"
	"log"
	"os"
)

const (
	profit   TileKind = "profit"
	loss     TileKind = "loss"
	quiz     TileKind = "quiz"
	branch   TileKind = "branch"
	overall  TileKind = "overallEffect"
	neighbor TileKind = "neighborEffect"
	require  TileKind = "require"
	gamble   TileKind = "gamnble"
)

const TilesJSONPath = "../../tiles.json"

type TileKind string

type Tile struct {
	prev   *Tile
	next   *Tile
	kind   TileKind
	id     int
	effect *Effect
	detail string
}

// JSONの構造に対応した一時的な構造体
type TileJSON struct {
	ID     int      `json:"id"`
	Kind   TileKind `json:"kind"`
	Detail string   `json:"detail"`
	Effect Effect   `json:"effect"`
	PrevID int      `json:"prev_id"`
	NextID int      `json:"next_id"`
}

//                                            __                                      __
//                                           |  \                                    |  \
//   _______   ______   _______    _______  _| $$_     ______   __    __   _______  _| $$_     ______    ______
//  /       \ /      \ |       \  /       \|   $$ \   /      \ |  \  |  \ /       \|   $$ \   /      \  /      \
// |  $$$$$$$|  $$$$$$\| $$$$$$$\|  $$$$$$$ \$$$$$$  |  $$$$$$\| $$  | $$|  $$$$$$$ \$$$$$$  |  $$$$$$\|  $$$$$$\
// | $$      | $$  | $$| $$  | $$ \$$    \   | $$ __ | $$   \$$| $$  | $$| $$        | $$ __ | $$  | $$| $$   \$$
// | $$_____ | $$__/ $$| $$  | $$ _\$$$$$$\  | $$|  \| $$      | $$__/ $$| $$_____   | $$|  \| $$__/ $$| $$
//  \$$     \ \$$    $$| $$  | $$|       $$   \$$  $$| $$       \$$    $$ \$$     \   \$$  $$ \$$    $$| $$
//   \$$$$$$$  \$$$$$$  \$$   \$$ \$$$$$$$     \$$$$  \$$        \$$$$$$   \$$$$$$$    \$$$$   \$$$$$$  \$$
//

// TODO: 完全コンストラクタ化を行うべき
func NewTile(prev *Tile, next *Tile, kind TileKind, id int, effect *Effect, detail string) *Tile {
	return &Tile{
		prev:   prev,
		next:   next,
		kind:   kind,
		id:     id,
		effect: effect,
		detail: detail,
	}
}

//  ______            __    __      __            __  __
// |      \          |  \  |  \    |  \          |  \|  \
//  \$$$$$$ _______   \$$ _| $$_    \$$  ______  | $$ \$$ ________   ______
//   | $$  |       \ |  \|   $$ \  |  \ |      \ | $$|  \|        \ /      \
//   | $$  | $$$$$$$\| $$ \$$$$$$  | $$  \$$$$$$\| $$| $$ \$$$$$$$$|  $$$$$$\
//   | $$  | $$  | $$| $$  | $$ __ | $$ /      $$| $$| $$  /    $$ | $$    $$
//  _| $$_ | $$  | $$| $$  | $$|  \| $$|  $$$$$$$| $$| $$ /  $$$$_ | $$$$$$$$
// |   $$ \| $$  | $$| $$   \$$  $$| $$ \$$    $$| $$| $$|  $$    \ \$$     \
//  \$$$$$$ \$$   \$$ \$$    \$$$$  \$$  \$$$$$$$ \$$ \$$ \$$$$$$$$  \$$$$$$$

func InitTiles() []*Tile {
	file, err := os.Open(TilesJSONPath)
	if err != nil {
		log.Fatalf("File open error: %v", err)
	}
	defer file.Close()

	var tilesJSON []TileJSON

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tilesJSON); err != nil {
		log.Fatalf("JSON decode error: %v", err)
	}

	tileMap := make(map[int]*Tile)

	tiles := make([]*Tile, 0, len(tilesJSON))

	// タイルを生成(この時点ではタイル同士はつながっていない)
	for _, tj := range tilesJSON {
		// TODO エフェクトの記述をしないといけない
		tile := NewTile(nil, nil, tj.Kind, tj.ID, nil, tj.Detail)
		tileMap[tile.id] = tile
		tiles = append(tiles, tile)
	}

	for _, tj := range tilesJSON {
		currentTile := tileMap[tj.ID]

		currentTile.prev = tileMap[tj.PrevID]
		currentTile.next = tileMap[tj.NextID]
	}

	return tiles
}
