package service

/*import (
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

	accessToken, err := token.SignedString([]byte("KEY"))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}*/
