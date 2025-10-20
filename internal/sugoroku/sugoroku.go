package sugoroku

import (
	"fmt"
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

func (g *Game) AddPlayer(id string) (*Player, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.players[id]; exists {
		return nil, fmt.Errorf("player with id %s already exists", id)
	}

	player := NewPlayer(id, g.tileMap[InitialTileID])
	g.players[id] = player

	return player, nil
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
	if p.position.prev != nil {
		targetTiles = append(targetTiles, p.position.prev)
	}
	if p.position.next != nil {
		targetTiles = append(targetTiles, p.position.next)
	}
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

func DistributeMoney(players []*Player, amount int) int {
	playerNum := len(players)
	amountPerPlayers := amount / playerNum
	return amountPerPlayers
}
