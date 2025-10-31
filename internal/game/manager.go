package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/service"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
	log "github.com/sirupsen/logrus"
)

type GameManager struct {
	game          *sugoroku.Game
	hub           *hub.Hub
	playerClients map[string]*hub.Client
	firestore     *firestore.Client
	mu            sync.RWMutex
}

func NewGameManager(g *sugoroku.Game, h *hub.Hub) *GameManager {
	fs, err := service.GetFirestoreClient()
	if err != nil {
		log.WithError(err).Fatal("failed to get firestore client")
	}
	return &GameManager{
		game:          g,
		hub:           h,
		playerClients: make(map[string]*hub.Client),
		firestore:     fs,
	}
}
func (gm *GameManager) MoveByDiceRoll(playerID string, steps int) error {
	player, err := gm.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("failed to get player: %w", err)
	}

	// 1. 移動前の状態を記録
	initialPosition := player.GetPosition().GetID()
	initialMoney := player.GetMoney()
	initialIsMarried := player.GetIsMarried()
	initialChildren := player.GetChildren()
	initialJob := player.GetJob()

	// 2. プレイヤーを移動させる
	flag := player.Move(steps) //めんどくさくなったのでフラグで実装してる。Effect型で比較するなどもっといいやり方はあると思う

	// 効果を判定
	log.WithFields(log.Fields{
		"playerID":    playerID,
		"newPosition": player.GetPosition().GetID(),
	}).Info("Player moved")
	// 3. マス効果を判定・適用

	finalPosition := player.GetPosition().GetID()

	if initialPosition != finalPosition {
		gm.broadcastPlayerMoved(playerID, finalPosition)
	}

	currentTile := player.GetPosition()
	effect := currentTile.GetEffect()

	if effect.RequiresUserInput() || flag == "GOAL" {
		// ユーザからの入力が必要な場合、もしくはゴールの場合こちらで処理
		switch e := effect.(type) {
		case sugoroku.BranchEffect:
			return gm.sendBranchSelection(player, currentTile, e)
		case sugoroku.QuizEffect:
			return gm.sendQuizInfo(player, currentTile, e)
		case sugoroku.GambleEffect:
			return gm.sendGambleRequire(player, currentTile)
		case sugoroku.GoalEffect:
			if err := gm.Goal(playerID, gm.playerClients[playerID]); err != nil { // TODO: ゴールした際に行う処理(clientとの接続解除など)を行ったほうが良いと思う
				return err
			}
			return nil //ゲーム終了するのでここで関数を脱出！
		default:
			return fmt.Errorf("unhandled user input required for effect type %T", e)
		}
	} else {
		// 即時効果を適用
		if err := effect.Apply(player, gm.game, nil); err != nil {
			return err
		}
	}

	// 4. ステータスの変更を検知して通知
	finalMoney := player.GetMoney()
	if initialMoney != finalMoney {
		gm.broadcastMoneyChanged(playerID, finalMoney)
	}

	finalIsMarried := player.GetIsMarried()
	if initialIsMarried != finalIsMarried {
		gm.broadcastPlayerStatusChanged(playerID, "isMarried", finalIsMarried)
	}

	finalChildren := player.GetChildren()
	if initialChildren != finalChildren {
		gm.broadcastPlayerStatusChanged(playerID, "children", finalChildren)
	}

	finalJob := player.GetJob()
	if initialJob != finalJob {
		gm.broadcastPlayerStatusChanged(playerID, "job", finalJob)
	}

	return nil
}

func (gm *GameManager) RegisterPlayerClient(playerID string, c *hub.Client) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	_, err := gm.game.AddPlayer(playerID)
	if err != nil {
		return err
	}
	gm.playerClients[playerID] = c
	return nil
}

func (gm *GameManager) UnregisterPlayerClient(playerID string, c *hub.Client) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	
	return gm.unregisterPlayerClientLocked(playerID, c)
}

// unregisterPlayerClientLocked は、既にロックが取得されている状態でプレイヤーを登録解除します
func (gm *GameManager) unregisterPlayerClientLocked(playerID string, c *hub.Client) error {
	log.WithField("playerID", playerID).Info("UnregisterPlayerClient: Starting")
	
	// ゲームからプレイヤーを削除
	log.WithField("playerID", playerID).Info("UnregisterPlayerClient: Deleting player from game")
	err := gm.game.DeletePlayer(playerID)
	if err != nil {
		log.WithError(err).WithField("playerID", playerID).Error("UnregisterPlayerClient: Failed to delete player from game")
		return err
	}
	log.WithField("playerID", playerID).Info("UnregisterPlayerClient: Player deleted from game")
	
	// GameManagerからプレイヤーを削除
	delete(gm.playerClients, playerID)
	log.WithField("playerID", playerID).Info("UnregisterPlayerClient: Player deleted from playerClients map")

	// Hubにクライアントの登録解除を通知
	// これにより、Hubはクライアントの接続を閉じ、リソースを解放します
	log.WithField("playerID", playerID).Info("UnregisterPlayerClient: Calling Hub.Unregister")
	c.Hub.Unregister(c)
	log.WithField("playerID", playerID).Info("UnregisterPlayerClient: Hub.Unregister completed")

	return nil
}

func (gm *GameManager) Goal(playerID string, c *hub.Client) error {
	log.WithField("playerID", playerID).Info("Goal function called")
	
	player, err := gm.game.GetPlayer(playerID)
	if err != nil {
		log.WithError(err).Error("failed to get player in Goal")
		return fmt.Errorf("failed to get player: %w", err)
	}
	money := player.GetMoney()
	log.WithFields(log.Fields{"playerID": playerID, "money": money}).Info("Player retrieved successfully")

	fmt.Printf("%s goal!!!!!!!", playerID)

	// Firestoreに保存するデータを作成
	data := map[string]interface{}{
		"playerID":   playerID,
		"money":      money,
		"finishedAt": time.Now(),
	}

	// Firestoreにデータを保存
	log.Info("Starting Firestore save")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = gm.firestore.Collection("playerClearData").Doc(playerID).Set(ctx, data)
	if err != nil {
		log.WithError(err).Error("failed to save to firestore")
		return fmt.Errorf("failed to save player data to firestore: %w", err)
	}
	log.Info("Firestore save completed")

	log.Info("Broadcasting player finished")
	gm.broadcastPlayerFinished(playerID, money)
	log.Info("Broadcast completed")

	log.Info("Calling UnregisterPlayerClient asynchronously")
	// 非同期で登録解除を行うことで、Broadcast処理との競合を回避
	go func() {
		log.Info("UnregisterPlayerClient goroutine started")
		if err := gm.UnregisterPlayerClient(playerID, c); err != nil {
			log.WithError(err).Error("UnregisterPlayerClient failed")
		} else {
			log.Info("UnregisterPlayerClient completed successfully")
		}
	}()
	
	return nil
}
