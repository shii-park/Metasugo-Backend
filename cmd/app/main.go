package main

import (
	"github.com/shii-park/Metasugo-Backend/internal/sugoroku"
)

func main() {
	sugoroku.NewGame() // ハンドラができた際に、gameにAddplayerができるようになる
}
