package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
	}
	res := hanshu(a)
	fmt.Println(res)
}
func hanshu(a []int) map[int]int {
	res := map[int]int{}
	for _, v := range a {
		res[v]++
	}
	return res
}
