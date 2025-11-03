package handler

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/shii-park/Metasugo-Backend/internal/service"
	"google.golang.org/api/iterator"
)

// RankingHandler handles ranking related requests.
type RankingHandler struct {
	firestore *firestore.Client
}

// NewRankingHandler creates a new RankingHandler.
func NewRankingHandler() (*RankingHandler, error) {
	fs, err := service.GetFirestoreClient()
	if err != nil {
		return nil, err
	}
	return &RankingHandler{firestore: fs}, nil
}

// GetRanking retrieves the ranking from Firestore.
func (h *RankingHandler) GetRanking(c *gin.Context) {
	ctx := context.Background()
	iter := h.firestore.Collection("playerClearData").OrderBy("money", firestore.Desc).Documents(ctx)
	defer iter.Stop()

	var players []map[string]any
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ranking"})
			return
		}
		players = append(players, doc.Data())
	}

	c.JSON(http.StatusOK, players)
}
