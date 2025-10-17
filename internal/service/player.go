package sugoroku

import "errors"

type command string

const (
	next command = "next"
	prev command = "prev"
)

type Player struct {
	onTile  *Tile
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

// TODO: エラー文の追加
func (p *Player) moveNextTile() {
	if p.onTile.next != nil {
		p.onTile = p.onTile.next
	}
}

// TODO: エラー文の追加
func (p *Player) movePrevTile() {
	if p.onTile.next != nil {
		p.onTile = p.onTile.prev
	}
}

func (p *Player) addMoney(amount int) error {
	if amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	p.money += amount
	return nil
}

func (p *Player) decreaseMoney(amount int) error {
	if amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	p.money -= amount
	return nil
}
