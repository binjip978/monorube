package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

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

func boarding(boardPass string) int {
	first := boardPass[:7]
	second := boardPass[7:]

	cf := 0
	power := 0
	for i := 6; i >= 0; i-- {
		if first[i] == 'B' {
			cf += int(math.Pow(2.0, float64(power)))
		}
		power++
	}

	rf := 0
	power = 0
	for i := 2; i >= 0; i-- {
		if second[i] == 'R' {
			rf += int(math.Pow(2.0, float64(power)))
		}
		power++
	}

	return cf*8 + rf
}

func maxBoardPass(filename string) int {
	lines := readLines(filename)
	max := 0
	for _, line := range lines {
		r := boarding(line)
		if r > max {
			max = r
		}
	}

	return max
}

func boardIDs(filename string) {
	lines := readLines(filename)
	var res []int
	for _, line := range lines {
		r := boarding(line)
		res = append(res, r)
	}
	sort.Ints(res)
	prev := res[0]
	for i := 1; i < len(res); i++ {
		if res[i]-1 != prev {
			fmt.Println(res[i] - 1)
			return
		}
		prev = res[i]
	}

	fmt.Println("None")
}

// problem 6a/b
func customs(filename string) int {
	// var merge = func(line string, set map[rune]bool) map[rune]bool {
	// 	for _, c := range line {
	// 		set[c] = true
	// 	}

	// 	return set
	// }

	var create = func(line string) map[rune]bool {
		set := make(map[rune]bool)

		for _, c := range line {
			set[c] = true
		}

		return set
	}

	var intersect = func(set1 map[rune]bool, set2 map[rune]bool) map[rune]bool {
		nm := make(map[rune]bool)

		for k, ok1 := range set1 {
			_, ok2 := set2[k]
			if ok1 && ok2 {
				nm[k] = true
			}
		}

		return nm
	}

	lines := readLines(filename)
	var cnt int
	first := true
	set := make(map[rune]bool)

	for _, line := range lines {
		if line == "" {
			for _, v := range set {
				if v {
					cnt++
				}
			}
			set = make(map[rune]bool)
			first = true
			continue
		}
		newSet := create(line)
		if first {
			set = newSet
			first = false
		} else {
			set = intersect(set, newSet)
		}
	}

	for _, v := range set {
		if v {
			cnt++
		}
	}

	return cnt
}

// Problem 7 a/b

type bagEdge struct {
	to     string
	amount int
}

type bagGraph struct {
	g map[string][]bagEdge
}

func parseBagString(line string) (string, []bagEdge) {
	var edges []bagEdge
	s1 := strings.Split(line, "bags contain ")
	from := strings.TrimSpace(s1[0])
	if strings.Contains(s1[1], "no other bags") {
		return from, edges
	}

	s2 := strings.Split(s1[1], ", ")

	for _, s3 := range s2 {
		s4 := strings.Split(s3, " ")
		var edge bagEdge
		c, _ := strconv.Atoi(s4[0])
		edge.amount = c
		edge.to = strings.TrimSpace(s4[1] + " " + s4[2])
		edges = append(edges, edge)
	}

	return from, edges
}

func buildBagGraph(filename string) bagGraph {
	bg := bagGraph{
		g: make(map[string][]bagEdge),
	}

	lines := readLines(filename)
	for _, line := range lines {
		from, egs := parseBagString(line)
		edges, ok := bg.g[from]
		if !ok {
			bg.g[from] = egs
			continue
		}
		edges = append(edges, egs...)
		bg.g[from] = edges
	}

	return bg
}

func buildInverseBagGraph(filename string) bagGraph {
	bg := bagGraph{
		g: make(map[string][]bagEdge),
	}

	lines := readLines(filename)
	for _, line := range lines {
		from, egs := parseBagString(line)
		for _, e := range egs {
			edges, ok := bg.g[e.to]
			if !ok {
				bg.g[e.to] = []bagEdge{{from, e.amount}}
				continue
			}
			edges = append(edges, bagEdge{from, e.amount})
			bg.g[e.to] = edges
		}
	}

	return bg
}

