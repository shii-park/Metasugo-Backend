package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
)

var (
	loadOnce  sync.Once
	loadErr   error
	tilesData map[string]any

	firestoreClient *firestore.Client
	firestoreOnce   sync.Once
	firestoreErr    error
)

// NewFirestoreClient initializes the Firestore client
func NewFirestoreClient() (*firestore.Client, error) {
	firestoreOnce.Do(func() {
		ctx := context.Background()
		projectID := os.Getenv("FIREBASE_PROJECT_ID")

		client, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			log.Printf("error initializing firestore client with database ID: %v\n", err)
			firestoreErr = err
			return
		}
		firestoreClient = client
	})

	if firestoreErr != nil {
		return nil, firestoreErr
	}
	if firestoreClient == nil {
		return nil, errors.New("firestore client is not initialized")
	}

	return firestoreClient, nil
}

func GetFirestoreClient() (*firestore.Client, error) {
	return NewFirestoreClient()
}

func GetTiles() (any, error) {
	loadOnce.Do(func() {
		filePath := os.Getenv("TILES_JSON_PATH")
		if filePath == "" {
			loadErr = errors.New("TILES_JSON_PATH environment variable not set")
			return
		}
		file, err := os.Open(filePath)
		if err != nil {
			loadErr = err
			return
		}
		defer file.Close()
		var m map[string]any

		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&m); err != nil {
			loadErr = err
			return
		}
		tilesData = m
	})
	if loadErr != nil {
		return nil, loadErr
	}
	if tilesData == nil {
		return nil, errors.New("盤面データを読み込めませんでした")
	}
	return tilesData, nil
}
