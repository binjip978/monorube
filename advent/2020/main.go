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
	var cw, ch, cnt int

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

// problem 3b
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

type password struct {
	byr string
	iyr string
	eyr string
	hgt string
	hcl string
	ecl string
	pid string
	cid string
}

func scanPasswords(filename string) []password {
	var scanLine = func(line string, m map[string]string) map[string]string {
		sp := strings.Split(line, " ")
		for _, s := range sp {
			fs := strings.Split(s, ":")
			m[fs[0]] = fs[1]
		}

		return m
	}

	var pass = func(m map[string]string) password {
		return password{
			byr: m["byr"],
			iyr: m["iyr"],
			eyr: m["eyr"],
			hgt: m["hgt"],
			hcl: m["hcl"],
			ecl: m["ecl"],
			pid: m["pid"],
			cid: m["cid"],
		}
	}

	var passwords []password
	lines := readLines(filename)

	m1 := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			p := pass(m1)
			passwords = append(passwords, p)
			m1 = make(map[string]string)
			continue
		}
		m1 = scanLine(line, m1)
	}

	return passwords
}

func valid(pass password) bool {
	return pass.byr != "" && pass.ecl != "" && pass.eyr != "" && pass.hcl != "" &&
		pass.hgt != "" && pass.iyr != "" && pass.pid != ""
}

func valid2(pass password) bool {
	if pass.byr == "" {
		return false
	}
	b, _ := strconv.Atoi(pass.byr)
	if b < 1920 || b > 2002 {
		return false
	}

	if pass.iyr == "" {
		return false
	}
	i, _ := strconv.Atoi(pass.iyr)
	if i < 2010 || i > 2020 {
		return false
	}

	if pass.eyr == "" {
		return false
	}
	e, _ := strconv.Atoi(pass.eyr)
	if e < 2020 || e > 2030 {
		return false
	}

	if pass.hgt == "" {
		return false
	}
	if strings.Contains(pass.hgt, "cm") {
		sp := strings.Split(pass.hgt, "cm")
		h, _ := strconv.Atoi(sp[0])
		if h < 150 || h > 193 {
			return false
		}
	} else if strings.Contains(pass.hgt, "in") {
		sp := strings.Split(pass.hgt, "in")
		h, _ := strconv.Atoi(sp[0])
		if h < 59 || h > 76 {
			return false
		}
	} else {
		return false
	}

	if len(pass.hcl) != 7 || pass.hcl[0] != '#' {
		return false
	}

	for _, c := range pass.hcl[1:] {
		switch c {
		case 'a', 'b', 'c', 'd', 'e', 'f', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return false
		}
	}

	switch pass.ecl {
	case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
	default:
		return false
	}

	if len(pass.pid) != 9 {
		return false
	}

	for _, c := range pass.pid {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return false
		}
	}

	return true
}

func validCount(passwords []password) int {
	var cnt int
	for _, pass := range passwords {
		if valid2(pass) {
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
	passwords := scanPasswords("input/4.txt")
	fmt.Println(validCount(passwords))
}
