package middleware

import (
	"context"
	//"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

var firebaseAuth *auth.Client

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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("firebase_uid", token.UID)
		c.Set("user_email", token.Claims["email"]) //これは必要ないかも
		c.Next()
	}
}
