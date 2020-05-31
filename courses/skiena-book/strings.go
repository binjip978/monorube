package main

import (
	"fmt"
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
