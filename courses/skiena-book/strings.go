package main

import "strings"

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
