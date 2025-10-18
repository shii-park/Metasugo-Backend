package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Metasugo-Backend/internal/handler"
)

func main() {
	router := gin.Default()

	router.POST("/signup", handler.SignUp)
}
