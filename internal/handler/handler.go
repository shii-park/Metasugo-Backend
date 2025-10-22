package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/service"
)

var upgrater = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleRanking() {
	return
}

func HandleWebSocket( /*w http.ResponseWriter,r *http.Request,*/ c *gin.Context) {
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

	client := &service.Client{
		ID:     uuid.New().String(),
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
		UserID: userId,
	}
	client.Hub.Register <- client

}
