package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/shii-park/Metasugo-Backend/internal/event"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

type MoneyChangedResponse struct {
	UserID string `json:"userID"`
	Money  int    `json:"money"`
}

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
	player, err := gm.game.AddPlayer(userID)
	if err != nil {
		return err
	}
	player.OnEvent = gm.onEvent
	gm.playerClients[userID] = c
	return nil
}

func (gm *GameManager) onEvent(e event.Event) {
	var message []byte
	var err error
	switch e.Type {
	case event.MoneyChanged:
		money, ok := e.Data["money"].(int)
		if !ok {
			log.Printf("error: could not assert money to int in event data: %+v", e.Data)
			return
		}
		response := MoneyChangedResponse{
			UserID: e.PlayerID,
			Money:  money,
		}
		message, err = json.Marshal(response)
		if err != nil {
			log.Printf("error: could not marshal money changed response: %v", err)
			return
		}
	}
	gm.hub.Broadcast(message)
}

func (m *GameManager) MoveByDiceRoll(playerID string, steps int) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return errors.New("Invalid player id")
	}
	player.Move(steps)

	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()

	// プレイヤーからの入力が必要な場合
	if effect.RequiresUserInput() {
		options := effect.GetOptions(currentTile)

		event := map[string]interface{}{
			"type": "USER_CHOICE_REQUIRED",
			"payload": map[string]interface{}{
				"tile_id": currentTile.GetID(),
				"options": options,
			},
		}
		return m.hub.SendToPlayer(playerID, event)

	} else {
		if err := effect.Apply(player, m.game, nil); err != nil {
			return err
		}
	}

	return nil
}

func (m *GameManager) HandlePlayerChoice(playerID string, choiceData map[string]interface{}) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}
	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()

	choice := choiceData["selection"]

	if err := effect.Apply(player, m.game, choice); err != nil {
		return fmt.Errorf("failed to apply choice :%w", err)
	}
	return nil
}
