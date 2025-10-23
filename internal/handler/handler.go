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

func HandleRanking() {
	return
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	firebaseUID, exists := c.Get("firebase_uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインしていません"})
		return
	}
	userId := firebaseUID.(string)
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
