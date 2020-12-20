package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"regexp"
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

// problem 12

type shipCMD struct {
	action string
	value  int
}

func parseShipCMD(filename string) []shipCMD {
	lines := readLines(filename)
	var cmds []shipCMD
	for _, line := range lines {
		v, _ := strconv.Atoi(line[1:])
		cmds = append(cmds, shipCMD{line[0:1], v})
	}

	return cmds
}

type cruiseShip struct {
	cmds []shipCMD
	dir  string
	we   int
	ns   int
}

func (s *cruiseShip) exec() float64 {
	for _, cmd := range s.cmds {
		switch cmd.action {
		case "N":
			s.ns += cmd.value
		case "S":
			s.ns -= cmd.value
		case "W":
			s.we += cmd.value
		case "E":
			s.we -= cmd.value
		case "L":
			r := cmd.value
			for r != 0 {
				switch s.dir {
				case "N":
					s.dir = "W"
				case "W":
					s.dir = "S"
				case "S":
					s.dir = "E"
				case "E":
					s.dir = "N"
				}
				r -= 90
			}
		case "R":
			r := cmd.value
			for r != 0 {
				switch s.dir {
				case "N":
					s.dir = "E"
				case "W":
					s.dir = "N"
				case "S":
					s.dir = "W"
				case "E":
					s.dir = "S"
				}
				r -= 90
			}
		case "F":
			switch s.dir {
			case "N":
				s.ns += cmd.value
			case "S":
				s.ns -= cmd.value
			case "W":
				s.we += cmd.value
			case "E":
				s.we -= cmd.value
			}
		}
	}

	return math.Abs(float64(s.we)) + math.Abs(float64(s.ns))
}

type cruiseShip2 struct {
	cmds         []shipCMD
	shipLocEast  int
	shipLocNorth int
	pointEast    int
	pointNorth   int
}

func (s *cruiseShip2) exec() float64 {
	for _, cmd := range s.cmds {
		switch cmd.action {
		case "N":
			s.pointNorth += cmd.value
		case "S":
			s.pointNorth -= cmd.value
		case "W":
			s.pointEast -= cmd.value
		case "E":
			s.pointEast += cmd.value
		case "L":
			switch cmd.value {
			case 90:
				s.pointEast, s.pointNorth = -s.pointNorth, s.pointEast
			case 180:
				s.pointEast, s.pointNorth = -s.pointEast, -s.pointNorth
			case 270:
				s.pointEast, s.pointNorth = s.pointNorth, -s.pointEast
			}
		case "R":
			switch cmd.value {
			case 90:
				s.pointEast, s.pointNorth = s.pointNorth, -s.pointEast
			case 180:
				s.pointEast, s.pointNorth = -s.pointEast, -s.pointNorth
			case 270:
				s.pointEast, s.pointNorth = -s.pointNorth, s.pointEast
			}
		case "F":
			s.shipLocEast += cmd.value * s.pointEast
			s.shipLocNorth += cmd.value * s.pointNorth
		}
	}

	return math.Abs(float64(s.shipLocEast)) + math.Abs(float64(s.shipLocNorth))
}

// problem 13

func firstBus(time int, periods []int) int {
	type rec struct {
		id        int
		startTime int
	}

	var recs []rec

	for _, p := range periods {
		d := time / p
		recs = append(recs, rec{p, d*p + p})
	}

	recMin := recs[0]

	for i := 1; i < len(recs); i++ {
		if recs[i].startTime < recMin.startTime {
			recMin = recs[i]
		}
	}

	return (recMin.startTime - time) * recMin.id
}

func prepareBus(filename string) (int, []int) {
	lines := readLines(filename)
	var ids []int
	sp := strings.Split(lines[1], ",")
	for _, s := range sp {
		if s == "x" {
			ids = append(ids, 0)
		}
		v, _ := strconv.Atoi(s)
		ids = append(ids, v)
	}

	t, _ := strconv.Atoi(lines[0])
	return t, ids
}

