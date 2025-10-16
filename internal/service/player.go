package sugoroku

type command string

const (
	next command = "next"
	prev command = "prev"
)

type Player struct {
	onTile   Tile
	id       string
	commands chan command
}
