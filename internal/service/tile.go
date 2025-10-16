package sugoroku

type TileKind int

type Tile struct {
	prev   *Tile
	next   *Tile
	kind   TileKind
	ID     int
	effect Effect
}

