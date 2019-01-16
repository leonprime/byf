package game

import (
	"bytes"
	"fmt"
	"strings"
)

// converts a game into a coverage matrix for solving with DLX
// this flattens the 2d game board by listing one row after another
func Game2DCoverageMatrix(g *Game2D) *Grid {
	var rows [][]bool
	n := len(g.Pieces)
	for i, piece := range g.Pieces {
		grids := piece.Positions(g.w, g.h)
		for _, grid := range grids {
			row := make([]bool, n+g.w*g.h, n+g.w*g.h)
			row[i] = true // set piece at index i to 1
			for y := 0; y < g.h; y++ {
				for x := 0; x < g.w; x++ {
					row[n+y*g.w+x] = grid.Get(x, y)
				}
			}
			rows = append(rows, row)
		}
	}
	return &Grid{cells: rows, w: g.w, h: g.h}
}

// given the dimensions of a 2d game board, returns all uniquely
// oriented positions of a piece on the game board
func (p *Piece) Positions(w, h int) []*Grid {
	var grids []*Grid
	for _, shape := range p.Shapes {
		grids = append(grids, perms(w, h, shape)...)
		for i := 1; i < p.Rotate; i++ {
			shape = shape.Rotate()
			grids = append(grids, perms(w, h, shape)...)
		}
	}
	if debug {
		fmt.Printf("generated %d positions for %s\n", len(grids), p.Name)
		printgrids(grids)
	}
	return grids
}

// does the actual permutation work for a given shape
func perms(w, h int, shape *Grid) []*Grid {
	var grids []*Grid
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// make a grid with piece at this position
			grid := newEmptyGrid(w, h)
			// SetSubgrid handles oob conditions with a noop
			grid.SetSubgrid(x, y, shape)
			if !grid.IsEmpty() {
				grids = append(grids, grid)
			}
		}
	}
	return grids
}

func printgrids(grids []*Grid) {
	var str [][]string
	for _, grid := range grids {
		str = append(str, strings.Split(grid.String(), "\n"))
	}
	var b bytes.Buffer
	for i := 0; i < len(str); i += 5 {
		for y := 0; y < grids[0].h; y++ {
			for j := i; j < i+5 && j < len(str); j++ {
				b.WriteString(str[j][y])
				b.WriteString("    ")
			}
			b.WriteRune('\n')
		}
		b.WriteRune('\n')
	}
	fmt.Println(b.String())
}
