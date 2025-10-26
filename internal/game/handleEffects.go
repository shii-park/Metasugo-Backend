package game

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func (m *GameManager) HandleBranch(playerID string, choiceData map[string]interface{}) error {
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

func (m *GameManager) HandleGamble(playerID string, payload map[string]interface{}) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	effect := player.GetPosition().GetEffect()

	if err := effect.Apply(player, m.game, payload); err != nil {
		return fmt.Errorf("invalid gamble input: %w", err)
	}

	baseValue := 3
	bet := int(payload["bet"].(float64))
	choice := payload["choice"].(string)

	initialMoney := player.GetMoney()

	diceResult := sugoroku.RollDice()
	isHigh := diceResult >= baseValue

	playerWon := (choice == "High" && isHigh) || (choice == "Low" && !isHigh)

	amount := bet
	if playerWon {
		player.Profit(amount)
	} else {
		player.Loss(amount)
	}
	finalMoney := player.GetMoney()

	resultPayload := map[string]interface{}{
		"userID":     playerID,
		"diceResult": diceResult,
		"choice":     choice,
		"won":        playerWon,
		"amount":     amount,
		"newMoney":   finalMoney,
	}
	m.sendGambleResult(playerID, resultPayload)

	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}

	return nil
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
