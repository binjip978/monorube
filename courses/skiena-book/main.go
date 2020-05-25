package main

import "fmt"

func main() {
	ar := []int{0, 4, 2, 1, 1, 0, 2, 1, 0, 1, 2, 0}
	pivot(ar, 0)
	fmt.Println(ar)
}
