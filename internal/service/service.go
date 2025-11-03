package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt" // ★ インポート追加
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var (
	loadOnce  sync.Once
	loadErr   error
	tilesData map[string]any

	// --- Firebase 関連の変数を整理 ---
	firebaseApp *firebase.App
	appOnce     sync.Once
	appErr      error

	// --- Firebase 関連の変数を整理 ---
	firebaseApp *firebase.App
	appOnce     sync.Once
	appErr      error

	firestoreClient *firestore.Client
	firestoreOnce   sync.Once
	firestoreErr    error

	authClient *auth.Client // ★ Authクライアント用
	authOnce   sync.Once
	authErr    error
	// ---
)

func getFirebaseApp() (*firebase.App, error) {
	appOnce.Do(func() {
		ctx := context.Background()

		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		if projectID == "" {
			appErr = errors.New("FIREBASE_PROJECT_ID environment variable not set")
			return
		}

		conf := &firebase.Config{ProjectID: projectID}

		app, err := firebase.NewApp(ctx, conf)
		if err != nil {
			log.Printf("error initializing firebase app: %v\n", err)
			appErr = err
			return
		}
		firebaseApp = app
	})
	return firebaseApp, appErr
}

func NewFirestoreClient() (*firestore.Client, error) {
	firestoreOnce.Do(func() {
		//
		app, err := getFirebaseApp()
		if err != nil {
			firestoreErr = fmt.Errorf("failed to get firebase app: %w", err)
			return
		}

		// 2. App から Firestore クライアントを取得
		ctx := context.Background()
		client, err := app.Firestore(ctx)
		if err != nil {
			log.Printf("error initializing firestore client from app: %v\n", err)
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
	// この関数は変更なし
	return NewFirestoreClient()
}

// ★★★ ここから追加 ★★★

// GetAuthClient は App から Auth クライアントを取得します
func GetAuthClient() (*auth.Client, error) {
	authOnce.Do(func() {
		// 1. 共通の App を取得
		app, err := getFirebaseApp()
		if err != nil {
			authErr = fmt.Errorf("failed to get firebase app: %w", err)
			return
		}

		// 2. App から Auth クライアントを取得
		ctx := context.Background()
		client, err := app.Auth(ctx)
		if err != nil {
			log.Printf("error initializing auth client from app: %v\n", err)
			authErr = err
			return
		}
		authClient = client
	})

	if authErr != nil {
		return nil, authErr
	}
	if authClient == nil {
		return nil, errors.New("auth client is not initialized")
	}

	return authClient, nil
}

// ★★★ ここまで追加 ★★★

func GetTiles() (interface{}, error) {
	// この関数は変更なし
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
