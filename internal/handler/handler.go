package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	//"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 許可オリジンの設定
		return true
	},
}

type WebSocketHandler struct {
	hub *hub.Hub
}

type wsRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func NewWebSocketHandler(h *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: h}
}

// Websocket接続時のハンドラー
func (h *WebSocketHandler) HandleWebSocket(gm *game.GameManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("firebase_uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインしていません"})
			return
		}

		//HTTPをWebSocketに昇格
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		client := hub.NewClient(h.hub, conn, userID)

		h.hub.Register(client)
		gm.RegisterPlayerClient(userID, client)

		go client.WritePump()
		//go client.ReadPump()
		defer func() {
			h.hub.Unregister(client)
			conn.Close()
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket読み込みエラー: %v", err)
				return
			}
			var req wsRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				client.SendError(err)
				continue
			}
			switch req.Type {
			case "getTile", "getTiles":
				h.HandleGetTile(client, req.Payload)
			default:
				client.SendError(nil) //TODO:エラー内容修正
			}
		}
	}
}

func (h *WebSocketHandler) HandleGetTile(client *hub.Client, request map[string]interface{}) {
	tile, err := service.GetTiles()
	if err != nil {
		_ = client.SendJSON(gin.H{"type": "error", "message": "タイルの取得に失敗しました"})
		return
	}
	_ = client.SendJSON(gin.H{"type": "tile", "data": tile})
}
