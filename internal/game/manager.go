package game

import (
	"errors"
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
		return errors.New("invalid player id")
	}

	// 移動前の状態を記録
	initialPosition := player.GetPosition().GetID()
	initialMoney := player.GetMoney()

	// プレイヤーを移動させる
	player.Move(steps)

	// マス効果を適用
	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()

	if effect.RequiresUserInput() {
		// ユーザー入力が必要な場合は、選択肢をクライアントに送信
		options := effect.GetOptions(currentTile)
		event := map[string]interface{}{
			"type": "USER_CHOICE_REQUIRED",
			"payload": map[string]interface{}{
				"tile_id": currentTile.GetID(),
				"options": options,
			},
		}
		// TODO: SendToPlayer を SendToClient に修正する必要があるか確認
		return m.hub.SendToPlayer(playerID, event)
	}

	// 即時効果を適用
	if err := effect.Apply(player, m.game, nil); err != nil {
		return err
	}

	// 移動と効果適用後の最終的な状態を取得
	finalPosition := player.GetPosition().GetID()
	finalMoney := player.GetMoney()

	// 状態が変化していれば、全クライアントに通知
	if initialPosition != finalPosition {
		m.broadcastPlayerMoved(playerID, finalPosition)
	}
	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}

	return nil
}

func (m *GameManager) HandlePlayerChoice(playerID string, choiceData map[string]interface{}) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	// 適用前の状態を記録
	initialPosition := player.GetPosition().GetID()
	initialMoney := player.GetMoney()

	// 選択を適用
	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()
	choice := choiceData["selection"]
	if err := effect.Apply(player, m.game, choice); err != nil {
		return fmt.Errorf("failed to apply choice: %w", err)
	}

	// 適用後の最終的な状態を取得
	finalPosition := player.GetPosition().GetID()
	finalMoney := player.GetMoney()

	// 状態が変化していれば、全クライアントに通知
	if initialPosition != finalPosition {
		m.broadcastPlayerMoved(playerID, finalPosition)
	}
	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}

	return nil
}
