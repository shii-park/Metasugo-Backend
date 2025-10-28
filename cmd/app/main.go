package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

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

	router := gin.Default()
	err = middleware.InitFirebase()
	if err != nil {
		log.Fatal("Firebaseの初期化に失敗:", err)
	}

	//ゲームの初期化
	g := sugoroku.NewGame()
	log.Println("=== Game created ===")

	// ルーティングの設定
	handler.SetupRoutes(router, g)

	router.Run()
}
