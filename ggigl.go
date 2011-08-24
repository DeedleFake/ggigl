package main

import (
	"os"
	"fmt"
	"sdl"
	"flag"
)

const (
	// The window's title.
	WinCap = "GGIGL: Go Game In Go Lang"
)

// Stores game information and basically runs everything.
type game struct {
	running bool

	pieces map[string]*Piece
	board  *Board

	selX   int
	selY   int
	turn   *Piece
	score map[*Piece]float64
	passed bool

	screen *sdl.Surface
}

// Runs everything, calling the methods required to get the game
// running, runs the main loop, and then cleanly exits.
func (g *game) run() (err os.Error) {
	err = g.load()
	if err != nil {
		return
	}
	defer g.quit()

	if sdl.Init(sdl.INIT_EVERYTHING) < 0 {
		return os.NewError(sdl.GetError())
	}
	defer sdl.Quit()

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

// Runs the games main loop.
func (g *game) main() (err os.Error) {
	g.running = true

	for g.running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch ev := e.(type) {
			case *sdl.KeyboardEvent:
				err = g.onKeyboardEvent(ev)
				if err != nil {
					return
				}
			case *sdl.MouseMotionEvent:
				err = g.onMouseMotionEvent(ev)
				if err != nil {
					return
				}
			case *sdl.MouseButtonEvent:
				err = g.onMouseButtonEvent(ev)
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

// Handles key presses.
func (g *game) onKeyboardEvent(ev *sdl.KeyboardEvent) (err os.Error) {
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
			g.placeTurn()
		case 'p':
			if g.passTurn() {
				// The game just ended.
			}
		}
	case sdl.KEYUP:
	}

	return
}

// Handles the mouse moving.
func (g *game) onMouseMotionEvent(ev *sdl.MouseMotionEvent) (err os.Error) {
	g.selX, g.selY = g.board.XYToCoord(int(ev.X), int(ev.Y))

	return
}

// Handles mouse clicks.
func (g *game) onMouseButtonEvent(ev *sdl.MouseButtonEvent) (err os.Error) {
	switch ev.Type {
	case sdl.MOUSEBUTTONDOWN:
		switch ev.Button {
		case sdl.BUTTON_LEFT:
			g.placeTurn()
			//case sdl.BUTTON_RIGHT:
			//	g.board.remove(g.selX, g.selY)
		}
	}

	return
}

// Switches turns.
func (g *game) changeTurns() {
	switch g.turn {
	case g.pieces["black"]:
		g.turn = g.pieces["white"]
	case g.pieces["white"]:
		g.turn = g.pieces["black"]
	default:
		panic("Invalid turn")
	}

	g.passed = false
}

// Takes a normal turn, placing a piece at the predetermined coordinates.
func (g *game) placeTurn() {
	if c := g.board.Place(g.selX, g.selY, g.turn); c >= 0 {
		g.changeTurns()
		g.score[g.turn] -= float64(c)
	}
}

// Passes. Checks if it's the second pass in a row. If it is, it
// returns true.
func (g *game) passTurn() bool {
	if g.passed {
		return true
	}

	g.changeTurns()
	g.passed = true

	return false
}

// Draws everything.
func (g *game) draw() (err os.Error) {
	if g.screen.Blit(nil, g.board.Image(), nil) < 0 {
		return os.NewError(sdl.GetError())
	}

	if (g.selX >= 0) || (g.selY >= 0) {
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
	}

	return
}

// Runs all of the initial start-up stuff.
func (g *game) load() (err os.Error) {
	var (
		size     int
		handicap int
		komi     float64
		superko  bool
	)
	flag.IntVar(&size,
		"size",
		int(Size19x19),
		"Board size; accepts from list: (9, 19)",
	)
	flag.IntVar(&handicap,
		"handicap",
		0,
		fmt.Sprintf("Handicap; maximums: (%v: %v, %v: %v)",
			Size9x9,
			MaxHandicap(Size9x9),
			Size19x19,
			MaxHandicap(Size19x19),
		),
	)
	flag.Float64Var(&komi,
		"komi",
		-1,
		"Komi; -1 to set based on handicap; Default: 5.5",
	)
	flag.BoolVar(&superko,
		"superko",
		false,
		"Use super ko instead of simple ko.",
	)
	flag.Parse()

	switch BoardSize(size) {
	case Size9x9, Size19x19:
	default:
		return fmt.Errorf("Bad board size: %v", size)
	}

	if (handicap < 0) || (handicap > MaxHandicap(BoardSize(size))) {
		return fmt.Errorf("Bad handicap: %v", handicap)
	}

	if komi < 0 {
		komi = 5.5
		if handicap != 0 {
			komi = 0
		}
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

	g.score = make(map[*Piece]float64)
	g.score[g.pieces["black"]] += komi

	ko := SimpleKo
	if superko {
		ko = SuperKo
	}

	g.board, err = NewBoard(BoardSize(size), ko)
	if err != nil {
		return
	}

	handi, err := GetHandicap(BoardSize(size), handicap)
	if err != nil {
		return
	}
	g.board.ApplyHandicap(g.pieces["black"], handi)

	return
}

// Run when the game exits.
func (g *game) quit() {
	// This will probably eventually autosave game information.
}

func main() {
	var g game
	err := g.run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
