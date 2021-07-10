package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	A_INSTRUCTION = "A_INSTRUCTION"
	C_INSTRUCTION = "C_INSTRUCTION"
	L_INSTRUCTION = "L_INSTRUCTION"
)

type parser struct {
	lines              []string
	currentLine        int
	currentInstruction *string
	symbolTable        map[string]int
	programLine        int
}

func newParser(fileName string) *parser {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(bufio.NewReader(f))
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	st := map[string]int{
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KDB":    24576,
	}

	return &parser{
		lines:       lines,
		currentLine: -1,
		symbolTable: st,
	}
}

func (p *parser) reset() {
	p.currentInstruction = nil
	p.currentLine = -1
}

func (p *parser) hasMoreLines() bool {
	return p.currentLine < len(p.lines)-1
}

func (p *parser) advance() {
	p.currentLine++
	line := p.lines[p.currentLine]
	line = strings.TrimSpace(line)

	var b bytes.Buffer
	var commentStart bool
	for _, r := range line {
		if unicode.IsSpace(r) {
			commentStart = false
			continue
		}
		if r == '/' && commentStart {
			break
		}
		if r == '/' {
			commentStart = true
			continue
		}

		b.WriteRune(r)
	}

	s := b.String()
	if len(s) > 0 {
		p.currentInstruction = &s
		return
	}

	p.currentInstruction = nil
}

func (p *parser) instructionType() string {
	if p.currentInstruction == nil {
		panic("should not call on empty instruction")
	}

	instr := *p.currentInstruction

	if instr[0] == '@' {
		return A_INSTRUCTION
	}

	if instr[0] == '(' && instr[len(instr)-1] == ')' {
		return L_INSTRUCTION
	}

	if strings.Contains(instr, "=") || strings.Contains(instr, ";") {
		return C_INSTRUCTION
	}

	panic(fmt.Sprintf("parsing if wrong for instruction: %s", instr))
}

func (p *parser) symbol() string {
	instrType := p.instructionType()
	if instrType == C_INSTRUCTION {
		panic("should not be called on C_INSTRUCTION")
	}

	instr := *p.currentInstruction

	if instrType == A_INSTRUCTION {
		return instr[1:]
	}

	return instr[1 : len(instr)-1]
}

func (p *parser) dest() string {
	if p.instructionType() != "C_INSTRUCTION" {
		panic(fmt.Sprintf("should not be called on %s", p.instructionType()))
	}
	sp := strings.Split(*p.currentInstruction, "=")
	if len(sp) == 2 {
		return sp[0]
	}

	return "null"
}

func (p *parser) comp() string {
	if p.instructionType() != "C_INSTRUCTION" {
		panic(fmt.Sprintf("should not be called on %s", p.instructionType()))
	}

	sp := strings.Split(*p.currentInstruction, "=")
	if len(sp) == 2 {
		return sp[1]
	}

	sp = strings.Split(*p.currentInstruction, ";")
	if len(sp) == 2 {
		return sp[0]
	}

	panic(fmt.Sprintf("comp could not be parsed %s", *p.currentInstruction))
}

func (p *parser) jump() string {
	if p.instructionType() != "C_INSTRUCTION" {
		panic(fmt.Sprintf("should not be called on %s", p.instructionType()))
	}

	sp := strings.Split(*p.currentInstruction, ";")
	if len(sp) == 2 {
		return sp[1]
	}

	return "null"
}

