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
