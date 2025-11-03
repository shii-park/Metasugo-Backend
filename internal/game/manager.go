package game

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
	"github.com/shii-park/Metasugo-Backend/internal/service"
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
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
		log.Fatalf("failed to get firestore client: %v", err)
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
	initialPosition := player.Position.Id
	initialMoney := player.Money
	initialIsMarried := player.IsMarried
	initialHasChildren := player.HasChildren
	initialJob := player.Job

	// 2. プレイヤーを移動させる
	flag := player.Move(steps) //めんどくさくなったのでフラグで実装してる。Effect型で比較するなどもっといいやり方はあると思う

	// 効果を判定
	log.Printf("PlayerMoved: %s moved to %d", playerID, player.Position.Id)
	// 3. マス効果を判定・適用

	finalPosition := player.Position.Id

	if initialPosition != finalPosition {
		gm.broadcastPlayerMoved(playerID, finalPosition)
	}

	currentTile := player.Position
	effect := currentTile.Effect

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
	finalMoney := player.Money
	if initialMoney != finalMoney {
		gm.broadcastMoneyChanged(playerID, finalMoney)
	}

	finalIsMarried := player.IsMarried
	if initialIsMarried != finalIsMarried {
		gm.broadcastPlayerStatusChanged(playerID, "isMarried", finalIsMarried)
	}

	finalHasChildren := player.HasChildren
	if initialHasChildren != finalHasChildren {
		gm.broadcastPlayerStatusChanged(playerID, "hasChildren", finalHasChildren)
	}

	finalJob := player.Job
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
	// ゲームからプレイヤーを削除
	err := gm.game.DeletePlayer(playerID)
	if err != nil {
		return err
	}
	// GameManagerからプレイヤーを削除
	delete(gm.playerClients, playerID)

	// Hubにクライアントの登録解除を通知
	// これにより、Hubはクライアントの接続を閉じ、リソースを解放します
	c.Hub.Unregister(c)

	return nil

}

func (gm *GameManager) Goal(playerID string, c *hub.Client) error {
	player, err := gm.game.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("failed to get player: %w", err)
	}
	money := player.Money

	// Firestoreに保存するデータを作成
	data := map[string]interface{}{
		"playerID":   playerID,
		"money":      money,
		"finishedAt": time.Now(),
	}

	// Firestoreにデータを保存
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = gm.firestore.Collection("playerClearData").Doc(playerID).Set(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to save player data to firestore: %w", err)
	}

	gm.broadcastPlayerFinished(playerID, money)

	if err := gm.UnregisterPlayerClient(playerID, c); err != nil {
		return err
	}
	return nil
}
