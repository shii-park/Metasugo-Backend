package sugoroku

type TileKind int

type Tile struct {
	prev   *Tile
	next   *Tile
	kind   TileKind
	ID     int
	effect Effect
}

// TODO: 完全コンストラクタ化を行うべき
func NewTile(prev *Tile, next *Tile, kind TileKind, ID int, effect Effect) *Tile {
	return &Tile{
		prev:   prev,
		next:   next,
		kind:   kind,
		ID:     ID,
		effect: effect,
	}
}
