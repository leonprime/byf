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
	//game.SetDebug()
	//dlx.SetDebug()

	game.LoadPieces("data/pieces.txt")
	g := game.New2D(5, 3, "otzrI")
	m, names := game.Game2DCoverageMatrix(g)
	dl := dlx.New(m.Cells(), names)
	dl.Search(0)

	// TODO: now transform solution back into a visual representation
	for _, soln := range dl.Solutions {
		fmt.Println(names)
		for _, y := range soln {
			fmt.Printf("row %d\n", y)
			fmt.Println(m.Row(y))
		}
		break
	}
}
