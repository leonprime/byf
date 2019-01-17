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

type Piece struct {
	Name   string
	Shapes []*Grid
	// # of rotation symmetries
	Rotate int
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