func timestampBF(ids []int) int {
	first := ids[0]
	for i := 500000000000000; i < math.MaxInt64; i++ {
		if i%1000000000 == 0 {
			fmt.Println(i)
		}
		if i%first == 0 {
			currTimestamp := i
			var fail bool
			for j := 1; j < len(ids); j++ {
				currTimestamp++
				if ids[j] == 0 {
					continue
				}
				if currTimestamp%ids[j] != 0 {
					fail = true
					break
				}
			}
			if !fail {
				return i
			}
		}
	}

	panic("AAAAA!!!")
}

// problem 14

func dockingSum(filename string) int {
	var mask string
	mem := make(map[int]int)

	lines := readLines(filename)
	for _, line := range lines {
		switch line[:3] {
		case "mas":
			sp := strings.Split(line, " = ")
			mask = sp[1]
		case "mem":
			sp := strings.Split(line, "] = ")
			sid := strings.Split(sp[0], "mem[")
			id, _ := strconv.Atoi(sid[1])
			value, _ := strconv.Atoi(sp[1])
			for i := 35; i >= 0; i-- {
				d := 35 - i
				if mask[i] == '1' {
					m := 1 << d
					value = value | m
				}
				if mask[i] == '0' {
					m := -1 ^ (1 << d)
					value = value & m
				}
			}
			mem[id] = value
		}
	}

	var res int
	for _, v := range mem {
		res += v
	}

	return res
}

func dockingSum2(filename string) int {
	var mask string
	mem := make(map[int64]int)

	lines := readLines(filename)
	for _, line := range lines {
		switch line[:3] {
		case "mas":
			sp := strings.Split(line, " = ")
			mask = sp[1]
		case "mem":
			sp := strings.Split(line, "] = ")
			sid := strings.Split(sp[0], "mem[")
			id, _ := strconv.Atoi(sid[1])
			value, _ := strconv.Atoi(sp[1])
			var b bytes.Buffer
			for i := 0; i < 36; i++ {
				switch mask[i] {
				case 'X':
					b.WriteRune('X')
				case '0':
					r := 1 & (id >> (35 - i))
					if r == 0 {
						b.WriteRune('0')
					} else {
						b.WriteRune('1')
					}
				case '1':
					b.WriteRune('1')
				}
			}

			memAddrs := allMem([]string{b.String()})

			for _, addr := range memAddrs {
				a, _ := strconv.ParseInt(addr, 2, 64)
				mem[a] = value
			}
		}
	}

	var res int
	for _, v := range mem {
		res += v
	}

	return res
}

func allMem(xs []string) []string {
	var xxs []string
	done := true
	for _, x := range xs {
		if strings.Contains(x, "X") {
			done = false
			i := strings.Index(x, "X")
			xxs = append(xxs, x[0:i]+"0"+x[i+1:])
			xxs = append(xxs, x[0:i]+"1"+x[i+1:])
		}
	}

	if !done {
		return allMem(xxs)
	}

	return xs
}

// problem 15

func elvesGame(numbers []int) int {
	named := make(map[int]int)
	for i := 0; i < len(numbers)-1; i++ {
		named[numbers[i]] = i + 1
	}
	lastNumber := numbers[len(numbers)-1]
	turn := len(numbers)

	for turn != 30000000 {
		p, ok := named[lastNumber]
		if ok {
			diff := turn - p
			named[lastNumber] = turn
			lastNumber = diff
		} else {
			named[lastNumber] = turn
			lastNumber = 0
		}
		turn++
	}

	return lastNumber
}

// problem 16

type interval struct {
	begin int
	end   int
}

type ticketClass struct {
	name string
	i1   interval
	i2   interval
}

func (t *ticketClass) in(value int) bool {
	i1 := value >= t.i1.begin && value <= t.i1.end
	i2 := value >= t.i2.begin && value <= t.i2.end
	return i1 || i2
}

