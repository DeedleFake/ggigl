package main

import (
	"os"
	"fmt"
	"sdl"
	"path"
	"runtime"
)

var (
	BoardPath = path.Join("data", "boards")
)

type BoardSize int

const (
	Size9x9   BoardSize = 9
	Size19x19 BoardSize = 19
)

type Board struct {
	size   int
	pieces []*Piece

	bg  *sdl.Surface
	img *sdl.Surface

	p1 *Piece

	p1cap float64
	p2cap float64
}

func NewBoard(size BoardSize) (*Board, os.Error) {
	b := new(Board)

	b.size = int(size)
	b.pieces = make([]*Piece, b.size*b.size)

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

func (b *Board) free() {
	b.bg.Free()
}

func (b *Board) At(x, y int) *Piece {
	return b.pieces[(y*b.size)+x]
}

func (b *Board) place(x, y int, p *Piece) {
	b.pieces[(y*b.size)+x] = p
}

func (b *Board)checkLib(x, y int) [][2]int {
	var checked [][2]int
	inchecked := func(x, y int) bool {
		for _, v := range(checked) {
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

		up := b.At(x, y - 1)
		down := b.At(x, y + 1)
		left := b.At(x - 1, y)
		right := b.At(x + 1, y)

		if (up == nil) || (down == nil) || (left == nil) || (right == nil) {
			return true
		}

		checked = append(checked, [2]int{x, y})

		if ((up != p) || inchecked(x, y - 1)) && ((down != p) || inchecked(x, y + 1)) && ((left != p) || inchecked(x - 1, y)) && ((right != p) || inchecked(x + 1, y)) {
			return false
		}

		var ret bool
		if (up == p) && (!inchecked(x, y - 1)) {
			ret = hasfree(x, y - 1)
		}
		if ((down == p) && (!inchecked(x, y + 1))) || !ret {
			ret = hasfree(x, y + 1)
		}
		if ((left == p) && (!inchecked(x - 1, y))) || !ret {
			ret = hasfree(x - 1, y)
		}
		if ((right == p) && (!inchecked(x + 1, y))) || !ret {
			ret = hasfree(x + 1, y)
		}

		return ret
	}

	if !hasfree(x, y) {
		return checked
	}

	return nil
}

func (b *Board) Place(x, y int, p *Piece) bool {
	if b.At(x, y) != nil {
		return false
	}

	b.place(x, y, p)

	if b.checkLib(x, y) != nil {
		b.place(x, y, nil)
		return false
	}

	if c := b.checkLib(x - 1, y); c != nil {
		for _, v := range(c) {
			b.Remove(v[0], v[1])
		}
	}
	if c := b.checkLib(x + 1, y); c != nil {
		for _, v := range(c) {
			b.Remove(v[0], v[1])
		}
	}
	if c := b.checkLib(x, y - 1); c != nil {
		for _, v := range(c) {
			b.Remove(v[0], v[1])
		}
	}
	if c := b.checkLib(x, y + 1); c != nil {
		for _, v := range(c) {
			b.Remove(v[0], v[1])
		}
	}

	if b.p1 == nil {
		b.p1 = p
	}

	return true
}

func (b *Board)Remove(x, y int) {
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

func (b *Board) CoordToXY(x, y int) (int, int) {
	switch BoardSize(b.size) {
	case Size19x19:
		x = (x * 25) + 14
		y = (y * 25) + 14
	}

	return x, y
}

func (b *Board) drawPiece(x, y int, p *Piece) {
	x, y = b.CoordToXY(x, y)

	pimg := p.Image()

	x -= int(pimg.W / 2)
	y -= int(pimg.H / 2)

	b.img.Blit(&sdl.Rect{X: int16(x), Y: int16(y)}, pimg, nil)
}

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

func (b *Board) Size() BoardSize {
	return BoardSize(b.size)
}

func (b *Board) ApplyHandicap(p *Piece, h Handicap) {
	for _, v := range h {
		b.place(v[0], v[1], p)
	}
}
