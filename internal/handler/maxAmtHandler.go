package handler

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/shii-park/Metasugo-Backend/internal/service"
)

type MaxAmtHandler struct {
	firestore *firestore.Client
}

func NewMaxAmtHandler() (*MaxAmtHandler, error) {
	fs, err := service.GetFirestoreClient()
	if err != nil {
		return nil, err
	}
	return &MaxAmtHandler{firestore: fs}, nil
}

func (h *MaxAmtHandler) GetMaxAmt(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("firebase_uid")

	docRef := h.firestore.Collection("playerClearData").Doc(userID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ドキュメント取得に失敗しました"})
		return
	}
	money, err := docSnap.DataAt("money")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "money の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"money": money})
}
