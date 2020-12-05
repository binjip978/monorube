package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type grid struct {
	g      [][]rune
	width  int
	height int
}

func buildGrid(filename string) grid {
	var g grid

	lines := readLines(filename)
	for _, line := range lines {
		var row []rune
		for _, c := range line {
			row = append(row, c)
		}

		g.g = append(g.g, row)
	}

	g.width = len(g.g[0])
	g.height = len(g.g)
	return g
}

func (g grid) String() string {
	var buf bytes.Buffer

	for _, row := range g.g {
		for _, s := range row {
			buf.WriteRune(s)
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

// problem 3a
func countTrees(m grid, wd int, hd int) int {
	cw, ch := 0, 0
	var cnt int

	for ch < m.height {
		if m.g[ch][cw] == '#' {
			cnt++
		}
		cw = (cw + wd) % m.width
		ch += hd

	}

	return cnt
}

type slope struct {
	right int
	down  int
}

func countTrees2(m grid, slopes []slope) int {
	prod := 1

	for _, slope := range slopes {
		prod *= countTrees(m, slope.right, slope.down)
	}

	return prod
}

// problem 1a
func sum2(xs []int) int {
	for i := 0; i < len(xs); i++ {
		for j := i + 1; j < len(xs); j++ {
			if xs[i]+xs[j] == 2020 {
				return xs[i] * xs[j]
			}
		}
	}

	return -1
}

// problem 1b
// hashing can be much faster
func sum3(xs []int) int {
	for i := 0; i < len(xs); i++ {
		for j := i + 1; j < len(xs); j++ {
			for k := j + 1; k < len(xs); k++ {
				if xs[i]+xs[j]+xs[k] == 2020 {
					return xs[i] * xs[j] * xs[k]
				}
			}
		}
	}

	return -1

}

func tobogganSplit(str string) (int, int, string, string) {
	sp := strings.Split(str, " ")
	csp := strings.Split(sp[0], "-")
	l, _ := strconv.Atoi(csp[0])
	h, _ := strconv.Atoi(csp[1])
	return l, h, sp[1][:1], sp[2]
}

func validPasswords(xs []string) int {
	var cnt int

	for _, x := range xs {
		l, h, sym, pass := tobogganSplit(x)
		c := strings.Count(pass, sym)
		if c >= l && c <= h {
			cnt++
		}
	}

	return cnt
}

func validPasswords2(xs []string) int {
	var cnt int

	for _, x := range xs {
		l, h, sym, pass := tobogganSplit(x)
		if (sym == pass[l-1:l] && sym != pass[h-1:h]) || (sym != pass[l-1:l] && sym == pass[h-1:h]) {
			cnt++
		}
	}

	return cnt
}

func readList(filename string) []int {
	b, _ := ioutil.ReadFile(filename)
	r := bufio.NewReader(bytes.NewReader(b))
	var res []int
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		i, err := strconv.Atoi(s[:len(s)-1])
		res = append(res, i)
	}

	return res
}

func readLines(filename string) []string {
	b, _ := ioutil.ReadFile(filename)
	r := bufio.NewReader(bytes.NewReader(b))
	var res []string
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		res = append(res, s[:len(s)-1])
	}

	return res
}

func main() {
	g := buildGrid("input/3.txt")
	fmt.Println(countTrees2(g, []slope{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}}))
}