func (p *parser) parse() string {
	// first run
	programLine := -1
	for p.hasMoreLines() {
		p.advance()
		if p.currentInstruction == nil {
			continue
		}

		switch p.instructionType() {
		case C_INSTRUCTION:
			programLine++
		case A_INSTRUCTION:
			programLine++
		case L_INSTRUCTION:
			s := p.symbol()
			p.symbolTable[s] = programLine + 1
		default:
			panic("can't be true")
		}
	}

	// second run
	p.reset()

	var res []string
	nextFreeAddress := 16

	for p.hasMoreLines() {
		p.advance()
		if p.currentInstruction == nil {
			continue
		}
		switch p.instructionType() {
		case C_INSTRUCTION:
			var b bytes.Buffer
			b.WriteString("111")
			switch p.comp() {
			case "0":
				b.WriteString("0101010")
			case "1":
				b.WriteString("0111111")
			case "-1":
				b.WriteString("0111010")
			case "D":
				b.WriteString("0001100")
			case "A":
				b.WriteString("0110000")
			case "M":
				b.WriteString("1110000")
			case "!D":
				b.WriteString("0001101")
			case "!A":
				b.WriteString("0110001")
			case "!M":
				b.WriteString("1110001")
			case "-D":
				b.WriteString("0001111")
			case "-A":
				b.WriteString("0110011")
			case "-M":
				b.WriteString("1110001")
			case "D+1":
				b.WriteString("0011111")
			case "A+1":
				b.WriteString("0110111")
			case "M+1":
				b.WriteString("1110111")
			case "D-1":
				b.WriteString("0001110")
			case "A-1":
				b.WriteString("0110010")
			case "M-1":
				b.WriteString("1110010")
			case "D+A":
				b.WriteString("0000010")
			case "D+M":
				b.WriteString("1000010")
			case "D-A":
				b.WriteString("0010011")
			case "D-M":
				b.WriteString("1010011")
			case "A-D":
				b.WriteString("0000111")
			case "M-D":
				b.WriteString("1000111")
			case "D&A":
				b.WriteString("0000000")
			case "D&M":
				b.WriteString("1000000")
			case "D|A":
				b.WriteString("0010101")
			case "D|M":
				b.WriteString("1010101")
			}

			switch p.dest() {
			case "null":
				b.WriteString("000")
			case "M":
				b.WriteString("001")
			case "D":
				b.WriteString("010")
			case "DM", "MD":
				b.WriteString("011")
			case "A":
				b.WriteString("100")
			case "AM":
				b.WriteString("101")
			case "AD":
				b.WriteString("110")
			case "ADM":
				b.WriteString("111")
			default:
				panic(fmt.Sprintf("%s", *p.currentInstruction))
			}

			switch p.jump() {
			case "null":
				b.WriteString("000")
			case "JGT":
				b.WriteString("001")
			case "JEQ":
				b.WriteString("010")
			case "JGE":
				b.WriteString("011")
			case "JLT":
				b.WriteString("100")
			case "JNE":
				b.WriteString("101")
			case "JLE":
				b.WriteString("110")
			case "JMP":
				b.WriteString("111")
			default:
				panic(fmt.Sprintf("%s", *p.currentInstruction))
			}

			res = append(res, b.String())

		case A_INSTRUCTION:
			sym := p.symbol()
			i, err := strconv.Atoi(sym)
			if err == nil {
				res = append(res, aInstr(i))
				continue
			}

			v, ok := p.symbolTable[sym]
			if !ok {
				// new variable
				p.symbolTable[sym] = nextFreeAddress
				res = append(res, aInstr(nextFreeAddress))
				nextFreeAddress++
				continue
			}
			res = append(res, aInstr(v))

		case L_INSTRUCTION:
		}

	}

	var b bytes.Buffer

	for _, i := range res {
		b.WriteString(i)
		b.WriteString("\n")
	}

	b.Truncate(b.Len() - 1)
	return b.String()
}

func aInstr(value int) string {
	bits := fmt.Sprintf("%b", value)
	prep := 16 - len(bits)
	var pad bytes.Buffer
	for i := 0; i < prep; i++ {
		pad.WriteString("0")
	}

	return pad.String() + bits
}

func main() {
	if len(os.Args) < 2 {
		panic("please provide asm program as an argument")
	}

	p := newParser(os.Args[1])
	program := p.parse()

	sp := strings.Split(os.Args[1], ".")
	output, err := os.OpenFile(fmt.Sprintf("%s.hack", sp[0]), os.O_CREATE|os.O_RDWR, 0775)

	if err != nil {
		panic(err)
	}

	_, _ = output.WriteString(program)
}
