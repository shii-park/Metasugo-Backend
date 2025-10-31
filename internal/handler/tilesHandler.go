package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func TilesHandler(c *gin.Context) {
	file, err := os.ReadFile("tiles.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read file"})
		return
	}
	c.Data(http.StatusOK, "application/json", file)
}
