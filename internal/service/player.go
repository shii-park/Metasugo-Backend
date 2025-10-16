package sugoroku

type command string

const (
	next command = "next"
	prev command = "prev"
)

type Player struct {
	onTile  Tile
	id      string
	money   int
	command chan command
}

func NewPlayer(id string) *Player {
	return &Player{
		id:      id,
		command: make(chan command),
	}
}
