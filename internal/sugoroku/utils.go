package sugoroku

import "math/rand"

func DistributeMoney(players []*Player, amount int) int {
	playerNum := len(players)
	amountPerPlayers := amount / playerNum
	return amountPerPlayers
}

// 1~6までのランダムな数を返す関数
func RollDice() int {
	return rand.Intn(6) + 1
}
