package sugoroku

import "errors"

type Effect interface {
	Apply(player *Player) error
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
	Amount int `json:"amount"`
}

type LossEffect struct {
	Amount int `json:"amount"`
}

type QuizEffect struct {
	QuizID int `json:"quiz_id"`
}

// TODO 一時的な型の実装をしている。　また変更するかも
type BranchEffect struct {
	ChoseID int `json:"chose_id"`
}

type OverallEffect struct {
	Amount int `json:"amount"`
}

type NeighborEffect struct {
	Amount int `json:"amount"`
}

type RequireEffect struct {
	RequireValue int `json:"require_value"`
}

type GambleEffect struct {
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
		return errors.New("cannot decrease money by negative amount")
	}
	p.money -= e.Amount
	return nil
}

// TODO効果の実装
func (e QuizEffect) Apply(p *Player) error {

	return nil
}

// TODO効果の実装
func (e BranchEffect) Apply(p *Player) error {
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

func (e GambleEffect) Apply(p *Player) error {
	return nil
}
