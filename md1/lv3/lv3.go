package main

import "fmt"

func main() {
	var n int
	var factorial int = 1
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		factorial *= i
	}
	fmt.Println(factorial)
}
