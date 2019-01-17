package game

import (
	"bytes"
	"fmt"
	"unicode"
)

type Grid struct {
	cells [][]bool
	w, h  int
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
	return &Grid{cells: cells, w: w, h: len(cells)}
}

// returns grid set to false
func newEmptyGrid(w, h int) *Grid {
	cells := make([][]bool, h, h)
	for y := range cells {
		cells[y] = make([]bool, w, w)
	}
	return &Grid{cells: cells, h: h, w: w}
}

func (g *Grid) Cells() [][]bool {
	return g.cells
}

func (g *Grid) Height() int {
	return g.h
}

func (g *Grid) Width() int {
	return g.w
}

func (g *Grid) Get(x, y int) bool {
	if x < 0 || y < 0 || x >= g.w || y >= g.h {
		panic(fmt.Sprintf("Grid.Get(%d, %d) is oob: Grid(w=%d, h=%d)", x, y, g.w, g.h))
	}
	return g.cells[y][x]
}

func (g *Grid) Set(x, y int, b bool) {
	if x < 0 || y < 0 || x >= g.w || y >= g.h {
		panic(fmt.Sprintf("Grid.Set(%d, %d) is oob: Grid(w=%d, h=%d)", x, y, g.w, g.h))
	}
	g.cells[y][x] = b
}

func (g *Grid) IsEmpty() bool {
	for y := range g.cells {
		for x := range g.cells[y] {
			if g.cells[y][x] {
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
	w, h := grid.Width(), grid.Height()
	if x+w > g.Width() || y+h > g.Height() {
		return
	}
	for yy, j := 0, y; yy < h; yy++ {
		for xx, i := 0, x; xx < w; xx++ {
			// TODO: try g.Set() as GetSubgrid does
			g.cells[j][i] = grid.Get(xx, yy)
			i++
		}
		j++
	}
}

func (g *Grid) GetSubgrid(x, y, w, h int) *Grid {
	if x+w > g.w || y+h > g.h {
		panic(fmt.Sprintf("subgrid is oob: (%d, %d) w=%d, h=%d on grid w=%d, h=%d", x, y, w, h, g.w, g.h))
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
	for x := 0; x < g.w; x++ {
		// x is the new y
		row := make([]bool, g.h, g.h)
		for y := 0; y < g.h; y++ {
			// y is x from right to left
			row[g.h-y-1] = g.cells[y][x]
		}
		cells = append(cells, row)
	}
	return &Grid{
		cells: cells,
		w:     g.h,
		h:     g.w,
	}
}

func (g *Grid) Row(y int) []bool {
	row := make([]bool, g.w, g.w)
	for x := 0; x < g.w; x++ {
		row[x] = g.Get(x, y)
	}
	return row
}

func (g *Grid) IsRowEmpty(y int) bool {
	for x := 0; x < g.w; x++ {
		if g.Get(x, y) {
			return false
		}
	}
	return true
}

func (g *Grid) IsColEmpty(x int) bool {
	for y := 0; y < g.h; y++ {
		if g.Get(x, y) {
			return false
		}
	}
	return true
}

func (g *Grid) String() string {
	var s bytes.Buffer
	for y := range g.cells {
		for x := range g.cells[y] {
			r := '.'
			if g.cells[y][x] {
				r = '█'
			}
			s.WriteRune(r)
		}
		s.WriteRune('\n')
	}
	return s.String()
}
