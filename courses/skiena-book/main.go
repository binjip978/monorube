package main

import "fmt"

func main() {
	ints := []interval{
		{1, 2}, {3, 4}, {1, 7}, {3, 5},
	}
	fmt.Println(ints)
	res := nonOverlapingSubset(ints)
	fmt.Println(res)
}
