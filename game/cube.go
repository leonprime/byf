package game

import "fmt"

type Play3D struct {
	Piece   *Piece
	Grid    *Grid3D
	X, Y, Z int
}

func (p *Play3D) String() string {
	return fmt.Sprintf("%s: (%d, %d, %d) w=%d, h=%d, d=%d\n%s", p.Piece.Name, p.X, p.Y, p.Z, p.Grid.W, p.Grid.H, p.Grid.D, p.Grid)
}

type Cube struct {
	W, H, D  int
	pieces   []*Piece
	Coverage *Coverage
}

func NewCube(w, h, d int, piecesSpec string) *Cube {
	c := &Cube{
		pieces: parsePiecesSpec(piecesSpec),
		W:      w,
		H:      h,
		D:      d,
	}
	c.Coverage = newCubeCoverage(c)
	return c
}
