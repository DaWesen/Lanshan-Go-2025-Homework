package main

import "fmt"

func main() {
	const pi = 3.14
	var r int = 5
	r = 5
	area := pi * float64(r*r)
	fmt.Println(area)
}
