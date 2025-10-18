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
	overall  TileKind = "overall"
	neighbor TileKind = "neighbor"
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
		var effect Effect
		// Effectによるマスインスタンス生成の分岐
		if len(tj.Effect) > 0 {
			var ewt effectWithType
			if err := json.Unmarshal(tj.Effect, &ewt); err != nil {
				log.Fatalf("Effect type unmarshal error: %v", err)
			}
			switch ewt.Type {
			case profit:
				var profitEffect ProfitEffect
				if err := json.Unmarshal(tj.Effect, &profitEffect); err != nil {
					log.Fatalf("ProfitEffect unmarshal error: %v", err)
				}
				effect = profitEffect
			case loss:
				var lossEffect LossEffect
				if err := json.Unmarshal(tj.Effect, &lossEffect); err != nil {
					log.Fatalf("ProfitEffect unmarshal error: %v", err)
				}
				effect = lossEffect
			case quiz:
				var quizEffect QuizEffect
				if err := json.Unmarshal(tj.Effect, &quizEffect); err != nil {
					log.Fatalf("QuizEffect unmarshal error: %v", err)
				}
				effect = quizEffect
			case branch:
				var branchEffect BranchEffect
				if err := json.Unmarshal(tj.Effect, &branchEffect); err != nil {
					log.Fatalf("BranchEffect unmarshal error: %v", err)
				}
				effect = branchEffect
			case overall:
				var overallEffect OverallEffect
				if err := json.Unmarshal(tj.Effect, &overallEffect); err != nil {
					log.Fatalf("OverallEffect unmarshal error: %v", err)
				}
				effect = overallEffect
			case neighbor:
				var neighborEffect NeighborEffect
				if err := json.Unmarshal(tj.Effect, &neighborEffect); err != nil {
					log.Fatalf("NeighborEffect unmarshal error: %v", err)
				}
				effect = neighborEffect
			case require:
				var requireEffect RequireEffect
				if err := json.Unmarshal(tj.Effect, &requireEffect); err != nil {
					log.Fatalf("RequireEffect unmarshal error: %v", err)
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

	return tiles
}
