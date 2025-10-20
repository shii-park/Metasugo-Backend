package main

import (
	//"github.com/gin-gonic/gin"

	//"github.com/Metasugo-Backend/internal/handler"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	//router := gin.Default()

	//router.POST("/signup", handler.SignUp)

	sugoroku.NewGame() // ハンドラができた際に、gameにAddplayerができるようになる
}
