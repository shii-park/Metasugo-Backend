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

func NewPlayer(id string) *Player {
	return &Player{
		id:      id,
		command: make(chan command),
	}
}

//                           __      __                        __
//                          |  \    |  \                      |  \
//  ______ ____    ______  _| $$_   | $$____    ______    ____| $$  _______
// |      \    \  /      \|   $$ \  | $$    \  /      \  /      $$ /       \
// | $$$$$$\$$$$\|  $$$$$$\\$$$$$$  | $$$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$$
// | $$ | $$ | $$| $$    $$ | $$ __ | $$  | $$| $$  | $$| $$  | $$ \$$    \
// | $$ | $$ | $$| $$$$$$$$ | $$|  \| $$  | $$| $$__/ $$| $$__| $$ _\$$$$$$\
// | $$ | $$ | $$ \$$     \  \$$  $$| $$  | $$ \$$    $$ \$$    $$|       $$
//  \$$  \$$  \$$  \$$$$$$$   \$$$$  \$$   \$$  \$$$$$$   \$$$$$$$ \$$$$$$$
//

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
