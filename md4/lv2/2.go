package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type timestampWriter struct {
	writer io.Writer
}

func newTimestampWriter(w io.Writer) *timestampWriter {
	return &timestampWriter{writer: w}
}
func (tw *timestampWriter) Write(p []byte) (n int, err error) {
	now := time.Now()
	timestamp := now.Format("2006-01-02 15:04:05")
	unixTime := now.Unix()
	message := string(p)
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}
	logEntry := fmt.Sprintf("[%s] [%d] %s\n", timestamp, unixTime, message)
	return tw.writer.Write([]byte(logEntry))
}
func main() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("无法打开日志文件: %v\n", err)
		return
	}
	defer file.Close()
	logWriter := newTimestampWriter(io.MultiWriter(file, os.Stdout))
	fmt.Fprintln(logWriter, "用户登录")
	time.Sleep(2 * time.Second)
	fmt.Fprintln(logWriter, "用户执行操作A")
	time.Sleep(1 * time.Second)
	fmt.Fprintln(logWriter, "用户执行操作B")
}
