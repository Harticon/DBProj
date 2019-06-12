package DBproj

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email" gorm:"unique_index;not null"`
	Password  string `json:"password" gorm:"not null"`
}

type Task struct {
	gorm.Model
	UserId    int
	Name      string `json:"name"`
	ExecuteAt string `json:"executeAt"`
}
