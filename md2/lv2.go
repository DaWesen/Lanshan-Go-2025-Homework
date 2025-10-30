package main

import "fmt"

func main() {
	var a, b int
	var fuhao, choice string

	for {
		fmt.Scanf("%d %s %d\n", &a, &fuhao, &b)
		res := calulate(fuhao, a, b)
		fmt.Println("结果:", res)
		fmt.Print("是否继续计算? (输入exit退出): ")
		fmt.Scanln(&choice)
		if choice == "exit" {
			break
		}
	}
	fmt.Println("程序结束")
}
func calulate(fuhao string, a, b int) int {
	switch fuhao {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b == 0 {
			fmt.Println("除数不能为0")
			return 0
		}
		return a / b
	default:
		fmt.Println("不支持的运算符:", fuhao)
		return 0
	}
}
