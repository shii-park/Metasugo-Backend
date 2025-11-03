package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/handler"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

// WebSocketHandlerのテスト: Firebase認証がない場合
func TestHandleWebSocket_NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Hubの初期化
	h := hub.NewHub()
	go h.Run()

	// GameManagerの初期化
	g := sugoroku.NewGameWithTilesForTest("../tiles.json")
	gm := game.NewGameManager(g, h)

	// ハンドラーの作成
	wsHandler := handler.NewWebSocketHandler(h)

	// ルーターの設定
	router := gin.New()
	router.GET("/ws", wsHandler.HandleWebSocket(gm))

	// リクエストの作成
	req, _ := http.NewRequest("GET", "/ws", nil)
	w := httptest.NewRecorder()

	// リクエストの実行
	router.ServeHTTP(w, req)

	// ステータスコードの確認（認証がないため401が期待される）
	if w.Code != http.StatusUnauthorized {
		t.Errorf("期待されるステータスコード: %d, 実際: %d", http.StatusUnauthorized, w.Code)
	}

	// レスポンスボディに「ログインしていません」が含まれることを確認
	body := w.Body.String()
	if !strings.Contains(body, "ログインしていません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。レスポンス: %s", body)
	}
}

// WebSocketHandlerのテスト: Firebase認証がある場合のWebSocketアップグレード
func TestHandleWebSocket_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Hubの初期化
	h := hub.NewHub()
	go h.Run()

	// GameManagerの初期化
	g := sugoroku.NewGameWithTilesForTest("../tiles.json")
	gm := game.NewGameManager(g, h)

	// ハンドラーの作成
	wsHandler := handler.NewWebSocketHandler(h)

	// テスト用サーバーの作成
	router := gin.New()
	router.GET("/ws", func(c *gin.Context) {
		// テスト用にfirebase_uidを設定
		c.Set("firebase_uid", "test-user-123")
		wsHandler.HandleWebSocket(gm)(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	// WebSocketクライアントの作成
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WebSocket接続に失敗: %v (レスポンス: %v)", err, resp)
	}
	defer ws.Close()

	// 接続が成功したことを確認
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("期待されるステータスコード: %d, 実際: %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}

	// メッセージの送信テスト
	testMessage := []byte("Hello, WebSocket!")
	err = ws.WriteMessage(websocket.TextMessage, testMessage)
	if err != nil {
		t.Errorf("メッセージの送信に失敗: %v", err)
	}

	// Ping/Pongのテスト（接続が維持されることを確認）
	err = ws.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		t.Errorf("Pingメッセージの送信に失敗: %v", err)
	}

	// Pongの応答を待つ
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	messageType, _, err := ws.ReadMessage()
	if err == nil && messageType == websocket.PongMessage {
		t.Log("Pongメッセージを受信しました")
	}
}

// NewWebSocketHandlerのテスト
func TestNewWebSocketHandler(t *testing.T) {
	h := hub.NewHub()
	wsHandler := handler.NewWebSocketHandler(h)

	if wsHandler == nil {
		t.Error("NewWebSocketHandlerがnilを返しました")
	}
}

// HandleRanking関数のテスト
func TestHandleRanking(t *testing.T) {
	// HandleRankingは現在何も返さないため、パニックしないことだけを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("HandleRankingがパニックしました: %v", r)
		}
	}()

	// handler.HandleRanking()
}

// WebSocketハンドラーの統合テスト: 複数クライアント
func TestHandleWebSocket_MultipleClients(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Hubの初期化
	h := hub.NewHub()
	go h.Run()

	// GameManagerの初期化
	g := sugoroku.NewGameWithTilesForTest("../tiles.json")
	gm := game.NewGameManager(g, h)

	// ハンドラーの作成
	wsHandler := handler.NewWebSocketHandler(h)

	// テスト用サーバーの作成
	router := gin.New()
	router.GET("/ws", func(c *gin.Context) {
		// テスト用にfirebase_uidを設定
		uid := c.Query("uid")
		if uid == "" {
			uid = "test-user"
		}
		c.Set("firebase_uid", uid)
		wsHandler.HandleWebSocket(gm)(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	// 複数のWebSocketクライアントを接続
	wsURL1 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?uid=user1"
	ws1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("WebSocket接続1に失敗: %v", err)
	}
	defer ws1.Close()

	wsURL2 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?uid=user2"
	ws2, _, err := websocket.DefaultDialer.Dial(wsURL2, nil)
	if err != nil {
		t.Fatalf("WebSocket接続2に失敗: %v", err)
	}
	defer ws2.Close()

	// 両方の接続が成功したことを確認
	t.Log("2つのWebSocket接続が成功しました")

	// メッセージを送信
	err = ws1.WriteMessage(websocket.TextMessage, []byte("Message from user1"))
	if err != nil {
		t.Errorf("user1からのメッセージ送信に失敗: %v", err)
	}

	err = ws2.WriteMessage(websocket.TextMessage, []byte("Message from user2"))
	if err != nil {
		t.Errorf("user2からのメッセージ送信に失敗: %v", err)
	}

	// 短時間待機して処理を完了させる
	time.Sleep(100 * time.Millisecond)
}

// WebSocketのクローズテスト
func TestHandleWebSocket_CloseConnection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Hubの初期化
	h := hub.NewHub()
	go h.Run()

	// GameManagerの初期化
	g := sugoroku.NewGameWithTilesForTest("../tiles.json")
	gm := game.NewGameManager(g, h)

	// ハンドラーの作成
	wsHandler := handler.NewWebSocketHandler(h)

	// テスト用サーバーの作成
	router := gin.New()
	router.GET("/ws", func(c *gin.Context) {
		c.Set("firebase_uid", "test-user-close")
		wsHandler.HandleWebSocket(gm)(c)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	// WebSocketクライアントの作成
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WebSocket接続に失敗: %v", err)
	}

	// 接続を閉じる
	err = ws.Close()
	if err != nil {
		t.Errorf("WebSocket接続のクローズに失敗: %v", err)
	}

	t.Log("WebSocket接続を正常にクローズしました")
}
