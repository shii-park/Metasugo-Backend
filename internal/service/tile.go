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
)

const TilesJSONPath = "../../tiles.json"

type TileKind string

type Tile struct {
	prev   *Tile
	next   *Tile
	kind   TileKind
	id     int
	effect Effect
	detail string
}

// JSONの構造に対応した一時的な構造体
type TileJSON struct {
	ID     int             `json:"id"`
	Kind   TileKind        `json:"kind"`
	Detail string          `json:"detail"`
	Effect json.RawMessage `json:"effect"`
	PrevID int             `json:"prev_id"`
	NextID int             `json:"next_id"`
}

type effectWithType struct {
	Type TileKind `json:"type"`
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
func NewTile(prev *Tile, next *Tile, kind TileKind, id int, effect Effect, detail string) *Tile {
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
	tiles, err := InitTilesFromPath(TilesJSONPath)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize tiles: %v", err))
	}
	return tiles
}

func InitTilesFromPath(path string) ([]*Tile, error) {
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

	tiles := make([]*Tile, 0, len(tilesJSON))

	// タイルを生成(この時点ではタイル同士はつながっていない)
	for _, tj := range tilesJSON {
		var effect Effect
		// Effectによるマスインスタンス生成の分岐
		if len(tj.Effect) > 0 && string(tj.Effect) != "null" && string(tj.Effect) != "{}" {
			var ewt effectWithType
			if err := json.Unmarshal(tj.Effect, &ewt); err != nil {
				return nil, fmt.Errorf("effect type unmarshal error for tile id %d: %w", tj.ID, err)
			}
			if ewt.Type == "" {
				return nil, fmt.Errorf("effect type is missing for tile id %d", tj.ID)
			}
			switch ewt.Type {
			case profit:
				var profitEffect ProfitEffect
				if err := json.Unmarshal(tj.Effect, &profitEffect); err != nil {
					return nil, fmt.Errorf("ProfitEffect unmarshal error for tile id %d: %w", tj.ID, err)
				}
				effect = profitEffect
			case loss:
				var lossEffect LossEffect
				if err := json.Unmarshal(tj.Effect, &lossEffect); err != nil {
					return nil, fmt.Errorf("LossEffect unmarshal error for tile id %d: %w", tj.ID, err)
				}
				effect = lossEffect
			case quiz:
				var quizEffect QuizEffect
				if err := json.Unmarshal(tj.Effect, &quizEffect); err != nil {
					return nil, fmt.Errorf("QuizEffect unmarshal error for tile id %d: %w", tj.ID, err)
				}
				effect = quizEffect
			case branch:
				var branchEffect BranchEffect
				if err := json.Unmarshal(tj.Effect, &branchEffect); err != nil {
					return nil, fmt.Errorf("BranchEffect unmarshal error for tile id %d: %w", tj.ID, err)
				}
				effect = branchEffect
			case overall:
				var overallEffect OverallEffect
				if err := json.Unmarshal(tj.Effect, &overallEffect); err != nil {
					return nil, fmt.Errorf("OverallEffect unmarshal error for tile id %d: %w", tj.ID, err)
				}
				effect = overallEffect
			case neighbor:
				var neighborEffect NeighborEffect
				if err := json.Unmarshal(tj.Effect, &neighborEffect); err != nil {
					return nil, fmt.Errorf("NeighborEffect unmarshal error for tile id %d: %w", tj.ID, err)
				}
				effect = neighborEffect
			case require:
				var requireEffect RequireEffect
				if err := json.Unmarshal(tj.Effect, &requireEffect); err != nil {
					return nil, fmt.Errorf("RequireEffect unmarshal error for tile id %d: %w", tj.ID, err)
				}
				effect = requireEffect
			case gamble:
				effect = GambleEffect{}
			default:
				effect = nil
			}
		}
		tile := NewTile(nil, nil, tj.Kind, tj.ID, effect, tj.Detail)
		tileMap[tile.id] = tile
		tiles = append(tiles, tile)
	}

	for _, tj := range tilesJSON {
		currentTile := tileMap[tj.ID]

		currentTile.prev = tileMap[tj.PrevID]
		currentTile.next = tileMap[tj.NextID]
	}

	return tiles, nil
}
