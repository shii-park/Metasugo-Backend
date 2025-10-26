package game

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
	"github.com/stretchr/testify/assert"
)

// setupTestEnvironment はテスト用の共通セットアップを行います。
func setupTestEnvironment(t *testing.T) (*GameManager, *hub.Hub) {
	game := sugoroku.NewGameWithTilesForTest("../../tiles.json")
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

	time.Sleep(10 * time.Millisecond)

	return client
}

func TestGameManager_NewGameManager(t *testing.T) {
	gm, _ := setupTestEnvironment(t)
	assert.NotNil(t, gm)
	assert.NotNil(t, gm.game)
	assert.NotNil(t, gm.hub)
}

func TestGameManager_MoveByDiceRoll_ImmediateEffect(t *testing.T) {
	gm, h := setupTestEnvironment(t)
	playerID := "test_player_profit"
	_ = createAndRegisterClient(t, gm, h, playerID)

	player, _ := gm.game.GetPlayer(playerID)
	initialMoney := player.GetMoney()

	// 2マス進むと "profit" マス (ID: 3) に止まるはず
	err := gm.MoveByDiceRoll(playerID, 2)
	assert.NoError(t, err)

	// 移動後のタイルIDをログに出力
	t.Logf("Player landed on tile ID: %d", player.GetPosition().GetID())

	// 効果が即時適用され、所持金が増えていることを確認
	finalMoney := player.GetMoney()
	assert.Greater(t, finalMoney, initialMoney, "Money should increase on a profit tile")
}

func TestGameManager_MoveByDiceRoll_StopsAtBranch(t *testing.T) {
	gm, h := setupTestEnvironment(t)
	playerID := "test_player_branch"
	client := createAndRegisterClient(t, gm, h, playerID)

	// 4マス進むと "branch" マス (ID: 5) に止まる
	err := gm.MoveByDiceRoll(playerID, 4)
	assert.NoError(t, err)

	// プレイヤーの位置が分岐マス(ID: 5)であることを確認
	player, _ := gm.game.GetPlayer(playerID)
	assert.Equal(t, 5, player.GetPosition().GetID())

	// クライアントがUSER_CHOICE_REQUIREDイベントを受信することを確認
	select {
	case msg := <-client.Send:
		var event map[string]interface{}
		err := json.Unmarshal(msg, &event)
		assert.NoError(t, err)
		assert.Equal(t, "USER_CHOICE_REQUIRED", event["type"])

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for USER_CHOICE_REQUIRED event")
	}
}

func TestGameManager_HandlePlayerChoice(t *testing.T) {
	gm, h := setupTestEnvironment(t)
	playerID := "test_player_choice"
	client := createAndRegisterClient(t, gm, h, playerID)

	// --- セットアップ：先にプレイヤーを分岐マスまで移動させる ---
	err := gm.MoveByDiceRoll(playerID, 4) // 4マス進んで分岐マス(ID: 5)へ
	assert.NoError(t, err)

	// 念のため、位置とイベント受信を確認
	player, _ := gm.game.GetPlayer(playerID)
	assert.Equal(t, 5, player.GetPosition().GetID())
	// イベントチャネルをクリアする
	select {
	case <-client.Send:
	default:
	}

	// --- テスト本番：ユーザーの選択を処理 --- 
	choicePayload := map[string]interface{}{
		"selection": float64(6), // JSON経由の数値はfloat64になるため
	}

	err = gm.HandlePlayerChoice(playerID, choicePayload)
	assert.NoError(t, err)

	// プレイヤーが選択したタイルID 6 に移動したことを確認
	assert.Equal(t, 6, player.GetPosition().GetID(), "Player should have moved to the chosen tile")
}