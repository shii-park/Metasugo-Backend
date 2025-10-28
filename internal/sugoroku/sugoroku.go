package sugoroku

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

const InitialTileID = 1

type Game struct {
	players map[string]*Player
	tileMap map[int]*Tile

	mu sync.RWMutex
}

//                                            __                                      __
//                                           |  \                                    |  \
//   _______   ______   _______    _______  _| $$_     ______   __    __   _______  _| $$_     ______    ______
//  /       \ /      \ |       \  /       \|   $$ \   /      \ |  \  |  \ /       \|   $$ \   /      \  /      \
// |  $$$$$$$|  $$$$$$\| $$$$$$$\|  $$$$$$$ \$$$$$$  |  $$$$$$\| $$  | $$|  $$$$$$$ \$$$$$$  |  $$$$$$\|  $$$$$$\
// | $$      | $$  | $$| $$  | $$ \$$    \   | $$ __ | $$   \$$| $$  | $$| $$        | $$ __ | $$  | $$| $$   \$$
// | $$_____ | $$__/ $$| $$  | $$ _\$$$$$$\  | $$|  \| $$      | $$__/ $$| $$_____   | $$|  \| $$__/ $$| $$
//  \$$     \ \$$    $$| $$  | $$|       $$   \$$  $$| $$       \$$    $$ \$$     \   \$$  $$ \$$    $$| $$
//   \$$$$$$$  \$$$$$$  \$$   \$$ \$$$$$$$     \$$$$  \$$        \$$$$$$   \$$$$$$$    \$$$$   \$$$$$$  \$$
//

func NewGame() *Game {
	tileMap := InitTiles()
	InitQuiz()
	return &Game{
		tileMap: tileMap,
		players: make(map[string]*Player),
	}
}

// テスト用のラッパー関数
func NewGameWithTilesForTest(path string) *Game {
	tileMap, err := InitTilesFromPath(path)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize tiles: %v", err))
	}
	InitQuiz()
	return &Game{
		tileMap: tileMap,
		players: make(map[string]*Player),
	}
}

//                           __      __                        __
//                          |  \    |  \                      |  \
//  ______ ____    ______  _| $$_   | $$____    ______    ____| $$  _______
// |      \    \  /      \|   $$ \  | $$    \  /      \  /      $$ /       \
// | $$$$$$\$$$$\|  $$$$$$\\$$$$$$  | $$$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$$
// | $$ | $$ | $$| $$    $$ | $$ __ | $$  | $$| $$  | $$| $$  | $$ \$$    \
// | $$ | $$ | $$| $$$$$$$$ | $$|  \| $$  | $$| $$__/ $$| $$__| $$ _\$$$$$$\
// | $$ | $$ | $$ \$$     \  \$$  $$| $$  | $$ \$$    $$ \$$    $$|       $$
//  \$$  \$$  \$$  \$$$$$$$   \$$$$  \$$   \$$  \$$$$$$   \$$$$$$$ \$$$$$$$
//

func (g *Game) AddPlayer(playerID string) (*Player, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.players[playerID]; exists {
		return nil, fmt.Errorf("player with id %s already exists", playerID)
	}

	player := NewPlayer(playerID, g.tileMap[InitialTileID])
	g.players[playerID] = player

	return player, nil
}

func (g *Game) DeletePlayer(playerID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.players[playerID]; exists {
		delete(g.players, playerID)
		log.Printf("DeletePlayer: %s has deleted", playerID)
		return nil
	}
	return fmt.Errorf("player with id %s does not exist", playerID)
}

func (g *Game) GetAllPlayers() []*Player {
	g.mu.RLock()
	defer g.mu.RUnlock()

	playerList := make([]*Player, 0, len(g.players))

	for _, player := range g.players {
		playerList = append(playerList, player)
	}
	return playerList
}

func (g *Game) GetNeighbors(p *Player) []*Player { // 計算量がプレイヤー数になってしまうのでリファクタリングできる(Tileにプレイヤー情報をもたせるなど)
	targetTiles := []*Tile{}
	targetPlayers := []*Player{}

	targetTiles = append(targetTiles, p.position.prevs...)
	targetTiles = append(targetTiles, p.position.nexts...)

	if p.position != nil {
		targetTiles = append(targetTiles, p.position)
	}

	if len(targetTiles) == 0 {
		return nil
	}

	for _, Player := range g.GetAllPlayers() {
		if p.id == Player.id {
			continue
		}
		for _, tile := range targetTiles {
			if Player.position == tile {
				targetPlayers = append(targetPlayers, Player)
				break
			}
		}
	}
	return targetPlayers
}

func (g *Game) GetPlayer(playerID string) (*Player, error) {
	player, exist := g.players[playerID]
	if !exist {
		return nil, errors.New("the player does not exist")
	}
	return player, nil
}
