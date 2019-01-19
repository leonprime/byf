package game

import (
	"bytes"
	"fmt"
	"strings"
)

// for debugging 3D
type Debug struct {
	Name    string
	Plays   []*Play3D
	W, H, D int
}

// the coverage matrix, its column names, and any coverage debugging
type Coverage struct {
	Columns []string
	M       *Grid
	Debugs  []*Debug
}

// converts a board game into a coverage matrix for solving with DLX
// this flattens a 2D game board by listing one row after another
func newBoardCoverage(b *Board) *Coverage {
	var (
		rows  [][]bool
		names []string
	)
	n := len(b.pieces)
	for i, piece := range b.pieces {
		grids := piece.Positions(b.W, b.H)
		for _, grid := range grids {
			row := make([]bool, n, n)
			row[i] = true // set piece at index i to 1
			for y := 0; y < b.H; y++ {
				for x := 0; x < b.W; x++ {
					row = append(row, grid.Get(x, y))
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
	cov := &Coverage{
		M:       &Grid{Cells: rows, W: len(rows[0]), H: len(rows)},
		Columns: names,
	}
	if debug.coverage() {
		fmt.Println(cov)
	}
	return cov
}

// returns all uniquely oriented positions of the piece
// on a 2D game board reprsented by (w, h),
func (p *Piece) Positions(w, h int) []*Grid {
	var grids []*Grid
	for _, shape := range p.Shapes {
		grids = append(grids, perms(w, h, shape)...)
		for i := 1; i < p.Rotate; i++ {
			shape = shape.Rotate()
			grids = append(grids, perms(w, h, shape)...)
		}
	}
	if debug.piece(p) {
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
		for y := 0; y < grids[0].H; y++ {
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
// this flattens the 3D cube by listing the 2D planes XY as in the board game
// from front to back in Z
func newCubeCoverage(c *Cube) *Coverage {
	var (
		rows   [][]bool
		names  []string
		debugs []*Debug
	)
	n := len(c.pieces)
	for i, piece := range c.pieces {
		grids := piece.Positions3D(c.W, c.H, c.D)
		for _, grid := range grids {
			row := make([]bool, n, n)
			row[i] = true // set piece at index i to 1
			for y := 0; y < c.H; y++ {
				for x := 0; x < c.W; x++ {
					for z := 0; z < c.D; z++ {
						row = append(row, grid.Get(x, y, z))
					}
				}
			}
			rows = append(rows, row)
		}
		// also set the name
		names = append(names, piece.Name)
		if debug.piece(piece) {
			var plays []*Play3D
			for _, grid := range grids {
				plays = append(plays, &Play3D{Piece: piece, Grid: grid})
			}
			debugs = append(debugs, &Debug{Name: fmt.Sprintf("positions_%s", piece.Name), Plays: plays, W: c.W, H: c.H, D: c.D})
		}
	}
	// rest of the columns should be named sequentially z*w*h + y*h + x
	for y := 0; y < c.H; y++ {
		for x := 0; x < c.W; x++ {
			for z := 0; z < c.D; z++ {
				names = append(names, fmt.Sprintf("c%d", z*c.W*c.H+y*c.H+x))
			}
		}
	}
	cov := &Coverage{
		M:       &Grid{Cells: rows, W: len(rows[0]), H: len(rows)},
		Columns: names,
		Debugs:  debugs,
	}
	if debug.coverage() {
		fmt.Println(cov)
	}
	return cov
}

// returns all uniquely oriented positions of the piece
// on a 3d cube reprsented by (w, h, d),
func (p *Piece) Positions3D(w, h, d int) []*Grid3D {
	//
	// the key to this is to use the 2D position permutations and "project" them
	// down each dimension.  this results in duplicates, so must clean that up too
	var grids []*Grid3D
	grids2d := p.Positions(w, h)
	for _, grid2d := range grids2d {
		for z := 0; z < d; z++ {
			grid := newEmptyGrid3D(w, h, d)
			grid.SetPlaneZ(z, grid2d)
			found := false
			for _, test := range grids {
				if grid.Equals(test) {
					found = true
					break
				}
			}
			if !found {
				grids = append(grids, grid)
			}
		}
		for y := 0; y < h; y++ {
			grid := newEmptyGrid3D(w, h, d)
			grid.SetPlaneY(y, grid2d)
			found := false
			for _, test := range grids {
				if grid.Equals(test) {
					found = true
					break
				}
			}
			if !found {
				grids = append(grids, grid)
			}
		}
		for x := 0; x < w; x++ {
			grid := newEmptyGrid3D(w, h, d)
			grid.SetPlaneX(x, grid2d)
			found := false
			for _, test := range grids {
				if grid.Equals(test) {
					found = true
					break
				}
			}
			if !found {
				grids = append(grids, grid)
			}
		}
	}
	if debug.piece(p) {
		fmt.Printf("generated %d 3D positions for %s\n", len(grids), p.Name)
	}
	return grids
}

func (c *Coverage) String() string {
	var b bytes.Buffer
	b.WriteString("coverage matrix A:\n")
	for _, col := range c.Columns {
		b.WriteString(col)
		b.WriteString(" ")
	}
	b.WriteRune('\n')
	b.WriteString(c.M.String())
	return b.String()
}

func (c *Coverage) RowString(y int) string {
	var b bytes.Buffer
	for i, v := range c.M.Row(y) {
		if i == len(c.Columns) {
			b.WriteString(" : ")
		}
		if v {
			b.WriteRune('█')
		} else {
			b.WriteRune('.')
		}
	}
	return b.String()
}
