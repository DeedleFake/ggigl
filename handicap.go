package main

type Handicap [][2]int

var handicaps Handicap

func init() {
	handicaps = [][2]int{
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

func GetHandicap(num int) Handicap {
	return handicaps[:num]
}

func MaxHandicap() int {
	return len(handicaps)
}
