package game

import (
	"encoding/json"
	"sync"

	"github.com/shii-park/Metasugo-Backend/internal/event"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

type MoneyResponseJSON struct {
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
			// TODO: error handling
			return
		}
		response := MoneyResponseJSON{
			UserID: e.PlayerID,
			Money:  money,
		}
		message, err = json.Marshal(response)
		if err != nil {
			// TODO: error handling
			return
		}
	}
	gm.hub.Broadcast(message)
}