func newTicketClass(line string) *ticketClass {
	s1 := strings.Split(line, ": ")
	s2 := strings.Split(s1[1], " or ")
	s3 := strings.Split(s2[0], "-")
	s4 := strings.Split(s2[1], "-")

	b1, _ := strconv.Atoi(s3[0])
	e1, _ := strconv.Atoi(s3[1])
	b2, _ := strconv.Atoi(s4[0])
	e2, _ := strconv.Atoi(s4[1])

	return &ticketClass{
		name: s1[0],
		i1:   interval{b1, e1},
		i2:   interval{b2, e2},
	}
}

type trainProblem struct {
	classes      []*ticketClass
	yourTickets  []int
	otherTickets [][]int
}

func newTrainProblem(filename string) *trainProblem {
	lines := readLines(filename)
	mode := "class"

	var classes []*ticketClass
	var your []int
	var other [][]int

	for _, line := range lines {
		switch line {
		case "":
			if mode == "your" {
				mode = "other"
			}
			if mode == "class" {
				mode = "your"
			}
		default:
			switch mode {
			case "class":
				classes = append(classes, newTicketClass(line))
			case "your":
				sp := strings.Split(line, ",")
				for _, s := range sp {
					v, _ := strconv.Atoi(s)
					your = append(your, v)
				}
			case "other":
				sp := strings.Split(line, ",")
				var ot []int
				for _, s := range sp {
					v, _ := strconv.Atoi(s)
					ot = append(ot, v)
				}
				other = append(other, ot)
			}
		}
	}

	return &trainProblem{classes, your, other}
}

func (t *trainProblem) errorRate() int {
	var errorRate int
	for _, ticket := range t.otherTickets {
		for _, v := range ticket {
			var hit bool
			for _, class := range t.classes {
				if class.in(v) {
					hit = true
					break
				}
			}
			if !hit {
				errorRate += v
			}
		}
	}

	return errorRate
}

// 1: to avoid bag in parse
func (t *trainProblem) filter() *trainProblem {
	var filtered [][]int

	for _, ticket := range t.otherTickets[1:] {
		goodTicket := true
		for _, v := range ticket {
			ok := false
			for _, class := range t.classes {
				if class.in(v) {
					ok = true
					break
				}
			}
			if !ok {
				goodTicket = false
				break
			}
		}
		if goodTicket {
			filtered = append(filtered, ticket)
		}
	}

	return &trainProblem{t.classes, t.yourTickets[1:], filtered}
}

func (t *trainProblem) mapping() int {
	cols := make(map[int][]int)
	for i := 0; i < len(t.otherTickets[0]); i++ {
		var c []int
		for j := 0; j < len(t.otherTickets); j++ {
			c = append(c, t.otherTickets[j][i])
		}
		cols[i] = c
	}

	cls := make(map[*ticketClass]bool)

	for _, c := range t.classes {
		cls[c] = true
	}

	res := make(map[string]int)

	for len(cls) != 0 {
		for class := range cls {
			var hc []int

			for i, col := range cols {
				all := true
				for _, v := range col {
					if !class.in(v) {
						all = false
						break
					}
				}
				if all {
					hc = append(hc, i)
				}
			}

			if len(hc) == 1 {
				res[class.name] = hc[0]
				delete(cls, class)
				delete(cols, hc[0])
			}
		}
	}

	cnt := 1
	for k, v := range res {
		if strings.Contains(k, "departure") {
			cnt *= t.yourTickets[v]
		}
	}

	return cnt
}

// problem 17

type zp struct {
	x int
	y int
	z int
}

func (p zp) nboors() []zp {
	var n []zp
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			for k := -1; k <= 1; k++ {
				if i == 0 && j == 0 && k == 0 {
					continue
				}
				n = append(n, zp{p.x + i, p.y + j, p.z + k})
			}
		}
	}

	return n
}

func initZP(filename string) []zp {
	var n []zp
	lines := readLines(filename)
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[0]); j++ {
			if lines[i][j] == '#' {
				n = append(n, zp{i, j, 0})
			}
		}
	}

	return n
}

