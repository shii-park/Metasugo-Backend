package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func SetupRoutes(router *gin.Engine, sg *sugoroku.Game) {
	// Hubの初期化
	hub := hub.NewHub()
	go hub.Run()

	// GameManagerの初期化
	gm := game.NewGameManager(sg, hub)

	// WebSocketHandlerの初期化とルーティング
	wsHandler := NewWebSocketHandler(hub)
	router.GET("/ws", wsHandler.HandleWebSocket(gm))

	// RankingHandlerの初期化とルーティング
	rankingHandler, err := NewRankingHandler()
	if err != nil {
		log.Fatalf("failed to create ranking handler: %v", err)
	}
	router.GET("/ranking", rankingHandler.GetRanking)
}
