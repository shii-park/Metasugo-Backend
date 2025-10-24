package main

import (
	"log"

	//"github.com/Metasugo-Backend/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"

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

	//ゲームの初期化
	h := hub.NewHub()
	go h.Run()
	g := sugoroku.NewGame() // ハンドラができた際に、gameにAddplayerができるようになる
	gm := game.NewGameManager(g, h)

	wsHandler := handler.NewWebSocketHandler(h)

	/**********エンドポイント**********/
	ranking := router.Group("/ranking")
	{
		ranking.POST("/score")        //スコア追加
		ranking.GET("/top")           //トップランキング取得
		ranking.GET("/all")           //全体ランキング取得
		ranking.GET("/user/:user_id") //特定ユーザのスコア取得
		ranking.GET("/me")            //自分のランクを取得
	}
	router.GET("/ws/connection", middleware.AuthToken(), wsHandler.HandleWebSocket)
	/**********エンドポイントここまで**********/


	router.Run()
}
