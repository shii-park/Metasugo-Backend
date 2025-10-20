package service

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId string) (string, error) {
	//JWT生成処理
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte("KEY"))
	if err != nil {
		log.Printf("署名に失敗しました: %v", err)
		return "", err
	}

	log.Println(signed)
	return signed, nil
}
