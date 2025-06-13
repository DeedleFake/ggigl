package main

import "fmt"

// A type that represents the spaces that pieces are placed on at
// specific levels of handicap.
type Handicap [][2]int

var (
	// Handicap information for a 9x9 board.
	handicaps9 Handicap

	// Handicap information for a 19x19 board.
	handicaps19 Handicap
)

func init() {
	handicaps9 = [][2]int{
		[2]int{6, 2},
		[2]int{2, 6},
		[2]int{6, 6},
		[2]int{2, 2},
		[2]int{4, 4},
		[2]int{2, 4},
		[2]int{6, 4},
		[2]int{4, 2},
		[2]int{4, 6},
	}

	handicaps19 = [][2]int{
		[2]int{15, 3},
		[2]int{3, 15},
		[2]int{15, 15},
		[2]int{3, 3},
		[2]int{9, 9},
		[2]int{3, 9},
		[2]int{15, 9},
		[2]int{9, 3},
		[2]int{9, 15},
	}
}

// Returns the Handicap representing the given level for the specified
// board.
func GetHandicap(size BoardSize, num int) (Handicap, error) {
	switch size {
	case Size9x9, Size19x19:
	default:
		return nil, fmt.Errorf("handicaps not supported for board size: %v", size)
	}

	if num > MaxHandicap(size) {
		return nil, fmt.Errorf("Handicap exceeds max for board size: %v", num)
	}

	switch size {
	case Size9x9:
		return handicaps9[:num], nil
	case Size19x19:
		return handicaps19[:num], nil
	}

	panic("This should never reach this point...")
}

// Returns the maximum handicap available for the specified board
// size.
func MaxHandicap(size BoardSize) int {
	switch size {
	case Size9x9:
		return len(handicaps9)
	case Size19x19:
		return len(handicaps19)
	}

	return -1
}
