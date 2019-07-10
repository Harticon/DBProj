package DBproj

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	UserId    int
	Name      string `json:"name"`
	ExecuteAt int    `json:"executeAt"`
}

type Token struct {
	Token string `json:"token"`
}

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email" gorm:"unique_index;not null" valid:"email"`
	Password  string `json:"password" gorm:"not null"`
}
