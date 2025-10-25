package main

import (
	"log"
	"os"

	//"github.com/Metasugo-Backend/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"

	"github.com/shii-park/Metasugo-Backend/internal/handler"
	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルの読み込みに失敗: ", err)
	}
	credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credFile == "" {
		log.Fatal("環境変数GOOGLE_APPLICATION_CREDENTIALSが設定されていません")
	}
	tilesFile := os.Getenv("TILES_JSON_PATH")
	if tilesFile == "" {
		log.Fatal("環境変数TILES_JSON_PATHが設定されていません")
	}

	router := gin.Default()
	err = middleware.InitFirebase()
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
	router.GET("/ws/connection", middleware.AuthToken(), wsHandler.HandleWebSocket(gm))
	/**********エンドポイントここまで**********/

	router.Run()
}
