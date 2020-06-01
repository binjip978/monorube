package main

import (
	"fmt"
	"math"
	"strings"
)

func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	i := 0
	j := len(s) - 1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}

	return true
}

func parseString(s string) (int, error) {
	m := 1
	r := 0
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case '0':
			r += m * 0
		case '1':
			r += m * 1
		case '2':
			r += m * 2
		case '3':
			r += m * 3
		case '4':
			r += m * 4
		case '5':
			r += m * 5
		case '6':
			r += m * 6
		case '7':
			r += m * 7
		case '8':
			r += m * 8
		case '9':
			r += m * 9
		default:
			return 0, fmt.Errorf("can't parse int")
		}
		m *= 10
	}

	return r, nil
}

func parseInt(n int) (string, error) {
	if n == 0 {
		return "0", nil
	}
	n1 := n
	if n < 0 {
		n1 = n1 * -1
	}
	inv := make([]byte, 0)
	for n1 > 0 {
		c := byte(n1%10) + 48
		inv = append(inv, c)
		n1 /= 10
	}
	b := make([]byte, 0)
	for i := len(inv) - 1; i >= 0; i-- {
		b = append(b, inv[i])
	}
	rev := string(b)
	if n < 0 {
		fmt.Println(rev)
		rev = "-" + rev
	}

	return rev, nil
}

// baseConversion 2 <= b1, b2 <= 16
func baseConversion(n string, b1 int, b2 int) string {
	// convert to base10
	power := func(v int, p int) int {
		return int(math.Pow(float64(v), float64(p)))
	}
	var d10 int
	pi := 0
	for i := len(n) - 1; i >= 0; i-- {
		switch n[i] {
		case '0':
			d10 = 0
		case '1':
			d10 += 1 * power(b1, pi)
		case '2':
			d10 += 2 * power(b1, pi)
		case '3':
			d10 += 3 * power(b1, pi)
		case '4':
			d10 += 4 * power(b1, pi)
		case '5':
			d10 += 5 * power(b1, pi)
		case '6':
			d10 += 6 * power(b1, pi)
		case '7':
			d10 += 7 * power(b1, pi)
		case '8':
			d10 += 8 * power(b1, pi)
		case '9':
			d10 += 9 * power(b1, pi)
		case 'A':
			d10 += 10 * power(b1, pi)
		case 'B':
			d10 += 11 * power(b1, pi)
		case 'C':
			d10 += 12 * power(b1, pi)
		case 'D':
			d10 += 13 * power(b1, pi)
		case 'E':
			d10 += 14 * power(b1, pi)
		case 'F':
			d10 += 15 * power(b1, pi)
		default:
			panic("unknown symbols")
		}
		pi++
	}
	if d10 == 0 {
		return "0"
	}
	converted := make([]byte, 0)

	for d10 != 0 {
		c := d10 % b2
		switch c {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
			converted = append(converted, byte(c)+48)
		case 10:
			converted = append(converted, 'A')
		case 11:
			converted = append(converted, 'B')
		case 12:
			converted = append(converted, 'C')
		case 13:
			converted = append(converted, 'D')
		case 14:
			converted = append(converted, 'E')
		case 15:
			converted = append(converted, 'F')
		default:
			panic("d2 is to big")
		}
		d10 /= b2
	}

	reverse := make([]byte, 0)
	for i := len(converted) - 1; i >= 0; i-- {
		reverse = append(reverse, converted[i])
	}

	return string(reverse)
}
