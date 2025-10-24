package sugoroku

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/shii-park/Metasugo-Backend/internal/event"
)

type Player struct {
	position *Tile
	id       string
	money    int
	mu       sync.Mutex
	OnEvent  func(e event.Event)
}

// プレイヤーのインスタンスを生成する
func NewPlayer(id string, position *Tile) *Player {
	return &Player{
		position: position,
		id:       id,
	}
}

// TODO: エラー文の追加、一時的にnextsの1こ目のマスに進むようになっている、ゴールの処理を書かなければならない
func (p *Player) moveNextTile() {
	if len(p.position.nexts) > 0 {
		p.position = p.position.nexts[0]
	}
}

// プレイヤーのIDを返すメソッド
func (p *Player) GetID() string {
	return p.id
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
		if p.position.kind == branch {
			break
		}
	}

	if err := p.position.effect.Apply(p, g); err != nil {
		return err
	}
	return nil
}

// プレイヤーのお金を増やすメソッド
func (p *Player) Profit(amount int) error {
	if amount < 0 {
		return errors.New("cannot add money by negative amount")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.money += amount
	if p.OnEvent != nil {
		p.OnEvent(event.Event{
			Type:     event.MoneyChanged,
			PlayerID: p.id,
			Data: map[string]interface{}{
				"money": p.money,
			},
		})
	}
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
	if p.OnEvent != nil {
		p.OnEvent(event.Event{
			Type:     event.MoneyChanged,
			PlayerID: p.id,
			Data: map[string]interface{}{
				"money": p.money,
			},
		})
	}
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

// 1~6までのランダムな数を返す関数
func RollDice() int {
	return rand.Intn(6) + 1
}
