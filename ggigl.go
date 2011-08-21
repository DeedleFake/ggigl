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

	selX int
	selY int
	turn *Piece

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
			case *sdl.KeyboardEvent:
				err = g.onKeyEvent(ev)
				if err != nil {
					return
				}
			case *sdl.QuitEvent:
				g.running = false
			}
		}

		g.screen.FillRect(nil, 0)
		err = g.draw()
		if err != nil {
			return
		}
		g.screen.Flip()
	}

	return
}

func (g *game) onKeyEvent(ev *sdl.KeyboardEvent) (err os.Error) {
	switch ev.Type {
	case sdl.KEYDOWN:
		switch ev.Keysym.Sym {
		case sdl.K_UP, 'k':
			g.selY--
			if g.selY < 0 {
				g.selY = 0
			}
		case sdl.K_DOWN, 'j':
			g.selY++
			s := int(g.board.Size())
			if g.selY >= s {
				g.selY = s - 1
			}
		case sdl.K_LEFT, 'h':
			g.selX--
			if g.selX < 0 {
				g.selX = 0
			}
		case sdl.K_RIGHT, 'l':
			g.selX++
			s := int(g.board.Size())
			if g.selX >= s {
				g.selX = s - 1
			}
		case sdl.K_SPACE:
			if g.board.Place(g.selX, g.selY, g.turn) {
				g.changeTurns()
			}
		}
	case sdl.KEYUP:
	}

	return
}

func (g *game) draw() (err os.Error) {
	if g.screen.Blit(nil, g.board.Image(), nil) < 0 {
		return os.NewError(sdl.GetError())
	}

	sx, sy := g.board.CoordToXY(g.selX, g.selY)
	//timg := g.turn.Image()
	//sx -= int(timg.W / 2)
	//sy -= int(timg.H / 2)
	//timg.SetAlpha(sdl.SRCALPHA, 128)
	//g.screen.Blit(&sdl.Rect{X: int16(sx), Y: int16(sy)}, timg, nil)
	//timg.SetAlpha(sdl.SRCALPHA, 255)
	switch g.turn {
	case g.pieces["black"]:
		g.screen.FillRect(&sdl.Rect{int16(sx - 10), int16(sy - 10), 20, 20},
			sdl.MapRGBA(g.screen.Format, 0, 0, 0, 128),
		)
	case g.pieces["white"]:
		g.screen.FillRect(&sdl.Rect{int16(sx - 10), int16(sy - 10), 20, 20},
			sdl.MapRGBA(g.screen.Format, 255, 255, 255, 128),
		)
	}

	return
}

func (g *game) load() (err os.Error) {
	var (
		size int
		handicap int
	)
	flag.IntVar(&size,
		"size",
		int(Size19x19),
		"Board size; accepts from list: (9, 19)",
	)
	flag.IntVar(&handicap,
		"handicap",
		0,
		fmt.Sprintf("Handicap; maximum: %v", MaxHandicap()),
	)
	flag.Parse()

	switch BoardSize(size) {
	case Size9x9, Size19x19:
	default:
		return fmt.Errorf("Bad board size: %v", size)
	}

	if (handicap < 0) || (handicap > MaxHandicap()) {
		return fmt.Errorf("Bad handicap: %v", handicap)
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

	g.turn = g.pieces["black"]

	g.board, err = NewBoard(BoardSize(size))
	if err != nil {
		return
	}

	g.board.ApplyHandicap(g.pieces["black"], GetHandicap(handicap))

	return
}

func (g *game) changeTurns() {
	switch g.turn {
	case g.pieces["black"]:
		g.turn = g.pieces["white"]
	case g.pieces["white"]:
		g.turn = g.pieces["black"]
	default:
		panic("Invalid turn")
	}
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
