package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/shii-park/Metasugo-Backend/internal/middleware"
	"google.golang.org/api/option"
)

// テスト用のFirebaseクライアントを初期化
func setupFirebaseForTest(t *testing.T) *auth.Client {
	// firebase-service-account.jsonが存在するか確認
	if _, err := os.Stat("../firebase-service-account.json"); os.IsNotExist(err) {
		t.Skip("firebase-service-account.jsonが見つかりません。このテストをスキップします。")
	}

	opt := option.WithCredentialsFile("../firebase-service-account.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		t.Fatalf("Firebaseアプリの初期化に失敗: %v", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		t.Fatalf("Firebase Auth クライアントの初期化に失敗: %v", err)
	}

	return client
}

// テスト用のカスタムトークンを生成（実際のユーザーがいない場合のテスト用）
func createTestToken(t *testing.T, client *auth.Client) string {
	token, err := client.CustomToken(context.Background(), "test-uid")
	if err != nil {
		t.Fatalf("カスタムトークンの生成に失敗: %v", err)
	}
	return token
}

// InitFirebase関数のテスト
func TestInitFirebase(t *testing.T) {
	err := middleware.InitFirebase()
	if err != nil {
		t.Logf("InitFirebaseがエラーを返しました: %v", err)
		t.Logf("GOOGLE_APPLICATION_CREDENTIALS環境変数が設定されているか確認してください")
	} else {
		t.Log("InitFirebaseが正常に完了しました")
	}
}

// AuthToken ミドルウェアのテスト: Authorization headerが無い場合
func TestAuthToken_NoAuthHeader(t *testing.T) {
	// Ginのテストモードに設定
	gin.SetMode(gin.TestMode)

	// ルーターの作成
	router := gin.New()
	router.Use(middleware.AuthToken())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// リクエストの作成
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// リクエストの実行
	router.ServeHTTP(w, req)

	// ステータスコードの確認
	if w.Code != 401 {
		t.Errorf("期待されるステータスコード: 401, 実際: %d", w.Code)
	}

	// レスポンスボディの確認
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "Authorization header required" {
		t.Errorf("期待されるエラーメッセージと異なります: %v", response["error"])
	}
}

// AuthToken ミドルウェアのテスト: 不正なフォーマット
func TestAuthToken_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.AuthToken())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("期待されるステータスコード: 401, 実際: %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "Invalid authorization format" {
		t.Errorf("期待されるエラーメッセージと異なります: %v", response["error"])
	}
}

// AuthToken ミドルウェアのテスト: 不正なトークン
func TestAuthToken_InvalidToken(t *testing.T) {
	// Firebaseの初期化
	err := middleware.InitFirebase()
	if err != nil {
		t.Skip("Firebaseの初期化に失敗したため、このテストをスキップします")
	}

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.AuthToken())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-12345")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("期待されるステータスコード: 401, 実際: %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "Invalid token" {
		t.Errorf("期待されるエラーメッセージと異なります: %v", response["error"])
	}
}

// AuthToken ミドルウェアのテスト: 正常なトークン（実際のIDトークンが必要）
func TestAuthToken_ValidToken(t *testing.T) {
	// Firebaseの初期化
	err := middleware.InitFirebase()
	if err != nil {
		t.Skip("Firebaseの初期化に失敗したため、このテストをスキップします")
	}

	// 実際のIDトークンを環境変数から取得（テスト実行時に設定が必要）
	validToken := os.Getenv("FIREBASE_TEST_TOKEN")
	if validToken == "" {
		t.Skip("FIREBASE_TEST_TOKEN環境変数が設定されていません。有効なトークンでのテストをスキップします。")
	}

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.AuthToken())
	router.GET("/test", func(c *gin.Context) {
		uid, exists := c.Get("firebase_uid")
		if !exists {
			c.JSON(500, gin.H{"error": "firebase_uid not set"})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"uid":     uid,
		})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("期待されるステータスコード: 200, 実際: %d", w.Code)
		t.Logf("レスポンス: %s", w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["message"] != "success" {
		t.Errorf("期待されるレスポンスと異なります: %v", response)
	}
	if response["uid"] == nil {
		t.Error("UIDが設定されていません")
	}
}
