package main

import (
	"fmt"
	"torres.guru/gagne/dlx"
	"torres.guru/gagne/game"
)

var games = []struct {
	rows, cols int
	pieces     string
}{
	{5, 3, "otzrI"},
	{5, 3, "oi|LO"},
}

func main() {
	game.SetDebug()
	//dlx.SetDebug()

	game.LoadPieces("data/pieces.txt")
	g := game.NewBoardGame(5, 3, "otzrI")
	dl := dlx.New(g.Coverage.M.Cells(), g.Coverage.Columns)
	dl.Search(0)
	plays := g.Play(dl.Solutions[0])
	fmt.Println(plays)

	// TODO rendering
}
