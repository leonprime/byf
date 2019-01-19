package game

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// Game piece.  Instead of generating all symmetries programmatically,
// we note all pieces are 2D and in a 2D game, pices are either symmetric or
// have a chirality.  Then, to generate all positional permutations, the only
// other thing we need are the rotational symmetries.  Furthermore, in a 3D
// game with 2D pieces, the same symmetries hold along each dimension.
type Piece struct {
	Name   string
	Shapes []*Grid // one shape if symmetrical, two if chiral
	Rotate int     // # of rotation symmetries
	Color  []uint8
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

// parses pieces from a piece spec
// a piece definition starts with "piece x" where x is the single character name of the piece
// followed by a single grid representing the piece.
// rotation symmetries are specified with "rotate n".  default is 0 (no rotation symmetries)
// color is specified with "color c" where c is a hex RGB value like FF0000
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
		color := make([]uint8, 3, 3)
		rotate := 0
		for j := i + 1; j < len(lines) && !strings.HasPrefix(lines[j], "piece"); j++ {
			if strings.HasPrefix(lines[j], "rotate") {
				r, err := strconv.Atoi(lines[j][7:8])
				if err != nil {
					panic(fmt.Sprintf("error parsing rotate%s: %s", lines[j], err))
				}
				rotate = r
				continue
			}
			if strings.HasPrefix(lines[j], "color") {
				str := lines[j][6:12]
				for i := 0; i < 3; i++ {
					n, err := strconv.ParseUint("0x"+str[i*2:i*2+2], 0, 8)
					if err != nil {
						panic(fmt.Sprintf("error parsing color %s: %s", lines[j], err))
					}
					color[i] = uint8(n)
				}
				continue
			}
			shape.WriteString(lines[j])
			shape.WriteRune('\n')
		}
		grid := newGrid(shape.String())
		if piece, ok := pieces[name]; ok {
			piece.Shapes = append(piece.Shapes, grid)
		} else {
			pieces[name] = &Piece{
				Name:   name,
				Shapes: []*Grid{grid},
				Rotate: rotate,
				Color:  color,
			}
		}
	}
	return pieces
}

var allPieces map[string]*Piece

// parse pieces from a file
func LoadPieces(fileName string) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	allPieces = ParsePieces(bytes.NewReader(b))
}

func AllPieces() []*Piece {
	var pieces []*Piece
	for _, piece := range allPieces {
		pieces = append(pieces, piece)
	}
	return pieces
}

// gets the pieces represented by a string of consecutive one-character piece names
func parsePiecesSpec(piecesSpec string) []*Piece {
	if allPieces == nil {
		panic("ensure LoadPieces(file) is called first")
	}
	var pieces []*Piece
	for _, char := range piecesSpec {
		if piece, ok := allPieces[string(char)]; ok {
			pieces = append(pieces, piece)
		} else {
			panic(fmt.Sprintf("no piece \"%c\" defined", char))
		}
	}
	return pieces
}
