package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	write1 := []byte("无缓冲写\n")
	write2 := []byte("有缓冲写\n")
	file1, _ := os.Create("test1.txt")
	defer file1.Close()
	now := time.Now()
	for i := 1; i < 10000; i++ {
		_, _ = file1.Write(write1)
	}
	nobufiotime := time.Since(now)
	file2, _ := os.Create("test2.txt")
	defer file2.Close()
	writer := bufio.NewWriter(file2)
	now = time.Now()
	for i := 1; i < 10000; i++ {
		_, _ = writer.Write(write2)
	}
	writer.Flush()
	bufiotime := time.Since(now)
	fmt.Printf("无缓冲写耗时: %v\n", nobufiotime)
	fmt.Printf("有缓冲写耗时: %v\n", bufiotime)
	if bufiotime < nobufiotime {
		fmt.Printf("有缓冲写比无缓冲写快")
	} else {
		fmt.Printf("无缓冲写比有缓冲写快")
	}
}
