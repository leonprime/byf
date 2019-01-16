package main

import (
	"fmt"
	"torres.guru/gagne/game"
)

var games = []struct {
	rows, cols int
	pieces     string
}{
	{3, 5, "otzrI"},
	{3, 5, "oi|LO"},
}

func main() {
	game.LoadPieces("data/pieces.txt")
	g := game.New2D(3, 5, "otzrI")
	fmt.Println(g)
}
