package game

import (
	"fmt"
	"log"
	"sync"

	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

type GameManager struct {
	game          *sugoroku.Game
	hub           *hub.Hub
	playerClients map[string]*hub.Client
	mu            sync.RWMutex
}

func NewGameManager(g *sugoroku.Game, h *hub.Hub) *GameManager {
	return &GameManager{
		game:          g,
		hub:           h,
		playerClients: make(map[string]*hub.Client),
	}
}

func (gm *GameManager) RegisterPlayerClient(userID string, c *hub.Client) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	_, err := gm.game.AddPlayer(userID)
	if err != nil {
		return err
	}
	gm.playerClients[userID] = c
	return nil
}

func (m *GameManager) MoveByDiceRoll(playerID string, steps int) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("failed to get player: %w", err)
	}

	// 1. 移動前の状態を記録
	initialPosition := player.GetPosition().GetID()
	initialMoney := player.GetMoney()

	// 2. プレイヤーを移動させる
	flag := player.Move(steps) //めんどくさくなったのでフラグで実装してる。Effect型で比較するなどもっといいやり方はあると思う

	log.Printf("PlayerMoved: %s moved to %d", playerID, player.GetPosition().GetID())
	// 3. マス効果を判定・適用

	finalPosition := player.GetPosition().GetID()

	if initialPosition != finalPosition {
		m.broadcastPlayerMoved(playerID, finalPosition)
	}

	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()

	if effect.RequiresUserInput() || flag == "GOAL" {
		// 3a. ユーザー入力が必要な場合 (Applyはここでは呼ばない)
		switch e := effect.(type) {
		case sugoroku.BranchEffect:
			return m.sendBranchSelection(player, currentTile, e)
		case sugoroku.QuizEffect:
			return m.sendQuizInfo(player, currentTile, e)
		case sugoroku.GambleEffect:
			return m.sendGambleRequire(player, currentTile)
		case sugoroku.GoalEffect:
			return m.sendGoal(playerID) // TODO: ゴールした際に行う処理(clientとの接続解除など)を行ったほうが良いと思う
		default:
			return fmt.Errorf("unhandled user input required for effect type %T", e)
		}
	} else {
		// 3b. 即時効果の場合 (ここでApplyを呼ぶ)
		if err := effect.Apply(player, m.game, nil); err != nil {
			return err
		}
	}

	finalMoney := player.GetMoney()
	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}
	return nil
}
