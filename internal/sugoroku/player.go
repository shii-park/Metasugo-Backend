package sugoroku

import (
	"errors"
	"log"
	"sync"
)

const (
	JobProfessor = "professor"
	JobLecturer  = "lecturer"

	initialMoney = 1000000
)

type Player struct {
	Position    *Tile
	Id          string
	Money       int
	mu          sync.Mutex
	IsMarried   bool
	HasChildren bool
	Job         string
}

// プレイヤーのインスタンスを生成する
func NewPlayer(id string, position *Tile) *Player {
	return &Player{
		Position:  position,
		Id:        id,
		IsMarried: false,
		Money:     initialMoney,
	}
}

// TODO: nextsの1こ目のマスに進むようになっている、ゴールの処理を書かなければならない
func (p *Player) moveNextTile() {
	if len(p.Position.nexts) > 0 {
		p.Position = p.Position.nexts[0]
	}
}

// TODO: prevsの1こ目のマスに進むようになっている
func (p *Player) movePrevTile() {
	if len(p.Position.prevs) > 0 {
		p.Position = p.Position.prevs[0]
	}

}

// プレイヤーを指定されたマス分移動させるメソッド
func (p *Player) Move(steps int) string {
	for i := 0; i < steps; i++ {
		p.moveNextTile()
		if p.Position.kind == branch {
			return "BRANCH"
		} else if p.Position.kind == goal {
			return "GOAL"
		}
	}
	return ""
}

// プレイヤーのお金を増やすメソッド
func (p *Player) Profit(amount int) error {
	if amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Money += amount
	log.Printf("PlayerProfit: %s earned %d. Wallet: %d", p.Id, amount, p.Money)
	return nil
}

// プレイヤーのお金を減らすメソッド
func (p *Player) Loss(amount int) error {
	if amount < 0 {
		return errors.New("cannot decrease Money by negative amount")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Money -= amount
	log.Printf("PlayerLose: %s lose %d. Wallet: %d", p.Id, amount, p.Money)

	return nil
}

// 特定のプレイヤーのお金を増やす関数
func ProfitForTargetPlayers(players []*Player, amount int) error {
	if amount < 0 {
		return errors.New("cannot add Money by negative amount")
	}
	for _, p := range players {
		p.Profit(amount)
	}
	return nil
}

// 特定のプレイヤーのお金を減らす関数
func LossForTargetPlayers(players []*Player, amount int) error {
	if amount < 0 {
		return errors.New("cannot decrease Money by negative amount")
	}
	for _, p := range players {
		p.Loss(amount)
	}
	return nil
}

// プレイヤーを結婚させるメソッド
func (p *Player) marry() {
	p.IsMarried = true
}

// プレイヤーに子供を授けるメソッド
func (p *Player) haveChildren() {
	p.HasChildren = true
}

// プレイヤーの職業を設定するメソッド
func (p *Player) setJob(job string) {
	p.Job = job
}
