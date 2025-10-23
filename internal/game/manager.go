package game

import (
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

type GameManager struct {
	game          *sugoroku.Game
	hub           *hub.Hub
	playerClients map[string]*hub.Client
}

func NewGameManager(g *sugoroku.Game, h *hub.Hub) *GameManager {
	return &GameManager{
		game:          g,
		hub:           h,
		playerClients: make(map[string]*hub.Client),
	}
}

func (gm *GameManager) RegisterPlayerClient(userID string, c *hub.Client) {
	gm.game.AddPlayer(userID)
	gm.playerClients[userID] = c
}
