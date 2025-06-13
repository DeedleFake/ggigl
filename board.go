package main

import (
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// The path to the board data.
	BoardPath = path.Join("data", "boards")
)

// A type for differentiating board sizes from regular ints. Will
// most likely do something else eventually.
type BoardSize int

const (
	Size9x9   BoardSize = 9
	Size19x19 BoardSize = 19
)

func (size BoardSize) Valid() bool {
	return size == Size9x9 || size == Size19x19
}

func (size BoardSize) image() *ebiten.Image {
	switch size {
	case Size9x9:
		return board9x9Image
	case Size19x19:
		return board19x19Image
	default:
		panic(fmt.Errorf("invalid board size: %v", size))
	}
}

func (size BoardSize) scaleUp(x, y int) (int, int) {
	switch size {
	case Size9x9:
		x = (x * 52) + 31
		y = (y * 52) + 31
	case Size19x19:
		x = (x * 25) + 14
		y = (y * 25) + 14
	}

	return x, y
}

func (size BoardSize) scaleDown(x, y int) (int, int) {
	switch size {
	case Size9x9:
		x /= 52
		y /= 52
	case Size19x19:
		x /= 25
		y /= 25
	}

	return x, y
}

// Stores information about changes to the board.
type Placement struct {
	p   *Piece
	loc [2]int
}

// Stores board information. Also calculates score and keeps track of
// the basic rules.
type Board struct {
	size   int
	pieces []*Piece
	tmp    []*Piece
	turns  [][]Placement

	ko func() bool

	img *ebiten.Image
}

// Initializes a new board of the given size using the given ko rule.
func NewBoard(size BoardSize, superko bool) (*Board, error) {
	if !size.Valid() {
		return nil, fmt.Errorf("invalid board size: %v", size)
	}

	b := Board{
		size:   int(size),
		pieces: make([]*Piece, size*size),
		tmp:    make([]*Piece, size*size),
		img:    size.image(),
	}

	b.ko = b.simpleKo
	if superko {
		b.ko = b.superKo
	}

	return &b, nil
}

// Returns the piece at the specified coordinates.
func (b *Board) At(x, y int) *Piece {
	return b.pieces[(y*b.size)+x]
}

// Places a piece at the specified coordinates without running any
// rule checks.
func (b *Board) place(x, y int, p *Piece) {
	b.pieces[(y*b.size)+x] = p
}

// Checks the simple ko rule. Returns true if the rule has been
// violated.
func (b *Board) simpleKo() bool {
	prev := b.getTurn(-2)
	if prev == nil {
		return false
	}

	for i := range b.pieces {
		if b.pieces[i] != prev[i] {
			return false
		}
	}

	return true
}

// Checks the super ko rule. Returns true if the rule has been
// violated.
func (b *Board) superKo() bool {
	checkTurn := func(prev []*Piece) bool {
		for i := range b.pieces {
			if b.pieces[i] != prev[i] {
				return false
			}
		}

		return true
	}

	for i := range b.turns {
		prev := b.getTurn(i)
		if prev == nil {
			continue
		}

		if checkTurn(prev) {
			return true
		}
	}

	return false
}

// Checks for ko. Returns true if the rule has been violated.
func (b *Board) checkKo() bool {
	return b.ko()
}

// Recursively checks the liberties of a piece and the neibourghing
// pieces of the same color. Returns a slice of the coordinates of
// pieces that need to be removed, or nil if it finds empty liberties.
func (b *Board) checkLib(x, y int) [][2]int {
	var checked [][2]int
	inchecked := func(x, y int) bool {
		for _, v := range checked {
			if (v[0] == x) && (v[1] == y) {
				return true
			}
		}

		return false
	}

	var hasfree func(int, int) bool
	hasfree = func(x, y int) bool {
		p := b.At(x, y)
		if p == nil {
			return true
		}

		if inchecked(x, y) {
			return false
		}

		checked = append(checked, [...]int{x, y})

		var ret bool

		if y > 0 {
			up := b.At(x, y-1)
			if !ret && (up == p || up == nil) {
				ret = ret || hasfree(x, y-1)
			}
		}

		if y < b.size-1 {
			down := b.At(x, y+1)
			if !ret && (down == p || down == nil) {
				ret = ret || hasfree(x, y+1)
			}
		}

		if x > 0 {
			left := b.At(x-1, y)
			if !ret && (left == p || left == nil) {
				ret = ret || hasfree(x-1, y)
			}
		}

		if x < b.size-1 {
			right := b.At(x+1, y)
			if !ret && (right == p || right == nil) {
				ret = ret || hasfree(x+1, y)
			}
		}

		return ret
	}

	if !hasfree(x, y) {
		return checked
	}

	return nil
}

