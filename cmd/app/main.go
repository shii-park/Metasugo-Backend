package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/shii-park/Metasugo-Backend/internal/handler"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	router := gin.Default()
	err := middleware.InitFirebase()
	if err != nil {
		log.Fatal("Firebaseの初期化に失敗:", err)
	}
	h := hub.NewHub()
	go h.Run()
	wsHandler := handler.NewWebSocketHandler(h)

	//いろんなエンドポイントをつくろう
	router.GET("/ranking", handler.HandleRanking)
	router.GET("/ws/connection", middleware.AuthToken(), wsHandler.HandleWebSocket)

	sugoroku.NewGame() // ハンドラができた際に、gameにAddplayerができるようになる
	router.Run()
}
