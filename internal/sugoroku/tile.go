package sugoroku

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	profit   TileKind = "profit"
	loss     TileKind = "loss"
	quiz     TileKind = "quiz"
	branch   TileKind = "branch"
	overall  TileKind = "overall"
	neighbor TileKind = "neighbor"
	require  TileKind = "require"
	gamble   TileKind = "gamble"
	goal     TileKind = "goal"
	conditional TileKind = "conditional"
	setStatus TileKind = "setStatus"
	childBonus TileKind = "childBonus"
)

const TilesJSONPath = "./tiles.json"

type TileKind string

type Tile struct {
	prevs  []*Tile
	nexts  []*Tile
	kind   TileKind
	id     int
	effect Effect
	detail string
}

// JSONの構造に対応した一時的な構造体
type TileJSON struct {
	ID      int             `json:"id"`
	Kind    TileKind        `json:"kind"`
	Detail  string          `json:"detail"`
	Effect  json.RawMessage `json:"effect"`
	PrevIDs []int           `json:"prev_ids"`
	NextIDs []int           `json:"next_ids"`
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
func NewTile(prevs []*Tile, nexts []*Tile, kind TileKind, id int, effect Effect, detail string) *Tile {
	return &Tile{
		prevs:  prevs,
		nexts:  nexts,
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

func InitTiles() map[int]*Tile {
	tileMap, err := InitTilesFromPath(TilesJSONPath)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize tiles: %v", err))
	}
	return tileMap
}

func InitTilesFromPath(path string) (map[int]*Tile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}
	defer file.Close()

	var tilesJSON []TileJSON

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tilesJSON); err != nil {
		return nil, fmt.Errorf("JSON decode error: %w", err)
	}

	tileMap := make(map[int]*Tile)

	// タイルを生成(この時点ではタイル同士はつながっていない)
	for _, tj := range tilesJSON {
		var effect Effect
		var err error
		// Effectによるマスインスタンス生成の分岐
		if len(tj.Effect) > 0 && string(tj.Effect) != "null" && string(tj.Effect) != "{}" {
			effect, err = CreateEffectFromJSON(tj.Effect)
			if err != nil {
				return nil, fmt.Errorf("failed to create effect for tile id %d: %w", tj.ID, err)
			}
		} else {
			effect = NoEffect{}
		}

		tile := NewTile(nil, nil, tj.Kind, tj.ID, effect, tj.Detail)
		tileMap[tile.id] = tile
	}

	for _, tj := range tilesJSON {
		currentTile := tileMap[tj.ID]

		for _, prevID := range tj.PrevIDs {
			currentTile.prevs = append(currentTile.prevs, tileMap[prevID])
		}

		for _, nextID := range tj.NextIDs {
			currentTile.nexts = append(currentTile.nexts, tileMap[nextID])
		}
	}

	return tileMap, nil
}

func (t *Tile) GetEffect() Effect {
	return t.effect
}

func (t *Tile) GetID() int {
	return t.id
}
