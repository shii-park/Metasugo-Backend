package game

import (
	log "github.com/sirupsen/logrus"

	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

// broadcastMoneyChanged は所持金変動イベントを全クライアントに通知
func (gm *GameManager) broadcastMoneyChanged(userID string, newMoney int) {
	gm.hub.Broadcast(map[string]any{
		"type": "MONEY_CHANGED",
		"payload": map[string]any{
			"userID":   userID,
			"newMoney": newMoney,
		},
	})
}

// broadcastPlayerMoved はプレイヤー移動イベントを全クライアントに通知
func (gm *GameManager) broadcastPlayerMoved(userID string, newPosition int) {
	gm.hub.Broadcast(map[string]any{
		"type": "PLAYER_MOVED",
		"payload": map[string]any{
			"userID":      userID,
			"newPosition": newPosition,
		},
	})
}
func (gm *GameManager) sendBranchSelection(player *sugoroku.Player, tile *sugoroku.Tile, effect sugoroku.BranchEffect) error {
	options := effect.GetOptions(tile)
	event := map[string]any{
		"type": "BRANCH_CHOICE_REQUIRED",
		"payload": map[string]any{
			"tileID":  tile.Id,
			"options": options,
		},
	}
	return gm.hub.SendToPlayer(player.Id, event)
}

func (gm *GameManager) sendQuizInfo(player *sugoroku.Player, tile *sugoroku.Tile, effect sugoroku.QuizEffect) error {
	quizData := effect.GetOptions(tile)
	event := map[string]any{
		"type": "QUIZ_REQUIRED",
		"payload": map[string]any{
			"tileID":   tile.Id,
			"quizData": quizData,
		},
	}
	return gm.hub.SendToPlayer(player.Id, event)
}

func (gm *GameManager) sendGambleRequire(player *sugoroku.Player, tile *sugoroku.Tile) error {
	baseValue := 3
	event := map[string]any{
		"type": "GAMBLE_REQUIRED",
		"payload": map[string]any{
			"tileID":         tile.Id,
			"referenceValue": baseValue,
		},
	}
	return gm.hub.SendToPlayer(player.Id, event)
}

func (gm *GameManager) sendGambleResult(playerID string, payload map[string]any) {
	event := map[string]any{
		"type":    "GAMBLE_RESULT",
		"payload": payload,
	}
	if err := gm.hub.SendToPlayer(playerID, event); err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"playerID": playerID,
		}).Error("failed to send gamble result to player")
	}
}

func (gm *GameManager) sendDiceRollResult(playerID string, diceResult int) error {
	event := map[string]any{
		"type": "DICE_RESULT",
		"payload": map[string]any{
			"userID":     playerID,
			"diceResult": diceResult,
		},
	}
	return gm.hub.SendToPlayer(playerID, event)
}

// broadcastPlayerFinished はプレイヤーがゴールしたことを全クライアントに通知
func (gm *GameManager) broadcastPlayerFinished(userID string, money int) {
	gm.hub.Broadcast(map[string]any{
		"type": "PLAYER_FINISHED",
		"payload": map[string]any{
			"userID": userID,
			"money":  money,
		},
	})
}

// broadcastPlayerStatusChanged はプレイヤーステータス変更イベントを全クライアントに通知
func (gm *GameManager) broadcastPlayerStatusChanged(userID string, status string, value any) {
	gm.hub.Broadcast(map[string]any{
		"type": "PLAYER_STATUS_CHANGED",
		"payload": map[string]any{
			"userID": userID,
			"status": status,
			"value":  value,
		},
	})
}
