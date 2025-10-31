package handler

import (
	log "github.com/sirupsen/logrus"

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
		log.WithError(err).Fatal("failed to create ranking handler")
	}
	bestScoreHandler, err := NewBestScoreHandler()
	if err != nil {
		log.Fatalf("MaxAmtHandlerの生成に失敗: %v", err)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
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
		//最高金額取得のルーティング
		authRequired.GET("/bestscore", bestScoreHandler.GetBestScore)
	}
}
