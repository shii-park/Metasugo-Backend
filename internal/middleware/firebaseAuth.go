package middleware

import (
	"context"
	"net/http"

	//"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

var (
	firebaseAuth *auth.Client
)

func InitFirebase() error {
	//export GOOGLE_APPLICATION_CREDENTIALS="/home/path/to/Metasugo-Backend/firebase-service-account.json"
	//環境変数からJSONファイルを読み込む

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return err
	}

	firebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		return err
	}

	return nil

}

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if firebaseAuth == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "認証システムが初期化されていません"})
			return
		}

		// クエリパラメータからトークンを取得しようと試みる
		idToken := c.Query("token")

		// クエリにトークンがない場合、Authorizationヘッダーを確認
		if idToken == "" {
			authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorizationヘッダまたは'token'クエリパラメータが必要です"})
				return
			}
			lower := strings.ToLower(authHeader)
			if !strings.HasPrefix(lower, "bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効な認証形式です"})
				return
			}
			idToken = strings.TrimSpace(authHeader[len("Bearer "):])
		}


		if idToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証トークンが見つかりません"})
			return
		}

		token, err := firebaseAuth.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			return
		}

		userEmail, _ := token.Claims["email"].(string)
		c.Set("firebase_uid", token.UID)
		c.Set("user_email", userEmail) //これは必要ないかも
		c.Next()
	}
}
