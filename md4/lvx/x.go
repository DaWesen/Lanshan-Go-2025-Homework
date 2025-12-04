package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var sourceDir, targetDir string
	fmt.Print("源目录: ")
	fmt.Scan(&sourceDir)
	fmt.Print("目标目录: ")
	fmt.Scan(&targetDir)
	if checkDir(sourceDir) != nil {
		fmt.Println("源目录不存在")
		return
	}
	os.MkdirAll(targetDir, 0755)
	fmt.Println("开始同步...")
	for {
		syncFiles(sourceDir, targetDir)
		time.Sleep(2 * time.Second)
	}
}
func checkDir(path string) error {
	_, err := os.Stat(path)
	return err
}
func syncFiles(source, target string) {
	sourceFiles, err := os.ReadDir(source)
	if err != nil {
		fmt.Println("读取源目录失败:", err)
		return
	}
	targetFiles, err := os.ReadDir(target)
	if err != nil {
		fmt.Println("读取目标目录失败:", err)
		return
	}
	hasOperation := false
	for _, targetFile := range targetFiles {
		found := false
		for _, sourceFile := range sourceFiles {
			if sourceFile.Name() == targetFile.Name() {
				found = true
				break
			}
		}
		if !found && !targetFile.IsDir() {
			filePath := filepath.Join(target, targetFile.Name())
			err := os.Remove(filePath)
			if err == nil {
				fmt.Println("删除:", targetFile.Name())
				hasOperation = true
			}
		}
	}
	for _, sourceFile := range sourceFiles {
		if sourceFile.IsDir() {
			continue
		}
		sourcePath := filepath.Join(source, sourceFile.Name())
		targetPath := filepath.Join(target, sourceFile.Name())

		if needSync(sourcePath, targetPath) {
			if copyFile(sourcePath, targetPath) == nil {
				fmt.Println("同步:", sourceFile.Name())
				hasOperation = true
			}
		}
	}
	if !hasOperation {
		fmt.Println("等待文件变化...")
	}
}
func needSync(source, target string) bool {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false
	}
	targetInfo, err := os.Stat(target)
	if err != nil {
		return true
	}
	if sourceInfo.ModTime().After(targetInfo.ModTime()) || sourceInfo.Size() != targetInfo.Size() {
		return true
	}
	return false
}
func copyFile(source, target string) error {
	content, err := os.ReadFile(source)
	if err != nil {
		fmt.Println("读取文件失败:", source, err)
		return err
	}
	err = os.WriteFile(target, content, 0644)
	if err != nil {
		fmt.Println("写入文件失败:", target, err)
		return err
	}
	return nil
}
