package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func SetupRoutes(router *gin.Engine, sg *sugoroku.Game) {
	// Hubの初期化
	hub := hub.NewHub()
	go hub.Run()

	// GameManagerの初期化
	gm := game.NewGameManager(sg, hub)

	// WebSocketHandlerの初期化
	wsHandler := NewWebSocketHandler(hub)

	// RankingHandlerの初期化
	rankingHandler, err := NewRankingHandler()
	if err != nil {
		log.Fatalf("failed to create ranking handler: %v", err)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 認証が必要なルートのグループを作成
	authRequired := router.Group("/")
	authRequired.Use(middleware.AuthToken())
	{
		// WebSocketのルーティング
		authRequired.GET("/ws", wsHandler.HandleWebSocket(gm))
		// ランキングのルーティング
		authRequired.GET("/ranking", rankingHandler.GetRanking)
		// タイルのルーティング
		authRequired.GET("/tiles", TilesHandler)
	}
}
