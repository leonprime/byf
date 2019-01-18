package game

import (
	"bytes"
	"fmt"
	"unicode"
)

type Grid struct {
	Cells [][]bool
	W, H  int
}

// build a grid from a grid spec, which is a # or █ for true and a . for false
// put each row on a separate line
func newGrid(spec string) *Grid {
	cells := make([][]bool, 0, 0)
	row := make([]bool, 0, 0)
	for _, char := range spec {
		if unicode.IsSpace(char) {
			if char == '\n' && len(row) > 0 {
				cells = append(cells, row)
				row = make([]bool, 0, 0)
			}
			continue
		}
		switch char {
		case '#', '█':
			row = append(row, true)
		case '.':
			row = append(row, false)
		}
	}
	if len(row) > 0 {
		cells = append(cells, row)
	}
	w := 0
	for r := range cells {
		if w > 0 && w != len(cells[r]) {
			panic(fmt.Sprintf("grid spec is not rectangular:\n%s", spec))
		}
		w = len(cells[r])
	}
	return &Grid{Cells: cells, W: w, H: len(cells)}
}

// returns grid set to false
func newEmptyGrid(w, h int) *Grid {
	cells := make([][]bool, h, h)
	for y := range cells {
		cells[y] = make([]bool, w, w)
	}
	return &Grid{Cells: cells, H: h, W: w}
}

func (g *Grid) IsOOB(x, y int) bool {
	if x < 0 || y < 0 || x >= g.W || y >= g.H {
		return true
	}
	return false
}

// test if value is set.  doesn't panic on oob
func (g *Grid) IsSet(x, y int) bool {
	return !g.IsOOB(x, y) && g.Get(x, y)
}

// get value.  panics if oob
func (g *Grid) Get(x, y int) bool {
	if g.IsOOB(x, y) {
		panic(fmt.Sprintf("Grid.Get(%d, %d) is oob: Grid(w=%d, h=%d)", x, y, g.W, g.H))
	}
	return g.Cells[y][x]
}

func (g *Grid) Set(x, y int, b bool) {
	if g.IsOOB(x, y) {
		panic(fmt.Sprintf("Grid.Set(%d, %d) is oob: Grid(w=%d, h=%d)", x, y, g.W, g.H))
	}
	g.Cells[y][x] = b
}

func (g *Grid) IsEmpty() bool {
	for y := range g.Cells {
		for x := range g.Cells[y] {
			if g.Cells[y][x] {
				return false
			}
		}
	}
	return true
}

// Sets the values to the given subgrid values if and only if the subgrid
// is entirely contained at positions (x, y) to (x+w, y+h)
// If the subgrid is out of bounds, nothing is set.
func (g *Grid) SetSubgrid(x, y int, grid *Grid) {
	if x+grid.W > g.W || y+grid.H > g.H {
		return
	}
	for yy, j := 0, y; yy < grid.H; yy++ {
		for xx, i := 0, x; xx < grid.W; xx++ {
			g.Set(i, j, grid.Get(xx, yy))
			i++
		}
		j++
	}
}

func (g *Grid) GetSubgrid(x, y, w, h int) *Grid {
	if x+w > g.W || y+h > g.H {
		panic(fmt.Sprintf("subgrid is oob: (%d, %d) w=%d, h=%d on grid w=%d, h=%d", x, y, w, h, g.W, g.H))
	}
	grid := newEmptyGrid(w, h)
	for yy, j := 0, y; yy < h; yy++ {
		for xx, i := 0, x; xx < w; xx++ {
			grid.Set(xx, yy, g.Get(i, j))
			i++
		}
		j++
	}
	return grid
}

// returns a grid that's rotated 90 degrees clockwise
func (g *Grid) Rotate() *Grid {
	var cells [][]bool
	for x := 0; x < g.W; x++ {
		// x is the new y
		row := make([]bool, g.H, g.H)
		for y := 0; y < g.H; y++ {
			// y is x from right to left
			row[g.H-y-1] = g.Cells[y][x]
		}
		cells = append(cells, row)
	}
	return &Grid{
		Cells: cells,
		W:     g.H,
		H:     g.W,
	}
}

func (g *Grid) Row(y int) []bool {
	row := make([]bool, g.W, g.W)
	for x := 0; x < g.W; x++ {
		row[x] = g.Get(x, y)
	}
	return row
}

func (g *Grid) IsRowEmpty(y int) bool {
	for x := 0; x < g.W; x++ {
		if g.Get(x, y) {
			return false
		}
	}
	return true
}

func (g *Grid) IsColEmpty(x int) bool {
	for y := 0; y < g.H; y++ {
		if g.Get(x, y) {
			return false
		}
	}
	return true
}

func (g *Grid) String() string {
	var s bytes.Buffer
	for y := range g.Cells {
		for x := range g.Cells[y] {
			r := '.'
			if g.Cells[y][x] {
				r = '█'
			}
			s.WriteRune(r)
		}
		s.WriteRune('\n')
	}
	return s.String()
}
