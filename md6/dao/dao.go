package dao

import (
	"encoding/json"
	"md6/modle"
	"os"
	"sync"
)

var (
	students = make(map[string]modle.Student)
	Archive  = "Archive.json"
	mu       sync.RWMutex
)

func Init() {
	Loadstudent()
}
func Loadstudent() {
	mu.Lock()
	defer mu.Unlock()
	file, err := os.Open(Archive)
	if err != nil {
		if os.IsNotExist(err) {
			Writeinarchive()
			return
		}
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&students)
	if err != nil {
		students = make(map[string]modle.Student)
	}
}
func Writeinarchive() {
	file, err := os.Create(Archive)
	if err != nil {
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(students)
}
func Addstudent(name, password, school string) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := students[name]; exists {
		return false
	}
	students[name] = modle.Student{
		Name:     name,
		Password: password,
		School:   school,
	}
	Writeinarchive()
	return true
}
func Givestudent(name string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := students[name]
	return exists
}
func Givepassword(name string) string {
	mu.Lock()
	defer mu.Unlock()
	if student, exists := students[name]; exists {
		return student.Password
	}
	return ""
}
func UpdatePassword(name, newPassword string) bool {
	mu.Lock()
	defer mu.Unlock()

	if student, exists := students[name]; exists {
		students[name] = modle.Student{
			Name:     name,
			Password: newPassword,
			School:   student.School,
		}
		Writeinarchive()
		return true
	}
	return false
}
func Getallstudents() map[string]modle.Student {
	mu.Lock()
	defer mu.Unlock()
	res := make(map[string]modle.Student)
	for k, v := range students {
		res[k] = v
	}
	return res
}
