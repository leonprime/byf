package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
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
	path := flag.String("path", ".", "output path for game solutions.")
	max := flag.Int("max", 0, "maximum number of solutions to find.  the default, 0, means find all solutions")
	n := flag.Int("n", 10, "number of solutions to print")
	debug := flag.Bool("debug", false, "turn on debugging")
	pieces := flag.String("pieces", "data/pieces.txt", "path to pieces data file")
	show := flag.Bool("show", false, "print available pieces and quit")
	countOnly := flag.Bool("countOnly", false, "don't print, just count solutions")

	flag.Usage = func() {
		f := flag.CommandLine.Output()
		fmt.Fprintf(f, "Usage: %s [options] w h pieceSpec\n", os.Args[0])
		fmt.Fprintf(f, "  w and h are the board width and height\n")
		fmt.Fprintf(f, "  pieceSpec is the set of pieces to play with (see data/pieces.txt)\n")
		fmt.Fprintf(f, "Example: %s 5 3 otzrI\n", os.Args[0])
		fmt.Fprintf(f, "  the solutions are saved at ${path}/solutions/5x3_otzrI\n")
		fmt.Fprintf(f, "Options:\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	game.LoadPieces(*pieces)
	if *show {
		for _, piece := range game.AllPieces() {
			fmt.Println(piece)
		}
		return
	}

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
	g.run(*path, *n, *max, *countOnly)
}

func (g Game) run(path string, n, max int, countOnly bool) {
	board := game.NewBoardGame(g.w, g.h, g.pieceSpec)
	dl := dlx.New(board.Coverage.M.Cells(), board.Coverage.Columns, max, countOnly)
	t := time.Now()
	dl.Search(0)
	dur := time.Now().Sub(t)

	fmt.Printf("game \"%s\" has %d solutions\n", g, dl.N)
	fmt.Printf("\ttime taken: %s\n", dur)
	fmt.Printf("\tsteps: %d\n", dl.S)
	if dl.N == 0 || countOnly {
		return
	}

	gamePath := fmt.Sprintf("%s/solutions/%s", path, g)
	os.RemoveAll(gamePath)
	os.MkdirAll(gamePath, os.ModePerm)

	if dl.N < n {
		n = dl.N
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
	fmt.Printf("wrote the first %d solutions to %s\n", n, gamePath)
}
