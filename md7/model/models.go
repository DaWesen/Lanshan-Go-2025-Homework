package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Student struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	School    string    `json:"school"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Tasks     []Task    `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE;" json:"tasks"`
}

type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Status    bool      `gorm:"default:false" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	StudentID uint      `gorm:"not null" json:"student_id"`
}

func (s *Student) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	s.Password = string(hashedPassword)
	return nil
}
func (s *Student) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
	return err == nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Student{}, &Task{})
}
