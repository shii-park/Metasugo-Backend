package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	//"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
)

/**************************************************/
var upgrater = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 許可オリジンの設定
		return true
	},
}

type WebSocketHandler struct {
	hub *hub.Hub
}

/**************************************************/
func NewWebSocketHandler(h *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: h}
}

/**************************************************/
func HandleRanking() {
	return
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userId := c.GetString("firebase_uid")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインしていません"})
		return
	}
	//HTTPをWebSocketに昇格

	conn, err := upgrater.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := hub.NewClient(h.hub, conn, userId)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
