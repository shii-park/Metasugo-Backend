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

	// 1. 移動前の状態を記録
	initialPosition := player.GetPosition().GetID()
	initialMoney := player.GetMoney()

	// 2. プレイヤーを移動させる
	player.Move(steps)

	// 3. マス効果を判定・適用
	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()

	if effect.RequiresUserInput() {
		// 3a. ユーザー入力が必要な場合 (Applyはここでは呼ばない)
		switch e := effect.(type) {
		case sugoroku.BranchEffect:
			return m.handleBranchInput(player, currentTile, e)
		case sugoroku.QuizEffect:
			return m.handleQuizInput(player, currentTile, e)
		case sugoroku.GambleEffect:
			return m.handleGambleInput(player, currentTile, e)
		default:
			return fmt.Errorf("unhandled user input required for effect type %T", e)
		}
	} else {
		// 3b. 即時効果の場合 (ここでApplyを呼ぶ)
		if err := effect.Apply(player, m.game, nil); err != nil {
			return err
		}
	}

	// 4. 最終的な状態の変化を検知して通知
	finalPosition := player.GetPosition().GetID()
	finalMoney := player.GetMoney()

	if initialPosition != finalPosition {
		m.broadcastPlayerMoved(playerID, finalPosition)
	}
	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}

	return nil
}

func (m *GameManager) HandlePlayerChoice(playerID string, choiceData map[string]any) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	initialPosition := player.GetPosition().GetID()
	initialMoney := player.GetMoney()

	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()
	choice := choiceData["selection"]
	if err := effect.Apply(player, m.game, choice); err != nil {
		return fmt.Errorf("failed to apply choice: %w", err)
	}

	finalPosition := player.GetPosition().GetID()
	finalMoney := player.GetMoney()

	if initialPosition != finalPosition {
		m.broadcastPlayerMoved(playerID, finalPosition)
	}
	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}

	return nil
}
