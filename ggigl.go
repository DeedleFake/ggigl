package main

import(
	"os"
	"fmt"
	"flag"
)

type game struct {
	pieces map[string]*Piece
	board *Board
}

func (g *game)run() (err os.Error) {
	var(
		size int
	)
	flag.IntVar(&size, "size", int(Size19x19), "Board size; accepts from list: (9, 19)")
	flag.Parse()

	switch size {
		case int(Size9x9), int(Size19x19):
		default:
			return fmt.Errorf("Bad board size: %v", size)
	}

	g.pieces = make(map[string]*Piece)

	g.pieces["black"], err = NewPiece("black")
	if err != nil {
		return
	}

	g.pieces["white"], err = NewPiece("white")
	if err != nil {
		return
	}

	g.board, err = NewBoard(BoardSize(size))
	if err != nil {
		return
	}

	return
}

func main() {
	var g game
	err := g.run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	println("There is no response...")
}
