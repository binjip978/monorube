package main

import "fmt"

func main() {
	r, err := substringIndex("hello worl", "world")
	fmt.Println(r, err)
}