// Attempts to place the specified piece at the specified coordinates,
// checking whether or not it's a legal move. Also checks the
// surrounding pieces to see if they've been captured. If they have,
// it removes them. Also sets player one. Returns true if the piece
// was placed, and false if it wasn't.
func (b *Board) Place(x, y int, p *Piece) (ret int) {
	copy(b.tmp, b.pieces)
	defer func() {
		if ret < 0 {
			copy(b.pieces, b.tmp)
			b.turns = b.turns[:len(b.turns)-1]
		} else {
			chng := &b.turns[len(b.turns)-1]
			*chng = append(*chng, Placement{
				p:   p,
				loc: [...]int{x, y},
			})
		}
	}()

	b.turns = append(b.turns, nil)

	if (x < 0) || (x > b.size-1) || (y < 0) || (y > b.size-1) || (b.At(x, y) != nil) {
		return -1
	}

	b.place(x, y, p)

	if (x > 0) && (b.At(x-1, y) != p) {
		if c := b.checkLib(x-1, y); c != nil {
			for _, v := range c {
				b.remove(v[0], v[1])
				ret++
			}
		}
	}
	if (x < b.size-1) && (b.At(x+1, y) != p) {
		if c := b.checkLib(x+1, y); c != nil {
			for _, v := range c {
				b.remove(v[0], v[1])
				ret++
			}
		}
	}
	if (y > 0) && (b.At(x, y-1) != p) {
		if c := b.checkLib(x, y-1); c != nil {
			for _, v := range c {
				b.remove(v[0], v[1])
				ret++
			}
		}
	}
	if (y < b.size-1) && (b.At(x, y+1) != p) {
		if c := b.checkLib(x, y+1); c != nil {
			for _, v := range c {
				b.remove(v[0], v[1])
				ret++
			}
		}
	}

	if (b.checkLib(x, y) != nil) || b.checkKo() {
		return -1
	}

	return
}

// Removes a piece from the board, updating the capture scores.
func (b *Board) remove(x, y int) {
	b.place(x, y, nil)

	chng := &b.turns[len(b.turns)-1]
	*chng = append(*chng, Placement{
		p:   nil,
		loc: [...]int{x, y},
	})
}

// Returns the board at the given turn, or nil if there's a problem.
func (b *Board) getTurn(num int) []*Piece {
	if num < 0 {
		num += len(b.turns)
	}

	if (num < 0) || (num >= len(b.turns)) {
		return nil
	}

	t := make([]*Piece, b.size*b.size)

	for i := range num {
		turn := b.turns[i]
		for _, v := range turn {
			t[(v.loc[1]*b.size)+v.loc[0]] = v.p
		}
	}

	return t
}

// Converts board coordinates to on-screen coordinates.
func (b *Board) CoordToXY(x, y int) (int, int) {
	imgBounds := b.img.Bounds()
	if (x > imgBounds.Dx()) || (y > imgBounds.Dy()) {
		return -1, -1
	}
	return b.Size().scaleUp(x, y)
}

// Converts on-screen coordinates to board coordinates.
func (b *Board) XYToCoord(x, y int) (int, int) {
	imgBounds := b.img.Bounds()
	if (x > imgBounds.Dx()) || (y > imgBounds.Dy()) {
		return -1, -1
	}
	return b.Size().scaleDown(x, y)
}

// Draws the specified piece at the specified coordinates.
func (b *Board) drawPiece(dst *ebiten.Image, x, y int, p *Piece) {
	x, y = b.CoordToXY(x, y)

	pimg := p.Image()
	pbounds := pimg.Bounds()
	x -= pbounds.Dx() / 2
	y -= pbounds.Dy() / 2

	var geom ebiten.GeoM
	geom.Translate(float64(x), float64(y))
	dst.DrawImage(pimg, &ebiten.DrawImageOptions{GeoM: geom})
}

func (b *Board) Draw(dst *ebiten.Image, geom ebiten.GeoM) {
	dst.DrawImage(b.img, &ebiten.DrawImageOptions{
		GeoM: geom,
	})

	for y := range b.size {
		for x := range b.size {
			if p := b.At(x, y); p != nil {
				b.drawPiece(dst, x, y, p)
			}
		}
	}
}

// Returns the board's size.
func (b *Board) Size() BoardSize {
	return BoardSize(b.size)
}

// Places the specified piece at the locations specified by the
// handicap.
func (b *Board) ApplyHandicap(p *Piece, h Handicap) {
	for _, v := range h {
		b.place(v[0], v[1], p)
	}
}
