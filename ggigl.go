package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// The window's title.
	WinCap = "GGIGL: Go Game In Go Lang"
)

// Stores game information and basically runs everything.
type game struct {
	board *Board

	selX   int
	selY   int
	turn   *Piece
	score  map[*Piece]float64
	passed bool
}

// Runs everything, calling the methods required to get the game
// running, runs the main loop, and then cleanly exits.
func (g *game) run() error {
	err := g.init()
	if err != nil {
		return fmt.Errorf("load assets: %w", err)
	}
	defer g.quit()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(WinCap)

	err = ebiten.RunGame(g)
	if err != nil {
		return fmt.Errorf("run game: %w", err)
	}

	return nil
}

// Runs all of the initial start-up stuff.
func (g *game) init() (err error) {
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

	if (handicap < 0) || (handicap > MaxHandicap(BoardSize(size))) {
		return fmt.Errorf("bad handicap: %v", handicap)
	}

	if komi < 0 {
		komi = 5.5
		if handicap != 0 {
			komi = 0
		}
	}

	g.board, err = NewBoard(BoardSize(size), superko)
	if err != nil {
		return fmt.Errorf("initialize board: %w", err)
	}

	g.turn = Black
	g.score = make(map[*Piece]float64)
	g.score[Black] += komi

	handi, err := GetHandicap(BoardSize(size), handicap)
	if err != nil {
		return fmt.Errorf("initialize handicap: %w", err)
	}
	g.board.ApplyHandicap(Black, handi)

	return nil
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.selY--
		if g.selY < 0 {
			g.selY = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.selY++
		s := int(g.board.Size())
		if g.selY >= s {
			g.selY = s - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.selX--
		if g.selX < 0 {
			g.selX = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.selX++
		s := int(g.board.Size())
		if g.selX >= s {
			g.selX = s - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.placeTurn()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if g.passTurn() {
			// The game just ended.
		}
	}

	g.selX, g.selY = g.board.XYToCoord(ebiten.CursorPosition())
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.placeTurn()
	}
	//if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
	//	g.board.remove(g.selX, g.selY)
	//}
	return nil
}

// Switches turns.
func (g *game) changeTurns() {
	switch g.turn {
	case Black:
		g.turn = White
	case White:
		g.turn = Black
	default:
		panic("invalid turn")
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

func (g *game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.board.Draw(screen, ebiten.GeoM{})

	if (g.selX >= 0) || (g.selY >= 0) {
		sx, sy := g.board.CoordToXY(g.selX, g.selY)
		draw.Draw(
			screen,
			image.Rect(sx-10, sy-10, sx+10, sy+10),
			g.turn.HighlightColor(),
			image.Point{},
			draw.Over,
		)
		//timg := g.turn.Image()
		//sx -= int(timg.W / 2)
		//sy -= int(timg.H / 2)
		//timg.SetAlpha(sdl.SRCALPHA, 128)
		//g.screen.Blit(&sdl.Rect{X: int16(sx), Y: int16(sy)}, timg, nil)
		//timg.SetAlpha(sdl.SRCALPHA, 255)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

// Run when the game exits.
func (g *game) quit() {
	// This will probably eventually autosave game information.
}

func main() {
	var g game
	err := g.run()
	if err != nil {
		slog.Error("failed to run game", "err", err)
		os.Exit(1)
	}
}
