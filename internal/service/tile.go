package sugoroku

type SpaceType int

type Tile struct {
	prev   *Tile
	next   *Tile
	Type   SpaceType
	ID     int
	effect Effect
}
