package game

import (
	"encoding/json"
	"log"

	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

// broadcastMoneyChanged は所持金変動イベントを全クライアントに通知
func (gm *GameManager) broadcastMoneyChanged(userID string, newMoney int) {
	message, err := json.Marshal(map[string]interface{}{
		"type": "MONEY_CHANGED",
		"payload": map[string]interface{}{
			"userID":   userID,
			"newMoney": newMoney,
		},
	})

	if err != nil {
		log.Printf("error: could not marshal money changed event: %v", err)
		return
	}

	gm.hub.Broadcast(message)
}

// broadcastPlayerMoved はプレイヤー移動イベントを全クライアントに通知
func (gm *GameManager) broadcastPlayerMoved(userID string, newPosition int) {
	message, err := json.Marshal(map[string]interface{}{
		"type": "PLAYER_MOVED",
		"payload": map[string]interface{}{
			"userID":      userID,
			"newPosition": newPosition,
		},
	})

	if err != nil {
		log.Printf("error: could not marshal player moved event: %v", err)
		return
	}

	gm.hub.Broadcast(message)
}

func (gm *GameManager) handleBranchInput(p *sugoroku.Player, t *sugoroku.Tile, effect sugoroku.Effect) error {
	options := effect.GetOptions(t)
	event := map[string]interface{}{
		"type": "BRANCH_CHOICE_REQUIRED",
		"payload": map[string]interface{}{
			"tile_id": t.GetID(),
			"options": options,
		},
	}
	return gm.hub.SendToPlayer(p.GetID(), event)
}

func (gm *GameManager) handleGambleInput(p *sugoroku.Player, t *sugoroku.Tile, effect sugoroku.Effect) error {
	options := effect.GetOptions(t)
	event := map[string]interface{}{
		"type": "GAMBLE_CHOICE_REQUIRED",
		"payload": map[string]interface{}{
			"tile_id": t.GetID(),
			"options": options,
		},
	}
	return gm.hub.SendToPlayer(p.GetID(), event)
}
func (gm *GameManager) handleQuizInput(p *sugoroku.Player, t *sugoroku.Tile, effect sugoroku.Effect) error {
	options := effect.GetOptions(t)
	event := map[string]interface{}{
		"type": "QUIZ_CHOICE_REQUIRED",
		"payload": map[string]interface{}{
			"tile_id": t.GetID(),
			"options": options,
		},
	}
	return gm.hub.SendToPlayer(p.GetID(), event)
}
