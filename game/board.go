package game

import (
	"fmt"
)

// a play is a piece that's placed on the board
type Play struct {
	Piece *Piece
	Grid  *Grid // piece as oriented on board
	X, Y  int   // position on board
}

func (p *Play) String() string {
	return fmt.Sprintf("%s: (%d, %d) w=%d, h=%d\n%s", p.Piece.Name, p.X, p.Y, p.Grid.W, p.Grid.H, p.Grid)
}

type Board struct {
	W, H     int
	pieces   []*Piece
	Coverage *Coverage
}

func NewBoard(w, h int, piecesSpec string) *Board {
	b := &Board{
		pieces: parsePiecesSpec(piecesSpec),
		W:      w,
		H:      h,
	}
	b.Coverage = newBoardCoverage(b)
	return b
}

// play a DLX solution by reading the selected rows
// from the coverage data
// returns a list of plays that represent the placement of each piece
func (b *Board) Play(rows []int) (plays []*Play) {
	for _, y := range rows {
		plays = append(plays, b.play(y))
	}
	return
}

func (b *Board) play(y int) *Play {
	play := &Play{}
	p := len(b.pieces)
	row := b.Coverage.M.Row(y)
	//
	// scan for the piece
	for i := 0; i < p; i++ {
		if row[i] {
			play.Piece = b.pieces[i]
			break
		}
	}
	//
	// rebuild the grid from the coverage row
	grid := newEmptyGrid(b.W, b.H)
	for i := p; i < len(row); i++ {
		x := (i - p) % b.W
		y := (i - p) / b.W
		grid.Set(x, y, row[i])
	}
	//
	// scan for piece location and extents
	w, h := 0, 0
	found := false
	for y := 0; y < b.H; y++ {
		if found {
			if grid.IsRowEmpty(y) {
				break
			} else {
				h++
			}
		} else {
			if grid.IsRowEmpty(y) {
				play.Y++
			} else {
				h++
				found = true
			}
		}
	}
	found = false
	for x := 0; x < b.W; x++ {
		if found {
			if grid.IsColEmpty(x) {
				break
			} else {
				w++
			}
		} else {
			if grid.IsColEmpty(x) {
				play.X++
			} else {
				w++
				found = true
			}
		}
	}
	// trim the grid to the subgrid bounding the piece
	play.Grid = grid.GetSubgrid(play.X, play.Y, w, h)
	if debug.piece(play.Piece) {
		fmt.Println(play)
		fmt.Printf("subset (%d, %d) w=%d, h=%d of:\n", play.X, play.Y, w, h)
		fmt.Println(grid)
	}
	return play
}
