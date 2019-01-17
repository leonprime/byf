package main

import (
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
	//game.SetDebug()
	//dlx.SetDebug()

	game.LoadPieces("data/pieces.txt")
	b := game.NewBoardGame(5, 3, "otzrI")
	dl := dlx.New(b.Coverage.M.Cells(), b.Coverage.Columns)
	dl.Search(0)
}
