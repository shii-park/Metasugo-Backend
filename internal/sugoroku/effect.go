package sugoroku

import "errors"

type Effect interface {
	Apply(player *Player, game *Game) error
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
	ProfitAmount int `json:"profit_amount"`
	LossAmount   int `json:"loss_amount"`
}

type NeighborEffect struct {
	ProfitAmount int `json:"profit_amount"`
	LossAmount   int `json:"loss_amount"`
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

func (e ProfitEffect) Apply(p *Player, g *Game) error {
	err := p.Profit(e.Amount)
	return err
}

func (e LossEffect) Apply(p *Player, g *Game) error {
	err := p.Loss(e.Amount)
	return err
}

// TODO効果の実装
func (e QuizEffect) Apply(p *Player, g *Game) error {

	return nil
}

// TODO効果の実装
func (e BranchEffect) Apply(p *Player, g *Game) error {
	return nil
}

// TODO効果の実装
func (e OverallEffect) Apply(p *Player, g *Game) error {
	targetPlayers := g.GetAllPlayers()
	if e.ProfitAmount > 0 {
		// 全体にお金をもらう
		p.Profit(e.ProfitAmount)
		amount := DistributeMoney(targetPlayers, e.ProfitAmount)
		LossForTargetPlayers(targetPlayers, amount)
	} else if e.LossAmount > 0 {
		// 全員にお金を配る
		p.Loss(e.LossAmount)
		amount := DistributeMoney(targetPlayers, e.LossAmount)
		ProfitForTargetPlayers(targetPlayers, amount)
	} else {
		return errors.New("invalid amount for overall effect")
	}
	return nil
}

// TODO効果の実装
func (e NeighborEffect) Apply(p *Player, g *Game) error {
	targetPlayers := g.GetNeighbors(p)
	if e.ProfitAmount > 0 {
		// 全体にお金をもらう
		p.Profit(e.ProfitAmount)
		amount := DistributeMoney(targetPlayers, e.ProfitAmount)
		LossForTargetPlayers(targetPlayers, amount)
	} else if e.LossAmount > 0 {
		// 全員にお金を配る
		p.Loss(e.LossAmount)
		amount := DistributeMoney(targetPlayers, e.LossAmount)
		ProfitForTargetPlayers(targetPlayers, amount)
	} else {
		return errors.New("invalid amount for overall effect")
	}
	return nil
}

// TODO効果の実装
func (e RequireEffect) Apply(p *Player, g *Game) error {

	return nil
}

func (e GambleEffect) Apply(p *Player, g *Game) error {
	return nil
}
