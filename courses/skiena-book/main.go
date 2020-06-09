package main

import "fmt"

func main() {
	l1 := &list{5, nil}
	l1.add(10)
	l1.add(29)
	l1.add(39)
	l1.add(45)

	fmt.Println("initial", l1)
}
