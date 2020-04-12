package main

import "testing"

func TestPrime(t *testing.T) {
	var table = []struct {
		number  uint64
		isPrime bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
		{11, true},
	}

	for _, entry := range table {
		if isPrime(entry.number) != entry.isPrime {
			t.Errorf("case for number: %d is not correct", entry.number)
		}
	}
}
