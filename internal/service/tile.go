package sugoroku

type TileKind string

const (
	profit         TileKind = "profit"
	loss           TileKind = "loss"
	quiz           TileKind = "quiz"
	branch         TileKind = "branch"
	overallEffect  TileKind = "overallEffect"
	neighborEffect TileKind = "neighborEffect"
	require        TileKind = "require"
	gamble         TileKind = "gamnble"
)

type Tile struct {
	prev   *Tile
	next   *Tile
	kind   TileKind
	id     int
	effect Effect
	detail string
}

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