func outmostShinyGold(g bagGraph) int {
	frontier := []string{"shiny gold"}

	visited := make(map[string]bool)

	for len(frontier) != 0 {
		f := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		ff, _ := g.g[f]
		for _, n := range ff {
			_, ok := visited[n.to]
			if ok {
				continue
			}
			frontier = append(frontier, n.to)
		}
		visited[f] = true
	}

	return len(visited) - 1
}

func bagGraphContains(g bagGraph, node string) int {
	edges := g.g[node]
	if len(edges) == 0 {
		return 0
	}

	var cnt int
	for _, e := range edges {
		cnt = cnt + e.amount + e.amount*bagGraphContains(g, e.to)
	}

	return cnt
}

// problem 8

type instr struct {
	op      string
	operand int
}

func code(filename string) []instr {
	lines := readLines(filename)
	var c []instr
	for _, line := range lines {
		sp := strings.Split(line, " ")
		v, _ := strconv.Atoi(sp[1])
		c = append(c, instr{sp[0], v})
	}

	return c
}

func execute(cs []instr) (int, bool) {
	var acc int
	var ip int
	visisted := make(map[int]bool)

	for ip != len(cs) {
		_, ok := visisted[ip]
		if ok {
			return acc, false
		}

		next := cs[ip]
		visisted[ip] = true
		switch next.op {
		case "nop":
			ip++
		case "acc":
			acc += next.operand
			ip++
		case "jmp":
			ip += next.operand
		default:
			panic("AAAAAAA!!!!!")
		}
	}

	return acc, true
}

func fixBootloader(cs []instr) int {
	for i, ins := range cs {
		if ins.op == "jmp" {
			css := make([]instr, len(cs))
			copy(css, cs)
			css[i].op = "nop"
			r, s := execute(css)
			if s {
				return r
			}
			continue
		}
		if ins.op == "nop" {
			css := make([]instr, len(cs))
			copy(css, cs)
			css[i].op = "jmp"
			r, s := execute(css)
			if s {
				return r
			}
			continue

		}
	}

	panic("AAAA")
	return 42
}

// problem 9
func nonSum(size int, filename string) int {
	ns := readList(filename)
	m := make(map[int]bool)

	for i := 0; i < size; i++ {
		m[ns[i]] = true
	}

	for i := size; i < len(ns); i++ {
		cand := ns[i]
		var hit bool
		for j := i - size; j < i; j++ {
			_, ok := m[cand-ns[j]]
			if ok {
				hit = true
				continue
			}
		}
		if !hit {
			return cand
		}
		d := ns[i-size]
		delete(m, d)
		m[cand] = false
	}

	panic("AAAAA!!! can't break it!!!")
}

func contiguousSum(filename string, cand int) int {
	ns := readList(filename)
	for i := 0; i < len(ns); i++ {
		var acc int
		for j := i; j < len(ns); j++ {
			if acc == cand {
				min := ns[i]
				max := ns[i]
				for z := i + 1; z < j; z++ {
					if ns[z] < min {
						min = ns[z]
					}
					if ns[z] > max {
						max = ns[z]
					}
				}

				return min + max
			}
			if acc > cand {
				break
			}
			acc += ns[j]
		}
	}

	panic("AAAAAA!!!")
}

// problem 10

func diff(filename string) int {
	jolts := readList(filename)
	sort.Ints(jolts)
	var curr, cnt1, cnt3 int

	for _, j := range jolts {
		d := j - curr
		if d == 1 {
			cnt1++
		}
		if d == 3 {
			cnt3++
		}
		curr = j
	}

	cnt3++ // device

	return cnt1 * cnt3
}

func arrangements(filename string) int {
	jolts := readList(filename)
	sort.Ints(jolts)
	jolts = append([]int{0}, jolts...)
	jolts = append(jolts, jolts[len(jolts)-1]+3)
	memo := make(map[int]int)

	return rec(0, jolts, memo)
}

