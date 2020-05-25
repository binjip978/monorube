package main

import "sort"

func insertionSort(items []int) {
	for i := 1; i < len(items); i++ {
		j := i
		for j > 0 && items[j] < items[j-1] {
			items[j], items[j-1] = items[j-1], items[j]
			j--
		}
	}
}

type interval struct {
	start int
	end   int
}

func nonOverlapingSubset(intervals []interval) []interval {
	if len(intervals) <= 1 {
		return intervals
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].end < intervals[j].end
	})

	res := []interval{intervals[0]}
	currentEnd := res[0].end

	for i := 1; i < len(intervals); i++ {
		if intervals[i].start > currentEnd {
			res = append(res, intervals[i])
			currentEnd = intervals[i].end
		}
	}

	return res
}

func substring(s string, p string) bool {
	for i := 0; i < len(s)-len(p)+1; i++ {
		if s[i:i+len(p)] == p {
			return true
		}
	}

	return false
}

func power(a int, n int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return a
	}
	k := n / 2
	if n%2 == 1 {
		return a * power(a*a, k)
	}
	return power(a*a, k)
}
