package game

import (
	"strings"
	"testing"
)

func TestNewGrid(t *testing.T) {
	g := newGrid(`
piece k
██
█.
.█
..
`)
	expect := [][]bool{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}
	if g.Height() != len(expect) {
		t.Fatal("unexpected rows")
	}
	if g.Width() != len(expect[0]) {
		t.Fatal("unexpected # cols")
	}
	for x := 0; x < g.Width(); x++ {
		for y := 0; y < g.Height(); y++ {
			if expect[y][x] != g.Get(x, y) {
				t.Fatalf("expected %v at (%d, %d).  Got %v", expect[y][x], x, y, g.Get(x, y))
			}
		}
	}
}

func TestRotateGrid(t *testing.T) {
	g := newGrid(`
██
█.
.█
..
`)
	g = g.Rotate()
	expect := [][]bool{
		{false, false, true, true},
		{false, true, false, true},
	}
	if g.Height() != len(expect) {
		t.Fatal("unexpected rows")
	}
	if g.Width() != len(expect[0]) {
		t.Fatal("unexpected # cols")
	}
	for x := 0; x < g.Width(); x++ {
		for y := 0; y < g.Height(); y++ {
			if expect[y][x] != g.Get(x, y) {
				t.Fatalf("expected %v at (%d, %d).  Got %v", expect[y][x], x, y, g.Get(x, y))
			}
		}
	}

}

func TestGridNotRectangle(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	newGrid(`
██
.█
█.█
`)
}

func TestParsePieces(t *testing.T) {
	data := `
piece t
rotate 9
###
.#.
piece n
rotate 2
.██.
█..█
█..█
piece t
rotate 3
...
█.█
`
	pieces := ParsePieces(strings.NewReader(data))

	pt, ok := pieces["t"]
	if !ok {
		t.Fatal("missing t")
	}
	if len(pt.Shapes) != 2 {
		t.Error("expected 2 t pieces")
	}
	if pt.Rotate != 9 {
		t.Errorf("expected t to have 9 rotations, got %d", pt.Rotate)
	}

	n, ok := pieces["n"]
	if !ok {
		t.Fatal("missing n")
	}
	if len(n.Shapes) != 1 {
		t.Error("expected 1 n pieces")
	}
	if n.Rotate != 2 {
		t.Errorf("expected n to have 2 rotations, got %d", n.Rotate)
	}

	tests := []*Grid{
		pt.Shapes[0],
		pt.Shapes[1],
		n.Shapes[0],
	}
	expects := [][][]bool{
		{ // first t
			{true, true, true},
			{false, true, false},
		},
		{ // second t
			{false, false, false},
			{true, false, true},
		},
		{ // n
			{false, true, true, false},
			{true, false, false, true},
			{true, false, false, true},
		},
	}
	for i := 0; i < 3; i++ {
		g := tests[i]
		expect := expects[i]
		for x := 0; x < g.Width(); x++ {
			for y := 0; y < g.Height(); y++ {
				if expect[y][x] != g.Get(x, y) {
					t.Fatalf("expected %v at (%d, %d).  Got %v", expect[y][x], x, y, g.Get(x, y))
				}
			}
		}
	}
}

func TestEmpty(t *testing.T) {
	grid := newEmptyGrid(5, 3)
	if !grid.IsEmpty() {
		t.Fatal("expected empty grid, got non-empty")
	}
	grid.cells[2][4] = true
	if grid.IsEmpty() {
		t.Fatal("expected non-empty grid, got empty")
	}
}

func TestSetSubgrid(t *testing.T) {
	g := newEmptyGrid(5, 3)
	piece := newGrid("███\n.█.")

	// attempt oob at various positions
	g.SetSubgrid(100, 100, piece)
	if !g.IsEmpty() {
		t.Fatal("expected empty grid, got non-empty")
	}
	// too far down
	g.SetSubgrid(0, 2, piece)
	if !g.IsEmpty() {
		t.Fatal("expected empty grid, got non-empty")
	}
	// too far to the right
	g.SetSubgrid(3, 0, piece)
	if !g.IsEmpty() {
		t.Fatal("expected empty grid, got non-empty")
	}

	// set it cleanly at lower right
	g.SetSubgrid(2, 1, piece)
	if g.IsEmpty() {
		t.Fatal("expected non-empty grid, got empty")
	}
	expect := [][]bool{
		{false, false, false, false, false},
		{false, false, true, true, true},
		{false, false, false, true, false},
	}
	for x := 0; x < g.Width(); x++ {
		for y := 0; y < g.Height(); y++ {
			if expect[y][x] != g.Get(x, y) {
				t.Fatalf("expected %v at (%d, %d).  Got %v", expect[y][x], x, y, g.Get(x, y))
			}
		}
	}
}
