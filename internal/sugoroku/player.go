package sugoroku

import (
	"errors"
	"log"
	"sync"
)

type Player struct {
	position  *Tile
	id        string
	money     int
	mu        sync.Mutex
	isMarried bool
}

// プレイヤーのインスタンスを生成する
func NewPlayer(id string, position *Tile) *Player {
	return &Player{
		position:  position,
		id:        id,
		isMarried: false,
	}
}

// TODO: nextsの1こ目のマスに進むようになっている、ゴールの処理を書かなければならない
func (p *Player) moveNextTile() {
	if len(p.position.nexts) > 0 {
		p.position = p.position.nexts[0]
	}
}

// TODO: prevsの1こ目のマスに進むようになっている
func (p *Player) movePrevTile() {
	if len(p.position.prevs) > 0 {
		p.position = p.position.prevs[0]
	}

}

// プレイヤーを指定されたマス分移動させるメソッド
func (p *Player) Move(steps int) string {
	for i := 0; i < steps; i++ {
		p.moveNextTile()
		if p.position.kind == branch {
			return "BRANCH"
		} else if p.position.kind == goal {
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
	p.money += amount
	log.Printf("PlayerProfit: %s earned %d. Wallet: %d", p.id, amount, p.money)
	return nil
}

// プレイヤーのお金を減らすメソッド
func (p *Player) Loss(amount int) error {
	if amount < 0 {
		return errors.New("cannot decrease money by negative amount")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.money -= amount
	log.Printf("PlayerLose: %s lose %d. Wallet: %d", p.id, amount, p.money)

	return nil
}

// 特定のプレイヤーのお金を増やす関数
func ProfitForTargetPlayers(players []*Player, amount int) error {
	if amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	for _, p := range players {
		p.Profit(amount)
	}
	return nil
}

// 特定のプレイヤーのお金を減らす関数
func LossForTargetPlayers(players []*Player, amount int) error {
	if amount < 0 {
		return errors.New("cannot decrease money by negative amount")
	}
	for _, p := range players {
		p.Loss(amount)
	}
	return nil
}

// プレイヤーのIDを返すメソッド
func (p *Player) GetID() string {
	return p.id
}

// プレイヤーの現在地のマス情報を返すメソッド
func (p *Player) GetPosition() *Tile {
	return p.position
}

// プレイヤーの所持金を返すメソッド
func (p *Player) GetMoney() int {
	return p.money
}

// プレイヤーの位置を移動させるメソッド(テストに用いる)
func (p *Player) SetPosition(tile *Tile) {
	p.position = tile
}

// プレイヤーを結婚させるメソッド
func (p *Player) marry() {
	p.isMarried = true
}
