package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
		log.Fatal("環境変数 GOOGLE_APPLICATION_CREDENTIALS が設定されていません")
	}

	router := gin.Default()

	// CORS 設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // フロントエンドURL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Firebase初期化
	err = middleware.InitFirebase()
	if err != nil {
		log.Fatal("Firebaseの初期化に失敗:", err)
	}

	// ゲームの初期化
	g := sugoroku.NewGame()
	log.Println("=== Game created ===")

	// ルーティング設定
	handler.SetupRoutes(router, g)

	router.Run() // デフォルトで :8080
}