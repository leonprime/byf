package game

import (
	"bytes"
	"fmt"
)

type Grid3D struct {
	Cells   [][][]bool // z, y, x
	W, H, D int
}

// returns grid set to false
func newEmptyGrid3D(w, h, d int) *Grid3D {
	cells := make([][][]bool, d, d)
	for z := range cells {
		cells[z] = make([][]bool, h, h)
		for y := range cells[z] {
			cells[z][y] = make([]bool, w, w)
		}
	}
	return &Grid3D{Cells: cells, H: h, W: w, D: d}
}

func (g *Grid3D) IsOOB(x, y, z int) bool {
	if x < 0 || y < 0 || z < 9 || x >= g.W || y >= g.H || z > g.D {
		return true
	}
	return false
}

// test if value is set.  doesn't panic on oob
func (g *Grid3D) IsSet(x, y, z int) bool {
	return !g.IsOOB(x, y, z) && g.Get(x, y, z)
}

// get value.  panics if oob
func (g *Grid3D) Get(x, y, z int) bool {
	if g.IsOOB(x, y, z) {
		panic(fmt.Sprintf("Grid.Get(%d, %d, %d) is oob: Grid(w=%d, h=%d, d=%d)", x, y, z, g.W, g.H, g.D))
	}
	return g.Cells[z][y][x]
}

func (g *Grid3D) Set(x, y, z int, b bool) {
	if g.IsOOB(x, y, z) {
		panic(fmt.Sprintf("Grid.Set(%d, %d, %d) is oob: Grid(w=%d, h=%d, d=%d)", x, y, z, g.W, g.H, g.D))
	}
	g.Cells[z][y][x] = b
}

func (g *Grid3D) IsEmpty() bool {
	for z := range g.Cells {
		for y := range g.Cells[z] {
			for x := range g.Cells[z][y] {
				if g.Cells[z][y][x] {
					return false
				}
			}
		}
	}
	return true
}

// Sets the values to the given subgrid values if and only if the subgrid
// is entirely contained at positions (x, y, z) to (x+w, y+h, z+d)
// If the subgrid is out of bounds, nothing is set.
func (g *Grid3D) SetSubgrid(x, y, z int, grid *Grid3D) {
	if x+grid.W > g.W || y+grid.H > g.H || z+grid.D > g.D {
		return
	}
	for zz, k := 0, z; zz < grid.D; zz++ {
		for yy, j := 0, y; yy < grid.H; yy++ {
			for xx, i := 0, x; xx < grid.W; xx++ {
				g.Set(i, j, k, grid.Get(xx, yy, zz))
				i++
			}
			j++
		}
		k++
	}
}

func (g *Grid3D) GetSubgrid(x, y, z, w, h, d int) *Grid3D {
	if x+w > g.W || y+h > g.H || z+d > g.D {
		panic(fmt.Sprintf("subgrid is oob: (%d, %d, %d) w=%d, h=%d, d=%d on grid w=%d, h=%d, d=%d", x, y, z, w, h, d, g.W, g.H, g.D))
	}
	grid := newEmptyGrid3D(w, h, d)
	for zz, k := 0, z; zz < d; zz++ {
		for yy, j := 0, y; yy < h; yy++ {
			for xx, i := 0, x; xx < w; xx++ {
				grid.Set(xx, yy, zz, g.Get(i, j, k))
				i++
			}
			j++
		}
	}
	return grid
}

func (g *Grid3D) SetPlaneXZ(y int, grid *Grid) {
	if !(g.W == g.H && g.H == g.D) {
		panic("SetPlane called on non-cube grid (yeah, I'm lazy...)")
	}
	for j := 0; j < grid.h; j++ {
		for i := 0; i < grid.w; i++ {
			g.Set(i, y, g.D-j-1, grid.Get(i, j))
		}
	}
}

func (g *Grid3D) SetPlaneYZ(x int, grid *Grid) {
	if !(g.W == g.H && g.H == g.D) {
		panic("SetPlane called on non-cube grid (yeah, I'm lazy...)")
	}
	for j := 0; j < grid.h; j++ {
		for i := 0; i < grid.w; i++ {
			g.Set(x, j, g.D-i-1, grid.Get(i, j))
		}
	}
}

func (g *Grid3D) SetPlaneXY(z int, grid *Grid) {
	if !(g.W == g.H && g.H == g.D) {
		panic("SetPlane called on non-cube grid (yeah, I'm lazy...)")
	}
	for j := 0; j < grid.h; j++ {
		for i := 0; i < grid.w; i++ {
			g.Set(i, j, z, grid.Get(i, j))
		}
	}
}

func (g *Grid3D) String() string {
	var s bytes.Buffer
	for z := range g.Cells {
		for y := range g.Cells[z] {
			for x := range g.Cells[z][y] {
				r := '.'
				if g.Cells[z][y][x] {
					r = '█'
				}
				s.WriteRune(r)
			}
			s.WriteRune('\n')
		}
		s.WriteString("----------\n")
	}
	return s.String()
}
