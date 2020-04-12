package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
		w0	   w1	  w2	 w3
	h0	(0, 0) (0, 1) (0, 2) (0, 3)
	h1	(1, 0) (1, 1) (1, 2) (1, 3)
	h2	(2, 0) (2, 1) (2, 2) (2, 3)
*/

type life struct {
	width  int
	height int
	grid   [][]bool
}

type coord struct {
	h, w int
}

func newGrid(width int, height int) [][]bool {
	grid := make([][]bool, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]bool, width)
	}

	return grid
}

// New creates new game of life
func new(width int, height int, prob float64) *life {
	grid := newGrid(width, height)

	// init seed
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if r.Float64() < prob {
				grid[h][w] = true
			}
		}
	}

	return &life{
		width:  width,
		height: height,
		grid:   grid,
	}
}

func (l *life) print() {
	fmt.Println(strings.Repeat("-", l.width))
	for h := 0; h < l.height; h++ {
		for w := 0; w < l.width; w++ {
			if l.grid[h][w] {
				fmt.Printf("%s", "o")
			} else {
				fmt.Printf("%s", " ")
			}
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", l.width))
}

func (l *life) inGrid(h, w int) bool {
	return w >= 0 && w < l.width && h >= 0 && h < l.height
}

func (l *life) closeCoord(h, w int) []coord {
	var candidates []coord
	pos := []int{-1, 0, 1}
	for _, p1 := range pos {
		for _, p2 := range pos {
			if p1 != 0 || p2 != 0 {
				c := coord{h + p1, w + p2}
				if l.inGrid(c.h, c.w) {
					candidates = append(candidates, c)
				}
			}
		}
	}

	return candidates
}

func (l *life) liveNeighbors(h, w int) []coord {
	cand := l.closeCoord(h, w)
	var nb []coord
	for _, n := range cand {
		if l.grid[n.h][n.w] {
			nb = append(nb, n)
		}
	}

	return nb
}

func (l *life) step() {
	nextGrid := newGrid(l.width, l.height)
	for h := 0; h < l.height; h++ {
		for w := 0; w < l.width; w++ {
			cnt := len(l.liveNeighbors(h, w))
			if l.grid[h][w] {
				if cnt == 2 || cnt == 3 {
					nextGrid[h][w] = true
				}
				if cnt < 2 {
					nextGrid[h][w] = false
				}
				if cnt > 3 {
					nextGrid[h][w] = false
				}
			} else {
				if cnt == 3 {
					nextGrid[h][w] = true
				}
			}
		}
	}
	l.grid = nextGrid
}

func main() {
	game := new(120, 26, 0.1)
	for i := 0; i < 120; i++ {
		game.print()
		game.step()
		time.Sleep(2 * time.Second)
	}
}
