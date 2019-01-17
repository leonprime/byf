package game

import (
	"fmt"
)

var debug bool

func SetDebug() {
	debug = true
}

type Play struct {
	piece *Piece
	grid  *Grid
	x, y  int
}

type BoardGame struct {
	w, h     int
	pieces   []*Piece
	Coverage *Coverage
	plays    []*Play
}

func NewBoardGame(w, h int, pieces_spec string) *BoardGame {
	if allPieces == nil {
		panic("ensure LoadPieces(file) is called first")
	}
	var pieces []*Piece
	for _, char := range pieces_spec {
		if piece, ok := allPieces[string(char)]; ok {
			pieces = append(pieces, piece)
		} else {
			panic(fmt.Sprintf("no piece \"%c\" defined", char))
		}
	}
	b := &BoardGame{
		pieces: pieces,
		w:      w,
		h:      h,
	}
	b.Coverage = newCoverage(b)
	return b
}

// play the solution by reading the selected rows
// from the coverage matrix and related data
func (b *BoardGame) Play(rows []int) {
}
