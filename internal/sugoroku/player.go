package sugoroku

import (
	"errors"
	"math/rand"
	"sync"
)

type Player struct {
	position *Tile
	id       string
	money    int
	mu       sync.Mutex
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

func NewPlayer(id string, position *Tile) *Player {
	return &Player{
		position: position,
		id:       id,
	}
}

//	                         __      __                        __
//	                        |  \    |  \                      |  \
//	______ ____    ______  _| $$_   | $$____    ______    ____| $$  _______
//
// |      \    \  /      \|   $$ \  | $$    \  /      \  /      $$ /       \
// | $$$$$$\$$$$\|  $$$$$$\\$$$$$$  | $$$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$$
// | $$ | $$ | $$| $$    $$ | $$ __ | $$  | $$| $$  | $$| $$  | $$ \$$    \
// | $$ | $$ | $$| $$$$$$$$ | $$|  \| $$  | $$| $$__/ $$| $$__| $$ _\$$$$$$\
// | $$ | $$ | $$ \$$     \  \$$  $$| $$  | $$ \$$    $$ \$$    $$|       $$
//	\$$  \$$  \$$  \$$$$$$$   \$$$$  \$$   \$$  \$$$$$$   \$$$$$$$ \$$$$$$$

// TODO: エラー文の追加、一時的にnextsの1こ目のマスに進むようになっている
func (p *Player) moveNextTile() {
	if len(p.position.nexts) > 0 {
		p.position = p.position.nexts[0]
	}
}

// TODO: エラー文の追加、一時的にprevsの1こ目のマスに進むようになっている
func (p *Player) movePrevTile() {
	if len(p.position.prevs) > 0 {
		p.position = p.position.prevs[0]
	}
}

func (p *Player) MoveByDiceRoll(steps int, g *Game) error {
	for i := 0; i < steps; i++ {
		p.moveNextTile()
	}
	if err := p.position.effect.Apply(p, g); err != nil {
		return err
	}
	return nil
}

func (p *Player) Profit(amount int) error {
	if amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.money += amount
	return nil
}

func (p *Player) Loss(amount int) error {
	if amount < 0 {
		return errors.New("cannot decrease money by negative amount")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.money -= amount
	return nil
}

func ProfitForTargetPlayers(players []*Player, amount int) error {
	if amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	for _, p := range players {
		p.Profit(amount)
	}
	return nil
}

func LossForTargetPlayers(players []*Player, amount int) error {
	if amount < 0 {
		return errors.New("cannot decrease money by negative amount")
	}
	for _, p := range players {
		p.Loss(amount)
	}
	return nil
}

func RollDice() int {
	return rand.Intn(6) + 1
}
