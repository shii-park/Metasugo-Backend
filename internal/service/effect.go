package sugoroku

import "errors"

type Effect interface {
	Apply(player *Player)
}

//  ________                                     __             ______
// |        \                                   |  \           /      \
//  \$$$$$$$$__    __   ______    ______    ____| $$  ______  |  $$$$$$\
//    | $$  |  \  |  \ /      \  /      \  /      $$ /      \ | $$_  \$$
//    | $$  | $$  | $$|  $$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$\| $$ \
//    | $$  | $$  | $$| $$  | $$| $$    $$| $$  | $$| $$    $$| $$$$
//    | $$  | $$__/ $$| $$__/ $$| $$$$$$$$| $$__| $$| $$$$$$$$| $$
//    | $$   \$$    $$| $$    $$ \$$     \ \$$    $$ \$$     \| $$
//     \$$   _\$$$$$$$| $$$$$$$   \$$$$$$$  \$$$$$$$  \$$$$$$$ \$$
//          |  \__| $$| $$
//           \$$    $$| $$
//            \$$$$$$  \$$

type ProfitEffect struct {
	Amount int `json: "amount"`
}

type LossEffect struct {
	Amount int `json: "amount"`
}

type QuizEffect struct {
	QuizID int `json:"quiz_id"`
}

type OverallEffect struct {
	Amount int `json: "amount"`
}

type NeighborEffect struct {
	Amount int `json: "amount"`
}

type RequireEffect struct {
	RequireValue int `json: "require_value"`
}

//  __       __             __      __                        __
// |  \     /  \           |  \    |  \                      |  \
// | $$\   /  $$  ______  _| $$_   | $$____    ______    ____| $$  _______
// | $$$\ /  $$$ /      \|   $$ \  | $$    \  /      \  /      $$ /       \
// | $$$$\  $$$$|  $$$$$$\\$$$$$$  | $$$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$$
// | $$\$$ $$ $$| $$    $$ | $$ __ | $$  | $$| $$  | $$| $$  | $$ \$$    \
// | $$ \$$$| $$| $$$$$$$$ | $$|  \| $$  | $$| $$__/ $$| $$__| $$ _\$$$$$$\
// | $$  \$ | $$ \$$     \  \$$  $$| $$  | $$ \$$    $$ \$$    $$|       $$
//  \$$      \$$  \$$$$$$$   \$$$$  \$$   \$$  \$$$$$$   \$$$$$$$ \$$$$$$$

func (e ProfitEffect) Apply(p *Player) error {
	if e.Amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	p.money += e.Amount
	return nil
}

func (e LossEffect) Apply(p *Player) error {
	if e.Amount < 0 {
		return errors.New("cannot decreace money by negative amount")
	}
	p.money -= e.Amount
	return nil
}

// TODO効果の実装
func (e QuizEffect) Apply(p *Player) error {

	return nil
}

// TODO効果の実装
func (e OverallEffect) Apply(p *Player) error {

	return nil
}

// TODO効果の実装
func (e NeighborEffect) Apply(p *Player) error {

	return nil
}

// TODO効果の実装
func (e RequireEffect) Apply(p *Player) error {

	return nil
}
