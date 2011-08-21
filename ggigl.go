package main

import (
	"os"
	"fmt"
	"sdl"
	"flag"
)

const (
	WinCap = "GGIGL: Go Game In Go Lang"
)

type game struct {
	running bool

	pieces map[string]*Piece
	board  *Board

	screen *sdl.Surface
}

func (g *game) run() (err os.Error) {
	err = g.load()
	if err != nil {
		return
	}
	defer g.quit()

	if sdl.Init(sdl.INIT_EVERYTHING) < 0 {
		return os.NewError(sdl.GetError())
	}

	g.screen = sdl.SetVideoMode(640, 480, 32, sdl.DOUBLEBUF)
	if g.screen == nil {
		return os.NewError(sdl.GetError())
	}

	sdl.WM_SetCaption(WinCap, "")

	err = g.main()
	if err != nil {
		return
	}

	return
}

func (g *game) main() (err os.Error) {
	g.running = true

	for g.running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch ev := e.(type) {
			case *sdl.QuitEvent:
				g.running = false
			}
		}

		g.screen.Flip()
	}

	return
}

func (g *game) load() (err os.Error) {
	var (
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

func (g *game) quit() {
	sdl.Quit()
}

func main() {
	var g game
	err := g.run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
