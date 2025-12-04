package sql

import (
	"fmt"
	"log"
	"md7/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func InitDatabase() (*Database, error) {
	var db *gorm.DB
	var err error
	dsn := "root:*******@tcp(127.0.0.1:3306)/MD7?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("mysqlwrong%v", err)
		db, err = gorm.Open(sqlite.Open("student.db"), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("sqlitewrong%v", err)
		}
		log.Println("sqlite success")
	} else {
		log.Println("mysql success")
	}
	if err := model.AutoMigrate(db); err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}
