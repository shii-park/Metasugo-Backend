package game

import (
	"errors"
	"fmt"
	"log"

	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func (gm *GameManager) HandleMove(playerID string) error {
	diceRollResult := sugoroku.RollDice()
	if err := gm.sendDiceRollResult(playerID, diceRollResult); err != nil {
		return fmt.Errorf("failed to send dice result: %w", err)
	}
	if err := gm.MoveByDiceRoll(playerID, diceRollResult); err != nil {
		return fmt.Errorf("failed to move player: %w", err)
	}
	return nil
}

// SUBMIT_BRANCHリクエスト時に発火する関数。
// 選んだタイルIDの方向へ移動させる。
func (m *GameManager) HandleBranch(playerID string, choiceData map[string]interface{}) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	// 適用前の状態を記録
	initialPosition := player.Position.GetID()
	initialMoney := player.Money

	// 選択を適用
	currentTile := player.Position
	effect := currentTile.GetEffect()
	choice := choiceData["selection"]

	if err := effect.Apply(player, m.game, choice); err != nil {
		return fmt.Errorf("failed to apply choice: %w", err)
	}

	// 適用後の最終的な状態を取得
	finalPosition := player.Position.GetID()
	finalMoney := player.Money

	// 状態が変化していれば、全クライアントに通知
	if initialPosition != finalPosition {
		m.broadcastPlayerMoved(playerID, finalPosition)
		log.Printf("PlayerMoved: %s moved to %d", playerID, player.Position.GetID())

	}
	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}

	return nil
}

// SUBMIT_GAMBLEリクエスト時に発火する関数。
// ペイロードからbetとHigh or Lowを読み込みギャンブルを行う。
// Gambleの結果をプレイヤーに返す。
func (m *GameManager) HandleGamble(playerID string, payload map[string]interface{}) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	effect := player.Position.GetEffect()

	if err := effect.Apply(player, m.game, payload); err != nil {
		return fmt.Errorf("failed to apply gamble choice: %w", err)
	}

	baseValue := 3
	bet := int(payload["bet"].(float64))
	choice := payload["choice"].(string)

	initialMoney := player.Money

	diceResult := sugoroku.RollDice()
	isHigh := diceResult >= baseValue

	playerWon := (choice == "High" && isHigh) || (choice == "Low" && !isHigh)

	amount := bet
	if playerWon {
		player.Profit(amount)
	} else {
		player.Loss(amount)
	}
	finalMoney := player.Money

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

// SUBMIT_QUIZリクエスト時に発火する関数。
// ペイロードからクイズIDと答えを読み取る。
func (m *GameManager) HandleQuiz(playerID string, payload map[string]interface{}) error {
	player, err := m.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}
	initialMoney := player.Money

	selection, ok := payload["selection"]
	if !ok {
		return errors.New("selection not found in payload")
	}

	currentTile := player.Position
	effect := currentTile.GetEffect()

	if err := effect.Apply(player, m.game, selection); err != nil {
		return fmt.Errorf("failed to apply quiz choice: %w", err)
	}

	finalMoney := player.Money

	if initialMoney != finalMoney {
		m.broadcastMoneyChanged(playerID, finalMoney)
	}
	return nil
}