func simulate(xs []zp) int {
	univ := make(map[zp]rune)
	for _, n := range xs {
		univ[n] = '#'
	}

	for i := 0; i < 6; i++ {
		cand := make(map[zp]bool)
		for z := range univ {
			nb := z.nboors()
			for _, b := range nb {
				cand[b] = true
			}
		}
		newUniv := make(map[zp]rune)

		for c := range cand {
			_, active := univ[c]
			if active {
				nb := c.nboors()
				var cnt int
				for _, n := range nb {
					_, ok := univ[n]
					if ok {
						cnt++
					}
				}
				if cnt == 2 || cnt == 3 {
					newUniv[c] = '#'
				}
			} else {
				nb := c.nboors()
				var cnt int
				for _, n := range nb {
					_, ok := univ[n]
					if ok {
						cnt++
					}
				}
				if cnt == 3 {
					newUniv[c] = '#'
				}
			}
		}

		univ = newUniv
	}

	return len(univ)
}

type zwp struct {
	x int
	y int
	z int
	w int
}

func (p zwp) nboors() []zwp {
	var n []zwp
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			for k := -1; k <= 1; k++ {
				for l := -1; l <= 1; l++ {
					if i == 0 && j == 0 && k == 0 && l == 0 {
						continue
					}
					n = append(n, zwp{p.x + i, p.y + j, p.z + k, p.w + l})
				}
			}
		}
	}

	return n
}

func initZWP(filename string) []zwp {
	var n []zwp
	lines := readLines(filename)
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[0]); j++ {
			if lines[i][j] == '#' {
				n = append(n, zwp{i, j, 0, 0})
			}
		}
	}

	return n
}

func simulate2(xs []zwp) int {
	univ := make(map[zwp]rune)
	for _, n := range xs {
		univ[n] = '#'
	}

	for i := 0; i < 6; i++ {
		cand := make(map[zwp]bool)
		for z := range univ {
			nb := z.nboors()
			for _, b := range nb {
				cand[b] = true
			}
		}
		newUniv := make(map[zwp]rune)

		for c := range cand {
			_, active := univ[c]
			if active {
				nb := c.nboors()
				var cnt int
				for _, n := range nb {
					_, ok := univ[n]
					if ok {
						cnt++
					}
				}
				if cnt == 2 || cnt == 3 {
					newUniv[c] = '#'
				}
			} else {
				nb := c.nboors()
				var cnt int
				for _, n := range nb {
					_, ok := univ[n]
					if ok {
						cnt++
					}
				}
				if cnt == 3 {
					newUniv[c] = '#'
				}
			}
		}

		univ = newUniv
	}

	return len(univ)
}

// problem 18

type lexem struct {
	operator rune
	value    int
}

func (l lexem) String() string {
	if l.operator == 0 {
		return fmt.Sprintf("%d", l.value)
	} else {
		return string(l.operator)
	}
}

func lexer(line string) []lexem {
	var res []lexem

	var b bytes.Buffer
	for _, c := range line {
		switch c {
		case '+':
			res = append(res, lexem{operator: '+'})
		case '*':
			res = append(res, lexem{operator: '*'})
		case ' ':
			if b.Len() != 0 {
				v, _ := strconv.Atoi(b.String())
				res = append(res, lexem{value: v})
				b.Reset()
			}
		case '(':
			res = append(res, lexem{operator: '('})
		case ')':
			if b.Len() != 0 {
				v, _ := strconv.Atoi(b.String())
				res = append(res, lexem{value: v})
				b.Reset()
			}
			res = append(res, lexem{operator: ')'})
		default:
			b.WriteRune(c)
		}
	}

	if b.Len() != 0 {
		v, _ := strconv.Atoi(b.String())
		res = append(res, lexem{value: v})
	}

	return res
}

type stack struct {
	s []lexem
}

func (s *stack) isEmpty() bool {
	return len(s.s) == 0
}

func (s *stack) push(l lexem) {
	s.s = append(s.s, l)
}

func (s *stack) peek() (lexem, error) {
	if s.isEmpty() {
		return lexem{}, fmt.Errorf("stack is empty")
	}

	return s.s[len(s.s)-1], nil
}

