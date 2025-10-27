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
func (gm *GameManager) sendBranchSelection(player *sugoroku.Player, tile *sugoroku.Tile, effect sugoroku.BranchEffect) error {
	options := effect.GetOptions(tile)
	event := map[string]interface{}{
		"type": "BRANCH_CHOICE_REQUIRED",
		"payload": map[string]interface{}{
			"tileID":  tile.GetID(),
			"options": options,
		},
	}
	return gm.hub.SendToPlayer(player.GetID(), event)
}

func (gm *GameManager) sendQuizInfo(player *sugoroku.Player, tile *sugoroku.Tile, effect sugoroku.QuizEffect) error {
	quizData := effect.GetOptions(tile)
	event := map[string]interface{}{
		"type": "QUIZ_REQUIRED",
		"payload": map[string]interface{}{
			"tileID":   tile.GetID(),
			"quizData": quizData,
		},
	}
	return gm.hub.SendToPlayer(player.GetID(), event)
}

func (gm *GameManager) sendGambleRequire(player *sugoroku.Player, tile *sugoroku.Tile) error {
	baseValue := 3
	event := map[string]interface{}{
		"type": "GAMBLE_REQUIRED",
		"payload": map[string]interface{}{
			"tileID":         tile.GetID(),
			"referenceValue": baseValue,
		},
	}
	return gm.hub.SendToPlayer(player.GetID(), event)
}

func (gm *GameManager) sendGambleResult(playerID string, payload map[string]interface{}) {
	message, err := json.Marshal(map[string]interface{}{
		"type":    "GAMBLE_RESULT",
		"payload": payload,
	})
	if err != nil {
		log.Printf("error: could not marshal gamble result event: %v", err)
		return
	}
	if err := gm.hub.SendToPlayer(playerID, message); err != nil {
		log.Printf("error: failed to send gamble result to player %s: %v", playerID, err)
	}
}

func (gm *GameManager) sendDiceRollResult(playerID string, diceResult int) error {
	baseValue := 3
	event := map[string]interface{}{
		"type": "DICE_RESULT",
		"payload": map[string]interface{}{
			"userID":     playerID,
			"diceResult": baseValue,
		},
	}
	return gm.hub.SendToPlayer(playerID, event)
}
