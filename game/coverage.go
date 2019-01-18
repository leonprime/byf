package game

import (
	"bytes"
	"fmt"
	"strings"
)

type Coverage struct {
	Columns []string
	M       *Grid
}

// converts a board game into a coverage matrix for solving with DLX
// this flattens the 2d game board by listing one row after another
func newBoardCoverage(b *Board) *Coverage {
	var (
		rows  [][]bool
		names []string
	)
	n := len(b.pieces)
	for i, piece := range b.pieces {
		grids := piece.Positions(b.W, b.H)
		for _, grid := range grids {
			row := make([]bool, n+b.W*b.H, n+b.W*b.H)
			row[i] = true // set piece at index i to 1
			for y := 0; y < b.H; y++ {
				for x := 0; x < b.W; x++ {
					row[n+y*b.W+x] = grid.Get(x, y)
				}
			}
			rows = append(rows, row)
		}
		// also set the name
		names = append(names, piece.Name)
	}
	// rest of the columns should be named sequentially y*h + x
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			names = append(names, fmt.Sprintf("c%d", y*b.H+x))
		}
	}
	return &Coverage{
		M:       &Grid{cells: rows, w: len(rows[0]), h: len(rows)},
		Columns: names,
	}
}

// returns all uniquely oriented positions of the piece
// on a 2d game board reprsented by (w, h),
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

// converts a cube game into a coverage matrix for solving with DLX
// this flattens the 3d cube by listing each 2d slice similar to the board
// coverage, then working front to back in depth
func newCubeCoverage(c *Cube) *Coverage {
	var (
		rows  [][]bool
		names []string
	)
	n := len(c.pieces)
	for i, piece := range c.pieces {
		grids := piece.Positions3D(c.W, c.H, c.D)
		for _, grid := range grids {
			row := make([]bool, n+c.W*c.H*c.D, n+c.W*c.H*c.D)
			row[i] = true // set piece at index i to 1
			for y := 0; y < c.H; y++ {
				for x := 0; x < c.W; x++ {
					for z := 0; z < c.D; z++ {
						row[n+x*c.W*c.H+y*c.W+x] = grid.Get(x, y, z)
					}
				}
			}
			rows = append(rows, row)
		}
		// also set the name
		names = append(names, piece.Name)
	}
	// rest of the columns should be named sequentially z*w*h + y*h + x
	for y := 0; y < c.H; y++ {
		for x := 0; x < c.W; x++ {
			for z := 0; z < c.D; z++ {
				names = append(names, fmt.Sprintf("c%d", z*c.W*c.H+y*c.H+x))
			}
		}
	}
	return &Coverage{
		M:       &Grid{cells: rows, w: len(rows[0]), h: len(rows)},
		Columns: names,
	}
}

// returns all uniquely oriented positions of the piece
// on a 3d cube reprsented by (w, h, d),
func (p *Piece) Positions3D(w, h, d int) []*Grid3D {
	//
	// the key to this is to use the 2d position permutations and "project" them
	// down each dimension
	var grids []*Grid3D
	grids2d := p.Positions(w, h)
	for _, grid2d := range grids2d {
		for x := 0; x < w; x++ {
			grid := newEmptyGrid3D(w, h, d)
			grid.SetPlaneYZ(x, grid2d)
			grids = append(grids, grid)
		}
		for y := 0; y < h; y++ {
			grid := newEmptyGrid3D(w, h, d)
			grid.SetPlaneXZ(y, grid2d)
			grids = append(grids, grid)
		}
		for z := 0; z < d; z++ {
			grid := newEmptyGrid3D(w, h, d)
			grid.SetPlaneXY(z, grid2d)
			grids = append(grids, grid)
		}
	}
	return grids
}
