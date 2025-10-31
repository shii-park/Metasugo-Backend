package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/shii-park/Metasugo-Backend/internal/handler"
	"github.com/shii-park/Metasugo-Backend/internal/logger"
	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	logger.Init()

	err := godotenv.Load()
	if err != nil {
		log.Warn(".envファイルの読み込みに失敗: ", err)
	}
	credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credFile == "" {
		log.Fatal("環境変数GOOGLE_APPLICATION_CREDENTIALSが設定されていません")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Recovery())
	err = middleware.InitFirebase()
	if err != nil {
		log.Fatal("Firebaseの初期化に失敗:", err)
	}

	//ゲームの初期化
	g := sugoroku.NewGame()
	log.Info("=== Game created ===")

	// ルーティングの設定
	handler.SetupRoutes(router, g)

	router.Run()
}
