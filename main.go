package main

import (
	"fmt"
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
	game.LoadPieces("data/pieces.txt")
	g := game.New2D(5, 3, "otzrI")
	m := game.Game2DCoverageMatrix(g)
	fmt.Println(m)
}
