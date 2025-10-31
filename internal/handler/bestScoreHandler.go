package handler

import (
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/shii-park/Metasugo-Backend/internal/service"
	"google.golang.org/api/iterator"
)

type BestScoreHandler struct {
	firestore *firestore.Client
}

func NewBestScoreHandler() (*BestScoreHandler, error) {
	fs, err := service.GetFirestoreClient()
	if err != nil {
		return nil, err
	}
	return &BestScoreHandler{firestore: fs}, nil
}

func (h *BestScoreHandler) GetBestScore(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("firebase_uid")

	iter := h.firestore.
		Collection("playerClearData").
		Where("playerID", "==", userID).
		Documents(ctx)
	defer iter.Stop()
	found := false
	var bestMoney int64

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("query error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ドキュメント取得に失敗しました"})
			return
		}
		v, err := doc.DataAt("money")
		if err != nil {
			continue
		}
		m, ok := v.(int64)
		if !ok {
			log.Printf("unexpected money type: %T", v)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "money の型が不正です"})
			return
		}

		if !found || m > bestMoney {
			bestMoney = m
			found = true
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "ユーザの記録が見つかりません"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"money":  bestMoney,
	})
}
