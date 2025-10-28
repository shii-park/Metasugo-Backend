package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/service"
)

//ハンドラを分割予定

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

		client := h.hub.NewClient(conn, userID)

		h.hub.Register(client)
		gm.RegisterPlayerClient(userID, client)

		go client.WritePump()
		go client.ReadPump()

		go h.processMessage(gm, client, userID)

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

func (h *WebSocketHandler) processMessage(gm *game.GameManager, client *hub.Client, userID string) {
	for message := range client.Receive {
		var req wsRequest
		if err := json.Unmarshal(message, &req); err != nil {
			_ = client.SendJSON(gin.H{ //TODO: sendErrorをつかうようにする
				"type": "error", "code": "invalid_json", "message": "JSONの解析に失敗しました"})
			continue
		}
		switch req.Type {
		case "ROLL_DICE":
			if err := gm.HandleMove(userID); err != nil {
				log.Printf("Error during handleRollDice: %v", err)
			}
		case "SUBMIT_CHOICE":
			if err := gm.HandleBranch(userID, req.Payload); err != nil {
				log.Printf("Error during handleSubmitChoice: %v", err)
			}
		case "SUBMIT_GAMBLE":
			if err := gm.HandleGamble(userID, req.Payload); err != nil {
				log.Printf("Error during handleGamble: %v", err)
			}
		case "SUBMIT_QUIZ":
			if err := gm.HandleQuiz(userID, req.Payload); err != nil {
				log.Printf("Error during handleQuiz: %v", err)
			}
		default:
			_ = client.SendJSON(gin.H{ //TODO: sendErrorをつかうようにする
				"type": "error", "code": "unknown_request", "message": "未対応のリクエストです",
			})
		}
	}
}
