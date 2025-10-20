package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/shii-park/Metasugo-Backend/internal/handler"
	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	router := gin.Default()
	err := middleware.InitFirebase()
	if err != nil {
		log.Fatal("Firebaseの初期化に失敗:", err)
	}

	//いろんなエンドポイントをつくろう
	router.GET("/ranking", handler.Ranking())
	router.GET("/ws/connection"handler.WebSocket())

	sugoroku.NewGame() // ハンドラができた際に、gameにAddplayerができるようになる
}