func (s *stack) pop() (lexem, error) {
	if s.isEmpty() {
		return lexem{}, fmt.Errorf("stack is empty")
	}

	elem := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]

	return elem, nil
}

func evaluate(lexems []lexem) int {
	s := stack{}
	for _, next := range lexems {
		if next.operator == 0 {
			prev, _ := s.peek()
			if prev.operator == '(' {
				s.push(next)
				continue
			}

			prev, err := s.pop()
			if err != nil {
				s.push(next)
				continue
			}

			if prev.operator == '+' {
				pprev, _ := s.pop()
				s.push(lexem{value: next.value + pprev.value})
				continue
			}

			if prev.operator == '*' {
				pprev, _ := s.pop()
				s.push(lexem{value: next.value * pprev.value})
				continue
			}
		}
		if next.operator == '+' || next.operator == '*' || next.operator == '(' {
			s.push(next)
			continue
		}
		if next.operator == ')' {
			curr, err := s.pop()
			if err != nil {
				panic(err)
			}
			_, _ = s.pop()
			op, _ := s.peek()
			if op.operator == '*' || op.operator == '+' {
				op, _ := s.pop()
				v2, _ := s.pop()

				if op.operator == '*' {
					s.push(lexem{value: curr.value * v2.value})
					continue
				}
				if op.operator == '+' {
					s.push(lexem{value: curr.value + v2.value})
					continue
				}
			} else {
				s.push(curr)
			}
		}
	}
	r, _ := s.pop()
	if !s.isEmpty() {
		fmt.Println(lexems)
	}

	return r.value
}

func homework(filename string) uint64 {
	lines := readLines(filename)
	var sum uint64
	for _, line := range lines {
		ls := lexer(line)
		s := evaluate(ls)
		sum += uint64(s)
	}

	return sum
}

// advanced math -> use two stacks!
func evaluate2(lx []lexem) int {
	terms := stack{}
	op := stack{}

	for _, c := range lx {
		switch c.operator {
		case 0:
			terms.push(c)
		case '(':
			op.push(c)
		case ')':
			prev, _ := op.peek()
			for prev.operator != '(' && !op.isEmpty() {
				t1, _ := terms.pop()
				t2, _ := terms.pop()
				operand, _ := op.pop()
				if operand.operator == '+' {
					terms.push(lexem{value: t1.value + t2.value})
				}
				if operand.operator == '*' {
					terms.push(lexem{value: t1.value * t2.value})
				}
				prev, _ = op.peek()
			}
			op.pop()
		case '+':
			op.push(c)
		case '*':
			prev, _ := op.peek()
			for !op.isEmpty() && prev.operator == '+' {
				t1, _ := terms.pop()
				t2, _ := terms.pop()
				_, _ = op.pop()
				terms.push(lexem{value: t1.value + t2.value})
				prev, _ = op.peek()
			}
			op.push(c)
		}
	}

	for !op.isEmpty() {
		t1, _ := terms.pop()
		t2, _ := terms.pop()
		operand, _ := op.pop()
		if operand.operator == '+' {
			terms.push(lexem{value: t1.value + t2.value})
		}
		if operand.operator == '*' {
			terms.push(lexem{value: t1.value * t2.value})
		}
	}

	r, _ := terms.pop()
	return r.value
}

func homework2(filename string) uint64 {
	lines := readLines(filename)
	var sum uint64
	for _, line := range lines {
		ls := lexer(line)
		s := evaluate2(ls)
		sum += uint64(s)
	}

	return sum
}

// problem 19

type term struct {
	final bool
	value string
	expr  [][]int
	id    int
}

func (t *term) terminate(m map[int]*term) string {
	if t.final {
		return t.value
	}

	if t.id == 8 {
		t42 := m[42]
		return t42.terminate(m) + "+"
	}

	if t.id == 11 {
		t42 := m[42]
		t31 := m[31]
		return "((" + t42.terminate(m) + "){x}(" + t31.terminate(m) + "){x})"
	}

	var buf bytes.Buffer
	buf.WriteString("(")
	for i := 0; i < len(t.expr); i++ {
		if i > 0 {
			buf.WriteString("|")
		}
		for _, tid := range t.expr[i] {
			tt := m[tid]
			buf.WriteString(tt.terminate(m))
		}
	}
	buf.WriteString(")")

	return buf.String()
}

