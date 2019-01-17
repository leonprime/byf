package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"torres.guru/gagne/display"
	"torres.guru/gagne/dlx"
	"torres.guru/gagne/game"
)

type Game struct {
	pieceSpec string
	w, h      int
}

func (g Game) String() string {
	return fmt.Sprintf("%dx%d_%s", g.w, g.h, g.pieceSpec)
}

func main() {
	path := flag.String("path", ".", "output path for game solutions")
	n := flag.Int("n", 10, "number of solutions to print")
	debug := flag.Bool("debug", false, "turn on debugging")

	flag.Usage = func() {
		f := flag.CommandLine.Output()
		fmt.Fprintf(f, "Usage: %s [options] w h pieceSpec\n", os.Args[0])
		fmt.Fprintf(f, "\tw and h are the board width and height\n")
		fmt.Fprintf(f, "\tpieceSpec is the set of pieces to play with (see data/pieces.txt)\n")
		fmt.Fprintf(f, "\texample: %s 5 3 otzrI\n", os.Args[0])
		fmt.Fprintf(f, "\tthe solutions are saved at ${path}/solutions/5x3_otzrI\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	if *debug {
		game.SetDebug()
		dlx.SetDebug()
	}

	if len(flag.Args()) != 3 {
		flag.Usage()
	}
	w, _ := strconv.Atoi(flag.Args()[0])
	h, _ := strconv.Atoi(flag.Args()[1])
	pieceSpec := flag.Args()[2]
	if w == 0 || h == 0 || len(pieceSpec) == 0 {
		flag.Usage()
	}

	g := Game{w: w, h: h, pieceSpec: pieceSpec}
	g.validate()
	g.run(*path, *n)
}

func (g Game) validate() {
}

func (g Game) run(path string, n int) {
	game.LoadPieces("data/pieces.txt")
	board := game.NewBoardGame(g.w, g.h, g.pieceSpec)
	dl := dlx.New(board.Coverage.M.Cells(), board.Coverage.Columns)
	dl.Search(0)

	fmt.Printf("game \"%s\" has %d solutions\n", g, len(dl.Solutions))

	gamePath := fmt.Sprintf("%s/solutions/%s", path, g)
	os.RemoveAll(gamePath)
	os.MkdirAll(gamePath, os.ModePerm)

	if len(dl.Solutions) < n {
		n = len(dl.Solutions)
	}
	for i := 0; i < n; i++ {
		plays := board.Play(dl.Solutions[i])
		filename := fmt.Sprintf("%s/%d.png", gamePath, i)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		display.Render(board.W, board.H, plays, f)
		f.Close()
	}
}
