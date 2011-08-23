package main

import (
	"os"
	"fmt"
	"sdl"
	"path"
	"runtime"
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

// Stores board information. Also calculates score and keeps track of
// the basic rules.
type Board struct {
	size   int
	pieces []*Piece
	prev [2][]*Piece
	tmp []*Piece

	bg  *sdl.Surface
	img *sdl.Surface

	p1 *Piece

	p1cap int
	p2cap int

	komi float64
}

// Initializes a new board of the given size.
func NewBoard(size BoardSize) (*Board, os.Error) {
	b := new(Board)

	b.size = int(size)
	b.pieces = make([]*Piece, b.size*b.size)
	for i := range(b.prev) {
		b.prev[i] = make([]*Piece, b.size*b.size)
	}
	b.tmp = make([]*Piece, b.size*b.size)

	b.bg = sdl.Load(path.Join(BoardPath, fmt.Sprintf("%v.png", b.size)))
	if b.bg == nil {
		return nil, os.NewError(sdl.GetError())
	}

	//b.img = sdl.CreateRGBSurface(sdl.HWSURFACE,
	//	int(b.bg.W),
	//	int(b.bg.H),
	//	int(b.bg.Format.BitsPerPixel),
	//	b.bg.Format.Rmask,
	//	b.bg.Format.Gmask,
	//	b.bg.Format.Bmask,
	//	b.bg.Format.Amask,
	//)
	//if b.img == nil {
	//	return nil, os.NewError(sdl.GetError())
	//}

	b.img = sdl.Load(path.Join(BoardPath, fmt.Sprintf("%v.png", b.size)))
	if b.img == nil {
		return nil, os.NewError(sdl.GetError())
	}

	runtime.SetFinalizer(b, (*Board).free)

	return b, nil
}

// Frees up the boards resources.
func (b *Board) free() {
	b.bg.Free()
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

// Checks whether or not the simple ko rule has been violated.
func (b *Board)checkKo() bool {
	for i := range(b.pieces) {
		if b.pieces[i] != b.prev[1][i] {
			return false
		}
	}

	return true
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
func (b *Board) Place(x, y int, p *Piece) (ret bool) {
	copy(b.tmp, b.pieces)
	defer func() {
		if ret == false {
			copy(b.pieces, b.tmp)
		}
	}()

	if (x < 0) || (x > b.size-1) || (y < 0) || (y > b.size-1) {
		return false
	}

	if b.At(x, y) != nil {
		return false
	}

	b.place(x, y, p)

	if (x > 0) && (b.At(x-1, y) != p) {
		if c := b.checkLib(x-1, y); c != nil {
			for _, v := range c {
				b.Remove(v[0], v[1])
			}
		}
	}
	if (x < b.size-1) && (b.At(x+1, y) != p) {
		if c := b.checkLib(x+1, y); c != nil {
			for _, v := range c {
				b.Remove(v[0], v[1])
			}
		}
	}
	if (y > 0) && (b.At(x, y-1) != p) {
		if c := b.checkLib(x, y-1); c != nil {
			for _, v := range c {
				b.Remove(v[0], v[1])
			}
		}
	}
	if (y < b.size-1) && (b.At(x, y+1) != p) {
		if c := b.checkLib(x, y+1); c != nil {
			for _, v := range c {
				b.Remove(v[0], v[1])
			}
		}
	}

	if (b.checkLib(x, y) != nil) || b.checkKo() {
		return false
	}

	copy(b.prev[1], b.prev[0])
	copy(b.prev[0], b.pieces)

	if b.p1 == nil {
		b.p1 = p
	}

	return true
}

// Removes a piece from the board, updating the capture scores.
func (b *Board) Remove(x, y int) {
	p := b.At(x, y)

	switch p {
	case nil:
	case b.p1:
		b.p2cap++
	default:
		b.p1cap++
	}

	b.place(x, y, nil)
}

// Converts board coordinates to on-screen coordinates.
func (b *Board) CoordToXY(x, y int) (int, int) {
	if (x > int(b.img.W)) || (y > int(b.img.H)) {
		return -1, -1
	}

	switch BoardSize(b.size) {
	case Size9x9:
		x = (x * 52) + 31
		y = (y * 52) + 31
	case Size19x19:
		x = (x * 25) + 14
		y = (y * 25) + 14
	}

	return x, y
}

// Converts on-screen coordinates to board coordinates.
func (b *Board) XYToCoord(x, y int) (int, int) {
	if (x > int(b.img.W)) || (y > int(b.img.H)) {
		return -1, -1
	}

	switch BoardSize(b.size) {
	case Size9x9:
		x /= 52
		y /= 52
	case Size19x19:
		x /= 25
		y /= 25
	}

	return x, y
}

// Draws the specified piece at the specified coordinates.
func (b *Board) drawPiece(x, y int, p *Piece) {
	x, y = b.CoordToXY(x, y)

	pimg := p.Image()

	x -= int(pimg.W / 2)
	y -= int(pimg.H / 2)

	b.img.Blit(&sdl.Rect{X: int16(x), Y: int16(y)}, pimg, nil)
}

// Returns the board's image, complete with all the pieces drawn onto
// it.
func (b *Board) Image() *sdl.Surface {
	b.img.FillRect(nil, sdl.MapRGB(b.img.Format, 0, 0, 0))

	b.img.Blit(nil, b.bg, nil)

	for y := 0; y < b.size; y++ {
		for x := 0; x < b.size; x++ {
			if p := b.At(x, y); p != nil {
				b.drawPiece(x, y, p)
			}
		}
	}

	return b.img
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

// Gives the appropriate komi.
func (b *Board) GiveKomi(komi float64) {
	b.komi += komi
}
