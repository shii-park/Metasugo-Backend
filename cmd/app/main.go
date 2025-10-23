package main

import (
	//"github.com/gin-gonic/gin"

	//"github.com/Metasugo-Backend/internal/handler"

	"github.com/shii-park/Metasugo-Backend/internal/game"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	// router := gin.Default()

	//router.POST("/signup", handler.SignUp)

	g := sugoroku.NewGame() // ハンドラができた際に、gameにAddplayerができるようになる
	h := hub.NewHub()
	game.NewGameManager(g, h)

}
