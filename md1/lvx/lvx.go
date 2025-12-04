package main

import "fmt"

func main() {
	var sum int = 0
	var count int = 0
	var number int
	for {
		fmt.Println("请输入一个整数(输入0结束):")
		fmt.Scanln(&number)
		if number == 0 {
			break
		}
		sum += number
		count++
	}
	average := average(sum, count)
	if average >= 60 {
		fmt.Printf("平均成绩为 %.2f，成绩合格", average)
	} else {
		fmt.Println("谢谢惠顾")
	}
}
func average(sum int, count int) float64 {
	return float64(sum) / float64(count)
}
