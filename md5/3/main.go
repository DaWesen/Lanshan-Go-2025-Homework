package main

import (
	"bufio"
	"fmt"
	task "md5/3/task1"
	"os"
	"strings"
	"sync"
)

func findFiles(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		path := dir + "/" + entry.Name()
		if entry.IsDir() {
			subFiles, err := findFiles(path)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, path)
		}
	}
	return files, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("用法: %s [目录] [关键词]\n", os.Args[0])
		return
	}
	dir, keyword := os.Args[1], os.Args[2]
	tasks := task.New(10)

	files, err := findFiles(dir)
	if err != nil {
		fmt.Printf("遍历目录出错: %v\n", err)
		return
	}

	var results []string
	var mu sync.Mutex
	for _, file := range files {
		tasks.Submit(func() {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			lineNum := 1
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, keyword) {
					mu.Lock()
					result := fmt.Sprintf("%s:%d: %s", file, lineNum, strings.TrimSpace(line))
					results = append(results, result)
					mu.Unlock()
				}
				lineNum++
			}
		})
	}
	tasks.Wait()
	for _, result := range results {
		fmt.Println(result)
	}
	if len(results) == 0 {
		fmt.Println("未找到匹配内容")
	} else {
		fmt.Printf("共找到 %d 处匹配\n", len(results))
	}
}
