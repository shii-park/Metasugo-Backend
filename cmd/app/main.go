package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/shii-park/Metasugo-Backend/internal/handler"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	//.env読み込み
	err := godotenv.Load()
	if err != nil {
		log.Printf("./envファイルの読み込みに失敗しました")
	}

	credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credFile == "" {
		log.Fatal("環境変数が設定されていません")
	}

	router := gin.Default()
	err = middleware.InitFirebase()
	if err != nil {
		log.Fatal("Firebaseの初期化に失敗:", err)
	}
	h := hub.NewHub()
	go h.Run()
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

	sugoroku.NewGame() // ハンドラができた際に、gameにAddplayerができるようになる
	router.Run()
}
