package game

import (
	"fmt"
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
func (gm *GameManager) MoveByDiceRoll(playerID string, steps int) error {
	player, err := gm.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("failed to get player: %w", err)
	}

	// 1. 移動前の状態を記録
	initialPosition := player.GetPosition().GetID()
	initialMoney := player.GetMoney()

	// 2. プレイヤーを移動させる
	flag := player.Move(steps) //めんどくさくなったのでフラグで実装してる。Effect型で比較するなどもっといいやり方はあると思う

	// 効果を判定
	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()

	if effect.RequiresUserInput() || flag == "GOAL" {
		// ユーザからの入力が必要な場合、もしくはゴールの場合こちらで処理
		switch e := effect.(type) {
		case sugoroku.BranchEffect:
			return gm.sendBranchSelection(player, currentTile, e)
		case sugoroku.QuizEffect:
			return gm.sendQuizInfo(player, currentTile, e)
		case sugoroku.GambleEffect:
			return gm.sendGambleRequire(player, currentTile)
		case sugoroku.GoalEffect:
			return gm.Goal(playerID, gm.playerClients[playerID]) // TODO: ゴールした際に行う処理(clientとの接続解除など)を行ったほうが良いと思う
		default:
			return fmt.Errorf("unhandled user input required for effect type %T", e)
		}
	} else {
		// 即時効果を適用
		if err := effect.Apply(player, gm.game, nil); err != nil {
			return err
		}
	}

	// 4. 最終的な状態の変化を検知して通知
	finalPosition := player.GetPosition().GetID()
	finalMoney := player.GetMoney()

	if initialPosition != finalPosition {
		gm.broadcastPlayerMoved(playerID, finalPosition)
	}
	if initialMoney != finalMoney {
		gm.broadcastMoneyChanged(playerID, finalMoney)
	}

	return nil
}

func (gm *GameManager) RegisterPlayerClient(playerID string, c *hub.Client) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	_, err := gm.game.AddPlayer(playerID)
	if err != nil {
		return err
	}
	gm.playerClients[playerID] = c
	return nil
}

func (gm *GameManager) UnregisterPlayerClient(playerID string, c *hub.Client) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	// ゲームからプレイヤーを削除
	err := gm.game.DeletePlayer(playerID)
	if err != nil {
		return err
	}
	// GameManagerからプレイヤーを削除
	delete(gm.playerClients, playerID)
	return nil

}

func (gm *GameManager) Goal(playerID string, c *hub.Client) error {
	if err := gm.sendGoal(playerID); err != nil {
		return err
	}
	if err := gm.UnregisterPlayerClient(playerID, c); err != nil {
		return err
	}
	return nil
}
