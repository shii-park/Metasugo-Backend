package sugoroku

type command string

const (
	next command = "next"
	prev command = "prev"
)

type Player struct {
	onTile  Tile
	id      string
	command chan command
}

func NewPlayer(id string) *Player {
	return &Player{
		id:      id,
		command: make(chan command),
	}
}
