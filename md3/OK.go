package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type User struct {
	Name     string
	Birthday string
	Source   string
	Age      int
	Momotalk int
}

type UserManager struct {
	passedInitialCheck bool
	userMap            map[string]User
}

func (AC UserManager) Adduser() {
	var BA User
	fmt.Printf("请输入学生信息(姓名 生日 来源 年龄 momotalk): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	_, err := fmt.Sscanf(input, "%s %s %s %d %d", &BA.Name, &BA.Birthday, &BA.Source, &BA.Age, &BA.Momotalk)
	if err != nil {
		fmt.Println("输入格式错误:", err)
		return
	}
	if AC.userMap == nil {
		AC.userMap = make(map[string]User)
	}

	AC.userMap[BA.Name] = BA
	fmt.Println("学生添加成功!")
}

func (AC UserManager) UpdateUser(oldName, newName, newBirthday, newSource string, newAge, newMomotalk int) {
	if !checkIQ(AC.passedInitialCheck) {
		return
	}

	BA, ok := AC.userMap[oldName]
	if !ok {
		fmt.Println("学生不存在:", oldName)
		return
	}

	BA.Name = newName
	BA.Birthday = newBirthday
	BA.Source = newSource
	BA.Age = newAge
	BA.Momotalk = newMomotalk

	AC.userMap[newName] = BA
	if newName != oldName {
		delete(AC.userMap, oldName)
	}

	fmt.Println("学生更新成功!")
}

func (AC UserManager) DeleteUser(name string) {
	fmt.Println("为什么要这样呢？")
}

type checker struct {
	IQ int
}

func (c checker) Check() bool {
	fmt.Println("少年，你是否身具荣耀的使命？")
	time.Sleep(1 * time.Second)
	fmt.Print("请输入你的智商（数字）: ")
	_, err := fmt.Scanf("%d", &c.IQ)
	if err != nil {
		fmt.Println("输入错误:", err)
		return false
	}
	if c.IQ >= 250 {
		fmt.Println("什么，竟然是蓝山go组组长!")
		return true
	}
	fmt.Println("蓝莓好像哪里出错了呢，快去看看吧")
	return false
}

func checkIQ(passedInitialCheck bool) bool {
	if passedInitialCheck {
		return true
	}
	fmt.Println("让我见证一下你真实的力量")
	return false
}

func main() {
	var action string
	AC := &UserManager{
		userMap: make(map[string]User),
	}

	fmt.Println("Good Luck!!!")
	c := &checker{}
	if c.Check() {
		AC.passedInitialCheck = true
		fmt.Println("如此看来，汝乃命定之人")
	} else {
		AC.passedInitialCheck = false
		fmt.Println("前面的区域，以后再来探索吧")
		time.Sleep(1 * time.Second)
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	fmt.Println("接下来请选择你的英雄：Add or Update or Delete or Show")
	for {
		fmt.Printf("请输入操作: ")
		action, _ = reader.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "Add":
			if !AC.passedInitialCheck {
				fmt.Println("汝是？")
				continue
			}
			AC.Adduser()
		case "Update":
			if !AC.passedInitialCheck {
				fmt.Println("汝是？")
				continue
			}

			fmt.Printf("请输入要更新的学生姓名: ")
			oldName, _ := reader.ReadString('\n')
			oldName = strings.TrimSpace(oldName)

			fmt.Printf("请输入更新后的学生信息(姓名 生日 来源 年龄 momotalk): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			var newName, newBirthday, newSource string
			var newAge, newMomotalk int
			_, err := fmt.Sscanf(input, "%s %s %s %d %d", &newName, &newBirthday, &newSource, &newAge, &newMomotalk)
			if err != nil {
				fmt.Println("输入格式错误:", err)
				continue
			}
			AC.UpdateUser(oldName, newName, newBirthday, newSource, newAge, newMomotalk)
		case "Delete":
			if !AC.passedInitialCheck {
				fmt.Println("汝是？")
				continue
			}
			fmt.Printf("为什么要这样做呢？")
		case "Show":
			fmt.Println("当前学生列表:")
			for k, v := range AC.userMap {
				fmt.Printf("key=%s value=%+v\n", k, v)
			}
		case "Exit":
			fmt.Println("什亭之箱退出")
			return
		default:
			fmt.Println("未知操作，请输入 Add / Update / Delete / Show / Exit")
		}
	}
}

