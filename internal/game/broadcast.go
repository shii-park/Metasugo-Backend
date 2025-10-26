package game

import (
	"encoding/json"
	"log"
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
