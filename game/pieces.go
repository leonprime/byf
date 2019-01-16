package game

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

type Grid struct {
	cells [][]bool
	w, h  int
}

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
			g.cells[j][i] = grid.Get(xx, yy)
			i++
		}
		j++
	}
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

type Piece struct {
	Name   string
	Shapes []*Grid
	// # of rotation symmetries
	Rotate int
}

func (p *Piece) String() string {
	var s bytes.Buffer
	s.WriteString(fmt.Sprintf("piece %s:\n", p.Name))
	for i, shape := range p.Shapes {
		s.WriteString(fmt.Sprintf("%d:\n", i))
		s.WriteString(shape.String())
	}
	return s.String()
}

func ParsePieces(r io.Reader) map[string]*Piece {
	pieces := make(map[string]*Piece)
	var lines []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		panic(s.Err())
	}
	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "piece") {
			continue
		}
		name := lines[i][6:7]
		var shape bytes.Buffer
		rotate := 0
		for j := i + 1; j < len(lines) && !strings.HasPrefix(lines[j], "piece"); j++ {
			if strings.HasPrefix(lines[j], "rotate") {
				rotate, _ = strconv.Atoi(lines[j][7:8])
				continue
			}
			shape.WriteString(lines[j])
			shape.WriteRune('\n')
		}
		grid := newGrid(shape.String())
		if piece, ok := pieces[name]; ok {
			piece.Shapes = append(piece.Shapes, grid)
		} else {
			pieces[name] = &Piece{Name: name, Shapes: []*Grid{grid}, Rotate: rotate}
		}
	}
	return pieces
}

var Pieces map[string]*Piece

func LoadPieces(fileName string) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	Pieces = ParsePieces(bytes.NewReader(b))
}
