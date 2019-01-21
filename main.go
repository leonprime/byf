package main

import (
	"flag"
	"fmt"
	"github.com/leonprime/byf/display"
	"github.com/leonprime/byf/dlx"
	"github.com/leonprime/byf/game"
	"io"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Game interface {
	Coverage() *game.Coverage
	Render(io.Writer, []int)
	String() string
}

type Game2D struct {
	pieceSpec string
	w, h      int
	board     *game.Board
}

func (g *Game2D) Coverage() *game.Coverage {
	g.board = game.NewBoard(g.w, g.h, g.pieceSpec)
	return g.board.Coverage
}

func (g *Game2D) Render(w io.Writer, rows []int) {
	plays := g.board.Play(rows)
	display.Render(g.board.W, g.board.H, plays, w)
}

func (g *Game2D) String() string {
	return fmt.Sprintf("%dx%d_%s", g.w, g.h, g.pieceSpec)
}

type Game3D struct {
	pieceSpec string
	w, h, d   int
	cube      *game.Cube
}

func (g *Game3D) Coverage() *game.Coverage {
	g.cube = game.NewCube(g.w, g.h, g.d, g.pieceSpec)
	return g.cube.Coverage
}

func (g *Game3D) Render(w io.Writer, rows []int) {
	plays := g.cube.Play(rows)
	display.Render3D(g.cube.W, g.cube.H, g.cube.D, plays, w)
}

func (g *Game3D) String() string {
	return fmt.Sprintf("%dx%dx%d_%s", g.w, g.h, g.d, g.pieceSpec)
}

func main() {
	max := flag.Int("max", 0, "max solutions to find.  0 means find all (default 0)")
	nprint := flag.Int("print", 10, "number of solutions to print")
	path := flag.String("path", ".", "output path for game solutions.")
	pieces := flag.String("pieces", "data/gagne.txt", "path to pieces data file")
	debug := flag.Bool("debug", false, "turn on all debugging")
	debugPiece := flag.String("debugPiece", "", "debug a specific piece")
	debugAllPieces := flag.Bool("debugAllPieces", false, "debug all pieces")
	debugCoverage := flag.Bool("debugCoverage", false, "debug coverage matrix")
	debugDLX := flag.Bool("debugDLX", false, "debug DLX algorithm")
	show := flag.Bool("show", false, "print available pieces and quit")
	nochiral := flag.Bool("nochiral", false, "don't include the chiral reflections. the first one found in data file is used")

	flag.Usage = func() {
		f := flag.CommandLine.Output()
		fmt.Fprintf(f, "Usage: %s [options] w h pieceSpec\n", os.Args[0])
		fmt.Fprintf(f, "  w and h are the board width and height\n")
		fmt.Fprintf(f, "  pieceSpec is the set of pieces to play with (see data/gagne.txt)\n")
		fmt.Fprintf(f, "Example: %s 5 3 otzvI\n", os.Args[0])
		fmt.Fprintf(f, "  the solutions are saved at ${path}/solutions/5x3_otzvI\n")
		fmt.Fprintf(f, "Options:\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	game.LoadPieces(*pieces, !*nochiral)
	if *show {
		for _, piece := range game.AllPieces() {
			fmt.Println(piece)
		}
		return
	}

	args := flag.Args()
	var (
		dim       int
		w, h, d   int
		pieceSpec string
		err       error
	)
	if len(args) == 3 {
		dim = 2
		pieceSpec = args[2]
	} else if len(args) == 4 {
		dim = 3
		d, err = strconv.Atoi(args[2])
		if err != nil {
			flag.Usage()
		}
		pieceSpec = args[3]
	} else {
		flag.Usage()
	}
	w, err = strconv.Atoi(args[0])
	if err != nil {
		flag.Usage()
	}
	h, err = strconv.Atoi(args[1])
	if err != nil {
		flag.Usage()
	}
	if w == 0 || h == 0 || len(pieceSpec) == 0 {
		flag.Usage()
	}
	if dim == 3 && d == 0 {
		flag.Usage()
	}

	if *debugPiece != "" {
		game.SetDebugPiece(*debugPiece)
	}
	if *debugAllPieces {
		game.SetDebugAllPieces()
	}
	if *debugCoverage {
		game.SetDebugCoverage()
	}
	if *debugDLX {
		dlx.SetDebug()
	}
	if *debug {
		game.SetDebugAllPieces()
		game.SetDebugCoverage()
		dlx.SetDebug()
	}

	var g Game
	if dim == 2 {
		g = &Game2D{w: w, h: h, pieceSpec: pieceSpec}
	} else {
		g = &Game3D{w: w, h: h, d: d, pieceSpec: pieceSpec}
	}
	run(g, *path, *nprint, *max)
}

func run(g Game, path string, nprint, max int) {
	cov := g.Coverage()
	renderDebugs(cov.Debugs, g.String(), path)

	dl := dlx.New(cov.M.Cells, cov.Columns, max, nprint)

	start := time.Now()

	// print if interrupted
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		printSolutions(g, dl, path, nprint, start)
		os.Exit(0)
	}()

	dl.Search(0)

	printSolutions(g, dl, path, nprint, start)
}

func printSolutions(g Game, dl *dlx.DancingLinks, path string, nprint int, start time.Time) {
	if dl.N >= 1000 {
		fmt.Print("\r") // clear out the count feedback
	}
	fmt.Printf("found %d solutions for game \"%s\"\n", dl.N, g)
	fmt.Printf("\ttime taken: %s\n", time.Now().Sub(start))
	fmt.Printf("\tsteps: %d\n", dl.S)

	if len(dl.Solutions) == 0 {
		return
	}
	gamePath := fmt.Sprintf("%s/solutions/%s", path, g)
	os.RemoveAll(gamePath)
	os.MkdirAll(gamePath, os.ModePerm)

	for i, solution := range dl.Solutions {
		filename := fmt.Sprintf("%s/%d.png", gamePath, i)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		g.Render(f, solution)
		f.Close()
	}
	quant := "the first"
	if dl.N == nprint {
		quant = "all"
	}
	fmt.Printf("wrote %s %d solutions to %s\n", quant, len(dl.Solutions), gamePath)
}

func renderDebugs(debugs []*game.Debug, gameName, path string) {
	if len(debugs) == 0 {
		return
	}
	debugPath := fmt.Sprintf("%s/debug/%s", path, gameName)
	os.RemoveAll(debugPath)
	os.MkdirAll(debugPath, os.ModePerm)
	for i, debug := range debugs {
		// put a number in the name because mac is case insensitive
		playPath := fmt.Sprintf("%s/%2d_%s", debugPath, i, debug.Name)
		os.RemoveAll(playPath)
		os.MkdirAll(playPath, os.ModePerm)
		for i := range debug.Plays {
			filename := fmt.Sprintf("%s/play%d.png", playPath, i)
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			display.Render3D(debug.W, debug.H, debug.D, debug.Plays[i:i+1], f)
			f.Close()
		}
	}
	fmt.Printf("generated %d debuging plays in %s\n", len(debugs), debugPath)
}