func rec(i int, jolts []int, memo map[int]int) int {
	r, ok := memo[i]
	if ok {
		return r
	}

	if i == len(jolts)-1 {
		return 1
	}
	var cnt int

	for j := 1; j <= 3; j++ {
		if i+j < len(jolts) && jolts[i+j]-jolts[i] <= 3 {
			cnt += rec(i+j, jolts, memo)
		}
	}

	memo[i] = cnt
	return cnt
}

// problem 11

type SeatModel struct {
	grid   [][]rune
	width  int
	height int
}

func (s *SeatModel) String() string {
	var buf bytes.Buffer
	for _, line := range s.grid {
		for _, s := range line {
			buf.WriteRune(s)
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

type pos struct {
	x int
	y int
}

func (s *SeatModel) move() bool {
	var cp [][]rune
	for _, l := range s.grid {
		c := make([]rune, len(l))
		copy(c, l)
		cp = append(cp, c)
	}

	var change bool

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			switch s.grid[y][x] {
			case 'L':
				n := s.pos1(x, y)
				occ := false
				for _, n1 := range n {
					if s.grid[n1.y][n1.x] == '#' {
						occ = true
						break
					}
				}
				if !occ {
					change = true
					cp[y][x] = '#'
				}
			case '#':
				n := s.pos1(x, y)
				occCnt := 0
				for _, n1 := range n {
					if s.grid[n1.y][n1.x] == '#' {
						occCnt++
					}
				}
				if occCnt >= 5 {
					change = true
					cp[y][x] = 'L'
				}
			}
		}
	}

	s.grid = cp
	return change
}

func (s *SeatModel) pos(x int, y int) []pos {
	// x width y height
	var in = func(xi int, yi int) bool {
		return xi >= 0 && xi < s.width && yi >= 0 && yi < s.height
	}

	var res []pos

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			if in(x+i, y+j) {
				res = append(res, pos{x + i, y + j})
			}
		}
	}

	return res
}

func (s *SeatModel) pos1(x int, y int) []pos {
	// x width y height
	var in = func(xi int, yi int) bool {
		return xi >= 0 && xi < s.width && yi >= 0 && yi < s.height
	}

	var res []pos
	var cx, cy int

	// 1) -1, -1
	cx = x
	cy = y
	for {
		cx = cx - 1
		cy = cy - 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	// 2) -1, 0
	cx = x
	cy = y
	for {
		cx = cx - 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	// 3) -1, +1
	cx = x
	cy = y
	for {
		cx = cx - 1
		cy = cy + 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	// 4) 0, -1
	cx = x
	cy = y
	for {
		cy = cy - 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	// 5) 0, +1
	cx = x
	cy = y
	for {
		cy = cy + 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	// 6) 1, -1
	cx = x
	cy = y
	for {
		cx = cx + 1
		cy = cy - 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	// 7) 1, 0
	cx = x
	cy = y
	for {
		cx = cx + 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	// 8) +1, +1
	cx = x
	cy = y
	for {
		cx = cx + 1
		cy = cy + 1

		if !in(cx, cy) {
			break
		}

		if s.grid[cy][cx] == 'L' || s.grid[cy][cx] == '#' {
			res = append(res, pos{cx, cy})
			break
		}
	}

	return res
}

func NewSeatModel(filename string) *SeatModel {
	var grid [][]rune
	lines := readLines(filename)
	for _, line := range lines {
		var g []rune
		for _, c := range line {
			g = append(g, c)
		}
		grid = append(grid, g)
	}

	return &SeatModel{
		grid:   grid,
		height: len(grid),
		width:  len(grid[0]),
	}
}

func (s *SeatModel) simulate() int {
	for s.move() {
	}

	var cnt int
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			if s.grid[y][x] == '#' {
				cnt++
			}
		}
	}

	return cnt
}

func main() {
	g := NewSeatModel("input/11.txt")
	fmt.Println(g.simulate())

}
