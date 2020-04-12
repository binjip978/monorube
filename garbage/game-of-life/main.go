package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type life struct {
	width  int // x
	height int // y
	grid   [][]bool
}

type coord struct {
	x, y int
}

// New creates new game of life
func new(width int, height int, prob float64) *life {
	grid := make([][]bool, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]bool, width)
	}

	// init seed
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if rand.Float64() < prob {
				grid[i][j] = true
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
	for i := 0; i < l.height; i++ {
		for j := 0; j < l.width; j++ {
			if l.grid[i][j] {
				fmt.Printf("%s", "o")
			} else {
				fmt.Printf("%s", " ")
			}
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", l.width))
}

func (l *life) inGrid(x, y int) bool {
	return x >= 0 && x < l.width && y >= 0 && y < l.height
}

func (l *life) closeCoord(x, y int) []coord {
	var candidates []coord
	pos := []int{-1, 0, 1}
	for _, p1 := range pos {
		for _, p2 := range pos {
			if p1 != 0 || p2 != 0 {
				c := coord{x + p1, y + p2}
				if l.inGrid(c.x, c.y) {
					candidates = append(candidates, c)
				}
			}
		}
	}

	return candidates
}

func (l *life) liveNeighbors(x, y int) []coord {
	// if x == 24 && y == 0 {
	// 	fmt.Println("BOOM")
	// }
	cand := l.closeCoord(x, y)
	// if x == 24 && y == 0 {
	// 	fmt.Println(l.width, l.height)
	// 	fmt.Println(x, y)
	// 	fmt.Println(cand)
	// }
	var nb []coord
	for _, n := range cand {
		if l.grid[n.y][n.x] {
			nb = append(nb, n)
		}
	}

	return nb
}

func (l *life) step() {
	for i := 0; i < l.height; i++ {
		for j := 0; j < l.width; j++ {
			cnt := len(l.liveNeighbors(i, j))
			if l.grid[i][j] {
				if cnt < 2 {
					l.grid[i][j] = false
				}
				if cnt > 3 {
					l.grid[i][j] = false
				}
			} else {
				if cnt == 3 {
					l.grid[i][j] = true
				}
			}
		}
	}
}

func main() {
	game := new(25, 25, 0.15)
	for i := 0; i < 30; i++ {
		game.print()
		game.step()
		time.Sleep(1 * time.Second)
	}
}
