package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
	"github.com/stretchr/testify/assert"
)

// setupTestEnvironment はテスト用の共通セットアップを行います。
func setupTestEnvironment(t *testing.T, tilePath string) (*GameManager, *hub.Hub) {
	// クイズデータをテスト用に初期化
	originalQuizJSONPath := sugoroku.QuizJSONPath
	sugoroku.QuizJSONPath = getTestFilePath(t, "test/test_quizzes.json")
	err := sugoroku.InitQuiz()
	assert.NoError(t, err)
	// この関数の終わりに元のパスに戻す
	t.Cleanup(func() {
		sugoroku.QuizJSONPath = originalQuizJSONPath
	})

	game := sugoroku.NewGameWithTilesForTest(tilePath)
	h := hub.NewHub()
	go h.Run()
	gm := NewGameManager(game, h)
	return gm, h
}

// createAndRegisterClient はテスト用のクライアントを作成し、HubとGameManagerに登録します。
func createAndRegisterClient(t *testing.T, gm *GameManager, hub *hub.Hub, playerID string) *hub.Client {
	// websocket.Connはテストに不要なためnilのままとする
	client := hub.NewClient(nil, playerID)

	hub.Register(client)
	err := gm.RegisterPlayerClient(playerID, client)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond) // 登録処理を待つ

	return client
}

// assertEventReceived はクライアントが特定のイベントを受信したことを表明します。
func assertEventReceived(t *testing.T, client *hub.Client, expectedEventType string) map[string]any {
	select {
	case msg := <-client.Send:
		t.Logf("Received JSON: %s", msg)
		var event map[string]any
		err := json.Unmarshal(msg, &event)
		assert.NoError(t, err, "Failed to unmarshal event message")
		assert.Equal(t, expectedEventType, event["type"], "Received event type mismatch")
		payload, ok := event["payload"].(map[string]any)
		assert.True(t, ok, "Payload is not a map")
		return payload
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Timed out waiting for %s event", expectedEventType)
	}
	return nil
}

// getTestFilePath はテストファイルの相対パスを絶対パスに変換します。
func getTestFilePath(t *testing.T, relativePath string) string {
	dir, err := os.Getwd()
	assert.NoError(t, err)
	return filepath.Join(dir, "..", "..", relativePath)
}

func TestGameManager_BroadcastsPlayerMoved(t *testing.T) {
	tilePath := getTestFilePath(t, "test/test_tiles.json")
	gm, h := setupTestEnvironment(t, tilePath)

	player1ID := "player1"
	player2ID := "player2"
	_ = createAndRegisterClient(t, gm, h, player1ID)
	player2 := createAndRegisterClient(t, gm, h, player2ID)

	// player1を1マス移動させる
	err := gm.MoveByDiceRoll(player1ID, 1)
	assert.NoError(t, err)

	// player2がPLAYER_MOVEDイベントを受信することを確認
	payload := assertEventReceived(t, player2, "PLAYER_MOVED")
	assert.Equal(t, "player1", payload["userID"])
	assert.Equal(t, float64(2), payload["newPosition"], "Player should have moved to tile 2")
}

func TestGameManager_BroadcastsMoneyChanged(t *testing.T) {
	tilePath := getTestFilePath(t, "test/test_tiles.json")
	gm, h := setupTestEnvironment(t, tilePath)

	player1ID := "player1"
	player2ID := "player2"
	_ = createAndRegisterClient(t, gm, h, player1ID)
	player2 := createAndRegisterClient(t, gm, h, player2ID)

	// player1を利益マス(ID:2)に移動させる (1マス進む)
	err := gm.MoveByDiceRoll(player1ID, 1)
	assert.NoError(t, err)

	// player2がMONEY_CHANGEDイベントを受信することを確認
	// 移動イベントもブロードキャストされるため、チャネルから読み飛ばす
	<-player2.Send
	payload := assertEventReceived(t, player2, "MONEY_CHANGED")
	assert.Equal(t, "player1", payload["userID"])
	assert.Equal(t, float64(10), payload["newMoney"], "Player money should be 10")
}

func TestGameManager_SendsQuizRequired(t *testing.T) {
	tilePath := getTestFilePath(t, "test/test_tiles.json")
	gm, h := setupTestEnvironment(t, tilePath)
	player1ID := "player1"
	player1 := createAndRegisterClient(t, gm, h, player1ID)

	// player1をクイズマス(ID:3)に移動させる (2マス進む)
	err := gm.MoveByDiceRoll(player1ID, 2)
	assert.NoError(t, err)

	// player1がQUIZ_REQUIREDイベントを受信することを確認
	payload := assertEventReceived(t, player1, "QUIZ_REQUIRED")
	assert.Equal(t, float64(3), payload["tileID"])
	quizData, ok := payload["quizData"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "1 + 1は？", quizData["question"])
	options, ok := quizData["options"].([]any)
	assert.True(t, ok)
	assert.ElementsMatch(t, []any{"1", "2", "3", "4"}, options)
}

func TestGameManager_SendsBranchChoiceRequired(t *testing.T) {
	tilePath := getTestFilePath(t, "test/test_tiles.json")
	gm, h := setupTestEnvironment(t, tilePath)
	player1ID := "player1"
	player1 := createAndRegisterClient(t, gm, h, player1ID)

	// player1を分岐マス(ID:4)に移動させる (3マス進む)
	err := gm.MoveByDiceRoll(player1ID, 3)
	assert.NoError(t, err)

	// player1がBRANCH_CHOICE_REQUIREDイベントを受信することを確認
	payload := assertEventReceived(t, player1, "BRANCH_CHOICE_REQUIRED")
	assert.Equal(t, float64(4), payload["tileID"])
	options, ok := payload["options"].([]any)
	assert.True(t, ok)
	assert.ElementsMatch(t, []any{float64(5), float64(6)}, options)
}
