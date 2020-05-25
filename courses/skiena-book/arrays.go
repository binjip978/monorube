package main

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
