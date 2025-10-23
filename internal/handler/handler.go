package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	//"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
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

func NewWebSocketHandler(h *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: h}
}

func HandleRanking(c *gin.Context) {
	return
}

// Websocket接続時のハンドラー
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context, gm *game.GameManager) {
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
	go client.ReadPump()
}