func parseGramar(filename string) (map[int]*term, []string) {
	mode := "term"
	lines := readLines(filename)
	var xs []string
	m := make(map[int]*term)

	for _, line := range lines {
		if line == "" {
			mode = "string"
		}
		switch mode {
		case "term":
			s1 := strings.Split(line, ": ")
			id, _ := strconv.Atoi(s1[0])

			if strings.Contains(s1[1], "\"") {
				f := strings.Trim(s1[1], "\"")
				m[id] = &term{true, f, nil, id}
				continue
			}

			s2 := strings.Split(s1[1], " | ")
			var expr [][]int
			for _, s3 := range s2 {
				s4 := strings.Split(s3, " ")
				var e1 []int
				for _, s5 := range s4 {
					v, _ := strconv.Atoi(s5)
					e1 = append(e1, v)
				}
				expr = append(expr, e1)
			}
			m[id] = &term{false, "", expr, id}
		case "string":
			xs = append(xs, line)
		}
	}

	return m, xs
}

func elvishImage(m map[int]*term, xs []string) int {
	t1 := m[0]
	rx := t1.terminate(m)
	rx = strings.ReplaceAll(rx, "x", "5")
	r := regexp.MustCompile("^" + rx + "$")

	var cnt int
	for _, s := range xs {
		if r.MatchString(s) {
			cnt++
		}
	}

	return cnt
}

func elvishImageX(m map[int]*term, xs []string, x int) int {
	t1 := m[0]
	rx := t1.terminate(m)
	xx := strconv.Itoa(x)
	rx = strings.ReplaceAll(rx, "x", xx)
	r := regexp.MustCompile("^" + rx + "$")

	var cnt int
	for _, s := range xs {
		if r.MatchString(s) {
			cnt++
		}
	}

	return cnt
}

// problem 20

type square struct {
	i  []string
	id int
}

func (s square) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("id: %d\n", s.id))
	for _, l := range s.i {
		b.WriteString(fmt.Sprintf("%s\n", l))
	}

	return b.String()
}

func (s square) lines() []string {
	var l bytes.Buffer
	var r bytes.Buffer

	for i := 0; i < 10; i++ {
		l.WriteByte(s.i[i][0])
		r.WriteByte(s.i[i][9])
	}

	return []string{
		s.i[0],
		reverse(s.i[0]),
		s.i[9],
		reverse(s.i[9]),
		l.String(),
		reverse(l.String()),
		r.String(),
		reverse(r.String()),
	}
}

func reverse(s string) string {
	var b bytes.Buffer
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}

	return b.String()
}

func parseSquares(filename string) []square {
	var sq []square
	lines := readLines(filename)
	var curr square
	for _, line := range lines {
		if strings.Contains(line, "Tile") {
			s := strings.Split(line, " ")
			id, _ := strconv.Atoi(s[1][:4])
			curr.i = nil
			curr.id = id
			continue
		}
		if line == "" {
			sq = append(sq, curr)
			continue
		}
		curr.i = append(curr.i, line)
	}

	return sq
}

func countEdges(sqx []square) map[string][]int {
	m := make(map[string][]int)
	for _, s := range sqx {
		ls := s.lines()
		for _, l := range ls {
			v, _ := m[l]
			v = append(v, s.id)
			m[l] = v
		}
	}

	return m
}

func main() {
	sqx := parseSquares("input/20.txt")
	m := countEdges(sqx)
	cnt := make(map[int]int)
	for _, v := range m {
		if len(v) == 1 {
			cnt[v[0]] = cnt[v[0]] + 1
		}
	}

	res := 1
	for k, v := range cnt {
		if v == 4 {
			res *= k
		}
	}

	fmt.Println(res)
}
