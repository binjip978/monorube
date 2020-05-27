package main

import "math"

func pivot(array []int, index int) {
	array[0], array[index] = array[index], array[0]

	pivot := array[0]
	j := 0
	for i := 1; i < len(array); i++ {
		if array[i] <= pivot {
			j++
			array[i], array[j] = array[j], array[i]
		}
	}

	array[0], array[j] = array[j], array[0]
}

// dutchFlag should order array with elements 0, 1, 2
func dutchFlag(array []int) {
	j := 0
	for i := 0; i < len(array); i++ {
		if array[i] == 0 {
			array[i], array[j] = array[j], array[i]
			j++
		}
	}
	for i := j; i < len(array); i++ {
		if array[i] == 1 {
			array[i], array[j] = array[j], array[i]
			j++
		}
	}
}

func advance(array []int) bool {
	return false
}

func deleteDubsFromSortedArray1(array []int) []int {
	if len(array) <= 1 {
		return array
	}

	var newArray []int
	newArray = append(newArray, array[0])
	prev := newArray[0]

	for i := 1; i < len(array); i++ {
		if array[i] != prev {
			newArray = append(newArray, array[i])
			prev = array[i]
		}
	}

	return newArray
}

// deleteDubsFromSortedArray2 returns number of valid elements
func deleteDubsFromSortedArray2(array []int) int {
	writeIndex := 1
	for i := 0; i < len(array); i++ {
		if array[i] != array[writeIndex-1] {
			array[writeIndex] = array[i]
			writeIndex++
		}
	}

	return writeIndex
}

func sellStockOnce(price []int) int {
	min := price[0]
	diff := 0

	for i := 1; i < len(price); i++ {
		cd := price[i] - min
		if cd > diff {
			diff = cd
		}
	}

	return diff
}

func primes(n int) []int {
	var p = func(v int) bool {
		for i := 2; i <= int(math.Sqrt(float64(v))); i++ {
			if v%i == 0 {
				return false
			}
		}

		return true
	}

	var res []int
	for i := 1; i <= n; i++ {
		if p(i) {
			res = append(res, i)
		}
	}

	return res
}

func incInteger(integer []int) []int {
	carry := 1
	for i := len(integer) - 1; i >= 0; i-- {
		res := carry + integer[i]
		if res <= 9 {
			integer[i] = res
			return integer
		}
		integer[i] = 0
	}

	// need carry on forward
	newInteger := []int{1}
	newInteger = append(newInteger, integer...)
	return newInteger
}
